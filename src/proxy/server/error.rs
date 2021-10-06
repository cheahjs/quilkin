/*
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

use crate::proxy::sessions::error::Error as SessionError;

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("failed to startup properly: {0}")]
    Initialize(String),
    #[error("session error: {0}")]
    Session(SessionError),
    #[error("i/o Error: {0}")]
    Bind(std::io::Error),
    #[error("tokio specific error: {0}")]
    Tokio(tokio::io::Error),
    #[error("Error occurred in worker threads: {0}")]
    Join(tokio::task::JoinError),
    #[error("receive loop exited with an error: {0}")]
    RecvLoop(String),
}

macro_rules! from {
    ($from_ty:ty : $($typ:ty => $path:ident);+ $(;)?) => {
        $(
            impl From<$typ> for $from_ty {
                fn from(value: $typ) -> Self {
                    Self::$path(value)
                }
            }
        )+
    }
}

from! {
    Error:
    SessionError => Session;
    std::io::Error => Bind;
}
