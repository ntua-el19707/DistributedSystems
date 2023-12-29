package  services
//Transaction Manager Test 
import (
    "fmt"
    "testing"
)

//Build The World

/**
    serviceCreater  create  the  service  and return  istances 
    @Returns  TransactionManagerService
*/
func  serviceCreateTransactionManager() (TransactionManagerService , * mockLogger , *  mockWallet , error) {
    
    logger := &mockLogger {}
    wallet := &mockWallet {}
    transactionManager := &TransactionManager{loggerService: logger , walletService: wallet }
    return  transactionManager ,  logger , wallet , transactionManager.construct()
}

/**
    TestCreateService  - test  if  inatce  is  created
*/
func  TestCreateSevice(t * testing.T) {
    _ , logger ,_ ,err := serviceCreateTransactionManager()
    if err != nil {
        t.Errorf("It should get no err  but  got %s" , err.Error())
    }
    if  len(logger.logs) !=  1 {
        t.Errorf("It should log 1 message  but instead log  %d" ,len(logger.logs))
    }
    const expected = "Service  created"

    if logger.logs[0] != expected {

        t.Errorf("It should log %s  but instead log  %s" ,expected , logger.logs[0])
    }
    fmt.Println("it should  create TransactionManagerService")
}
/**
    TestCreateServiceFailed  - test  if  instance  is  not  created 
*/
func  TestCreateSeviceFailed(t * testing.T) {
    logger := &mockLogger {}
    transactionManager := &TransactionManager{loggerService: logger , }
    err := transactionManager.construct()
    if err == nil {
       t.Errorf("It should get err but  got  nothing" )
    }
    const expected = "Provider  for  walletService  should  be  given"
    if  err.Error() != expected {
       t.Errorf("It should get err %s but  got %s" ,  expected,  err.Error() )
    } 

    fmt.Println("it should  not  create TransactionManagerService")
}

