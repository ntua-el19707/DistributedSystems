package  services
import (
    "errors"
    "crypto/rsa"
    "fmt"
    "entitys"
)


const TransactionManagerServiceName = "TransactionManagerService"
// TransactionMangerService  interface  Rules 
type TransactionManagerService interface {
    Service // it  is  indeed a  service 
    TransferMoney(to , validator  * rsa.PublicKey ,  ammount  float64)( [] TransactionService , error )  
    SendMessage( to , validatorMsg ,validatorCoins *  rsa.PublicKey ,  msg string   ) ([]  TransactionService  ,  error  ) 
    unValidService() error 
}

// TransactionManger  struct 

type  TransactionManager struct {
    walletService WalletService 
    loggerService LogerService 
    balanceService BalanceService 
}
// Construct 
func ( transactionManager *  TransactionManager ) construct () error {
    // if  logger service  nil  =>  create  autmaticly 
    if  transactionManager.loggerService == nil{
        transactionManager.loggerService = &Logger{ServiceName:TransactionManagerServiceName}
         err := transactionManager.loggerService.construct() 
         if  err != nil {
            return  err
         }
    }

    if transactionManager.walletService == nil {
        //error  message 
        const  errMsg =  "Provider  for  walletService  should  be  given"
        errmsg := transactionManager.loggerService.SprintErrorf(errMsg)
        return  errors.New(errmsg)
    }
    transactionManager.loggerService.Log("Service  created")
    return  nil
}

// Transfer  Money  
func  (transactionManager   TransactionManager ) TransferMoney( to ,  validator *  rsa.PublicKey   ,  amount float64) ( []  TransactionService ,  error )   {
    err :=  transactionManager.unValidService()
    if err != nil {
        return nil , err
    }
    type  transactionErrorPair struct {
        Service   TransactionService  
        Err  error 
    } 
    channel :=  make  (chan  transactionErrorPair ,  2 )
    createTransaction :=  func (money float64 ,  Receiver  *  rsa.PublicKey  ,  reason string ,  channel  chan transactionErrorPair )  {
        // Create  Transaction
        receiver := entitys.Client{Address:*Receiver}
        pair := entitys.BillingInfo{To :receiver  ,  }
        standard := TransactionsStandard{serviceName:transactionServiceName  ,balanceService:transactionManager.balanceService , walletService :transactionManager.walletService  } 
        transactionInfo := entitys.TransactionCoins{BillDetails:entitys.TransactionDetails{Bill:pair} , Amount:money , Reason:reason}
        transaction  := TransactionCoins{Transaction:  entitys.TransactionCoinEntityRoot{Transaction:transactionInfo}  , services:standard  ,  }
        var rsp transactionErrorPair
        err := transaction.construct() 
        if  err  != nil {
            // Failed  to  create  a transaction   service           
            rsp=  transactionErrorPair {Service:nil , Err:err ,}
            channel <- rsp  
            return
        } 
        err  = transaction.CreateTransaction() 
        if  err  != nil{
            // Failed  to  create  a transaction    
            rsp =  transactionErrorPair {Service:nil , Err:err ,}
            channel <- rsp  
            return 
        } 
   
        //Sign Transaction 
        err = transactionManager.walletService.sign(&transaction )
        
            rsp =  transactionErrorPair {Service:&transaction , Err:err ,}
            channel <-  rsp
    }

    //calculate  money 
    const taxConst float64 =  3.0
    howMuch ,  tax :=  Tax(amount ,  taxConst , transactionManager.loggerService )
    //? IF  i  make  the  lock  of  the  money(Full instead  of  partial ) here in  block  state  will  impove  the  next  routines  
    // create  transactions Parallel 
    go  createTransaction(howMuch  ,  to ,  "Transfer" , channel ) 
    go  createTransaction(tax  , validator, "fee"  , channel) 
    transactions := make([] TransactionService ,  0 )   

    var errTransaction error 
    for  i:= 0 ; i< 2; i++ {
        tpair := <-channel 
    
        if  tpair.Err == nil {
            fmt.Println(tpair)
            transactions = append(transactions ,  tpair.Service)
            
        }else {
            errTransaction = tpair.Err 
        }
    }
    if  len(transactions) != 2 {
        for _ , t := range transactions {
            //? Fall Back one  or more  transactions  failed  for those who succeed cancel the Frozen money
            depend := t.getAmount()
         
            transactionManager.walletService.UnFreeze(depend)
        }
     }
     if len(transactions) == 2  {
        list := make( []  entitys.TransactionCoinEntityRoot  ,0  )
        for  _ , tr:=  range  transactions {
            entity  , ok := tr.getInterface().(entitys.TransactionCoinEntityRoot)
            if  ok {
                list = append(list , entity)
            }
        }
 
        BlockChainCoinsService.InsertTransaction(list)

     }


    return  transactions , errTransaction 

}    
 

func (transactionManager   TransactionManager) unValidService() error {
    if  transactionManager.loggerService == nil {
        return  errors.New("Service  has  no  logger  service")
    } 
    if  transactionManager.walletService == nil  {
        return  errors.New("Service  has  no  wallet  service")
    }
    return  nil 

} 
func (transactionManager    TransactionManager) SendMessage( to , validatorMsg ,validatorCoins *  rsa.PublicKey ,  msg string   ) ([]  TransactionService  ,  error  ) {
    transactions := make([] TransactionService ,0 )   
    
    type  transactionErrorPair struct {
        Service   TransactionService  
        Err  error 
    } 
    
    msgExec := func( msg  string , tpair chan  transactionErrorPair  ) {
        //create msg transactio

     receiver := entitys.Client{Address: * to}
        pair := entitys.BillingInfo{To :receiver  ,  }
        standard := TransactionsStandard{serviceName:transactionServiceName  ,balanceService:transactionManager.balanceService , walletService :transactionManager.walletService  } 
        transactionInfo := entitys.TransactionMsg{BillDetails:entitys.TransactionDetails{Bill:pair} , Msg:msg}

        t := &TransactionMsg{services: standard , Transaction:entitys.TransactionMsgEntityRoot{Transaction:transactionInfo}} 
        err  := t.construct() 
        if err != nil {
            tpair <- transactionErrorPair{Service:nil , Err:err , }
            return 
        }
        err  = t.CreateTransaction() 
        if err != nil {
            tpair <- transactionErrorPair{Service:nil , Err:err , }
            return 
        }
        err =transactionManager.walletService.sign(t)
        tpair <- transactionErrorPair{Service:t , Err:err , }

    } 
    type  transactionsErrorPairList struct {
        Service  []  TransactionService
        Err  error 
    }
    enforcerParalell := func (pay float64  ,  tpair chan  transactionsErrorPairList ) {
        list , err := transactionManager.TransferMoney(validatorMsg ,  validatorCoins  ,   pay  )
        if err == nil{
            tpair <- transactionsErrorPairList{Service:list ,Err:nil ,}
        } else {
            tpair <- transactionsErrorPairList{Service:nil ,Err:err ,}
        }
    }

    pay :=  float64(len(msg)) 
    singleTransaction  := make (chan transactionErrorPair  , 1 )
    multipleTransaction := make (chan transactionsErrorPairList  , 1 )
    go  enforcerParalell(pay ,multipleTransaction)
    go  msgExec(msg  ,singleTransaction )
    var errInTransactions error 
    for  i := 0 ; i<2 ;i++{
        select {
            case single := <-singleTransaction : 
            if single.Err != nil{
                   errInTransactions = single.Err 
            }else {
                transactions = append(transactions, single.Service )}
            case multiple :=  <-  multipleTransaction:
            if multiple.Err != nil{
                   errInTransactions = multiple.Err 
            }else {
                transactions = append(transactions,   multiple.Service ...  )}
        }
    } 
    if errInTransactions != nil {
        return  nil , errInTransactions
    }
    return  transactions ,  nil  
}
// Helpfull  Functionality 
/**  
   Tax -  calculates  tax and  pay 
   @Param amount float64 
   @Param tax float64  (%)
   @logger  LogerService 
*/
func Tax(amount  ,  tax   float64 ,  logger LogerService) ( float64 ,  float64) {
    logger.Log(fmt.Sprintf("Start  Calculate  Tax for %.3f at  %.3f   %s" ,  amount ,  tax , "%"  ))
    tax  =   tax / 100 // regulate(%) => 1  ex.  3.5% =>  0.035 
    pay :=  amount * tax 
    amount = amount - pay 
    logger.Log(fmt.Sprintf("Commit  Calculate  Tax  %.3f coins value     , %.3f coins tax   " ,  amount ,  pay  ))
    return  amount,  pay 

}
