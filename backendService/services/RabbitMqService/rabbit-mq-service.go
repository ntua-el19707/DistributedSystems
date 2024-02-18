package RabbitMqService

import (
	"Logger"
	"Service"
	"entitys"
	"fmt"

	"MessageSystem"
)

type RabbitMqService interface {
	Service.Service
	ConsumerTransactionsCoins()
	ConsumerTransactionsMsg()

	GetChannelTransactionCoin() chan entitys.TransactionCoinSet
	GetChannelTransactionMsg() chan entitys.TransactionMessageSet
	PublishBlockCoin(block entitys.BlockCoinMessageRabbitMq) error
	PublishBlockMsg(block entitys.BlockMessageMessageRabbitMq) error
	ConsumeNextBlockCoin() entitys.BlockCoinMessageRabbitMq
	ConsumeNextBlockMsg() entitys.BlockMessageMessageRabbitMq
	PublishTractioncoinSet(t entitys.TransactionCoinSet) error
	PublishTractionMsgSet(t entitys.TransactionMessageSet) error
	BroadCastSystemInfo(p entitys.RabbitMqSystemInfoPack) error
	ConsumeNextSystemInfo() entitys.RabbitMqSystemInfoPack
}

const serviceName = "RabbitMqService"

type QueueAndExchange struct {
	Queue    string
	Exchange string
}
type RabbitMqProviders struct {
	LoggerService                   Logger.LoggerService
	consumerTransactionCoin         func(chan MessageSystem.ConsumerMsgResp[entitys.TransactionCoinSet], string, string, string, Logger.LoggerService)
	consumerTransactionMsg          func(chan MessageSystem.ConsumerMsgResp[entitys.TransactionMessageSet], string, string, string, Logger.LoggerService)
	consumerBlockMsg                func(string, string, Logger.LoggerService) (entitys.BlockMessageMessageRabbitMq, error)
	consumerBlockCoin               func(string, string, Logger.LoggerService) (entitys.BlockCoinMessageRabbitMq, error)
	consumerSystemInfo              func(string, string, Logger.LoggerService) (entitys.RabbitMqSystemInfoPack, error)
	channelTransactionCoinSet       chan entitys.TransactionCoinSet
	channelTransactionMsg           chan entitys.TransactionMessageSet
	ctr                             bool
	RabbitMqUri                     string
	TransactionCoinSetQueueExchange QueueAndExchange
	TransactionMsgSetQueueExchange  QueueAndExchange
	BlockMsgQueueExchange           QueueAndExchange
	BlockCoinQueueExchange          QueueAndExchange
	SystemInfoQueue                 QueueAndExchange
}

func (p *RabbitMqProviders) Construct() error {
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: serviceName}
		err := p.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	p.consumerTransactionCoin = MessageSystem.Consumer[entitys.TransactionCoinSet]
	p.consumerTransactionMsg = MessageSystem.Consumer[entitys.TransactionMessageSet]
	p.consumerBlockMsg = MessageSystem.ConsumeOne[entitys.BlockMessageMessageRabbitMq]
	p.consumerBlockCoin = MessageSystem.ConsumeOne[entitys.BlockCoinMessageRabbitMq]
	p.consumerSystemInfo = MessageSystem.ConsumeOne[entitys.RabbitMqSystemInfoPack]
	p.channelTransactionCoinSet = make(chan entitys.TransactionCoinSet)
	p.channelTransactionMsg = make(chan entitys.TransactionMessageSet)

	err := MessageSystem.CreateAndBind(p.RabbitMqUri, p.TransactionCoinSetQueueExchange.Queue, p.TransactionCoinSetQueueExchange.Exchange, p.LoggerService)
	if err != nil {
		return err
	}
	err = MessageSystem.CreateAndBind(p.RabbitMqUri, p.TransactionMsgSetQueueExchange.Queue, p.TransactionMsgSetQueueExchange.Exchange, p.LoggerService)
	if err != nil {
		return err
	}
	err = MessageSystem.CreateAndBind(p.RabbitMqUri, p.BlockCoinQueueExchange.Queue, p.BlockCoinQueueExchange.Exchange, p.LoggerService)
	if err != nil {
		return err
	}
	err = MessageSystem.CreateAndBind(p.RabbitMqUri, p.BlockMsgQueueExchange.Queue, p.BlockMsgQueueExchange.Exchange, p.LoggerService)
	if err != nil {
		return err
	}
	err = MessageSystem.CreateAndBind(p.RabbitMqUri, p.SystemInfoQueue.Queue, p.SystemInfoQueue.Exchange, p.LoggerService)
	if err != nil {
		return err
	}
	p.ctr = true

	return nil
}

// RabbitMq  service
type RabbitMqImpl struct {
	Providers RabbitMqProviders
}

func (service *RabbitMqImpl) Construct() error {
	return service.Providers.Construct()
}

func (service *RabbitMqImpl) ConsumerTransactionsCoins() {
	providers := &service.Providers
	if providers.ctr == true {
		channel := make(chan MessageSystem.ConsumerMsgResp[entitys.TransactionCoinSet])
		go providers.consumerTransactionCoin(channel, providers.RabbitMqUri, providers.TransactionCoinSetQueueExchange.Queue, providers.TransactionCoinSetQueueExchange.Exchange, providers.LoggerService)
		for {
			pack := <-channel
			fmt.Println(pack)
			if pack.Err == nil {
				providers.channelTransactionCoinSet <- pack.Payload
			}
		}
	}

}
func (service *RabbitMqImpl) ConsumerTransactionsMsg() {
	providers := &service.Providers
	if providers.ctr == true {
		channel := make(chan MessageSystem.ConsumerMsgResp[entitys.TransactionMessageSet])
		go providers.consumerTransactionMsg(channel, providers.RabbitMqUri, providers.TransactionMsgSetQueueExchange.Queue, providers.TransactionMsgSetQueueExchange.Exchange, providers.LoggerService)
		for {
			pack := <-channel
			if pack.Err == nil {
				providers.channelTransactionMsg <- pack.Payload
			}
		}
	}

}

func (service *RabbitMqImpl) PublishTractioncoinSet(t entitys.TransactionCoinSet) error {
	providers := &service.Providers
	topic := providers.TransactionCoinSetQueueExchange.Exchange
	return MessageSystem.ProducerBroadCast(t, providers.RabbitMqUri, topic, providers.LoggerService)
}
func (service *RabbitMqImpl) PublishTractionMsgSet(t entitys.TransactionMessageSet) error {
	providers := &service.Providers
	topic := providers.TransactionMsgSetQueueExchange.Exchange
	return MessageSystem.ProducerBroadCast(t, providers.RabbitMqUri, topic, providers.LoggerService)
}
func (service *RabbitMqImpl) PublishBlockCoin(block entitys.BlockCoinMessageRabbitMq) error {
	providers := &service.Providers
	topic := providers.BlockCoinQueueExchange.Exchange
	return MessageSystem.ProducerBroadCast(block, providers.RabbitMqUri, topic, providers.LoggerService)
}
func (service *RabbitMqImpl) PublishBlockMsg(block entitys.BlockMessageMessageRabbitMq) error {
	providers := &service.Providers
	topic := providers.BlockMsgQueueExchange.Exchange
	return MessageSystem.ProducerBroadCast(block, providers.RabbitMqUri, topic, providers.LoggerService)
}
func (service *RabbitMqImpl) BroadCastSystemInfo(payload entitys.RabbitMqSystemInfoPack) error {

	providers := &service.Providers
	topic := providers.SystemInfoQueue.Exchange
	return MessageSystem.ProducerBroadCast(payload, providers.RabbitMqUri, topic, providers.LoggerService)
}

func (service *RabbitMqImpl) GetChannelTransactionCoin() chan entitys.TransactionCoinSet {
	return service.Providers.channelTransactionCoinSet
}
func (service *RabbitMqImpl) GetChannelTransactionMsg() chan entitys.TransactionMessageSet {
	return service.Providers.channelTransactionMsg
}
func (service *RabbitMqImpl) ConsumeNextBlockCoin() entitys.BlockCoinMessageRabbitMq {
	var block entitys.BlockCoinMessageRabbitMq
	var err error
	providers := &service.Providers
	if providers.ctr == true {
		block, err = providers.consumerBlockCoin(providers.RabbitMqUri, providers.BlockCoinQueueExchange.Queue, providers.LoggerService)
		if err != nil {
			providers.LoggerService.Fatal(err.Error())
		}
	}
	fmt.Println(block)
	return block
}
func (service *RabbitMqImpl) ConsumeNextBlockMsg() entitys.BlockMessageMessageRabbitMq {
	var block entitys.BlockMessageMessageRabbitMq
	var err error
	providers := &service.Providers
	if providers.ctr == true {
		block, err = providers.consumerBlockMsg(providers.RabbitMqUri, providers.BlockMsgQueueExchange.Queue, providers.LoggerService)

		if err != nil {
			providers.LoggerService.Fatal(err.Error())
		}
	}
	fmt.Println(block)
	return block
}
func (service *RabbitMqImpl) ConsumeNextSystemInfo() entitys.RabbitMqSystemInfoPack {
	var systemInfo entitys.RabbitMqSystemInfoPack
	var err error
	providers := &service.Providers
	if providers.ctr == true {
		systemInfo, err = providers.consumerSystemInfo(providers.RabbitMqUri, providers.SystemInfoQueue.Queue, providers.LoggerService)
		if err != nil {
			providers.LoggerService.Fatal(err.Error())
		}
	}
	return systemInfo
}

// Mock
type MockRabbitMqImpl struct {
	Channel                       chan entitys.TransactionCoinSet
	Blocks                        []entitys.BlockCoinMessageRabbitMq
	Block                         entitys.BlockCoinMessageRabbitMq
	BlockMsg                      entitys.BlockMessageMessageRabbitMq
	BlockMsgRsp                   entitys.BlockMessageMessageRabbitMq
	ErrPublishBlockMsg            error
	ErrPublishBlockCoin           error
	ErrPublishTransactionCoin     error
	TransactionSetCoin            entitys.TransactionCoinSet
	index                         int
	CallPublishBlock              int
	CallConsumeBlock              int
	CallPublishBlockMsg           int
	CallConsumeBlockMsg           int
	CallPublishTransactionCoinSet int
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
func (mock *MockRabbitMqImpl) PublishBlockCoin(block entitys.BlockCoinMessageRabbitMq) error {
	mock.Block = block

	mock.CallPublishBlock++
	return mock.ErrPublishBlockCoin
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
func (mock *MockRabbitMqImpl) PublishBlockMsg(block entitys.BlockMessageMessageRabbitMq) error {
	mock.BlockMsg = block
	mock.CallPublishBlockMsg++
	return mock.ErrPublishBlockMsg
}
func (mock *MockRabbitMqImpl) ConsumeNextBlockMsg() entitys.BlockMessageMessageRabbitMq {
	mock.CallConsumeBlockMsg++
	return mock.BlockMsgRsp
}
func (mock *MockRabbitMqImpl) PublishTractioncoinSet(t entitys.TransactionCoinSet) error {
	mock.CallPublishTransactionCoinSet++
	mock.TransactionSetCoin = t
	return mock.ErrPublishTransactionCoin
}
