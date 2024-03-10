package services

import (
	"Hasher"
	"Logger"
	"Lottery"
	"RabbitMqService"
	"Service"
	"TransactionManager"
	"WalletAndTransactions"
	"crypto/rsa"
	"entitys"
	"fmt"
	"log"

	"Register"
	"SystemInfo"
	"asyncLoad"

	"Inbox"
)

var TransactionManagerService TransactionManager.TransactionManagerService

var RabbitMqS RabbitMqService.RabbitMqService
var WalletService WalletAndTransactions.WalletStructV1Implementation
var SystemInfoService SystemInfo.SystemInfoService
var blockChainCoinService WalletAndTransactions.BlockChainCoinsService
var blockChainMsgService WalletAndTransactions.BlockChainMsgService
var InboxService Inbox.InboxService
var FindBalanceService WalletAndTransactions.BalanceService
var scaleFactorMsg float64
var scaleFactorCoin float64

var capicity_Msg, capicity_Coin, expectWorkers int
var per_Node float64
var stakeCoinQueue RabbitMqService.QueueAndExchange
var stakeMsgQueue RabbitMqService.QueueAndExchange

func setQueues(node, rabbitMqUri string) {
	logger := Logger.Logger{ServiceName: "Set Queues And Constuct Wallet"}
	err := logger.Construct()
	if err != nil {
		log.Fatal(err.Error())

	}
	logger.Log("Start  Set Queues and create wallet ")

	genSet := func(queue, exchange, node string) RabbitMqService.QueueAndExchange {
		queue = fmt.Sprintf("%s-%s", queue, node)
		return RabbitMqService.QueueAndExchange{Queue: queue, Exchange: exchange}
	}
	RabbitMqUri := rabbitMqUri
	TransactionCoinSetQueueExchange := genSet("transactionCoins", "TCOINS", node)
	TransactionMsgSetQueueExchange := genSet("transactionMsg", "TMSG", node)
	BlockMsgQueueExchange := genSet("BlockCoins", "BCOIN", node)
	BlockCoinQueueExchange := genSet("BlockMsg", "BMSG", node)
	SystemInfoQueue := genSet("SystemInfo", "SINFO", node)
	stakeCoinQueue = genSet("StakeCoins", "STCOIN", node)
	stakeMsgQueue = genSet("StakeMsg", "STMSG", node)
	providers := RabbitMqService.RabbitMqProviders{RabbitMqUri: RabbitMqUri, StakeBlockMsgQueue: stakeMsgQueue, StakeBlockCoinQueue: stakeCoinQueue, SystemInfoQueue: SystemInfoQueue, TransactionCoinSetQueueExchange: TransactionCoinSetQueueExchange, TransactionMsgSetQueueExchange: TransactionMsgSetQueueExchange, BlockCoinQueueExchange: BlockCoinQueueExchange, BlockMsgQueueExchange: BlockMsgQueueExchange}
	WalletService = WalletAndTransactions.WalletStructV1Implementation{}
	bootStrapOrDie(&WalletService, &logger)
	RabbitMqS = &RabbitMqService.RabbitMqImpl{Providers: providers}

	bootStrapOrDie(RabbitMqS, &logger)
	logger.Log("Commit Set Queues and create wallet ")

}

func registerAndSystemInfo(coordinator bool, ExpectedWorkers int, Me, hostCoordinator, node, publicUri string) {
	logger := Logger.Logger{ServiceName: "register and system info Loader"}
	err := logger.Construct()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Log("Start Creating  register and system info")
	logger.Log("Start Creating  System info")
	providers := SystemInfo.SystemInfoProviders{RabbitMqService: RabbitMqS}
	SystemInfoService = &SystemInfo.SystemInfoImpl{Coordinator: coordinator, ExpectedWorkers: ExpectedWorkers, Providers: providers}
	bootStrapOrDie(SystemInfoService, &logger)
	logger.Log("Commit Creating  System info")
	if !coordinator {
		logger.Log("Start Creating  register  service")
		register := Register.RegisterImpl{}
		register.Who = hostCoordinator
		register.Me = Me
		register.MyPk = WalletService.GetPub()
		register.MyId = node
		register.UriPublic = publicUri
		bootStrapOrDie(&register, &logger)
		logger.Log("commit Creating  register  service")
		register.Register()
		err, scaleFactorMsg, scaleFactorCoin = SystemInfoService.Consume()
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		var params entitys.ClientRequestBody
		params.PublicKey = WalletService.GetPub()
		clientInfo := entitys.ClientInfo{Id: node, Uri: Me, UriPublic: publicUri}
		params.Client = clientInfo
		SystemInfoService.AddWorker(params)
	}
	logger.Log("Commit Creating  register and system info")

}

/** Providers - will create  the servies  and if  the servoice  can  be created  fall the system
 */
func providers(c bool) {
	logger := Logger.Logger{ServiceName: "provider Loader"}
	err := logger.Construct()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Log("Start creatig services 'dependecy injection'")
	// -- WalletService --

	//for now  use  mockFindBalance
	hashService := Hasher.HashImpl{}
	bootStrapOrDie(&hashService, &logger)
	spinProviders := Lottery.LotteryProviders{
		HasherService: &hashService,
	}
	lottery1 := Lottery.LotteryImpl{Services: spinProviders}
	bootStrapOrDie(&lottery1, &logger)
	blockProviders1 := WalletAndTransactions.BlockServiceProviders{
		RabbitMqService:       RabbitMqS,
		HashService:           &hashService,
		WalletServiceInstance: &WalletService,
		LotteryService:        &lottery1,
	}

	blockChainCoinService = &WalletAndTransactions.BlockChainCoinsImpl{ScaleFactor: scaleFactorCoin, Services: blockProviders1, Workers: []rsa.PublicKey{WalletService.GetPub()}}
	bootStrapOrDie(blockChainCoinService, &logger)
	lottery2 := Lottery.LotteryImpl{Services: spinProviders}
	bootStrapOrDie(&lottery2, &logger)
	blockProviders2 := WalletAndTransactions.BlockServiceProviders{
		RabbitMqService:       RabbitMqS,
		HashService:           &hashService,
		WalletServiceInstance: &WalletService,
		LotteryService:        &lottery2,
	}
	blockChainMsgService = &WalletAndTransactions.BlockChainMsgImpl{ScaleFactor: scaleFactorMsg, Services: blockProviders2, Workers: []rsa.PublicKey{WalletService.GetPub()}}
	bootStrapOrDie(blockChainMsgService, &logger)
	if c {
		err = blockChainCoinService.Genesis(capicity_Coin, expectWorkers, per_Node)

		if err != nil {
			logger.Fatal(err.Error())
		}
		err = blockChainMsgService.Genesis(capicity_Msg)
		if err != nil {
			logger.Fatal(err.Error())
		}

		broadCastBlockCoin := entitys.BlockCoinMessageRabbitMq{
			BlockCoin: blockChainCoinService.RetriveChain()[0],
		}
		err = RabbitMqS.PublishBlockCoin(broadCastBlockCoin)
		if err != nil {
			logger.Fatal(err.Error())
		}
		broadCastBlockMsg := entitys.BlockMessageMessageRabbitMq{
			BlockMsg: blockChainMsgService.RetriveChain()[0],
		}
		err = RabbitMqS.PublishBlockMsg(broadCastBlockMsg)
		if err != nil {
			logger.Fatal(err.Error())
		}
		RabbitMqS.ConsumeNextBlockCoin()
		RabbitMqS.ConsumeNextBlockMsg()
	} else {
		blockCoin := RabbitMqS.ConsumeNextBlockCoin()
		blockMsg := RabbitMqS.ConsumeNextBlockMsg()
		blockChainMsgService.InsertNewBlock(blockMsg.BlockMsg)
		blockChainCoinService.InsertNewBlock(blockCoin.BlockCoin)
		blockChainCoinService.SetWorkers(SystemInfoService.GetWorkers())
		blockChainMsgService.SetWorkers(SystemInfoService.GetWorkers())
		//NOW The re is  litle  a chance to fail Mine only if  internal error if err =>  commit  harakiri
	}
	thisNode, _ := SystemInfoService.NodeDetails(WalletService.GetPub())

	blockChainCoinService.SetWhoAndQueue(thisNode.IndexId, stakeCoinQueue)
	blockChainMsgService.SetWhoAndQueue(thisNode.IndexId, stakeMsgQueue)
	FindBalanceService = &WalletAndTransactions.BalanceImplementation{
		BlockChainService: blockChainCoinService,
		SystemInfoService: SystemInfoService,
	}
	bootStrapOrDie(FindBalanceService, &logger)

	TransactionManagerService = &TransactionManager.TransactionManager{
		WalletServiceInstance:      &WalletService,
		FindBalanceServiceInstance: FindBalanceService,
	}
	bootStrapOrDie(TransactionManagerService, &logger)
	InboxService = &Inbox.InboxImpl{Providers: Inbox.InboxProviders{BlockChainService: blockChainMsgService, SystemInfoService: SystemInfoService}}
	bootStrapOrDie(InboxService, &logger)

	asyncProviders := asyncLoad.AsyncLoadProviders{
		RabbitMqService:  RabbitMqS,
		BlockCoinService: blockChainCoinService,
		BlockMsgService:  blockChainMsgService,
	}

	asyncService := asyncLoad.AsyncLoadImpl{Providers: asyncProviders}
	bootStrapOrDie(&asyncService, &logger)
	go asyncService.Consumer()
	logger.Log("Commit creatig services 'dependecy injection'")
}
func equalPublicKeys(key1, key2 *rsa.PublicKey) bool {
	return key1.N.Cmp(key2.N) == 0 && key1.E == key2.E
}

func bootStrapOrDie(s Service.Service, loggerService Logger.LoggerService) {
	err := s.Construct()
	if err != nil {
		loggerService.Fatal(err.Error())
	}
}
func SetUp() {
	logger := Logger.Logger{ServiceName: "set up "}
	err := logger.Construct()
	if err != nil {
		log.Fatal(err.Error())
	}
	blockChainCoinService.SetWorkers(SystemInfoService.GetWorkers())
	blockChainMsgService.SetWorkers(SystemInfoService.GetWorkers())
	pk := WalletService.GetPub()
	for _, key := range SystemInfoService.GetWorkers() {
		if !equalPublicKeys(&key, &pk) {
			list, err := TransactionManagerService.TransferMoney(key, float64(1000))
			if err != nil {
				logger.Fatal(err.Error())
			}
			err = RabbitMqS.PublishTractioncoinSet(list)
			if err != nil {
				logger.Fatal(err.Error())
			}

		}
	}

}
func BootOrDie(node, hostC, Me, rabbitMqUri, publicUri string, coordinator bool, ExpectedWorkers, capicityMsg, capicityCoin int, sFm, sFc, perNode float64) {
	scaleFactorMsg = sFm
	scaleFactorCoin = sFc
	setQueues(node, rabbitMqUri)
	capicity_Msg = capicityMsg
	capicity_Coin = capicityCoin
	per_Node = perNode
	expectWorkers = ExpectedWorkers
	registerAndSystemInfo(coordinator, ExpectedWorkers, Me, hostC, node, publicUri)
	providers(coordinator)
}
