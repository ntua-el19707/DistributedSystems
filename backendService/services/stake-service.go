package services 

import  (
	"crypto/rsa"
	"errors"
	"fmt"
)

type  stakeService interface{
	Service
	distributionOfStake(weight  float64) (map[rsa.PublicKey] float64 , float64 )  
	MapOfDistibutesRoundUp(weight  float64)(map[rsa.PublicKey] int  , int )
	getCurrentHash() string
	getWorkers() [] rsa.PublicKey
}
// -- SERVICE  -- Providers 
type stakeProviders struct  {
	loggerService  LogerService 
}
func  (providers  *  stakeProviders ) construct(serviceName  string ) error{
	if  providers.loggerService == nil {
		 providers.loggerService =  &Logger{ServiceName :  serviceName}
		 err :=   providers.loggerService.construct() 
		 if err != nil {
			return  err
		 }
	} 
	return nil 
}  
type stakeCoinBlockChain  struct{
	block blockCoin
	workers []  rsa.PublicKey
	services  stakeProviders
}
/** 
    construct -- Service  Construct  @Service stakeService   @Implementation  stakeCoinBlockChain
    @Returns  error  
 */
func  ( service * stakeCoinBlockChain  )  construct() error {
	
	
	if service.block.b.capicity != len(service.block.Transactions) {
		return  errors.New(fmt.Sprintf("Block  is  not  full cappicity  is  %d but  has  %d " , service.block.b.capicity , len(service.block.Transactions)  ) )
	} 
	serviceName :=  fmt.Sprintf("StakeBlockCoinService_%s"   ,  service.block.b.currentHash )  
	err := service.services.construct(serviceName)
	if  err !=  nil{
		return  err
	} 
	service.services.loggerService.Log("Service  created")
	return nil 
}

type  sumResponseStruct struct {
	sum  float64 
	key  rsa.PublicKey 
}
func  sum(block blockCoin ,key rsa.PublicKey  ,  scaleFactor  float64   , notify  chan sumResponseStruct  ){
    var  sum float64 
	for _, t := range block.Transactions {
		if t.BillDetails.Bill.From.Address ==  key {
			sum +=  t.Amount
		} 
	    if t.BillDetails.Bill.To.Address ==  key {
			sum +=  scaleFactor * t.Amount
		} 
	}
	resp :=  sumResponseStruct{sum:sum  ,  key:key}
	notify <- resp 
}

/** 
     distributionOfStake() -- Distribute Stake   @Service stakeService   @Implementation  stakeCoinBlockChain
     @Return  map[rsa.PublicKey] float64 , float64 
*/
func  (service * stakeCoinBlockChain  ) 	distributionOfStake(weight float64 ) (map[rsa.PublicKey] float64 , float64 ){
    //loger  service  insatnce 
    logger := service.services.loggerService 
    //*  Log  The  Start of  Service 
	logger.Log("Start  DistibutionOfStake ")
	
    var  total float64  //  holds  the  Total amount of Block     
	totalWorkers :=  len(service.workers) //How  many  clients  
    //create  a channel to  collect the weight  of  each one  
	collector := make(chan sumResponseStruct,  totalWorkers)
	for  i:=0  ;i < totalWorkers ; i++ {
		go sum(service.block , service.workers[i] , weight , collector )
	}
    //Create  Distribution Map 
	distributionMap := make(map[rsa.PublicKey]float64)
	for  i:=0  ;i < totalWorkers ; i++ {
        //Collect  Each One  
		record :=  <- collector 
		distributionMap[record.key] = record.sum
		total +=  record.sum 
	}
    //*  Log  The commit  
	logger.Log("Commit  DistibutionOfStake ")
	return distributionMap ,  total
    //*  NOTES  :  distributionOfStake() -(main routine) will  Break  n (len  workers  ) routines
    //?  Why  There  is  no semaphore  for  distributionMap[record.key] = record.sum and  total
    //&  Beacause  the collection  and  save  is  happening in main routine (Blocking  style  ) 
}
/** 
     MapOfDistibutesRoundUp() -- create  Weight  ineteger  for distribution map      @Service stakeService   @Implementation  stakeCoinBlockChain
     @Return  map[rsa.PublicKey] float64 , float64 
*/
func (service stakeCoinBlockChain) MapOfDistibutesRoundUp(weight  float64) (map[rsa.PublicKey]int , int ) {
	logger := service.services.loggerService
	logger.Log("Start MapOfDistibutesRoundUp ")

	amounts ,  total := service.distributionOfStake(weight) 

	roundedMap := make(map[rsa.PublicKey]int)
	sum := 0 
	for key, amount := range amounts {
		roundedMap[key] = int((amount/total) * 100000) 
		sum += roundedMap[key]
	}
    //? why  1000 intead of 100 
    //&  increase  accuracy  of  % etc  0.4832 with 100 =>  48% , 1000 => 48.3% , => 10000 => 48.32% 
    logger.Log("Commit MapOfDistibutesRoundUp ")
	return roundedMap ,sum
}

/** 
      getCurrentHash() -- get block  current  hash    @Service stakeService   @Implementation  stakeCoinBlockChain
     @Return  string 
*/
func (service stakeCoinBlockChain) 	getCurrentHash() string  {
	return service.block.b.currentHash //  return  hash  of  block
}

/** 
     getWorkers() -- get Workers    @Service stakeService   @Implementation  stakeCoinBlockChain
     @Return  [] rsa.PublicKey
*/
func (service stakeCoinBlockChain) 	getWorkers() []  rsa.PublicKey {
    //*  VERY IMPORTANT Law : all  nodes must have the same  order in service worker and same  list 
    return service.workers //  return  hash  of  block
}
type mockStake  struct  {
	hash  string 
	distributedRoundUp map[rsa.PublicKey] int 
	workers []  rsa.PublicKey 
    total  int 
	callGetCurrentHash int 
	callMapOfDistibutesRoundUp int 

 }

func  	( service  * mockStake )  construct() error  {
	return nil
}
func  	( service  * mockStake )    distributionOfStake(weight  float64) (map[rsa.PublicKey] float64 , float64 )  {
	return  nil , 0 
}
func  	( service  * mockStake )	MapOfDistibutesRoundUp(weight float64)(map[rsa.PublicKey] int  , int ) {
	service.callMapOfDistibutesRoundUp++
	return 	service.distributedRoundUp ,service.total 
}
func  	( service  * mockStake )	getCurrentHash() string {
	service.callGetCurrentHash++
	return  service.hash
}
func 	( service  * mockStake ) 	getWorkers() [] rsa.PublicKey{
	return service.workers
}
