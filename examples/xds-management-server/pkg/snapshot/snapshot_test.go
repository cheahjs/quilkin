package snapshot

import (
	"context"
	"os"
	"testing"
	"time"

	envoylistener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/util/clock"
	"quilkin.dev/xds-management-server/pkg/cluster"
	"quilkin.dev/xds-management-server/pkg/filterchain"
	"quilkin.dev/xds-management-server/pkg/filters"
	debugfilterv1alpha "quilkin.dev/xds-management-server/pkg/filters/debug/v1alpha1"
)

// defaultUpdateInterval is how often to check for updates in tests.
const defaultUpdateInterval = 1 * time.Millisecond

func getDefaultFilterChain(t *testing.T, snapshot cache.Snapshot) types.Resource {
	listeners := snapshot.GetResources(resource.ListenerType)
	require.Len(t, listeners, 1)
	defaultListener, found := listeners[""]
	require.True(t, found)
	return defaultListener
}

func getProxySnapshot(t *testing.T, snapshotCache cache.SnapshotCache, proxyID string, version string) cache.Snapshot {
	snapshot, err := snapshotCache.GetSnapshot(proxyID)
	require.NoError(t, err)
	require.NoError(t, snapshot.Consistent())
	require.EqualValues(t, version, snapshot.GetVersion(resource.ListenerType))
	return snapshot
}

func waitForSnapshotUpdate(t *testing.T, snapshotCache cache.SnapshotCache, proxyIDs []string, version string) {
	require.Eventually(t, func() bool {
		for _, proxyID := range proxyIDs {
			snapshot, err := snapshotCache.GetSnapshot(proxyID)
			if err != nil {
				return false
			}
			if snapshot.GetVersion(resource.ClusterType) != version {
				return false
			}
		}
		return true
	}, 1*time.Second, 1*time.Millisecond)
}

func makeTestUpdater(
	ctx context.Context,
) (
	*Updater,
	chan<- []cluster.Cluster,
	chan<- filterchain.ProxyFilterChain,
	*clock.FakeClock,
) {
	logger := &log.Logger{}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.WarnLevel)

	clusterCh := make(chan []cluster.Cluster)
	filterChainCh := make(chan filterchain.ProxyFilterChain)

	fakeClock := clock.NewFakeClock(time.Now())
	u := NewUpdater(
		logger,
		clusterCh,
		filterChainCh,
		defaultUpdateInterval,
		fakeClock)

	go u.Run(ctx)

	return u, clusterCh, filterChainCh, fakeClock
}

func TestSnapshotUpdaterClusterUpdate(t *testing.T) {
	// Test if we get a cluster update, all proxies are updated.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, clusterCh, filterChainCh, fakeClock := makeTestUpdater(ctx)

	proxyIDs := []string{"proxy-1", "proxy-2", "proxy-3"}
	for _, proxyID := range proxyIDs {
		filterChainCh <- filterchain.ProxyFilterChain{
			ProxyID: proxyID,
			FilterChains: []*envoylistener.FilterChain{{
				FilterChainMatch: &envoylistener.FilterChainMatch{
					ApplicationProtocols: []string{"AA=="},
				},
			}},
		}
	}

	// Wait for the first round of updates.
	fakeClock.Step(defaultUpdateInterval)
	waitForSnapshotUpdate(t, u.snapshotCache, proxyIDs, "1")

	// Send cluster update.
	clusterCh <- []cluster.Cluster{{
		Name: "cluster-1", Endpoints: []cluster.Endpoint{{
			IP:   "127.0.0.2",
			Port: 32,
		}}}}

	// Wait for a v2 snapshot.
	fakeClock.Step(defaultUpdateInterval)
	waitForSnapshotUpdate(t, u.snapshotCache, proxyIDs, "2")

	for _, proxyID := range proxyIDs {
		proxySnapshot, err := u.snapshotCache.GetSnapshot(proxyID)
		require.NoError(t, err)
		require.NoError(t, proxySnapshot.Consistent())
		requireV0FilterChainInSnapshot(t, proxySnapshot)

		clusters := proxySnapshot.GetResources(resource.ClusterType)
		require.Len(t, clusters, 1)
		testCluster, found := clusters["cluster-1"]
		require.True(t, found)
		require.Contains(t, testCluster.String(), "127.0.0.2")
	}
}

func TestSnapshotUpdaterProxyFilterChainUpdates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, _, filterChainCh, fakeClock := makeTestUpdater(ctx)

	proxyIDs := []string{"proxy-1", "proxy-2", "proxy-3"}
	for _, proxyID := range proxyIDs {
		_, err := u.snapshotCache.GetSnapshot(proxyID)
		require.Error(t, err, "found unexpected snapshot for proxy")
	}

	for _, proxyID := range proxyIDs {
		filterChainCh <- filterchain.ProxyFilterChain{
			ProxyID: proxyID,
			FilterChains: []*envoylistener.FilterChain{{
				FilterChainMatch: &envoylistener.FilterChainMatch{
					ApplicationProtocols: []string{"AA=="},
				},
			}},
		}
	}

	fakeClock.Step(defaultUpdateInterval)
	waitForSnapshotUpdate(t, u.snapshotCache, proxyIDs, "1")

	for _, proxyID := range proxyIDs {
		proxySnapshot, err := u.snapshotCache.GetSnapshot(proxyID)
		require.NoError(t, err)
		require.NoError(t, proxySnapshot.Consistent())
		requireV0FilterChainInSnapshot(t, proxySnapshot)

		require.EqualValues(t, "1", proxySnapshot.GetVersion(resource.ClusterType))
		require.EqualValues(t, "1", proxySnapshot.GetVersion(resource.ListenerType))

		clusters := proxySnapshot.GetResources(resource.ClusterType)
		require.Empty(t, clusters)

		listeners := proxySnapshot.GetResources(resource.ListenerType)
		require.Len(t, listeners, 1)
	}
}

func TestSnapshotUpdaterContinuousProxyFilterChainUpdates(t *testing.T) {
	// Test that proxies are updated independently and continuously.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, _, filterChainCh, fakeClock := makeTestUpdater(ctx)

	proxyIDs := []string{"proxy-1", "proxy-2"}

	// Start with empty an filter chain for all proxies.
	for _, proxyID := range proxyIDs {
		filterChainCh <- filterchain.ProxyFilterChain{
			ProxyID: proxyID,
			FilterChains: []*envoylistener.FilterChain{{
				FilterChainMatch: &envoylistener.FilterChainMatch{
					ApplicationProtocols: []string{"AA=="},
				},
			}},
		}
	}

	fakeClock.Step(defaultUpdateInterval)
	waitForSnapshotUpdate(t, u.snapshotCache, proxyIDs, "1")

	// Send an updated filter chain to proxy-2
	debugFilter, err := filterchain.CreateXdsFilter(filters.DebugFilterName,
		&debugfilterv1alpha.Debug{
			Id: &wrapperspb.StringValue{Value: "hello"},
		})
	require.NoError(t, err)

	filterChainCh <- filterchain.ProxyFilterChain{
		ProxyID: "proxy-2",
		FilterChains: []*envoylistener.FilterChain{{
			FilterChainMatch: &envoylistener.FilterChainMatch{
				ApplicationProtocols: []string{"AA=="},
			},
			Filters: []*envoylistener.Filter{debugFilter},
		}},
	}
	fakeClock.Step(defaultUpdateInterval)
	waitForSnapshotUpdate(t, u.snapshotCache, []string{"proxy-2"}, "2")

	// Check that proxy-1 is on v1 while proxy-2 is on v2
	proxy1Snapshot := getProxySnapshot(t, u.snapshotCache, "proxy-1", "1")
	requireV0FilterChainInSnapshot(t, proxy1Snapshot)
	require.EqualValues(t, "filter_chains:{filter_chain_match:{application_protocols:\"AA==\"}}", getDefaultFilterChain(t, proxy1Snapshot).String())

	proxy2Snapshot := getProxySnapshot(t, u.snapshotCache, "proxy-2", "2")
	requireV0FilterChainInSnapshot(t, proxy2Snapshot)
	require.Contains(t, getDefaultFilterChain(t, proxy2Snapshot).String(), "hello")
}

func requireV0FilterChainInSnapshot(t *testing.T, snapshot cache.Snapshot) {
	listener := snapshot.GetResources(resource.ListenerType)
	defaultListener, found := listener[""]
	require.True(t, found)
	require.Contains(t, defaultListener.String(), "filter_chain_match:{application_protocols:\"AA==\"")
}
