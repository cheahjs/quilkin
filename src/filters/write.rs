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

use std::collections::HashMap;

use crate::{
    endpoint::{Endpoint, EndpointAddress},
    metadata::DynamicMetadata,
};

#[cfg(doc)]
use crate::filters::Filter;

/// The input arguments to [`Filter::write`].
#[non_exhaustive]
pub struct WriteContext<'a> {
    /// The upstream endpoint that we're expecting packets from.
    pub endpoint: &'a Endpoint,
    /// The source of the received packet.
    pub from: EndpointAddress,
    /// The destination of the received packet.
    pub to: EndpointAddress,
    /// Contents of the received packet.
    pub contents: Vec<u8>,
    /// Arbitrary values that can be passed from one filter to another
    pub metadata: DynamicMetadata,
}

/// The output of [`Filter::write`].
///
/// New instances are created from [`WriteContext`].
///
/// ```rust
/// # use quilkin::filters::{WriteContext, WriteResponse};
///   fn write(ctx: WriteContext) -> Option<WriteResponse> {
///       Some(ctx.into())
///   }
/// ```
#[non_exhaustive]
pub struct WriteResponse {
    /// Contents of the packet to be sent back to the original sender.
    pub contents: Vec<u8>,
    /// Arbitrary values that can be passed from one filter to another.
    pub metadata: DynamicMetadata,
}

impl WriteContext<'_> {
    /// Creates a new [`WriteContext`]
    pub fn new(
        endpoint: &Endpoint,
        from: EndpointAddress,
        to: EndpointAddress,
        contents: Vec<u8>,
    ) -> WriteContext {
        WriteContext {
            endpoint,
            from,
            to,
            contents,
            metadata: HashMap::new(),
        }
    }

    /// Creates a new [`WriteContext`] from a given [`WriteResponse`].
    pub fn with_response(
        endpoint: &Endpoint,
        from: EndpointAddress,
        to: EndpointAddress,
        response: WriteResponse,
    ) -> WriteContext {
        WriteContext {
            endpoint,
            from,
            to,
            contents: response.contents,
            metadata: response.metadata,
        }
    }
}

impl From<WriteContext<'_>> for WriteResponse {
    fn from(ctx: WriteContext) -> Self {
        Self {
            contents: ctx.contents,
            metadata: ctx.metadata,
        }
    }
}
