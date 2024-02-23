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
	BlockMsgService  *WalletAndTransactions.BlockChainMsgImpl
}

func (p *AsyncLoadProviders) Construct() error {
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: "Async-Load-Service"}
		return p.Construct()
	}
	return nil
}
func (s *AsyncLoadImpl) Construct() error {
	return s.Providers.Construct()
}

func (s *AsyncLoadImpl) Consumer() {
	go s.consumeTransactionsCoins()
	go s.consumeTransactionsMsg()
	for {
	}
}

func (s *AsyncLoadImpl) consumeTransactionsCoins() {
	//create  a channel to  collect Transactions
	channel := s.Providers.RabbitMqService.GetChannelTransactionCoin()
	go s.Providers.RabbitMqService.ConsumerTransactionsCoins()
	for {
		s.Providers.LoggerService.Log("Waiting for  Transaction coin")
		transaction := <-channel
		s.Providers.LoggerService.Log("Received Transaction Coins")

		s.Providers.BlockCoinService.InsertTransaction(transaction)
	}
}
func (s *AsyncLoadImpl) consumeTransactionsMsg() {
	channel := s.Providers.RabbitMqService.GetChannelTransactionMsg()
	go s.Providers.RabbitMqService.ConsumerTransactionsMsg()
	for {
		transaction := <-channel
		s.Providers.BlockMsgService.InsertTransaction(transaction)
	}
}
