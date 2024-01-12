package  services 

import (
	"testing"
	"fmt"
    "entitys"
	"crypto/rand"
	"crypto/rsa"
)
/** 
    create  stake  service  
*/ 
func  createStackService() ( stakeService ,  * mockLogger  ,  * stakeCoinBlockChain , error  ) {
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    err :=  impl.construct() 
    return   impl,  mockLogger , impl , err
} 

func TestCreateStakeService(t * testing.T) {
    _,logger,_,err:= createStackService() 
    if  err  != nil  {
        t.Errorf("Expcected  to not  get err  but  got %v" , err)
    } 
    const expected string = "Service  created"
    if  len(logger.logs) != 1 {
        t.Errorf("Expcected  to  get 1  log message   but  got %d" , len(logger.logs))
    }
    if  logger.logs[0] !=  expected {
        t.Errorf("Expcected  to not  get msg %s  but  got %s" , expected , logger.logs[0])
    }
    fmt.Println("it  should  CreateStakeService")
}
func TestFailCreateStakeService(t * testing.T) {
   
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    impl.block.b.capicity = 1 
    err :=  impl.construct() 
    const expected  string  =  "Block  is  not  full cappicity  is  1 but  has  0 "

    if  err.Error()   != expected   {
        t.Errorf("Expcected  to not  get  err : %s  but  got %s" , expected  ,  err.Error())
    } 
    fmt.Println("it  should  failed   CreateStakeService not  full block")
}
func  TestDistribution(t * testing.T){
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    keyGen := func   (n int  ) ([]  rsa.PublicKey ,error){
		var publicKeys [] rsa.PublicKey

        for i := 0; i < n; i++ {
            privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
            if err != nil {
                return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
            }

		publicKeys = append(publicKeys,  privateKey.PublicKey)
           
        }
        return  publicKeys , nil
    }
    impl.workers , _ =  keyGen(2)
    impl.block.b.capicity = 5 

    billDetails1 := entitys.TransactionDetails{}
    billDetails2 := entitys.TransactionDetails{}
    billDetails1.Bill.From.Address =  impl.workers[0] 
    billDetails1.Bill.To.Address =   impl.workers[1]
    
    billDetails2.Bill.To.Address =  impl.workers[0] 
    billDetails2.Bill.From.Address =  impl.workers[1]

    transactions  :=  []  entitys.TransactionCoins{ entitys.TransactionCoins{BillDetails: billDetails1 ,Amount:20.0 } ,  entitys.TransactionCoins{BillDetails: billDetails1 , Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails1 ,  Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails2 , Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails2 ,  Amount:20.0}  }
    impl.block.Transactions  =   transactions
    err := impl.construct()
    if  err != nil {
        t.Errorf("Expected no err  but  got %v" ,err)
    }
    dMap , total := impl.distributionOfStake(0)
    if  dMap[impl.workers[0]] != 60.0 ||  dMap[impl.workers[1]] != 40.0  || total != 100.0 {
        t.Errorf("Expected worker0:%3f worker1:%3f , total %3f but  got  worker0:%3f worker1:%3f , total %3f " , 60.0 ,40.0,100.0 ,   dMap[impl.workers[0]] ,  dMap[impl.workers[1]] ,total)
    }
    fmt.Println("Create Distribution Map ")
}
func  TestDistributionWeight1Half(t * testing.T){
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    keyGen := func   (n int  ) ([]  rsa.PublicKey ,error){
		var publicKeys [] rsa.PublicKey

        for i := 0; i < n; i++ {
            privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
            if err != nil {
                return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
            }

		publicKeys = append(publicKeys,  privateKey.PublicKey)
           
        }
        return  publicKeys , nil
    }
    impl.workers , _ =  keyGen(2)
    impl.block.b.capicity = 5 

    billDetails1 := entitys.TransactionDetails{}
    billDetails2 := entitys.TransactionDetails{}
    billDetails1.Bill.From.Address =  impl.workers[0] 
    billDetails1.Bill.To.Address =   impl.workers[1]
    
    billDetails2.Bill.To.Address =  impl.workers[0] 
    billDetails2.Bill.From.Address =  impl.workers[1]

    transactions  :=  []  entitys.TransactionCoins{ entitys.TransactionCoins{BillDetails: billDetails1 ,Amount:20.0 } ,  entitys.TransactionCoins{BillDetails: billDetails1 , Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails1 ,  Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails2 , Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails2 ,  Amount:20.0}  }
    impl.block.Transactions  =   transactions
    err := impl.construct()
    if  err != nil {
        t.Errorf("Expected no err  but  got %v" ,err)
    }
    dMap , total := impl.distributionOfStake(1.5)
    if  dMap[impl.workers[0]] != 120.0 ||  dMap[impl.workers[1]] != 130.0  || total != 250.0 {
        t.Errorf("Expected worker0:%3f worker1:%3f , total %3f but  got  worker0:%3f worker1:%3f , total %3f " , 120.0 ,130.0,250.0 ,   dMap[impl.workers[0]] ,  dMap[impl.workers[1]] ,total)
    }
    fmt.Println("Create Distribution Map weight 1.5  ")
}

func  TestDistributionWeight1HalfIntMap(t * testing.T){
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    keyGen := func   (n int  ) ([]  rsa.PublicKey ,error){
		var publicKeys [] rsa.PublicKey

        for i := 0; i < n; i++ {
            privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
            if err != nil {
                return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
            }

		publicKeys = append(publicKeys,  privateKey.PublicKey)
           
        }
        return  publicKeys , nil
    }
    impl.workers , _ =  keyGen(2)
    impl.block.b.capicity = 5 

    billDetails1 := entitys.TransactionDetails{}
    billDetails2 := entitys.TransactionDetails{}
    billDetails1.Bill.From.Address =  impl.workers[0] 
    billDetails1.Bill.To.Address =   impl.workers[1]
    
    billDetails2.Bill.To.Address =  impl.workers[0] 
    billDetails2.Bill.From.Address =  impl.workers[1]

    transactions  :=  []  entitys.TransactionCoins{ entitys.TransactionCoins{BillDetails: billDetails1 ,Amount:20.0 } ,  entitys.TransactionCoins{BillDetails: billDetails1 , Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails1 ,  Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails2 , Amount:20.0} , entitys.TransactionCoins{BillDetails: billDetails2 ,  Amount:20.0}  }
    impl.block.Transactions  =   transactions
    err := impl.construct()
    if  err != nil {
        t.Errorf("Expected no err  but  got %v" ,err)
    }
    dMap , total :=impl.MapOfDistibutesRoundUp (1.5)
    if  dMap[impl.workers[0]] != 48000 ||  dMap[impl.workers[1]] != 52000 || total != 100000 {
        t.Errorf("Expected iworker0:%d worker1:%d , total %d but  got  worker0:%d worker1:%d , total %d " , 48000,52000,100000 ,   dMap[impl.workers[0]] ,  dMap[impl.workers[1]] ,total)
    }
    fmt.Println("Create Distribution Mapi Rouned  Up  weight 1.5  ")
}
func  TestGetCurrentHash(t * testing.T) {
   
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    err := impl.construct() 
    if  err != nil  {
        t.Errorf("Expected to  not  get  err  but  got  %v" ,err)
    } 
    impl.block.b.currentHash = "aa"
    actual := impl.getCurrentHash()
    if  actual != "aa"{
        t.Errorf("Expected to  not  get  'aa' but  got  '%s'" ,actual)
    }

    fmt.Println("it  should  get  correct  hash aa")

}

func  TestGetCurrentHash2(t * testing.T) {
   
    mockLogger := &mockLogger{} 
    impl  := &stakeCoinBlockChain{services:stakeProviders{loggerService:mockLogger}}
    err := impl.construct() 
    if  err != nil  {
        t.Errorf("Expected to  not  get  err  but  got  %v" ,err)
    } 
    impl.block.b.currentHash = "bb"
    actual := impl.getCurrentHash()
    if  actual != "bb"{
        t.Errorf("Expected to  not  get  'bb' but  got  '%s'" ,actual)
    }

    fmt.Println("it  should  get  correct  hash bb")

}
