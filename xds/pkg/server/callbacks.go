/*
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"fmt"

	discoveryservice "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	log "github.com/sirupsen/logrus"
)

// callbacks implements callbacks for the go-control-plane xds server.
type callbacks struct {
	log      *log.Logger
	nodeIDCh chan<- string
}

func (c *callbacks) OnStreamOpen(
	_ context.Context,
	streamID int64,
	typeURL string,
) error {
	c.log.WithFields(log.Fields{
		"streamId": streamID,
		"type_url": typeURL,
	}).Debugf("OnStreamOpen")
	return nil
}

func (c *callbacks) OnStreamClosed(streamID int64) {
	c.log.WithFields(log.Fields{
		"stream_id": streamID,
	}).Debugf("OnStreamClosed")
}

// OnStreamRequest is called whenever a new DiscoveryRequest is received from a proxy.
//  we use this event to track proxies that are connected to the server.
func (c *callbacks) OnStreamRequest(streamID int64, request *discoveryservice.DiscoveryRequest) error {
	c.log.WithFields(log.Fields{
		"stream_id":            streamID,
		"request_version_info": request.VersionInfo,
		"request_nonce":        request.ResponseNonce,
	}).Debugf("OnStreamRequest")

	if c.nodeIDCh != nil {
		c.nodeIDCh <- request.Node.Id
	}
	return nil
}

func (c *callbacks) OnStreamResponse(
	streamID int64,
	request *discoveryservice.DiscoveryRequest,
	response *discoveryservice.DiscoveryResponse,
) {
	c.log.WithFields(log.Fields{
		"stream_id":             streamID,
		"request_version_info":  request.VersionInfo,
		"request_nonce":         request.ResponseNonce,
		"response_version_info": response.VersionInfo,
		"response_nonce":        response.Nonce,
	}).Debugf("OnStreamResponse")
}

func (c *callbacks) OnFetchRequest(
	context.Context,
	*discoveryservice.DiscoveryRequest,
) error {
	return nil
}

func (c *callbacks) OnFetchResponse(
	*discoveryservice.DiscoveryRequest,
	*discoveryservice.DiscoveryResponse,
) {
}

func (c *callbacks) OnDeltaStreamOpen(context.Context, int64, string) error {
	return fmt.Errorf("delta XDS is not supported")
}

func (c *callbacks) OnDeltaStreamClosed(int64) {
}

func (c *callbacks) OnStreamDeltaRequest(
	int64, *discoveryservice.DeltaDiscoveryRequest,
) error {
	return fmt.Errorf("delta XDS is not supported")
}

func (c *callbacks) OnStreamDeltaResponse(
	int64,
	*discoveryservice.DeltaDiscoveryRequest,
	*discoveryservice.DeltaDiscoveryResponse,
) {
}
