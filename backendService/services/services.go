package  services
import "log"
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
//Logger  Services 
var LogerServiceBalance LogerService
var LogerServiceWallet LogerService
//Services  of Application
var  Wallet_Service  WalletService 
var  Balance_Service BalanceService
var  Generator_Service GeneratorService
/** Providers - will create  the servies  and if  the servoice  can  be created  fall the system 
*/
func Providers(){
	list := make ([]  Service , 0) 
	//first Logger services 
	LogerServiceBalance = &Logger{ServiceName:balanceServiceName}
	list = append(list  , LogerServiceBalance)
	LogerServiceWallet = &Logger{ServiceName:walletServiceName}
	list = append(list  , LogerServiceWallet)
	
	Wallet_Service = &walletStructV1Service{}
	list = append(list  , Wallet_Service)
	Balance_Service = &balanceImplementation{}
	list = append(list  , Balance_Service)
	Generator_Service = &generatorImplementation{ServiceName:generatorServiceName ,CharSet:"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"}
		list = append(list  , Generator_Service)
	for  _,service := range  list {
		err := construct(service)
		if err != nil {
			log.Fatal(err.Error())
		}

	}
	
}

