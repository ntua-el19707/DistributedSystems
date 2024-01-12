package  services

type  Service interface {
	construct()error  
}
/**
	construct -  for construction of a  service
	@Param  s Service 
*/
func construct(s   Service) error  {
	return s.construct()
}
const allChars ="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//names
const balanceServiceName string = "balanceService"
const walletServiceName = "walletService"
const generatorServiceName = "generatorService"
const transactionServiceName = "transactionService"

//Services  of Application
var TransactionManagerInstanceService TransactionManagerService
var BalanceInstanceService BalanceService
var BlockChainCoinsService blockChainService
var WalletServiceInstance WalletService
/** Providers - will create  the servies  and if  the servoice  can  be created  fall the system 
*/
func Providers(){
	logger :=  &Logger{ServiceName:"Providers"}
	logger.Log("Start  Loadding  Providers")
	WalletServiceInstance  = &walletStructV1Service{}
    err :=  WalletServiceInstance.construct() 
	if err != nil {
		logger.Fatal(err.Error())
	}
	BalanceInstanceService = &balanceImplementation{walletService:WalletServiceInstance}
	err = BalanceInstanceService.construct()
		if err != nil {
		logger.Fatal(err.Error())
	}
	TransactionManagerInstanceService =&TransactionManager{ walletService:WalletServiceInstance  , balanceService:BalanceInstanceService}
	err = TransactionManagerInstanceService.construct()
		if err != nil {
		logger.Fatal(err.Error())
	}


	hash := &hashIpmpl{}
	err = hash.construct() 
	if err != nil{
			logger.Fatal(err.Error())
	}

	BlockChainCoinsService =  &blockChainCoinsImpl{services:blockServiceProviders { walletService:WalletServiceInstance , hashService:hash}}
	err = 	BlockChainCoinsService.construct()
		if err != nil{
			logger.Fatal(err.Error())
	}
	err =  	BlockChainCoinsService.genesis()
	if  err != nil {
		logger.Fatal(err.Error())
	}
	

	logger.Log("Commit  Loading  Providers")
	
	
	
	
}

