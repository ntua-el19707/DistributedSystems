package  services

import (	
	"crypto/rsa"
	"crypto"
	"crypto/sha256"
	"time"
	"fmt"
	"errors"
	"entitys"
)
type  TransactionService interface {
	Service
	CreateTransaction() error 
	VerifySignature(PublicKey  *  rsa.PublicKey   ) error
	setSign(signiture  [] byte)
    getTransaction() ([] byte ,error) 
	getSigniture() [] byte
	getAmount() float64
}
type  Person struct{
	Address rsa.PublicKey `json:"address"`
} 
type  SendReceivePair struct {
	Sender Person `json:"sender"`
	Receiver Person `json:"receiver"`
}
type  TransactionsStandard struct {
	serviceName string 
	loggerService LogerService 
	balanceService BalanceService
	generatorService  GeneratorService
    walletService  WalletService
}
type  TransactionCoins struct {
	Transaction	 entitys.TransactionCoinEntityRoot
	services	TransactionsStandard

}
// TransactionStandard
func (t  *  TransactionsStandard ) construct () (string , error)  {
    var  err error 
    //create  genaratorService  if  not exist 
	if  t.generatorService == nil {
		t.generatorService =  &generatorImplementation{ServiceName: generatorServiceName , CharSet:allChars ,}
        err:= t.generatorService.construct() 
		if err != nil {
			return  "" , err
		}
	} 
	
    //ISSUE  id  for  Transaction 
	id :=  t.generatorService.generateId(10)
    
    if  t.loggerService  == nil {

		t.serviceName = fmt.Sprintf("%s_%s" , t.serviceName ,  id)
        t.loggerService  = &Logger{ServiceName:t.serviceName}
        err =  t.loggerService.construct() 
        if  err != nil {
            return  "" , err
        } 
    }
    logger :=  t.loggerService 
    logger.Log("Start checking  services ")
    err =  t.valid()
    if  err != nil {
        errmsg := fmt.Sprintf("Abbort   services not  valid %s" , err.Error())
        logger.Error(errmsg)
        return  "" , errors.New(logger.SprintErrorf(errmsg))
    }
    logger.Log("Commit Checking  services")
    return id ,  nil 
} 
func  (t *  TransactionsStandard ) valid ()  error  {
    if  t.loggerService  == nil   {
        const errmsg string =  "transaction standard has  no  loggerService "
        return errors.New(errmsg )
    }
    if  t.balanceService  == nil   {
        const errmsg string =  "transaction standard has  no  balanceService "
        return errors.New(errmsg )
    }
    if  t.generatorService  == nil   {
        const errmsg string =  "transaction standard has  no  generatorService "
        return errors.New(errmsg )
    }
    if  t.walletService  == nil   {
        const errmsg string =  "transaction standard has  no  walletService "
        return errors.New(errmsg )
    }

    return  nil 
}  
/**
	CreateTransaction -  crreate  a  transaction 
	@Returns  error   
*/
func  (transaction *  TransactionCoins ) construct() error{
	
    id , err := transaction.services.construct()
    if err != nil {
        return err 
    }
	transaction.Transaction.Transaction.BillDetails.Transaction_id = id
	createdMessage := fmt.Sprintf("Created  service : %s \n" , transaction.services.
	serviceName )
    loggerService := transaction.services.loggerService
	loggerService.Log(createdMessage)
	return  nil 
}
/**
	CreateTransaction -  create  a  transaction 
	@Returns  error   
*/
func  (t *  TransactionCoins ) CreateTransaction()  error {
	transaction := &t.Transaction.Transaction
	services  :=  &t.services
	logger := t.services.loggerService
    err := services.valid()
    if err !=  nil {
        return err
    }
	
   	transaction.BillDetails.Bill.From.Address = * services.walletService.getPub() 
	bill := transaction.BillDetails.Bill
	transactionDetails :=  fmt.Sprintf("%s From %v to %v  For  %f" , transaction.BillDetails.Transaction_id , bill.From.Address , bill.To.Address , transaction.Amount )
	
	logger.Log(fmt.Sprintf("Start creating  Transaction %s" , transactionDetails))
	//issue  time to transaction
	transaction.BillDetails.Created_at = time.Now().Unix()

	amount := transaction.Amount
    //check if acount has the amount  

	sender := transaction.BillDetails.Bill.From.Address
	//? i need  a service  to find  the balance  of the sender
    
    services.balanceService.LockBalance() // ensure  the  balance  will not  change 
    //wallet will have  balance -frozenMoney coins  that 
    //trnsction  is  not yet  added  in chain 
    //Be  optimist and  froze  the  coins
    services.walletService.Freeze(amount)
    //find  balance 
	balance ,err := services.balanceService.findBalance(sender)
	if err != nil {
    //failed  load  balance  so transaction will  error =>  unfroze  the  money trnsaction will not  happen 
    services.balanceService.UnLockBalance() 
    services.walletService.UnFreeze(amount)
		return err
	}
    frozenMoney := services.walletService.getFreeze()
	if balance - frozenMoney    <  0 {
		message := fmt.Sprintf("Request To  sent  %.3f  from  %.3f  balance Failed  due to total Money Froze(for wallet  ) %.3f\n" , amount ,  balance , frozenMoney)
   // not  valid  trancation the total 
        services.balanceService.UnLockBalance() 
        services.walletService.UnFreeze(amount)

		return errors.New(logger.SprintErrorf(message))
       
	}
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s" , transactionDetails))
    services.balanceService.UnLockBalance() 
	
	return  nil 
}
/**
	setSign  -  set the signiture
	@Param  siginiture  string

*/
func  (t *  TransactionCoins ) setSign(signature  [] byte){
	t.Transaction.Signiture = signature
}
/**
	getAmount  -  get amount 
	@Param  siginiture  string

*/
func  (t  TransactionCoins ) getAmount() float64{
	return t.Transaction.Transaction.Amount
}
/**
	getTransaction  -  get the transaction
	@Param  siginiture  string

*/
func  (t  TransactionCoins ) getTransaction() ([]  byte , error)  {
   	//*The  signature  must not  be in document IMPORTANT
	return  entitys.JsonStringfy(t.Transaction.Transaction)
}
/**
	VerifySignature  -  verify the transaction
	@Param  siginiture  string

*/
func  (t  TransactionCoins ) VerifySignature(PublicKey  *  rsa.PublicKey   ) error  {
	//define  verify 
	verify :=   func(transaction []  byte ) error {
		hashed := sha256.Sum256(transaction)
    
		err := rsa.VerifyPKCS1v15(PublicKey, crypto.SHA256, hashed[:], t.getSigniture())
		if  err != nil {
			return  err
		}
		return  nil
	}
	data ,  err := t.getTransaction()
	if  err != nil {
		return err
	}
	return  verify(data)
}
func  (t  TransactionCoins ) getSigniture() ([]  byte )  {
	return t.Transaction.Signiture}

type  TransactionMsg struct {	
	Transaction	 entitys.TransactionMsgEntityRoot
	services	TransactionsStandard
 }
/**
	construct - create  transactionMsg  service 
	@Returns  error   
*/
func  (transaction *  TransactionMsg) construct() error{
    id , err := transaction.services.construct()
    if err != nil {
        return err 
    }
	transaction.Transaction.Transaction.BillDetails.Transaction_id = id
	createdMessage := fmt.Sprintf("Created  service : %s \n" , transaction.services.
	serviceName )
    loggerService := transaction.services.loggerService
	loggerService.Log(createdMessage)
	return  nil 
}
/**
	CreateTransaction -  create  a  transaction 
	@Returns  error   
*/
func   (t *TransactionMsg) CreateTransaction()  error {
	transaction := &t.Transaction.Transaction
	services  :=  &t.services
	logger := t.services.loggerService
    err := services.valid()
    if err !=  nil {
        return err
    }
	
   	transaction.BillDetails.Bill.From.Address = * services.walletService.getPub() 
	bill := transaction.BillDetails.Bill
	transactionDetails :=  fmt.Sprintf("%s From %v to %v  For  %s" , transaction.BillDetails.Transaction_id , bill.From.Address , bill.To.Address , transaction.Msg )
	
	logger.Log(fmt.Sprintf("Start creating  Transaction %s" , transactionDetails))
	//issue  time to transaction
	transaction.BillDetails.Created_at = time.Now().Unix()
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s" , transactionDetails))
	return  nil 
}
/**
	setSign  -  set the signiture
	@Param  siginiture  string

*/
func  (t *  TransactionMsg ) setSign(signature  [] byte){
	t.Transaction.Signiture = signature}
/**
	getTransaction  -  get the transaction
	@Param  siginiture  string

*/
func  (t  TransactionMsg ) getTransaction() ([]  byte , error)  {
   	//*The  signature  must not  be in document IMPORTANT
return  entitys.JsonStringfy(t.Transaction.Transaction)
}
/**
	getAmount  -  get amount 
	

*/
func  (t   TransactionMsg  ) getAmount() float64{
	return float64(0)
}
/**
	VerifySignature  -  verify the transaction
	@Param  siginiture  string

*/
func  (t  TransactionMsg ) VerifySignature(PublicKey  *  rsa.PublicKey   ) error  {
	//define  verify 
	verify :=   func(transaction []  byte ) error {
		hashed := sha256.Sum256(transaction)
    
		err := rsa.VerifyPKCS1v15(PublicKey, crypto.SHA256, hashed[:], t.getSigniture())
		if  err != nil {
			return  err
		}
		return  nil
	}
	data ,  err := t.getTransaction()
	if  err != nil {
		return err
	}
	return  verify(data)
}
func  (t  TransactionMsg ) getSigniture() ([]  byte )  {
	return t.Transaction.Signiture
}


