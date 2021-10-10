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

mod capture;
mod config;
mod metrics;
mod proto;

use std::sync::Arc;

use slog::o;

use crate::filters::prelude::*;

use capture::Capture;
use metrics::Metrics;
use proto::quilkin::extensions::filters::capture_bytes::v1alpha1::CaptureBytes as ProtoConfig;

use crate::log::SharedLogger;
use crate::warn;
pub use config::{Config, Strategy};

pub const NAME: &str = "quilkin.extensions.filters.capture_bytes.v1alpha1.CaptureBytes";

/// Creates a new factory for generating capture filters.
pub fn factory(base: &SharedLogger) -> DynFilterFactory {
    Box::from(CaptureBytesFactory::new(base))
}

struct CaptureBytes {
    log: SharedLogger,
    capture: Box<dyn Capture + Sync + Send>,
    /// metrics reporter for this filter.
    metrics: Metrics,
    metadata_key: Arc<String>,
    size: usize,
    remove: bool,
}

impl CaptureBytes {
    fn new(base: &SharedLogger, config: Config, metrics: Metrics) -> Self {
        CaptureBytes {
            log: base.child(o!("source" => "extensions::CaptureBytes")),
            capture: config.strategy.as_capture(),
            metrics,
            metadata_key: Arc::new(config.metadata_key),
            size: config.size,
            remove: config.remove,
        }
    }
}

impl Filter for CaptureBytes {
    fn read(&self, mut ctx: ReadContext) -> Option<ReadResponse> {
        // if the capture size is bigger than the packet size, then we drop the packet,
        // and occasionally warn
        if ctx.contents.len() < self.size {
            warn!(
                self.log,
                "Packets are being dropped due to their length being less than {} bytes",
                self.size; "count" => self.metrics.packets_dropped_total.get()
            );
            self.metrics.packets_dropped_total.inc();
            return None;
        }
        let token = self
            .capture
            .capture(&mut ctx.contents, self.size, self.remove);

        ctx.metadata
            .insert(self.metadata_key.clone(), Box::new(token));

        Some(ctx.into())
    }
}

struct CaptureBytesFactory {
    log: SharedLogger,
}

impl CaptureBytesFactory {
    pub fn new(base: &SharedLogger) -> Self {
        CaptureBytesFactory { log: base.clone() }
    }
}

impl FilterFactory for CaptureBytesFactory {
    fn name(&self) -> &'static str {
        NAME
    }

    fn create_filter(&self, args: CreateFilterArgs) -> Result<FilterInstance, Error> {
        let (config_json, config) = self
            .require_config(args.config)?
            .deserialize::<Config, ProtoConfig>(self.name())?;
        let filter = CaptureBytes::new(&self.log, config, Metrics::new(&args.metrics_registry)?);
        Ok(FilterInstance::new(
            config_json,
            Box::new(filter) as Box<dyn Filter>,
        ))
    }
}

#[cfg(test)]
mod tests {
    use std::sync::Arc;

    use prometheus::Registry;
    use serde_yaml::{Mapping, Value};

    use crate::endpoint::{Endpoint, Endpoints};
    use crate::test_utils::assert_write_no_change;

    use super::{CaptureBytes, CaptureBytesFactory, Config, Metrics, Strategy};

    use super::capture::{Capture, Prefix, Suffix};

    use crate::filters::{
        metadata::CAPTURED_BYTES, CreateFilterArgs, Filter, FilterFactory, ReadContext,
    };
    use crate::log::test_logger;

    const TOKEN_KEY: &str = "TOKEN";

    fn capture_bytes(config: Config) -> CaptureBytes {
        CaptureBytes::new(
            &test_logger(),
            config,
            Metrics::new(&Registry::default()).unwrap(),
        )
    }

    #[test]
    fn factory_valid_config_all() {
        let factory = CaptureBytesFactory::new(&test_logger());
        let mut map = Mapping::new();
        map.insert(
            Value::String("strategy".into()),
            Value::String("SUFFIX".into()),
        );
        map.insert(
            Value::String("metadataKey".into()),
            Value::String(TOKEN_KEY.into()),
        );
        map.insert(Value::String("size".into()), Value::Number(3.into()));
        map.insert(Value::String("remove".into()), Value::Bool(true));

        let filter = factory
            .create_filter(CreateFilterArgs::fixed(
                Registry::default(),
                Some(&Value::Mapping(map)),
            ))
            .unwrap()
            .filter;
        assert_end_strategy(filter.as_ref(), TOKEN_KEY, true);
    }

    #[test]
    fn factory_valid_config_defaults() {
        let factory = CaptureBytesFactory::new(&test_logger());
        let mut map = Mapping::new();
        map.insert(Value::String("size".into()), Value::Number(3.into()));
        let filter = factory
            .create_filter(CreateFilterArgs::fixed(
                Registry::default(),
                Some(&Value::Mapping(map)),
            ))
            .unwrap()
            .filter;
        assert_end_strategy(filter.as_ref(), CAPTURED_BYTES, false);
    }

    #[test]
    fn factory_invalid_config() {
        let factory = CaptureBytesFactory::new(&test_logger());
        let mut map = Mapping::new();
        map.insert(Value::String("size".into()), Value::String("WRONG".into()));

        let result = factory.create_filter(CreateFilterArgs::fixed(
            Registry::default(),
            Some(&Value::Mapping(map)),
        ));
        assert!(result.is_err(), "Should be an error");
    }

    #[test]
    fn read() {
        let config = Config {
            strategy: Strategy::Suffix,
            metadata_key: TOKEN_KEY.into(),
            size: 3,
            remove: true,
        };
        let filter = capture_bytes(config);
        assert_end_strategy(&filter, TOKEN_KEY, true);
    }

    #[test]
    fn read_overflow_capture_size() {
        let config = Config {
            strategy: Strategy::Suffix,
            metadata_key: TOKEN_KEY.into(),
            size: 99,
            remove: true,
        };
        let filter = capture_bytes(config);
        let endpoints = vec![Endpoint::new("127.0.0.1:81".parse().unwrap())];
        let response = filter.read(ReadContext::new(
            Endpoints::new(endpoints).unwrap().into(),
            "127.0.0.1:80".parse().unwrap(),
            "abc".to_string().into_bytes(),
        ));

        assert!(response.is_none());
        let count = filter.metrics.packets_dropped_total.get();
        assert_eq!(1, count);
    }

    #[test]
    fn write() {
        let config = Config {
            strategy: Strategy::Suffix,
            metadata_key: TOKEN_KEY.into(),
            size: 0,
            remove: false,
        };
        let filter = capture_bytes(config);
        assert_write_no_change(&filter);
    }

    #[test]
    fn end_capture() {
        let end = Suffix {};
        let mut contents = b"helloabc".to_vec();
        let result = end.capture(&mut contents, 3, false);
        assert_eq!(b"abc".to_vec(), result);
        assert_eq!(b"helloabc".to_vec(), contents);

        let result = end.capture(&mut contents, 3, true);
        assert_eq!(b"abc".to_vec(), result);
        assert_eq!(b"hello".to_vec(), contents);
    }

    #[test]
    fn beginning_capture() {
        let beg = Prefix {};
        let mut contents = b"abchello".to_vec();

        let result = beg.capture(&mut contents, 3, false);
        assert_eq!(b"abc".to_vec(), result);
        assert_eq!(b"abchello".to_vec(), contents);

        let result = beg.capture(&mut contents, 3, true);
        assert_eq!(b"abc".to_vec(), result);
        assert_eq!(b"hello".to_vec(), contents);
    }

    fn assert_end_strategy<F>(filter: &F, key: &str, remove: bool)
    where
        F: Filter + ?Sized,
    {
        let endpoints = vec![Endpoint::new("127.0.0.1:81".parse().unwrap())];
        let response = filter
            .read(ReadContext::new(
                Endpoints::new(endpoints).unwrap().into(),
                "127.0.0.1:80".parse().unwrap(),
                "helloabc".to_string().into_bytes(),
            ))
            .unwrap();

        if remove {
            assert_eq!(b"hello".to_vec(), response.contents);
        } else {
            assert_eq!(b"helloabc".to_vec(), response.contents);
        }

        let token = response
            .metadata
            .get(&Arc::new(key.into()))
            .unwrap()
            .downcast_ref::<Vec<u8>>()
            .unwrap();
        assert_eq!(b"abc", token.as_slice());
    }
}
