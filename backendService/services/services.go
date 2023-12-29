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
/** Providers - will create  the servies  and if  the servoice  can  be created  fall the system 
*/
func Providers(){
	logger :=  &Logger{ServiceName:"Providers"}
	logger.Log("Start  Loadding  Providers")
	wallet := &walletStructV1Service{}
    err :=  wallet.construct() 
	if err != nil {
		logger.Fatal(err.Error())
	}
	BalanceInstanceService = &balanceImplementation{walletService:wallet}
	err = BalanceInstanceService.construct()
		if err != nil {
		logger.Fatal(err.Error())
	}
	TransactionManagerInstanceService =&TransactionManager{ walletService:wallet , balanceService:BalanceInstanceService}
	err = TransactionManagerInstanceService.construct()
		if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Log("Commit  Loading  Providers")
	
	
	
	
}

