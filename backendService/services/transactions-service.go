package  services

import (	
	"crypto/rsa"
	"crypto"
	"crypto/sha256"
	"time"
	"fmt"
	"errors"
	"encoding/json"
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
	SenderRecieverPair  SendReceivePair `json:"fromTo"`
	Nonce  int   `json:"nonce"`
	Transaction_id string `json:"transaction_id"`
	signedTransaction  [] byte 
	Created_at  int64 `json:"created_at"`
	loggerService LogerService 
	balanceService BalanceService
	generatorService  GeneratorService
    walletService  WalletService
}
type  TransactionCoins struct {
	
	TransactionStandard   TransactionsStandard `json:"transaction"`
	Amount  float64 `json:"ammount"`
    Reason  string  `json:"reason"`
}
// TransactionStandard
func (t  *  TransactionsStandard ) construct () error  {
    var  err error 
    //create  genaratorService  if  not exist 
	if  t.generatorService == nil {
		t.generatorService =  &generatorImplementation{ServiceName: generatorServiceName , CharSet:allChars ,}
        err:= t.generatorService.construct() 
		if err != nil {
			return  err
		}
	} 
	
    //ISSUE  id  for  Transaction 
	id :=  t.generatorService.generateId(10)
	t.Transaction_id = id
    
    if  t.loggerService  == nil {

		t.serviceName = fmt.Sprintf("%s_%s" , t.serviceName ,  t.Transaction_id )
        t.loggerService  = &Logger{ServiceName:t.serviceName}
        err =  t.loggerService.construct() 
        if  err != nil {
            return  err
        } 
    }
    logger :=  t.loggerService 
    logger.Log("Start checking  services ")
    err =  t.valid()
    if  err != nil {
        errmsg := fmt.Sprintf("Abbort   services not  valid %s" , err.Error())
        logger.Error(errmsg)
        return  errors.New(logger.SprintErrorf(errmsg))
    }
    logger.Log("Commit Checking  services")
    return nil 
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
    if len(t.Transaction_id ) != 10 {
        const errmsg string =  "transaction id  has  no  length of  10 probably not  created "
        return errors.New(errmsg )
    }
    return  nil 
}  
/**
	CreateTransaction -  crreate  a  transaction 
	@Returns  error   
*/
func  (transaction *  TransactionCoins ) construct() error{
	
    err := transaction.TransactionStandard.construct()
    if err != nil {
        return err 
    }
	createdMessage := fmt.Sprintf("Created  service : %s \n" , transaction.TransactionStandard.
	serviceName )
    loggerService := transaction.TransactionStandard.loggerService
	loggerService.Log(createdMessage)
	return  nil 
}
/**
	CreateTransaction -  create  a  transaction 
	@Returns  error   
*/
func  (t *  TransactionCoins ) CreateTransaction()  error {
	transaction := &t.TransactionStandard
	logger := transaction.loggerService
    err := transaction.valid()
    if err !=  nil {
        return err
    }
	
    t.TransactionStandard.SenderRecieverPair.Sender.Address = * transaction.walletService.getPub() 

	transactionDetails :=  fmt.Sprintf("%s From %v to %v  For  %f" , transaction.Transaction_id , transaction.SenderRecieverPair.Sender.Address , transaction.SenderRecieverPair.Receiver.Address , t.Amount )
	
	logger.Log(fmt.Sprintf("Start creating  Transaction %s" , transactionDetails))
	//issue  time to transaction
	transaction.Created_at = time.Now().Unix()

	amount := t.Amount
    //check if acount has the amount  

	sender := transaction.SenderRecieverPair.Sender.Address
	//? i need  a service  to find  the balance  of the sender
    
    transaction.balanceService.LockBalance() // ensure  the  balance  will not  change 
    //wallet will have  balance -frozenMoney coins  that 
    //trnsction  is  not yet  added  in chain 
    //Be  optimist and  froze  the  coins
    transaction.walletService.Freeze(amount)
    //find  balance 
	balance ,err := transaction.balanceService.findBalance(sender)
	if err != nil {
    //failed  load  balance  so transaction will  error =>  unfroze  the  money trnsaction will not  happen 
    transaction.balanceService.UnLockBalance() 
    transaction.walletService.UnFreeze(amount)
		return err
	}
    frozenMoney := transaction.walletService.getFreeze()
	if balance - frozenMoney    <  0 {
		message := fmt.Sprintf("Request To  sent  %.3f  from  %.3f  balance Failed  due to total Money Froze(for wallet  ) %.3f\n" , amount ,  balance , frozenMoney)
   // not  valid  trancation the total 
        transaction.balanceService.UnLockBalance() 
        transaction.walletService.UnFreeze(amount)

		return errors.New(logger.SprintErrorf(message))
       
	}
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s" , transactionDetails))
    transaction.balanceService.UnLockBalance() 
	
	return  nil 
}
/**
	setSign  -  set the signiture
	@Param  siginiture  string

*/
func  (t *  TransactionCoins ) setSign(signature  [] byte){
	t.TransactionStandard.signedTransaction  = signature
}
/**
	getAmount  -  get amount 
	@Param  siginiture  string

*/
func  (t  TransactionCoins ) getAmount() float64{
	return t.Amount
}
/**
	getTransaction  -  get the transaction
	@Param  siginiture  string

*/
func  (t  TransactionCoins ) getTransaction() ([]  byte , error)  {
   	//*The  signature  must not  be in document IMPORTANT
	data ,  err := json.Marshal(t) 
	if  err != nil {
		return nil ,err
	}

	return data, nil
}
/**
	VerifySignature  -  verify the transaction
	@Param  siginiture  string

*/
func  (t  TransactionCoins ) VerifySignature(PublicKey  *  rsa.PublicKey   ) error  {
	//define  verify 
	verify :=   func(transaction []  byte ) error {
		hashed := sha256.Sum256(transaction)
    
		err := rsa.VerifyPKCS1v15(PublicKey, crypto.SHA256, hashed[:], t.TransactionStandard.signedTransaction)
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
	return t.TransactionStandard.signedTransaction}

type  TransactionMsg struct {	
	TransactionStandard   TransactionsStandard `json:"transaction"`
	Msg  string `json:"msg"`
 }
/**
	construct - create  transactionMsg  service 
	@Returns  error   
*/
func  (transaction *  TransactionMsg) construct() error{
	
    err := transaction.TransactionStandard.construct()
    if err != nil {
        return err 
    }
	createdMessage := fmt.Sprintf("Created  service : %s \n" , transaction.TransactionStandard.
	serviceName )
    loggerService := transaction.TransactionStandard.loggerService
	loggerService.Log(createdMessage)
	return  nil 
}
/**
	CreateTransaction -  create  a  transaction 
	@Returns  error   
*/
func   (t *TransactionMsg) CreateTransaction()  error {
	transaction := &t.TransactionStandard
    
    err := transaction.valid()
    if err !=  nil {
        return err
    }
	logger := transaction.loggerService
    t.TransactionStandard.SenderRecieverPair.Sender.Address = * transaction.walletService.getPub() 
    transactionDetails :=  fmt.Sprintf("%s From %v to %v  For  %s" , transaction.Transaction_id , transaction.SenderRecieverPair.Sender.Address , transaction.SenderRecieverPair.Receiver.Address , t.Msg )
	
	logger.Log(fmt.Sprintf("Start creating  Transaction %s" , transactionDetails))
	//issue  time to transaction
	transaction.Created_at = time.Now().Unix()
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s" , transactionDetails))
	return  nil 
}
/**
	setSign  -  set the signiture
	@Param  siginiture  string

*/
func  (t *  TransactionMsg ) setSign(signature  [] byte){
	t.TransactionStandard.signedTransaction  = signature
}
/**
	getTransaction  -  get the transaction
	@Param  siginiture  string

*/
func  (t  TransactionMsg ) getTransaction() ([]  byte , error)  {
   	//*The  signature  must not  be in document IMPORTANT
	data ,  err := json.Marshal(t) 
	if  err != nil {
		return nil ,err
	}

	return data, nil
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
    
		err := rsa.VerifyPKCS1v15(PublicKey, crypto.SHA256, hashed[:], t.TransactionStandard.signedTransaction)
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
	return t.TransactionStandard.signedTransaction
}


