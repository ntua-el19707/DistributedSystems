package services 
import  (
	"crypto/rsa"
	"time"
	"errors"
"fmt"
)

type  blockChainCoins  []  blockCoin

func  (chain *  blockChainCoins ) chainGenesis(providers blockServiceProviders, validator rsa.PublicKey ) {
	logger := providers.loggerService
	hasher := providers.hashService
	logger.Log("Start  creating a  new  chain -- GENESIS --  ")
	parent := hasher.ParrentOFall()
	current := hasher.instantHash(-1000)
	var empty blockChainCoins
	*chain = empty
	genessisBlock := &blockCoin{}
	err :=  genessisBlock.genesis(validator ,parent , current) 
	if err  != nil {
		logger.Error("Abbort chain  creation  ")
		return 
	} 
 	*chain = append( * chain , * genessisBlock)
	logger.Log("Commit  creating a  new  chain -- GENESIS --  ")
}
func  (chain   * blockChainCoins ) insertNewBlock(logger LogerService , validator rsa.PublicKey , blockDetails  blockCoin ) error  {
	logger.Log("Start insert  a  new  block in chain")
	expectedIndex := len(*chain )
	if expectedIndex !=  blockDetails.b.index {
		const  errmsgTemplate string  = "expecting the next block  would  have  index %d  but  had  %d "
		errmsg :=  fmt.Sprintf(errmsgTemplate , expectedIndex ,  blockDetails.b.index)
		logger.Error(fmt.Sprintf("Abbort: %s" , errmsg))
		return  errors.New(errmsg)
	}
	logger.Log("Start  validation of  block ")
	err :=  blockDetails.b.validateBlock()
	if err  !=  nil {
		errmsg :=  err.Error()
		logger.Error(fmt.Sprintf("Abbort: Failed valiadtion  due to %s" , errmsg))
		return  errors.New(errmsg)
	}
	logger.Log("Commit  validation of  block ")

	*chain = append( *chain , blockDetails)
	logger.Log("Commit  insert a  new  block in chain ")
	return  nil
}

type  block  struct {
	index int 
	createdAt  int64
	validator  rsa.PublicKey
	capicity int 
	currentHash string
	purrentHash string 
} 

func (b * block)  genesis(validator rsa.PublicKey ,  parrent ,  current string ) error {
	b.index = 0 //first  block 
	b.createdAt = time.Now().Unix() //creation  time  stamp 
	b.validator = validator
	b.capicity = 5 
	b.purrentHash  = parrent 
	b.currentHash = current //later	
	return nil
} 
func (b * block)  validateBlock() error {
	return nil
}

type  blockCoin struct {
	b block
}
func (b * blockCoin ) genesis(validator rsa.PublicKey ,  parrent ,  current string ) error  {
	return  b.b.genesis(validator ,parrent , current)
}


type  blockChainBlock interface  {
	genesis(validator  rsa.PublicKey ,parrent ,  current string )
	validateBlock() error

}
type  blockChainService  interface {
	Service
	genesis() error 

}
type  blockChainCoinsImpl  struct {

	chain  blockChainCoins 
	services blockServiceProviders

} 
type  blockServiceProviders struct {

	loggerService  LogerService
	walletService  WalletService // i will  wand  my rsa for  the brodcasted  block
	hashService  hashService 
}

const  blockChainServiceName = "blockChainService"
func ( p *  blockServiceProviders )  construct () error  {
	var err error 
	if  p.loggerService == nil {
		p.loggerService = &Logger{ServiceName: blockChainServiceName}
		err = p.loggerService.construct()
		if  err != nil {
			return  err
		}
	 }
   return p.valid()


   }   
func ( p *  blockServiceProviders )  valid () error  {

	if  p.loggerService == nil {
		const  errmsg =  "The are  is no loggerService"
		return errors.New(errmsg)
	}
	
	if  p.walletService == nil {
		const  errmsg =  "The are  is no walletService"
		return errors.New(errmsg)
	}
	if  p.hashService == nil {
		const  errmsg =  "The are  is no hashService"
		return errors.New(errmsg)
	}
	return  nil

}   

func (service  *  blockChainCoinsImpl ) construct() error  {
	err := service.services.construct()
	if err != nil {
		return err
	}  
	logger := service.services.loggerService 
	logger.Log("Service  Created")
    return nil
}
func (service  *  blockChainCoinsImpl )  genesis() error   {
	err := service.services.valid() 
	if err != nil {
		return err
	}

	service.chain.chainGenesis(service.services, * service.services.walletService.getPub())
	return  nil
}
