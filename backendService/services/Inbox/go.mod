module github.com/ntua-el19707/DistributedSystems/backendService/services/Inbox

go 1.21.3

replace SystemInfo v0.0.0 => ../SystemInfo

require SystemInfo v0.0.0

replace MessageSystem v0.0.0 => ../../MessageSystem

require MessageSystem v0.0.0

replace Service v0.0.0 => ../service

require Service v0.0.0

replace Logger v0.0.0 => ../Logger

require Logger v0.0.0

require entitys v0.0.0

replace entitys v0.0.0 => ../../entitys

replace WalletAndTransactions v0.0.0 => ../WalletAndTransactions

require WalletAndTransactions v0.0.0

replace Hasher v0.0.0 => ../Hash

require Hasher v0.0.0

replace Generator v0.0.0 => ../Generator

require Generator v0.0.0

replace FindBalance v0.0.0 => ../FindBalance

require FindBalance v0.0.0

replace RabbitMqService v0.0.0 => ../RabbitMqService

require RabbitMqService v0.0.0

replace Stake v0.0.0 => ../Stake

require Stake v0.0.0

replace Lottery v0.0.0 => ../Lottery

require Lottery v0.0.0
