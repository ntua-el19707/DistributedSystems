package Stake

import (
	"Logger"
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
	var sum float64 //  sum = 0   initial
	for _, t := range Transactions {
		//Scan  The  block
		if t.BillDetails.Bill.From.Address == key {
			sum += t.Amount
		}
		if t.BillDetails.Bill.To.Address == key {
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
func (service StakeCoinBlockChain) MapOfDistibutesRoundUp(scaleFactor float64) (map[rsa.PublicKey]int, int) {
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
func (service StakeCoinBlockChain) GetCurrentHash() string {
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
