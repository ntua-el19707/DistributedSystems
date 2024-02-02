module github.com/ntua-el19707/DistributedSystems/backendService/services

go 1.21.3

replace Logger v0.0.0 => ./Logger

require Logger v0.0.0

replace entitys v0.0.0 => ../entitys

require entitys v0.0.0
