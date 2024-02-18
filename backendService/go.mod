module github.com/ntua-el19707/DistributedSystems/backendService

go 1.21.3

replace Inbox v0.0.0 => ./services/Inbox

require Inbox v0.0.0

replace SystemInfo v0.0.0 => ./services/SystemInfo

require SystemInfo v0.0.0

replace Register v0.0.0 => ./services/Register

require Register v0.0.0

replace MessageSystem v0.0.0 => ./MessageSystem

require MessageSystem v0.0.0

replace asyncLoad v0.0.0 => ./services/asyncLoad

require asyncLoad v0.0.0

replace services v0.0.0 => ./services

require services v0.0.0

replace entitys v0.0.0 => ./entitys

require entitys v0.0.0

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
)

replace TransactionManager v0.0.0 => ./services/TransactionManager

require TransactionManager v0.0.0

replace Service v0.0.0 => ./services/service

require Service v0.0.0

replace WalletAndTransactions v0.0.0 => ./services/WalletAndTransactions

require WalletAndTransactions v0.0.0

replace Hasher v0.0.0 => ./services/Hash

require Hasher v0.0.0

replace Generator v0.0.0 => ./services/Generator

require Generator v0.0.0

replace FindBalance v0.0.0 => ./services/FindBalance

require FindBalance v0.0.0

replace RabbitMqService v0.0.0 => ./services/RabbitMqService

require RabbitMqService v0.0.0

replace Lottery v0.0.0 => ./services/Lottery

require Lottery v0.0.0

replace Stake v0.0.0 => ./services/Stake

require Stake v0.0.0

replace Logger v0.0.0 => ./services/Logger

require Logger v0.0.0
