package Stake

import (
	"Logger"
	"RabbitMqService"
	"Service"
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
)

type StakeService interface {
	Service.Service
	distributionOfStake(weight float64) (map[rsa.PublicKey]float64, float64)
	MapOfDistibutesRoundUp(weight float64) (map[rsa.PublicKey]int, int)
	GetCurrentHash() string
	GetWorkers() []rsa.PublicKey
}

// -- SERVICE  -- Providers
type StakeProviders struct {
	Service.Service
	LoggerService Logger.LoggerService
}

func (providers *StakeProviders) Construct(serviceName string) error {
	if providers.LoggerService == nil {
		providers.LoggerService = &Logger.Logger{ServiceName: serviceName}
		err := providers.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	return nil
}

type StakeCoinBlockChain struct {
	Block    entitys.BlockCoinEntity
	Workers  []rsa.PublicKey
	Services StakeProviders
}

/*
*

	construct -- Service  Construct  @Service stakeService   @Implementation  stakeCoinBlockChain
	@Returns  error
*/
func (service *StakeCoinBlockChain) Construct() error {
	block := service.Block
	expectedSize := block.BlockEntity.Capicity
	actualSize := len(block.Transactions)
	if expectedSize != actualSize {
		return errors.New(fmt.Sprintf("Block  is  not  full cappicity  is  %d but  has  %d ", expectedSize, actualSize))
	}
	serviceName := fmt.Sprintf("StakeBlockCoinService_%s", block.BlockEntity.CurrentHash)
	err := service.Services.Construct(serviceName)
	if err != nil {
		return err
	}
	service.Services.LoggerService.Log("Service  created")
	return nil
}

type sumResponseStruct struct {
	sum float64
	key rsa.PublicKey
}

func EqualPublicKeys(key1, key2 *rsa.PublicKey) bool {
	return key1.N.Cmp(key2.N) == 0 && key1.E == key2.E
}

/*
*

	sum -  local  funvtions  that  sum  the  transactions of  a  client 'key'
	@Param  Transactions []  entitys.TransactionCoins
	@Param  key  rsa.PublicKey
	@Param  scaleFactor  float64
	@Param  notify   chan sumResponseStruct
*/
func sum(Transactions []entitys.TransactionCoins, key rsa.PublicKey, scaleFactor float64, notify chan sumResponseStruct) {
	//?  what  is  scaleFactor  ?
	//&  The  scaleFactor  set how much value the coins  tha someone receives
	// for example  let  say  that in the  block  client A => + 50  (received ) ,  client B => -50  (send)
	// with scale  factor  0 => Total sum  50   ,   A: 50  ,  B:0
	// with scale  factor  0.5 => Total sum  75 ,   A: 50  ,  B:25
	// with scale  factor  1 => Total sum  100  ,   A: 50  ,  B:50
	// with scale  factor  2 => Total sum  150  ,   A: 50  ,  B:100
	// 0 <= scaleFactor <1  'From'  have  more chance
	// scaleFactor > 1  'To'  have  more chance
	// scaleFactor = 1  equal  distibution
	sum := float64(0) //  sum = 0   initial
	for _, t := range Transactions {
		//Scan  The  block
		from := t.BillDetails.Bill.From.Address
		to := t.BillDetails.Bill.To.Address
		var zero rsa.PublicKey
		if zero != from {
			if EqualPublicKeys(&from, &key) {
				sum += t.Amount
			}
		}
		if EqualPublicKeys(&to, &key) {
			sum += scaleFactor * t.Amount
		}
	}
	resp := sumResponseStruct{sum: sum, key: key}
	notify <- resp
}

/*
*

	distributionOfStake() -- Distribute Stake   @Service stakeService   @Implementation  stakeCoinBlockChain
	@Param   scaleFactor  float64
	@Return  map[rsa.PublicKey] float64 , float64
*/
func (service *StakeCoinBlockChain) distributionOfStake(scaleFactor float64) (map[rsa.PublicKey]float64, float64) {
	//loger  service  insatnce
	logger := service.Services.LoggerService
	block := service.Block
	//*  Log  The  Start of  Service
	logger.Log("Start  DistibutionOfStake ")

	var total float64                    //  holds  the  Total amount of Block
	totalWorkers := len(service.Workers) //How  many  clients
	//create  a channel to  collect the weight  of  each one
	collector := make(chan sumResponseStruct, totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		go sum(block.Transactions, service.Workers[i], scaleFactor, collector)
	}
	//Create  Distribution Map
	distributionMap := make(map[rsa.PublicKey]float64)
	for i := 0; i < totalWorkers; i++ {
		//Collect  Each One
		record := <-collector
		distributionMap[record.key] = record.sum
		total += record.sum
	}
	//*  Log  The commit
	logger.Log("Commit  DistibutionOfStake ")
	return distributionMap, total
	//*  NOTES  :  distributionOfStake() -(main routine) will  Break  n (len  workers  ) routines
	//?  Why  There  is  no semaphore  for  distributionMap[record.key] = record.sum and  total
	//&  Beacause  the collection  and  save  is  happening in main routine (Blocking  style  )
}

/*
*

	MapOfDistibutesRoundUp() -- create  Weight  ineteger  for distribution map      @Service stakeService   @Implementation  stakeCoinBlockChain
	@Param   scaleFactor  float64
	@Return  map[rsa.PublicKey] float64 , float64
*/
func (service *StakeCoinBlockChain) MapOfDistibutesRoundUp(scaleFactor float64) (map[rsa.PublicKey]int, int) {
	logger := service.Services.LoggerService
	logger.Log("Start MapOfDistibutesRoundUp ")

	amounts, total := service.distributionOfStake(scaleFactor)
	roundedMap := make(map[rsa.PublicKey]int)
	sum := 0
	for key, amount := range amounts {
		roundedMap[key] = int((amount / total) * 100000)
		sum += roundedMap[key]
	}
	//? why  1000 intead of 100
	//&  increase  accuracy  of  % etc  0.4832 with 100 =>  48% , 1000 => 48.3% , => 10000 => 48.32%
	logger.Log("Commit MapOfDistibutesRoundUp ")
	return roundedMap, sum
}

/*
*

	GetCurrentHash() -- get block  current  hash    @Service stakeService   @Implementation  stakeCoinBlockChain
	@Return  string
*/
func (service *StakeCoinBlockChain) GetCurrentHash() string {
	return service.Block.BlockEntity.CurrentHash //  return  hash  of  block
}

/*
*

	getWorkers() -- get Workers    @Service stakeService   @Implementation  stakeCoinBlockChain
	@Return  [] rsa.PublicKey
*/
func (service StakeCoinBlockChain) GetWorkers() []rsa.PublicKey {
	//*  VERY IMPORTANT Law : all  nodes must have the same  order in service worker and same  list
	return service.Workers //  return  workers
}

type StakeMesageBlockChain struct {
	Block    entitys.BlockMessage
	Workers  []rsa.PublicKey
	Services StakeProviders
}

/*
*

	sumMsGLen -  local  funvtions  that  sum  the  transactions of  a  client 'key'
	@Param  Transactions []  entitys.TransactionMsg
	@Param  key  rsa.PublicKey
	@Param  scaleFactor  float64
	@Param  notify   chan sumResponseStruct
*/
func sumMsgLen(Transactions []entitys.TransactionMsg, key rsa.PublicKey, scaleFactor float64, notify chan sumResponseStruct) {
	//?  what  is  scaleFactor  ?
	// 0 <= scaleFactor <1  'From'  have  more chance
	// scaleFactor > 1  'To'  have  more chance
	// scaleFactor = 1  equal  distibution
	var sum float64 //  sum = 0   initial
	for _, t := range Transactions {
		//Scan  The  block
		from := t.BillDetails.Bill.From.Address
		to := t.BillDetails.Bill.To.Address
		var zero rsa.PublicKey
		if zero != from {
			if EqualPublicKeys(&from, &key) {
				sum += float64(len(t.Msg))
			}
		}
		if EqualPublicKeys(&to, &key) {
			sum += scaleFactor * float64(len(t.Msg))
		}
	}
	resp := sumResponseStruct{sum: sum, key: key}
	notify <- resp
}

/*
*

	construct -- Service  Construct  @Service stakeService   @Implementation  stakeMessageBlockChain
	@Returns  error
*/
func (service *StakeMesageBlockChain) Construct() error {
	block := service.Block
	expectedSize := block.BlockEntity.Capicity
	actualSize := len(block.Transactions)
	if expectedSize != actualSize {
		return errors.New(fmt.Sprintf("Block  is  not  full cappicity  is  %d but  has  %d ", expectedSize, actualSize))
	}
	serviceName := fmt.Sprintf("StakeBlockMessageService_%s", block.BlockEntity.CurrentHash)
	err := service.Services.Construct(serviceName)
	if err != nil {
		return err
	}
	service.Services.LoggerService.Log("Service  created")
	return nil
}

/*
*

	distributionOfStake() -- Distribute Stake   @Service stakeService   @Implementation  stakeMessageBlockChain
	@Param   scaleFactor  float64
	@Return  map[rsa.PublicKey] float64 , float64
*/
func (service *StakeMesageBlockChain) distributionOfStake(scaleFactor float64) (map[rsa.PublicKey]float64, float64) {
	//loger  service  insatnce
	logger := service.Services.LoggerService
	block := service.Block
	//*  Log  The  Start of  Service
	logger.Log("Start  DistibutionOfStake ")

	var total float64                    //  holds  the  Total amount of Block
	totalWorkers := len(service.Workers) //How  many  clients
	//create  a channel to  collect the weight  of  each one
	collector := make(chan sumResponseStruct, totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		go sumMsgLen(block.Transactions, service.Workers[i], scaleFactor, collector)
	}
	//Create  Distribution Map
	distributionMap := make(map[rsa.PublicKey]float64)
	for i := 0; i < totalWorkers; i++ {
		//Collect  Each One
		record := <-collector
		distributionMap[record.key] = record.sum
		total += record.sum
	}
	//*  Log  The commit
	logger.Log("Commit  DistibutionOfStake ")
	return distributionMap, total
	//*  NOTES  :  distributionOfStake() -(main routine) will  Break  n (len  workers  ) routines
	//?  Why  There  is  no semaphore  for  distributionMap[record.key] = record.sum and  total
	//&  Beacause  the collection  and  save  is  happening in main routine (Blocking  style  )
}

/*
*

	MapOfDistibutesRoundUp() -- create  Weight  ineteger  for distribution map      @Service stakeService   @Implementation  stakeMsgBlockChain
	@Param   scaleFactor  float64
	@Return  map[rsa.PublicKey] float64 , float64
*/
func (service *StakeMesageBlockChain) MapOfDistibutesRoundUp(scaleFactor float64) (map[rsa.PublicKey]int, int) {
	logger := service.Services.LoggerService
	logger.Log("Start MapOfDistibutesRoundUp ")

	amounts, total := service.distributionOfStake(scaleFactor)

	roundedMap := make(map[rsa.PublicKey]int)
	sum := 0
	for key, amount := range amounts {
		roundedMap[key] = int((amount / total) * 100000)
		sum += roundedMap[key]
	}
	//? why  1000 intead of 100
	//&  increase  accuracy  of  % etc  0.4832 with 100 =>  48% , 1000 => 48.3% , => 10000 => 48.32%
	logger.Log("Commit MapOfDistibutesRoundUp ")
	return roundedMap, sum
}

/*
*

	GetCurrentHash() -- get block  current  hash    @Service stakeService   @Implementation  StakeMesageBlockChain
	@Return  string
*/
func (service *StakeMesageBlockChain) GetCurrentHash() string {
	return service.Block.BlockEntity.CurrentHash //  return  hash  of  block
}

/*
*

	getWorkers() -- get Workers    @Service stakeService   @Implementation  StakeMesageBlockChain
	@Return  [] rsa.PublicKey
*/
func (service StakeMesageBlockChain) GetWorkers() []rsa.PublicKey {
	//*  VERY IMPORTANT Law : all  nodes must have the same  order in service worker and same  list
	return service.Workers //  return  workers
}

type StakeProviders2 struct {
	LoggerService Logger.LoggerService
	RabbitMq      RabbitMqService.RabbitMqService
}

const errCouldNotFindProvider string = "could  not find  provider '%s'"
const errCurrentHashShouldBeGiven string = "should  give  give the current hash 'Spin Requirment'"
const errNotHaveQueueAndTopic string = "do not  have  queue  or a exchange topic"

func (providers *StakeProviders2) Construct() error {
	if providers.LoggerService == nil {
		providers.LoggerService = &Logger.Logger{ServiceName: "stake  service v3"}
		err := providers.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	if providers.RabbitMq == nil {
		errMsg := fmt.Sprintf(errCouldNotFindProvider, "RabbitMqService")
		return errors.New(errMsg)
	}
	return nil
}

type StakeBCCv3struct struct {
	Workers          []rsa.PublicKey
	Providers        StakeProviders2
	HashCurrent      string
	Who              int
	QueueAndExchange RabbitMqService.QueueAndExchange
	totalWorkers     int
	vld              bool
}

func (service *StakeBCCv3struct) Construct() error {
	err := service.Providers.Construct()
	if err != nil {
		return err
	}
	if service.HashCurrent == "" {
		return errors.New(errCurrentHashShouldBeGiven)
	}
	if service.QueueAndExchange.Queue == "" || service.QueueAndExchange.Exchange == "" {
		return errors.New(errNotHaveQueueAndTopic)
	}
	service.totalWorkers = len(service.Workers)
	service.vld = true
	service.Providers.LoggerService.Log("Service  created")
	return nil
}
func (service *StakeBCCv3struct) distributionOfStake(Bcc float64) (map[rsa.PublicKey]float64, float64) {

	if !service.vld {
		return nil, 0
	}
	logger := service.Providers.LoggerService
	logger.Log("Start distribution of  StakeV3 'BCC'")
	pack := entitys.StakePack{Node: service.Who, Bcc: Bcc}

	// -- BROADCAST --
	logger.Log("Start BroadCast Stake")
	err := service.Providers.RabbitMq.PublishStake(pack, service.QueueAndExchange)
	if err != nil {
		suicideNode := fmt.Sprintf("Fatal error due to %s", err.Error())
		logger.Fatal(suicideNode)
	}
	logger.Log("Commit BroadCast Stake")
	distributionMap := make(map[rsa.PublicKey]float64)
	sum := float64(0)
	logger.Log(fmt.Sprintf("Start Consuming Message waiting for a total  %d message", service.totalWorkers))
	// -- Consume Stake --
	for i := 0; i < service.totalWorkers; i++ {
		logger.Log(fmt.Sprintf("Start Consuming Message waiting for %d message", i))
		recieved := service.Providers.RabbitMq.ConsumeStake(service.QueueAndExchange)
		if recieved.Node >= service.totalWorkers {
			suicideNode := fmt.Sprintf("Fatal error due to recieved  node %d  but have total nodes %d", recieved.Node, service.totalWorkers)
			logger.Fatal(suicideNode)
		}
		sum += recieved.Bcc
		distributionMap[service.Workers[recieved.Node]] = recieved.Bcc
		logger.Log(fmt.Sprintf("Commit Consuming Message waiting for %d message", i))
	}
	logger.Log(fmt.Sprintf("Commit Consuming Message waiting for a total %d message", service.totalWorkers))

	logger.Log("Commit distribution of  StakeV3 'BCC'")
	return distributionMap, sum
}

/*
*

	MapOfDistibutesRoundUp() -- create  Weight  ineteger  for distribution map      @Service stakeService   @Implementation  StakeBCCv3struct
	@Param   scaleFactor  float64
	@Return  map[rsa.PublicKey] float64 , float64
*/
func (service *StakeBCCv3struct) MapOfDistibutesRoundUp(scaleFactor float64) (map[rsa.PublicKey]int, int) {
	if !service.vld {
		return nil, 0
	}
	logger := service.Providers.LoggerService
	logger.Log("Start MapOfDistibutesRoundUp ")

	amounts, total := service.distributionOfStake(scaleFactor)
	roundedMap := make(map[rsa.PublicKey]int)
	sum := 0
	for key, amount := range amounts {
		roundedMap[key] = int((amount / total) * 100000)
		sum += roundedMap[key]
	}
	//? why  1000 intead of 100
	//&  increase  accuracy  of  % etc  0.4832 with 100 =>  48% , 1000 => 48.3% , => 10000 => 48.32%
	logger.Log("Commit MapOfDistibutesRoundUp ")
	return roundedMap, sum
}
func (service *StakeBCCv3struct) GetCurrentHash() string {
	return service.HashCurrent
}
func (service *StakeBCCv3struct) GetWorkers() []rsa.PublicKey {
	return service.Workers
}

type MockStake struct {
	Hash                       string
	DistributedRoundUp         map[rsa.PublicKey]int
	Workers                    []rsa.PublicKey
	Total                      int
	CallGetCurrentHash         int
	CallMapOfDistibutesRoundUp int
}

func (service *MockStake) Construct() error {
	return nil
}
func (service *MockStake) distributionOfStake(weight float64) (map[rsa.PublicKey]float64, float64) {
	return nil, 0
}
func (service *MockStake) MapOfDistibutesRoundUp(weight float64) (map[rsa.PublicKey]int, int) {
	service.CallMapOfDistibutesRoundUp++
	return service.DistributedRoundUp, service.Total
}
func (service *MockStake) GetCurrentHash() string {
	service.CallGetCurrentHash++
	return service.Hash
}
func (service *MockStake) GetWorkers() []rsa.PublicKey {
	return service.Workers
}
