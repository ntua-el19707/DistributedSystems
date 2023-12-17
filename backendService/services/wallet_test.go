package  services 

import 
(
    "testing"
    "fmt"
)
func TestWalletCreateServiceVersionStruct1(t * testing.T){
   var walletService   walletStructV1Service 
   err :=  walletService.construct()
   if  err != nil {
    t.Errorf("%s" ,  err.Error())
   }
   fmt.Println("it  should  create  a new  wallet service vesrion struct 1 ")
    
}

