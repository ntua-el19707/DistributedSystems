package WalletAndTransactions

import (
	"Logger"
	"Service"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"sync"
)

const walletSize = 2048

// Wallet Success Messages
const createdWalletMessage = "Has  succesfully  created  a new  wallet with  address %s"

// Wallet Errors Template
const failedToGenerateKeys = "Failed  to generate Keys  due to : %s "

type WalletService interface {
	Service.Service
	// will  generate  a wallet
	GenerateWallet(size int, method func(io.Reader, int) (*rsa.PrivateKey, error)) error
	Sign(transaction TransactionService) error
	GetPub() rsa.PublicKey
	Freeze(coins float64) error
	UnFreeze(coins float64) error
	GetFreeze() float64
}

// -- WalletStructV1Implementation --
/**
  WalletStructV1Implementation - struct
  @Param PublicKey   rsa.PublicKey
  @Param privateKey  rsa.PrivateKey
*/
type WalletStructV1Implementation struct {
	PublicKey          rsa.PublicKey
	privateKey         *rsa.PrivateKey
	frozen             float64
	mu                 sync.Mutex
	happenConstruction bool
	nonce              int
	signMethod         func(io.Reader, *rsa.PrivateKey, crypto.Hash, []byte) ([]byte, error)
	LoggerService      Logger.LoggerService
}

/*
*

	Construct -  WalletStructV1Implementation all  service Required a  contruction
	Has  To geneate a wallet
*/
func (wallet *WalletStructV1Implementation) Construct() error {
	//? why create  a private  construction  In @Interface Service  i said  that iwill not  pass arguments in construct
	//& i want to pass an argument  to the constructor that will be const 'The genrate method  come from rsa lib'
	//in order to force  fail constuctor from  testing
	return wallet.privateConstructor(rsa.GenerateKey)
}

// Create  a method  to help Construct Tyhe walletService
func (wallet *WalletStructV1Implementation) privateConstructor(generateMethod func(io.Reader, int) (*rsa.PrivateKey, error)) error {
	if wallet.LoggerService == nil {
		wallet.LoggerService = &Logger.Logger{ServiceName: "wallet-service"}
	}

	logger := wallet.LoggerService
	logger.Log("Start Construction of  a new wallet")
	err := wallet.GenerateWallet(walletSize, generateMethod)
	if err != nil {
		logger.Error("Abbort Construction of  a new wallet")
		return err
	}
	wallet.signMethod = rsa.SignPKCS1v15
	wallet.happenConstruction = true // TODO later and a valid  in wallet service
	logger.Log("Commit Construction of  a new wallet")
	return nil
	// ? Note  for  thought  if i save  the  keys in  phisicall storage
	// ? i  can change the Construct and  make  it so evry time  that  the  process  rice  has  the same  wallet
}

/*
*

	GenerateWallet - generate  a new wallet
	@Param size  int
*/
func (wallet *WalletStructV1Implementation) GenerateWallet(size int, GenerateMethod func(io.Reader, int) (*rsa.PrivateKey, error)) error {
	//generate  rsa key
	logger := wallet.LoggerService
	logger.Log("Start create  a new wallet")
	var err error
	wallet.privateKey, err = GenerateMethod(rand.Reader, size)

	if err != nil {
		const errorTemplate = "Abbort  error :%s"
		errorMessage := fmt.Sprintf(failedToGenerateKeys, err.Error())
		message := fmt.Sprintf(errorTemplate, errorMessage)
		logger.Error(message)
		return errors.New(message)
	}
	wallet.PublicKey = wallet.privateKey.PublicKey
	if err != nil {
		return err
	}
	logger.Log("Commit  wallet created")
	return nil
}

/*
*

	Sign - Sign  a Transaction
	@Param  trasactionService
	@Return error
*/
func (wallet *WalletStructV1Implementation) Sign(transactionService TransactionService) error {
	logger := wallet.LoggerService
	logger.Log("Start  Sign  transaction")
	//define  SignTransaction
	setNonce := func() {
		var nonceLocal int
		wallet.mu.Lock()
		nonceLocal = wallet.nonce
		wallet.nonce++
		wallet.mu.Unlock()
		transactionService.SetNonce(nonceLocal)
	}
	SignDocument := func(transaction []byte) ([]byte, error) {
		// hashed  transaction
		hashed := sha256.Sum256(transaction)
		// Sign  hashed transaction
		Signature, err := wallet.signMethod(rand.Reader, wallet.privateKey, crypto.SHA256, hashed[:])

		if err != nil {
			return nil, err
		}
		return Signature, nil
	}
	logger.Log("start seting  nonce ")

	setNonce()
	logger.Log("commit seting  nonce")
	logger.Log("Start  getTransaction")
	transaction, err := transactionService.GetTransaction()
	if err != nil {
		const errorTemplate = "Abbort  error :  Failed  to   getTransaction  due to %s"
		message := fmt.Sprintf(errorTemplate, err.Error())
		logger.Log(message)
		err = errors.New(logger.Sprintf(message))
		return err

	}
	logger.Log("Commit getTransaction")

	// Sign the document
	logger.Log("Start SignDocument")
	Signature, err := SignDocument(transaction)
	if err != nil {
		const errorTemplate = "Abbort  error :  Failed  to   SignTransaction  due to %s"
		message := fmt.Sprintf(errorTemplate, err.Error())
		logger.Log(message)
		err = errors.New(logger.Sprintf(message))
		return err
	}
	logger.Log("Commit SignDocument")
	logger.Log("Start setSign")
	transactionService.setSign(Signature)
	logger.Log("Commit setSign")
	logger.Log("Commit  Sign  transaction")
	return nil
}

/*
*

	GetPub - get publickey
	@Return  rsa.PublicKey
*/
func (wallet *WalletStructV1Implementation) GetPub() rsa.PublicKey {
	return wallet.PublicKey
}

/*
*

	encodeToPemPublicKey - parse  rsa to pem
	@Param key rsa.PublicKey
	@Param  who string

func encodeToPemPublicKey(key rsa.PublicKey, logger Logger.LoggerService) (string, error) {

	logger.Log("will attempt to parse  rsa.PublicKey to pem format")
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		message := fmt.Sprintf("Error : Marshal Failed Due to %s", err.Error())
		return "", errors.New(logger.Sprintf(message))
	}
	publicKeyPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPEM,
	}
	publicKeyPEMBytes := pem.EncodeToMemory(publicKeyPEMBlock)
	logger.Log("parsed the  rsa.PublicKey to pem format")
	return string(publicKeyPEMBytes), nil

}
*/
// -- const  errTemplate  fot  freeze unfreeze
const errCannotFreezeUnFreezeNegativeCoins = "Cannot '%s' negative coins"
const errCannotUnFreezeAndGoNegativeCoins = "Cannot unfreeze that amount  not  enough frozen coins"

func (wallet *WalletStructV1Implementation) Freeze(coins float64) error {
	var err error
	wallet.mu.Lock()
	wallet.LoggerService.Log(fmt.Sprintf("start freeze  money %.3f ", coins))
	if coins < 0 {
		err = errors.New(fmt.Sprintf(errCannotFreezeUnFreezeNegativeCoins, "freeze"))
	} else {
		wallet.frozen += coins
	}
	wallet.LoggerService.Log(fmt.Sprintf("commit freeze  money %.3f ", coins))
	wallet.mu.Unlock()
	return err
}
func (wallet *WalletStructV1Implementation) UnFreeze(coins float64) error {
	var err error
	wallet.mu.Lock()
	wallet.LoggerService.Log(fmt.Sprintf("start unfreeze  money %.3f ", coins))
	const epsilon = 1e-9
	if coins < 0 {
		err = errors.New(fmt.Sprintf(errCannotFreezeUnFreezeNegativeCoins, "unfreeze"))
	} else if wallet.frozen-coins < -epsilon {
		wallet.LoggerService.Error(fmt.Sprintf("%f", wallet.frozen-coins))
		err = errors.New(errCannotUnFreezeAndGoNegativeCoins)
	} else {
		wallet.frozen -= coins
	}
	wallet.LoggerService.Log(fmt.Sprintf("commit unfreeze  money %.3f ", coins))
	wallet.mu.Unlock()
	return err
}
func (wallet *WalletStructV1Implementation) GetFreeze() float64 {
	wallet.mu.Lock()
	coins := wallet.frozen
	wallet.mu.Unlock()
	return coins
}

/*
*

	encodeToPemPrivateKey - parse  rsa to pem
	@Param key rsa.PublicKey
	@Param  who string
*/ /**
func encodeToPemPrivateKey(key *rsa.PrivateKey, logger Logger.LoggerService) string {
	logger.Log("will attempt to parse  rsa.PrivateKey to pem format  ")
	bytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPEMBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytes,
	}
	privateKey := pem.EncodeToMemory(privateKeyPEMBlock)
	logger.Log("parsed the  rsa.PrivateKey to pem format  ")
	return string(privateKey)

}
*/
// Mock Wallet
type MockWallet struct {
	ErrorGenerateWallet    error
	ErrorSignWallet        error
	CounterGeneratorWallet int
	CounterSign            int
	CounterGetPub          int
	CounterFreeze          int
	CounterUnFreeze        int
	CounterGetFreeze       int
	Frozen                 float64
	MockPublicKey          rsa.PublicKey
}

func (mock *MockWallet) Construct() error {
	return nil
}
func (mock *MockWallet) GenerateWallet(size int, method func(io.Reader, int) (*rsa.PrivateKey, error)) error {
	mock.CounterGeneratorWallet++
	return mock.ErrorGenerateWallet
}
func (mock *MockWallet) Sign(transation TransactionService) error {
	mock.CounterSign++
	return mock.ErrorSignWallet
}
func (mock *MockWallet) GetPub() rsa.PublicKey {
	mock.CounterGetPub++
	return mock.MockPublicKey
}
func (mock *MockWallet) Freeze(coins float64) error {

	mock.CounterFreeze++
	return nil
}
func (mock *MockWallet) UnFreeze(coins float64) error {
	mock.CounterUnFreeze++
	return nil
}
func (mock *MockWallet) GetFreeze() float64 {
	mock.CounterGetFreeze++
	return mock.Frozen
}
