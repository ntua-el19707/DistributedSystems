module github.com/ntua-el19707/DistributedSystems/backendService/enitys

go 1.21.3

replace Service v0.0.0 => ../services/service

require Service v0.0.0

replace Hasher v0.0.0 => ../services/Hash

require Hasher v0.0.0

replace Logger v0.0.0 => ../services/Logger

require Logger v0.0.0
