package RabbitMqService

import (
	"Service"
	"entitys"
)

type RabbitMqService interface {
	Service.Service
	ConsumerTransactionsCoins()
	GetChannelTransactionCoin() chan entitys.TransactionCoinSet
	PublishBlock(block entitys.BlockCoinMessageRabbitMq)
	ConsumeNextBlock() entitys.BlockCoinMessageRabbitMq
}

// Mock
type MockRabbitMqImpl struct {
	Channel          chan entitys.TransactionCoinSet
	Blocks           []entitys.BlockCoinMessageRabbitMq
	Block            entitys.BlockCoinMessageRabbitMq
	index            int
	CallPublishBlock int
	CallConsumeBlock int
}

func (mock *MockRabbitMqImpl) Construct() error {
	return nil
}
func (mock *MockRabbitMqImpl) ConsumerTransactionsCoins() {

}
func (mock *MockRabbitMqImpl) GetChannelTransactionCoin() chan entitys.TransactionCoinSet {
	return mock.Channel
}
func (mock *MockRabbitMqImpl) ConsumeBlock() chan entitys.TransactionCoinSet {
	return mock.Channel
}
func (mock *MockRabbitMqImpl) PublishBlock(block entitys.BlockCoinMessageRabbitMq) {
	mock.Block = block
	mock.CallPublishBlock++
}
func (mock *MockRabbitMqImpl) ConsumeNextBlock() entitys.BlockCoinMessageRabbitMq {
	mock.CallConsumeBlock++
	index := mock.index
	if index > len(mock.Blocks)-1 {
		index = 0
	}
	mock.index++
	return mock.Blocks[index]
}
