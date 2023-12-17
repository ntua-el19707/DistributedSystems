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

//Services  of Application
var  Wallet_Service  WalletService 

/** Providers - will create  the servies  and if  the servoice  can  be created  fall the system 
*/
func Providers(){
	list := make ([]  Service , 0) 
	
	Wallet_Service = &walletStructV1Service{}
	list = append(list  , Wallet_Service)
	for  _,service := range  list {
		err := construct(service)
		if err != nil {
			log.Fatal(err.Error())
		}

	}
	
}

