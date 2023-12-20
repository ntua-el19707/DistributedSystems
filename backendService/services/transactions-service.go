package  services

import (	
	"crypto/rsa"
	"crypto"
	"crypto/sha256"
	"time"
	"fmt"
	"errors"
	"encoding/json"
)
type  TransactionService interface {
	Service
	CreateTransaction() error 
	VerifySignature(PublicKey  *  rsa.PublicKey   ) error
	setSign(signiture  [] byte)
getTransaction() ([] byte ,error) 
	getSigniture() [] byte
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
}
type  TransactionCoins struct {
	
	TransactionStandard   TransactionsStandard `json:"transaction"`
	Amount  float64 `json:"ammount"`

}
/**
	CreateTransaction -  crreate  a  transaction 
	@Returns  error   
*/
func  (transaction *  TransactionCoins ) construct() error{
	//var  zero TransactionsStandard 
	/*if transaction.TransactionStandard == zero {
		transaction.TransactionStandard = TransactionsStandard{}
	}*/
	var loggerService LogerService = transaction.TransactionStandard.loggerService

	
	if  transaction.TransactionStandard.balanceService == nil {
		transaction.TransactionStandard.balanceService =  &balanceImplementation{}
		err := transaction.TransactionStandard.balanceService.construct() 
		if err != nil {
			return  err
		}

	} 
	
	if  transaction.TransactionStandard.generatorService == nil {
		transaction.TransactionStandard.generatorService =  &generatorImplementation{ServiceName: generatorServiceName , CharSet:allChars ,}
		err := transaction.TransactionStandard.generatorService.construct() 
		if err != nil {
			return  err
		}
	} 

	//ISSUE  id  for  Transaction 
	id :=  transaction.TransactionStandard.generatorService.generateId(10)
	transaction.TransactionStandard.Transaction_id = id
	
	if  transaction.TransactionStandard.loggerService == nil {
		transaction.TransactionStandard.serviceName = fmt.Sprintf("%s_%s" , transaction.TransactionStandard.serviceName ,  transaction.TransactionStandard.Transaction_id )
		transaction.TransactionStandard.loggerService = &Logger{ServiceName:transaction.TransactionStandard.serviceName}
		err := transaction.TransactionStandard.loggerService.construct()
		if err != nil {
			return err
		}
		loggerService = transaction.TransactionStandard.loggerService 
	}
	
	createdMessage := fmt.Sprintf("Created  service : %s \n" , transaction.TransactionStandard.
	serviceName )
	loggerService.Log(createdMessage)
	return  nil 
}
/**
	CreateTransaction -  create  a  transaction 
	@Returns  error   
*/
func  (t *  TransactionCoins ) CreateTransaction()  error {
	transaction := t.TransactionStandard
	logger := transaction.loggerService
	transactionDetails :=  fmt.Sprintf("%s From %v to %v  For  %f" , transaction.Transaction_id , transaction.SenderRecieverPair.Sender.Address , transaction.SenderRecieverPair.Receiver.Address , t.Amount )
	
	logger.Log(fmt.Sprintf("Start creating  Transaction %s" , transactionDetails))
	//issue  time to transaction
	transaction.Created_at = time.Now().Unix()
	
	amount := t.Amount
    //check if acount has the amount  

	sender := transaction.SenderRecieverPair.Sender.Address
	//? i need  a service  to find  the balance  of the sender
	balance ,err := transaction.balanceService.findBalance(sender)
	if err != nil {
		return err
	}
	if amount + balance  <  0 {
		message := fmt.Sprintf("Request To  sent  %f  from  %f  balance\n" , amount ,  balance)
		return errors.New(logger.Sprintf(message))
	}
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s" , transactionDetails))
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