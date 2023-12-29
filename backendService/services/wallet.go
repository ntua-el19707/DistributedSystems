package services 

import (
    "crypto"
    "crypto/rand"
	"crypto/rsa"
    "crypto/x509"
    "crypto/sha256"
    "encoding/pem"
    "errors"
    "fmt"
    "sync"

)


const walletSize = 2048
//Wallet Success Messages 
const createdWalletMessage = "Has  succesfully  created  a new  wallet with  address %s"
// Wallet Errors Template
const failedToGenerateKeys = "Failed  to generate Keys  due to : %s "

type WalletService  interface {
    construct()error  
    // will  generate  a wallet 
    generate_wallet(size int ) error
    sign(transaction TransactionService) error 
    getPub() *  rsa.PublicKey
    Freeze(coins  float64) 
    UnFreeze(coins float64)
    getFreeze() float64   
}
// -- walletStructV1Service --
/**
     walletStructV1Service - struct 
     @Param PublicKey   rsa.PublicKey
     @Param privateKey  rsa.PrivateKey
*/
type  walletStructV1Service  struct {
    PublicKey  *   rsa.PublicKey
    privateKey  * rsa.PrivateKey
    frozen  float64
    mu sync.Mutex
    loggerService  LogerService 
}

/**
    construct -  walletStructV1Service all  service Required a  contruction
    Has  To geneate a wallet 
 */
func  (wallet *  walletStructV1Service ) construct () error {

    if  wallet.loggerService == nil {
        wallet.loggerService = &Logger{ServiceName:walletServiceName}
    }
    
    logger := wallet.loggerService
    logger.Log("Start construction of  a new wallet")
    err := wallet.generate_wallet(walletSize)
    if err != nil {
         logger.Log("Abbort construction of  a new wallet")
         return err
    }
     logger.Log("Commit construction of  a new wallet")
    return nil
//? Note  for  thought  if i save  the  keys in  phisicall storage
//? i  can change the construct and  make  it so evry time  that  the  process  rice  has  the same  wallet
}
/**
    generate_wallet - generate  a new wallet 
    @Param size  int
*/
func (wallet *  walletStructV1Service) generate_wallet(size  int) error {
    //generate  rsa key
    logger := wallet.loggerService
    logger.Log("Start create  a new wallet")
    var err  error 
    wallet.privateKey ,  err =  rsa.GenerateKey(rand.Reader ,  size) 


    if  err != nil {
        errorMessage := fmt.Sprintf(failedToGenerateKeys ,  err.Error() )
        logErrorMessage := logger.Sprintf(errorMessage) 
        return  errors.New(logErrorMessage)
    }
    wallet.PublicKey = &wallet.privateKey.PublicKey
    if  err != nil {
        return  err
    }
   logger.Log("Commit  wallet created")
    return  nil
}
/**
    sign - sign  a Transaction
    @Param  trasactionService
    @Return error
*/
func (wallet *  walletStructV1Service)sign(transactionService TransactionService) error {
    logger := wallet.loggerService
    logger.Log("Start  sign  transaction")
    //define  signTransaction
    signDocument :=   func (transaction []  byte )([] byte , error){
        // hashed  transaction  
        hashed := sha256.Sum256(transaction)
        // sign  hashed transaction
	    signature, err := rsa.SignPKCS1v15(rand.Reader, wallet.privateKey, crypto.SHA256, hashed[:])
	    
        if err != nil {
		    return nil, err
	    }
	    return signature, nil
    }
    logger.Log("Start  getTransaction")
    transaction  , err := transactionService.getTransaction()
    if  err != nil {
        const  errorTemplate =  "Abbort  error :  Failed  to   getTransaction  due to %s"
        message :=  fmt.Sprintf(errorTemplate ,  err.Error())
        logger.Log(message)
        err = errors.New(logger.Sprintf(message))
        return  err

    }
    logger.Log("Commit getTransaction")
	// Sign the document
    logger.Log("Start signDocument")
	signature, err := signDocument(transaction)
	if err != nil {
        const  errorTemplate =  "Abbort  error :  Failed  to   signTransaction  due to %s"
        message :=  fmt.Sprintf(errorTemplate ,  err.Error())
        logger.Log(message)
        err = errors.New(logger.Sprintf(message))
		return err
	}
    logger.Log("Commit signDocument")
    logger.Log("Start setSign")
    transactionService.setSign(signature)
    logger.Log("Commit setSign")
    logger.Log("Commit  sign  transaction")
    return  nil
}
/**
    getPub - get publickey  
    @Return  rsa.PublicKey
*/
func (wallet   walletStructV1Service) getPub()  * rsa.PublicKey {
       return  wallet.PublicKey
}

/**
    encodeToPemPublicKey - parse  rsa to pem 
    @Param key rsa.PublicKey 
    @Param  who string
*/
func encodeToPemPublicKey(key rsa.PublicKey ,  logger LogerService)  ( string ,error) {
    
    logger.Log("will attempt to parse  rsa.PublicKey to pem format")
    publicKeyPEM, err := x509.MarshalPKIXPublicKey(&key)
    if err != nil {
        message := fmt.Sprintf("Error : Marshal Failed Due to %s" , err.Error())
        return "" , errors.New(logger.Sprintf(message))
    }
    publicKeyPEMBlock := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyPEM,
    }
    publicKeyPEMBytes := pem.EncodeToMemory(publicKeyPEMBlock)
    logger.Log("parsed the  rsa.PublicKey to pem format")
    return  string(publicKeyPEMBytes) ,nil 

}
func  (wallet *  walletStructV1Service ) Freeze(coins  float64) {
    wallet.mu.Lock()
    wallet.frozen += coins 
    wallet.mu.Unlock() 
} 
func  (wallet *  walletStructV1Service ) UnFreeze(coins  float64) {
    wallet.loggerService.Log("start unfreeze  money")
    wallet.mu.Lock()
    wallet.frozen -=  coins 
    wallet.mu.Unlock() 
    wallet.loggerService.Log("commit unfreeze  money")
}
func  (wallet *  walletStructV1Service ) getFreeze() float64 {
    wallet.mu.Lock()
    coins := wallet.frozen
    wallet.mu.Unlock() 
    return  coins 
}

/**
    encodeToPemPrivateKey - parse  rsa to pem 
    @Param key rsa.PublicKey 
    @Param  who string
*/
func encodeToPemPrivateKey(key * rsa.PrivateKey ,logger LogerService ) ( string) {
    logger.Log("will attempt to parse  rsa.PrivateKey to pem format  ")
    bytes := x509.MarshalPKCS1PrivateKey(key)
    privateKeyPEMBlock := &pem.Block{
        Type:  "PRIVATE KEY",
        Bytes: bytes ,
    }
    privateKey  := pem.EncodeToMemory(privateKeyPEMBlock)
    logger.Log("parsed the  rsa.PrivateKey to pem format  ")
    return  string(privateKey)


}

//Mock Wallet 
type  mockWallet struct {
    errorGenerateWallet error 
    errorSignWallet error 
    counterGeneratorWallet int 
    counterSign int 
    counterGetPub int
    counterFreeze int 
    counterUnFreeze int
    countergetFreeze int 
    frozen float64

}
func  (mock *  mockWallet ) construct() error  {
    return nil 
}
func  (mock *  mockWallet ) generate_wallet( size  int ) error  {
    mock.counterGeneratorWallet++
    return mock.errorGenerateWallet 
}
func  (mock *  mockWallet ) sign( transation TransactionService) error  {
    mock.counterSign++
    return mock.errorSignWallet
}
func  (mock *  mockWallet ) getPub()  * rsa.PublicKey  {
    mock.counterGetPub++
    return &rsa.PublicKey{}
}
func  (mock *  mockWallet ) Freeze( coins float64)    {

    mock.counterFreeze++ 
}
func  (mock *  mockWallet ) UnFreeze( coins float64)    {
    mock.counterUnFreeze++ 
}
func  (mock *  mockWallet ) getFreeze()   float64 {
    mock.countergetFreeze++ 
    return mock.frozen 
}



