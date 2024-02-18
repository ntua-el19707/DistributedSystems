module github.com/ntua-el19707/DistributedSystems/backendService/RabbitMqService

go 1.21.3

replace Service v0.0.0 => ../service

require Service v0.0.0

require Logger v0.0.0

replace Hasher v0.0.0 => ../Hash

replace entitys v0.0.0 => ../../entitys

require entitys v0.0.0

require Generator v0.0.0

replace Logger v0.0.0 => ../Logger

require Hasher v0.0.0

replace MessageSystem v0.0.0 => ../../MessageSystem

require MessageSystem v0.0.0

replace Generator v0.0.0 => ../Generator

require github.com/rabbitmq/amqp091-go v1.9.0 // indirect
