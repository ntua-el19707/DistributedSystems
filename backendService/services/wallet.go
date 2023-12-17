package services 

import (
    "crypto/rand"
	"crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "fmt"
    "log"
)

const  serviceName = "walletService"
const walletSize = 4096
//Wallet Success Messages 
const createdWalletMessage = "The ' %s ' service  has  succesfully  created  a new  wallet with  address \n%s"
// Wallet Errors Template
const failedToGenerateKeys = "The ' %s ' service  Failed  to generate Keys  due to : %s "

type WalletService  interface {
    	construct()error  
    // will  generate  a wallet 
    generate_wallet(size int ) error
}
// -- walletStructV1Service --
/**
     walletStructV1Service - struct 
     @Param PublicKey   rsa.PublicKey
     @Param privateKey  rsa.PrivateKey
*/
type  walletStructV1Service  struct {
    PublicKey   rsa.PublicKey
    privateKey  * rsa.PrivateKey
}

/**
    construct -  walletStructV1Service all  service Required a  contruction
    Has  To geneate a wallet 
 */
func  (wallet *  walletStructV1Service ) construct () error {
    err := wallet.generate_wallet(walletSize)
    return err
//? Note  for  thought  if i save  the  keys in  phisicall storage
//? i  can change the construct and  make  it so evry time  that  the  process  rice  has  the same  wallet
}
/**
    generate_wallet - generate  a new wallet 
    @Param size  int
*/
func (wallet *  walletStructV1Service) generate_wallet(size  int) error {
    //generate  rsa key
    var err  error 
    wallet.privateKey ,  err =  rsa.GenerateKey(rand.Reader ,  size) 
 
    if  err != nil {
        errorMessage := fmt.Sprintf(failedToGenerateKeys ,  serviceName ,  err.Error() )
        return  errors.New(errorMessage)
    }
    wallet.PublicKey = wallet.privateKey.PublicKey
   
    publicPemKey ,err :=  encodeToPemPublicKey(wallet.PublicKey ,serviceName)


    if  err != nil {
        return  err
    }
    log.Printf(createdWalletMessage ,  serviceName ,    publicPemKey )
    return  nil
}

/**
    encodeToPemPublicKey - parse  rsa to pem 
    @Param key rsa.PublicKey 
    @Param  who string
*/
func encodeToPemPublicKey(key rsa.PublicKey ,  who string)  ( string ,error) {
    
    log.Printf("['%s'] will attempt to parse  rsa.PublicKey to pem format  " , who )
    publicKeyPEM, err := x509.MarshalPKIXPublicKey(&key)
    if err != nil {
        message := fmt.Sprintf("['%s'] Error : Marshal Failed Due to %s" , who , err.Error())
        return "" , errors.New(message)
    }
    publicKeyPEMBlock := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyPEM,
    }
    publicKeyPEMBytes := pem.EncodeToMemory(publicKeyPEMBlock)
    log.Printf("['%s'] parsed the  rsa.PublicKey to pem format  " , who )
    return  string(publicKeyPEMBytes) ,nil 

}
/**
    encodeToPemPrivateKey - parse  rsa to pem 
    @Param key rsa.PublicKey 
    @Param  who string
*/
func encodeToPemPrivateKey(key * rsa.PrivateKey ,who string  ) ( string) {
    log.Printf("['%s'] will attempt to parse  rsa.PrivateKey to pem format  " , who )
    bytes := x509.MarshalPKCS1PrivateKey(key)
    privateKeyPEMBlock := &pem.Block{
        Type:  "PRIVATE KEY",
        Bytes: bytes ,
    }
    privateKey  := pem.EncodeToMemory(privateKeyPEMBlock)
     log.Printf("['%s'] parsed the  rsa.PrivateKey to pem format  " , who )
    return  string(privateKey)


}
