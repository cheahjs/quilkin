/*
 * Copyright 2020 Google LLC
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

use std::sync::Arc;

use prometheus::{exponential_buckets, Error as PrometheusError, Histogram, Registry};

use crate::config::{Filter as FilterConfig, ValidationError};
use crate::filters::{prelude::*, FilterRegistry};
use crate::metrics::{histogram_opts, CollectorExt};

const FILTER_LABEL: &str = "filter";

/// A chain of [`Filter`]s to be executed in order.
///
/// Executes each filter, passing the [`ReadContext`] and [`WriteContext`]
/// between each filter's execution, returning the result of data that has gone
/// through all of the filters in the chain. If any of the filters in the chain
/// return `None`, then the chain is broken, and `None` is returned.
pub struct FilterChain {
    filters: Vec<(String, FilterInstance)>,
    filter_read_duration_seconds: Vec<Histogram>,
    filter_write_duration_seconds: Vec<Histogram>,
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("{}", .0)]
    Prometheus(PrometheusError),
    #[error("failed to create filter {}: {}", filter_name, error)]
    Filter {
        filter_name: String,
        error: ValidationError,
    },
}

impl From<PrometheusError> for Error {
    fn from(error: PrometheusError) -> Self {
        Self::Prometheus(error)
    }
}

impl FilterChain {
    pub fn new(filters: Vec<(String, FilterInstance)>, registry: &Registry) -> Result<Self, Error> {
        let subsystem = "filter";

        Ok(Self {
            filter_read_duration_seconds: filters
                .iter()
                .map(|(name, _)| {
                    Histogram::with_opts(
                        histogram_opts(
                            "read_duration_seconds",
                            subsystem,
                            "Seconds taken to execute a given filter's `read`.",
                            Some(exponential_buckets(0.000125, 2.5, 11).unwrap()),
                        )
                        .const_label(FILTER_LABEL, name),
                    )
                    .and_then(|histogram| histogram.register_if_not_exists(registry))
                })
                .collect::<Result<_, prometheus::Error>>()?,
            filter_write_duration_seconds: filters
                .iter()
                .map(|(name, _)| {
                    Histogram::with_opts(
                        histogram_opts(
                            "write_duration_seconds",
                            subsystem,
                            "Seconds taken to execute a given filter's `write`.",
                            Some(exponential_buckets(0.000125, 2.5, 11).unwrap()),
                        )
                        .const_label(FILTER_LABEL, name),
                    )
                    .and_then(|histogram| histogram.register_if_not_exists(registry))
                })
                .collect::<Result<_, prometheus::Error>>()?,
            filters,
        })
    }

    /// Validates the filter configurations in the provided config and constructs
    /// a FilterChain if all configurations are valid.
    pub fn try_create(
        filter_configs: Vec<FilterConfig>,
        filter_registry: &FilterRegistry,
        metrics_registry: &Registry,
    ) -> Result<Self, Error> {
        let mut filters = Vec::new();

        for filter_config in filter_configs {
            match filter_registry.get(
                &filter_config.name,
                CreateFilterArgs::fixed(metrics_registry.clone(), filter_config.config.as_ref())
                    .with_metrics_registry(metrics_registry.clone()),
            ) {
                Ok(filter) => filters.push((filter_config.name, filter)),
                Err(err) => {
                    return Err(Error::Filter {
                        filter_name: filter_config.name.clone(),
                        error: err.into(),
                    });
                }
            }
        }

        FilterChain::new(filters, metrics_registry)
    }

    /// Returns an iterator over the current filters' configs.
    pub(crate) fn get_configs(&self) -> impl Iterator<Item = (&str, Arc<serde_json::Value>)> {
        self.filters
            .iter()
            .map(|(config_json, config)| (config_json.as_str(), config.config.clone()))
    }
}

impl Filter for FilterChain {
    fn read(&self, ctx: ReadContext) -> Option<ReadResponse> {
        self.filters
            .iter()
            .zip(self.filter_read_duration_seconds.iter())
            .try_fold(ctx, |ctx, ((_, instance), histogram)| {
                Some(ReadContext::with_response(
                    ctx.from.clone(),
                    histogram.observe_closure_duration(|| instance.filter.read(ctx))?,
                ))
            })
            .map(ReadResponse::from)
    }

    fn write(&self, ctx: WriteContext) -> Option<WriteResponse> {
        self.filters
            .iter()
            .rev()
            .zip(self.filter_write_duration_seconds.iter().rev())
            .try_fold(ctx, |ctx, ((_, instance), histogram)| {
                Some(WriteContext::with_response(
                    ctx.endpoint,
                    ctx.from.clone(),
                    ctx.to.clone(),
                    histogram.observe_closure_duration(|| instance.filter.write(ctx))?,
                ))
            })
            .map(WriteResponse::from)
    }
}

#[cfg(test)]
mod tests {
    use std::str::from_utf8;

    use crate::{
        config,
        endpoint::{Endpoint, Endpoints, UpstreamEndpoints},
        filters::{debug, FilterRegistry, FilterSet},
        test_utils::{new_test_chain, TestFilterFactory},
    };

    use super::*;

    #[test]
    fn from_config() {
        let provider = debug::factory();

        // everything is fine
        let filter_configs = vec![config::Filter {
            name: provider.name().into(),
            config: Default::default(),
        }];

        let registry = FilterRegistry::new(FilterSet::default());
        let chain =
            FilterChain::try_create(filter_configs, &registry, &Registry::default()).unwrap();
        assert_eq!(1, chain.filters.len());

        // uh oh, something went wrong
        let filter_configs = vec![config::Filter {
            name: "this is so wrong".into(),
            config: Default::default(),
        }];
        let result = FilterChain::try_create(filter_configs, &registry, &Registry::default());
        assert!(result.is_err());
    }

    fn endpoints() -> Vec<Endpoint> {
        vec![
            Endpoint::new("127.0.0.1:80".parse().unwrap()),
            Endpoint::new("127.0.0.1:90".parse().unwrap()),
        ]
    }

    fn upstream_endpoints(endpoints: Vec<Endpoint>) -> UpstreamEndpoints {
        Endpoints::new(endpoints).unwrap().into()
    }

    #[test]
    fn chain_single_test_filter() {
        let registry = prometheus::Registry::default();
        let chain = new_test_chain(&registry);
        let endpoints_fixture = endpoints();

        let response = chain
            .read(ReadContext::new(
                upstream_endpoints(endpoints_fixture.clone()),
                "127.0.0.1:70".parse().unwrap(),
                b"hello".to_vec(),
            ))
            .unwrap();

        let expected = endpoints_fixture.clone();
        assert_eq!(
            expected,
            response.endpoints.iter().cloned().collect::<Vec<_>>()
        );
        assert_eq!(
            "hello:odr:127.0.0.1:70",
            from_utf8(response.contents.as_slice()).unwrap()
        );
        assert_eq!(
            "receive",
            response.metadata[&"downstream".to_string()]
                .as_string()
                .unwrap()
        );

        let response = chain
            .write(WriteContext::new(
                &endpoints_fixture[0],
                endpoints_fixture[0].address.clone(),
                "127.0.0.1:70".parse().unwrap(),
                b"hello".to_vec(),
            ))
            .unwrap();

        assert_eq!(
            "receive",
            response.metadata[&"upstream".to_string()]
                .as_string()
                .unwrap()
        );
        assert_eq!(
            "hello:our:127.0.0.1:80:127.0.0.1:70",
            from_utf8(response.contents.as_slice()).unwrap()
        );
    }

    #[test]
    fn chain_double_test_filter() {
        let registry = prometheus::Registry::default();
        let chain = FilterChain::new(
            vec![
                (
                    "TestFilter".into(),
                    TestFilterFactory::create_empty_filter(),
                ),
                (
                    "TestFilter".into(),
                    TestFilterFactory::create_empty_filter(),
                ),
            ],
            &registry,
        )
        .unwrap();

        let endpoints_fixture = endpoints();

        let response = chain
            .read(ReadContext::new(
                upstream_endpoints(endpoints_fixture.clone()),
                "127.0.0.1:70".parse().unwrap(),
                b"hello".to_vec(),
            ))
            .unwrap();

        let expected = endpoints_fixture.clone();
        assert_eq!(
            expected,
            response.endpoints.iter().cloned().collect::<Vec<_>>()
        );
        assert_eq!(
            "hello:odr:127.0.0.1:70:odr:127.0.0.1:70",
            from_utf8(response.contents.as_slice()).unwrap()
        );
        assert_eq!(
            "receive:receive",
            response.metadata[&"downstream".to_string()]
                .as_string()
                .unwrap()
        );

        let response = chain
            .write(WriteContext::new(
                &endpoints_fixture[0],
                endpoints_fixture[0].address.clone(),
                "127.0.0.1:70".parse().unwrap(),
                b"hello".to_vec(),
            ))
            .unwrap();
        assert_eq!(
            "hello:our:127.0.0.1:80:127.0.0.1:70:our:127.0.0.1:80:127.0.0.1:70",
            from_utf8(response.contents.as_slice()).unwrap()
        );
        assert_eq!(
            "receive:receive",
            response.metadata[&"upstream".to_string()]
                .as_string()
                .unwrap()
        );
    }

    #[test]
    fn get_configs() {
        struct TestFilter2 {}
        impl Filter for TestFilter2 {}

        let registry = prometheus::Registry::default();
        let filter_chain = FilterChain::new(
            vec![
                (
                    "TestFilter".into(),
                    TestFilterFactory::create_empty_filter(),
                ),
                (
                    "TestFilter2".into(),
                    FilterInstance {
                        config: Arc::new(serde_json::json!({
                            "k1": "v1",
                            "k2": 2
                        })),
                        filter: Box::new(TestFilter2 {}),
                    },
                ),
            ],
            &registry,
        )
        .unwrap();

        let configs = filter_chain.get_configs().collect::<Vec<_>>();
        assert_eq!(
            vec![
                ("TestFilter", Arc::new(serde_json::Value::Null)),
                (
                    "TestFilter2",
                    Arc::new(serde_json::json!({
                        "k1": "v1",
                        "k2": 2
                    }))
                )
            ],
            configs
        )
    }
}
