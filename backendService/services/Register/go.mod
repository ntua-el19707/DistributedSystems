module github.com/ntua-el19707/DistributedSystems/backendService/services/Register

go 1.21.3

replace Service v0.0.0 => ../service

require Service v0.0.0

replace Hasher v0.0.0 => ../Hash

require Hasher v0.0.0

replace Logger v0.0.0 => ../Logger

require Logger v0.0.0

replace entitys v0.0.0 => ../../entitys

require entitys v0.0.0
