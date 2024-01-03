package services

import (
	"fmt"
	"testing"
)

func  buildTransactionMsgService()( TransactionService ,  * mockLogger , * mockGenerator  , * mockFindBalance , *TransactionMsg, *  mockWallet , error ){
	
	mockLogger :=  &mockLogger{}
	mockGenerator :=  &mockGenerator{response:"aaaaabbbbb"}
	mockFindBalance := &mockFindBalance{} 
    wallet := &mockWallet{} 
    transactionStandard := TransactionsStandard{serviceName:transactionServiceName ,walletService:wallet, loggerService:mockLogger ,  balanceService:mockFindBalance ,generatorService:mockGenerator }
	transaction :=  &TransactionMsg{services:transactionStandard, }
	transactionService := transaction
	err := transactionService.construct()
	if err != nil {
		return nil , nil , nil , nil , nil ,nil ,err
	}

	return  transactionService ,  	mockLogger ,  mockGenerator ,mockFindBalance , transaction , wallet, nil
}
func  buildTransactionService()( TransactionService ,  * mockLogger , * mockGenerator  , * mockFindBalance , *TransactionCoins, *  mockWallet , error ){
	
	mockLogger :=  &mockLogger{}
	mockGenerator :=  &mockGenerator{response:"aaaaabbbbb"}
	mockFindBalance := &mockFindBalance{} 
    wallet := &mockWallet{} 
    transactionStandard := TransactionsStandard{serviceName:transactionServiceName ,walletService:wallet, loggerService:mockLogger ,  balanceService:mockFindBalance ,generatorService:mockGenerator }
	transaction :=  &TransactionCoins{services:transactionStandard, }
	transactionService := transaction
	err := transactionService.construct()
	if err != nil {
		return nil , nil , nil , nil , nil ,nil ,err
	}

	return  transactionService ,  	mockLogger ,  mockGenerator ,mockFindBalance , transaction , wallet, nil


} 
func  TestCreateServiceTransactionMsg( t * testing.T){
    _, logger  ,_ , _ ,_ ,_ ,err := buildTransactionMsgService() 
    if  err != nil {
        t.Errorf("Expceted  no err  but  got %v" ,err)
    }
    if len(logger.logs) !=3 {
        t.Errorf("Expceted  to log 3  msg   but  got %d" ,len(logger.logs))
    }
	fmt.Println("created Transaction Msg Service")
}

func  TestCreateServiceTransactionMsgCreateValidTransaction( t * testing.T){
    tr, _  ,_ , _ ,_ ,_ ,err := buildTransactionMsgService() 
    err = tr.CreateTransaction()
    if  err != nil {
        t.Errorf("Expceted  no err  but  got %v" ,err)
    }
	fmt.Println("create transaction Msg Service")
}
func  TestCreateServiceTransaction( t * testing.T){
	_,logger,generator,_ ,_,_,err := buildTransactionService()
	if err != nil {
		t.Errorf("Failed  to create  trnsaction service  due  to %s" ,  err.Error())
	}
	if generator.timesCallgenerateId !=1  {
		t.Errorf("Generate  id function  should  be  call  1 time  but get called %d" ,generator.timesCallgenerateId  )
	}
	if len(logger.logs) != 3 {
		t.Errorf("Logger logs  should  have  3  message but  have  %d messages " , len(logger.logs)  )
	}
	msg := fmt.Sprintf("Created  service : %s \n" , transactionServiceName ,)
	if  logger.logs[2] != msg {
		t.Errorf("message  should  be  %s  but  got  %s  " , msg ,  logger.logs[2]  )
	}

	fmt.Println("Create Transaction Service")
}
func TestCreateTransactionInvalid( t * testing.T ){
	service,_,_,findBalance, transaction , wallet ,_ := buildTransactionService()
	findBalance.amount = 10 
    transaction.Transaction.Transaction.Amount = 15 
    wallet.frozen = 15 
    errmsg := fmt.Sprintf("Request To  sent  %.3f  from  %.3f  balance Failed  due to total Money Froze(for wallet  ) %.3f\n" , transaction.Transaction.Transaction.Amount,  findBalance.amount  , wallet.frozen )	
	err := service.CreateTransaction() 
	if err == nil {
		t.Errorf("It should be  invalid")
	}
	if err.Error() != errmsg {
		t.Errorf("It should get err: %s  but  got %s" ,  errmsg ,  err.Error())
	}
    if findBalance.locked  {
		t.Errorf("It should call first  lock and  then unlock " )

    }
    if  findBalance.lockedCall !=  1 {
		t.Errorf("It should locked balance once  but lock %d " , findBalance.lockedCall  )

    }
    if  findBalance.unlockedCall !=  1 {
		t.Errorf("It should unlocked balance once  but unlock %d " , findBalance.unlockedCall  )

    }
    if  wallet.counterFreeze != 1 {
		t.Errorf("It should freeze money  once  but freeze %d " , wallet.counterFreeze  )

    } 
    if  wallet.counterUnFreeze != 1 {
		t.Errorf("It should un freeze money  once  but un  freeze %d " , wallet.counterUnFreeze  )

    } 
    if  wallet.countergetFreeze != 1 {
		t.Errorf("It should get freeze money  once  but get  freeze %d " , wallet.countergetFreeze  )

    } 

    fmt.Println("Should  not  create  invalid transaction")

}
func TestCreateTransactionvalid( t * testing.T ){
	service,logger,_,findBalance, transaction , wallet ,_ := buildTransactionService()
	findBalance.amount = 100 
    transaction.Transaction.Transaction.Amount = 15 
    wallet.frozen  = transaction.Transaction.Transaction.Amount
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
    if findBalance.locked  {
		t.Errorf("It should call first  lock and  then unlock " )

    }
    if  findBalance.lockedCall !=  1 {
		t.Errorf("It should locked balance once  but lock %d " , findBalance.lockedCall  )

    }
    if  findBalance.unlockedCall !=  1 {
		t.Errorf("It should unlocked balance once  but unlock %d " , findBalance.unlockedCall  )

    }
    if  wallet.counterFreeze != 1 {
		t.Errorf("It should freeze money  once  but freeze %d " , wallet.counterFreeze  )

    } 
    if  wallet.counterUnFreeze != 0 {
		t.Errorf("It should not un freeze money  but did un   freeze %d " , wallet.counterUnFreeze  )

    } 
    if  wallet.countergetFreeze != 1 {
		t.Errorf("It should get freeze money  once  but get  freeze %d " , wallet.countergetFreeze  )

    } 
	fmt.Println("Should    create  valid transaction")

} 

