package  services 

import 
(
    "testing"
    "fmt"
)

func buildWalletService()(  WalletService ,  *mockLogger  ,  error ) {
    mockLogger :=  &mockLogger{} 

    walletService := &walletStructV1Service{loggerService:mockLogger}
    err := walletService.construct()
    if err !=nil {
        return  nil , nil , err
    }
    return  walletService ,  mockLogger ,  nil


}
func TestWalletCreateServiceVersionStruct1(t * testing.T){
 
   _ , logger ,  err :=  buildWalletService()
   if  err != nil {
    t.Errorf("%s" ,  err.Error())
   }
   if len(logger.logs) !=  4 {
        t.Errorf("It  should  have  been 4  messages  to logs  but  got %d" ,  len(logger.logs))
   }
   fmt.Println("it  should  create  a new  wallet service vesrion struct 1 ")
    
}


func TestWalletSign(t * testing.T){
 
   service , _,  _ :=  buildWalletService()
   tservice ,_ ,_,_ ,transaction ,_  ,_ :=  buildTransactionService()
   transaction.Transaction.Transaction.Amount = 10
   err := tservice.CreateTransaction()
   if err!=nil {
    t.Error("failed to create  transaction" + err.Error())
   }
   err = service.sign(tservice)
   if err != nil {
    t.Error("Failed to sign  Error " + err.Error())
   }
   err = tservice.VerifySignature(service.getPub())
   if err != nil {
    t.Error("Failed to verify Error " + err.Error())
   }

   fmt.Println("it  should  sign  transaction ")
    
}
