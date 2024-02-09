package asyncLoad

import (
	"Logger"
	"RabbitMqService"
	"Service"

	"WalletAndTransactions"
)

type AsyncLoadService interface {
	Service.Service
	Consumer()
}

type AsyncLoadImpl struct {
	Providers AsyncLoadProviders
}

type AsyncLoadProviders struct {
	LoggerService    Logger.LoggerService
	RabbitMqService  RabbitMqService.RabbitMqService
	BlockCoinService *WalletAndTransactions.BlockChainCoinsImpl
}

func (s *AsyncLoadImpl) Construct() error {
	return nil
}

func (s *AsyncLoadImpl) Consumer() {}

func (s *AsyncLoadImpl) consumeTransactions() {
	//create  a channel to  collect Transactions
	channel := s.Providers.RabbitMqService.GetChannelTransactionCoin()
	go s.Providers.RabbitMqService.ConsumerTransactionsCoins()
	for {
		transaction := <-channel
		s.Providers.BlockCoinService.InsertTransaction(transaction)
	}
}
