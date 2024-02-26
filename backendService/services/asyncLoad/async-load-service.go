package asyncLoad

import (
	"Logger"
	"RabbitMqService"
	"Service"
	"errors"

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
	BlockCoinService WalletAndTransactions.BlockChainCoinsService
	BlockMsgService  WalletAndTransactions.BlockChainMsgService
}

const ErrNoRabbitMqProviders = "Failed to create service  'no RabbitMqService'"
const ErrNoBlockCoinService = "Failed to create service  'no BlockCoinService'"
const ErrNoMsgService = "Failed to create service  'no BlockMsfService'"

func (p *AsyncLoadProviders) Construct() error {
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: "Async-Load-Service"}
		return p.Construct()
	}
	if p.RabbitMqService == nil {
		return errors.New(ErrNoRabbitMqProviders)
	}
	if p.BlockCoinService == nil {
		return errors.New(ErrNoBlockCoinService)
	}
	if p.BlockMsgService == nil {
		return errors.New(ErrNoMsgService)
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

		err := s.Providers.BlockCoinService.InsertTransaction(transaction)
		if err != nil {
			s.Providers.LoggerService.Error(err.Error())
		}
	}
}
func (s *AsyncLoadImpl) consumeTransactionsMsg() {
	channel := s.Providers.RabbitMqService.GetChannelTransactionMsg()
	go s.Providers.RabbitMqService.ConsumerTransactionsMsg()
	for {
		transaction := <-channel
		err := s.Providers.BlockMsgService.InsertTransaction(transaction)
		if err != nil {
			s.Providers.LoggerService.Error(err.Error())
		}
	}
}
