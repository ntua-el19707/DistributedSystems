package services

import (
	"fmt"
	"testing"
)

func  buildTransactionService()( TransactionService ,  * mockLogger , * mockGenerator  , * mockFindBalance , *TransactionCoins , error ){
	
	mockLogger :=  &mockLogger{}
	mockGenerator :=  &mockGenerator{response:"T1"}
	mockFindBalance := &mockFindBalance{}
	transactionStandard := TransactionsStandard{serviceName:transactionServiceName , loggerService:mockLogger ,  balanceService:mockFindBalance ,generatorService:mockGenerator }
	transaction :=  &TransactionCoins{TransactionStandard:transactionStandard, }
	transactionService := transaction
	err := transactionService.construct()
	if err != nil {
		return nil , nil , nil , nil , nil,err
	}

	return  transactionService ,  	mockLogger ,  mockGenerator ,mockFindBalance , transaction, nil


} 
func  TestCreateServiceTransaction( t * testing.T){
	_,logger,generator,_ ,_,err := buildTransactionService()
	if err != nil {
		t.Errorf("Failed  to create  trnsaction service  due  to %s" ,  err.Error())
	}
	if generator.timesCallgenerateId !=1  {
		t.Errorf("Generate  id function  should  be  call  1 time  but get called %d" ,generator.timesCallgenerateId  )
	}
	if len(logger.logs) != 1 {
		t.Errorf("Logger logs  should  have  1  message but  have  %d messages " , len(logger.logs)  )
	}
	msg := fmt.Sprintf("Created  service : %s \n" , transactionServiceName ,)
	if  logger.logs[0] != msg {
		t.Errorf("message  should  be  %s  but  got  %s  " , msg ,  logger.logs[0]  )
	}

	fmt.Println("Create Transaction Service")
}
func TestCreateTransactionInvalid( t * testing.T ){
	service,_,_,findBalance, transaction ,_ := buildTransactionService()
	findBalance.amount = 10 
    transaction.Amount = -15 
	errmsg:= fmt.Sprintf("Request To  sent  %f  from  %f  balance\n",  transaction.Amount ,findBalance.amount  )
	err := service.CreateTransaction() 
	if err == nil {
		t.Errorf("It should be  invalid")
	}
	if err.Error() != errmsg {
		t.Errorf("It should get err: %s  but  got %s" ,  errmsg ,  err.Error())
	}
	fmt.Println("Should  not  create  invalid transaction")

}
func TestCreateTransactionvalid( t * testing.T ){
	service,logger,_,findBalance, transaction ,_ := buildTransactionService()
	findBalance.amount = 100 
    transaction.Amount = -15 
	logger.logs = make([] string , 0)
	err := service.CreateTransaction()
	if err != nil {
		t.Errorf("It should be  valid not  get  err  %s" , err.Error())
	}
	if  len(logger.logs) != 2 {
		t.Errorf("Logger  should  receive 2  messages  but intead  got  %d " , len(logger.logs))
	}
	if findBalance.findBalanceCalledTimes !=1 {
		t.Errorf("findbalance should  be  called  once but intead  called  %d " ,findBalance.findBalanceCalledTimes)
	}
	fmt.Println("Should    create  valid transaction")

} 