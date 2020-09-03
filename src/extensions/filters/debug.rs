/*
 * Copyright 2020 Google LLC All Rights Reserved.
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

use std::net::SocketAddr;

use slog::{info, o, Logger};

use crate::config::EndPoint;
use crate::extensions::filter_registry::{CreateFilterArgs, Error, FilterFactory};
use crate::extensions::Filter;

/// Debug Filter logs all incoming and outgoing packets
///
/// # Configuration
///
/// ```yaml
/// local:
///   port: 7000 # the port to receive traffic to locally
/// filters:
///   - name: quilkin.extensions.filters.debug.v1alpha1.DebugFilter
///     config:
///       id: "debug-1"
/// client:
///   addresses:
///     - 127.0.0.1:7001
///   connection_id: 1x7ijy6
/// ```
///  `config.id` (optional) adds a "id" field with a given value to each log line.
///     This can be useful to identify debug log positioning within a filter config if you have
///     multiple DebugFilters configured.
///
pub struct DebugFilter {
    log: Logger,
}

impl DebugFilter {
    /// Constructor for the DebugFilter. Pass in a "id" to append a string to your log messages from this
    /// Filter.
    fn new(base: &Logger, id: Option<String>) -> Self {
        let log = match id {
            None => base.new(o!("source" => "extensions::DebugFilter")),
            Some(id) => base.new(o!("source" => "extensions::DebugFilter", "id" => id)),
        };

        DebugFilter { log }
    }
}

/// Factory for the DebugFilter
pub struct DebugFilterFactory {
    log: Logger,
}

impl DebugFilterFactory {
    pub fn new(base: &Logger) -> Self {
        DebugFilterFactory { log: base.clone() }
    }
}

impl FilterFactory for DebugFilterFactory {
    fn name(&self) -> String {
        "quilkin.extensions.filters.debug.v1alpha1.DebugFilter".into()
    }

    fn create_filter(&self, args: CreateFilterArgs) -> Result<Box<dyn Filter>, Error> {
        // pull out the Option<&Value>
        let prefix = match args.config {
            Some(serde_yaml::Value::Mapping(map)) => map.get(&serde_yaml::Value::from("id")),
            _ => None,
        };

        match prefix {
            // if no config value supplied, then no prefix, which is fine
            None => Ok(Box::new(DebugFilter::new(&self.log, None))),
            // return an Error if the id exists but is not a string.
            Some(value) => match value.as_str() {
                None => Err(Error::FieldInvalid {
                    field: "config.id".to_string(),
                    reason: "id value should be a string".to_string(),
                }),
                Some(prefix) => Ok(Box::new(DebugFilter::new(
                    &self.log,
                    Some(prefix.to_string()),
                ))),
            },
        }
    }
}

impl Filter for DebugFilter {
    fn on_downstream_receive(
        &self,
        endpoints: &[EndPoint],
        from: SocketAddr,
        contents: Vec<u8>,
    ) -> Option<(Vec<EndPoint>, Vec<u8>)> {
        info!(self.log, "on local receive"; "from" => from, "contents" => packet_to_string(contents.clone()));
        Some((endpoints.to_vec(), contents))
    }

    fn on_upstream_receive(
        &self,
        endpoint: &EndPoint,
        from: SocketAddr,
        to: SocketAddr,
        contents: Vec<u8>,
    ) -> Option<Vec<u8>> {
        info!(self.log, "received endpoint packet"; "endpoint" => endpoint.name.clone(),
        "from" => from,
        "to" => to,
        "contents" => packet_to_string(contents.clone()));
        Some(contents)
    }
}

/// packet_to_string takes the content, and attempts to convert it to a string.
/// Returns a string of "error decoding packet" on failure.
fn packet_to_string(contents: Vec<u8>) -> String {
    match String::from_utf8(contents) {
        Ok(str) => str,
        Err(_) => String::from("error decoding packet"),
    }
}

#[cfg(test)]
mod tests {
    use serde_yaml::Mapping;
    use serde_yaml::Value;

    use crate::config::ConnectionConfig::Server;
    use crate::test_utils::{
        assert_filter_on_downstream_receive_no_change, assert_filter_on_upstream_receive_no_change,
        logger,
    };

    use super::*;

    #[test]
    fn on_downstream_receive() {
        let df = DebugFilter::new(&logger(), None);
        assert_filter_on_downstream_receive_no_change(&df);
    }

    #[test]
    fn on_upstream_receive() {
        let df = DebugFilter::new(&logger(), None);
        assert_filter_on_upstream_receive_no_change(&df);
    }

    #[test]
    fn from_config_with_id() {
        let log = logger();
        let mut map = Mapping::new();
        let connection = Server { endpoints: vec![] };
        let factory = DebugFilterFactory::new(&log);

        map.insert(Value::from("id"), Value::from("name"));
        assert!(factory
            .create_filter(CreateFilterArgs::new(
                &connection,
                Some(&Value::Mapping(map)),
            ))
            .is_ok());
    }

    #[test]
    fn from_config_without_id() {
        let log = logger();
        let mut map = Mapping::new();
        let connection = Server { endpoints: vec![] };
        let factory = DebugFilterFactory::new(&log);

        map.insert(Value::from("id"), Value::from("name"));
        assert!(factory
            .create_filter(CreateFilterArgs::new(
                &connection,
                Some(&Value::Mapping(map)),
            ))
            .is_ok());
    }

    #[test]
    fn from_config_should_error() {
        let log = logger();
        let mut map = Mapping::new();
        let connection = Server { endpoints: vec![] };
        let factory = DebugFilterFactory::new(&log);

        map.insert(Value::from("id"), Value::from(false));
        match factory.create_filter(CreateFilterArgs::new(
            &connection,
            Some(&Value::Mapping(map)),
        )) {
            Ok(_) => assert!(false, "should be an error"),
            Err(err) => {
                assert_eq!(
                    Error::FieldInvalid {
                        field: "config.id".to_string(),
                        reason: "id value should be a string".to_string()
                    }
                    .to_string(),
                    err.to_string()
                );
            }
        }
    }
}
