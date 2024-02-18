module github.com/ntua-el19707/DistributedSystems/backendService/MessageSystem

go 1.21.3

require github.com/rabbitmq/amqp091-go v1.9.0 // indirect

replace Service v0.0.0 => ../services/service

require Service v0.0.0

replace Logger v0.0.0 => ../services/Logger

require Logger v0.0.0

replace entitys v0.0.0 => ../entitys

require entitys v0.0.0
