package WalletAndTransactions

import (
	"Generator"
	"Logger"
	"Service"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"entitys"
	"errors"
	"fmt"
	"time"
)

type TransactionService interface {
	Service.Service
	semiConstruct() error
	CreateTransaction() error
	VerifySignature() error
	setSign(signiture []byte)
	GetTransaction() ([]byte, error)
	GetSigniture() []byte
	GetAmount() float64
	GetInterface() interface{}
}
type TransactionsStandard struct {
	ServiceName            string
	LoggerService          Logger.LoggerService
	BalanceServiceInstance BalanceService
	GeneratorService       Generator.GeneratorService
	WalletService          WalletService
	jsonStringfy           func(entitys.TransactionRecord) ([]byte, error)
	verifyMethod           func(*rsa.PublicKey, crypto.Hash, []byte, []byte) error
}
type TransactionCoins struct {
	Transaction entitys.TransactionCoinEntityRoot
	Services    TransactionsStandard
}

// TransactionStandard
func (t *TransactionsStandard) Construct() (string, error) {
	const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var err error
	//create  genaratorService  if  not exist
	if t.GeneratorService == nil {
		t.GeneratorService = &Generator.GeneratorImplementation{ServiceName: "generator-service", CharSet: allChars}
		err := t.GeneratorService.Construct()
		if err != nil {
			return "", err
		}
	}

	//ISSUE  id  for  Transaction
	id := t.GeneratorService.GenerateId(10)

	if t.LoggerService == nil {

		t.ServiceName = fmt.Sprintf("%s_%s", t.ServiceName, id)
		t.LoggerService = &Logger.Logger{ServiceName: t.ServiceName}
		err = t.LoggerService.Construct()
		if err != nil {
			return "", err
		}
	}
	logger := t.LoggerService
	logger.Log("Start checking  Services ")
	err = t.valid()
	if err != nil {
		errmsg := fmt.Sprintf("Abbort   Services not  valid %s", err.Error())
		logger.Error(errmsg)
		return "", errors.New(logger.SprintErrorf(errmsg))
	}
	t.jsonStringfy = entitys.JsonStringfy
	t.verifyMethod = rsa.VerifyPKCS1v15
	logger.Log("Commit Checking  Services")
	return id, nil
}
func (t *TransactionsStandard) valid() error {
	if t.LoggerService == nil {
		const errmsg string = "transaction standard has  no  LoggerService "
		return errors.New(errmsg)
	}
	if t.BalanceServiceInstance == nil {
		const errmsg string = "transaction standard has  no  BalanceServiceInstance "
		return errors.New(errmsg)
	}
	if t.GeneratorService == nil {
		const errmsg string = "transaction standard has  no  GeneratorService "
		return errors.New(errmsg)
	}
	if t.WalletService == nil {
		const errmsg string = "transaction standard has  no  WalletService "
		return errors.New(errmsg)
	}

	return nil
}

/*
*

	CreateTransaction -  crreate  a  transaction
	@Returns  error
*/
func (transaction *TransactionCoins) Construct() error {

	id, err := transaction.Services.Construct()
	if err != nil {
		return err
	}
	transaction.Transaction.Transaction.BillDetails.Transaction_id = id
	createdMessage := fmt.Sprintf("Created  service : %s \n", transaction.Services.
		ServiceName)
	LoggerService := transaction.Services.LoggerService
	LoggerService.Log(createdMessage)
	return nil
}
func (t *TransactionsStandard) semiConstruct(id string) error {
	var err error
	if t.LoggerService == nil {

		t.ServiceName = fmt.Sprintf("%s_%s", t.ServiceName, id)
		t.LoggerService = &Logger.Logger{ServiceName: t.ServiceName}
		err = t.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	/*	logger := t.LoggerService
		logger.Log("Start checking  Services ")
		err = t.valid()
		if err != nil {
			errmsg := fmt.Sprintf("Abbort   Services not  valid %s", err.Error())
			logger.Error(errmsg)
			return errors.New(logger.SprintErrorf(errmsg))
		}*/
	t.jsonStringfy = entitys.JsonStringfy
	t.verifyMethod = rsa.VerifyPKCS1v15
	return nil
}

// only if  use  verifyTransaction
func (transaction *TransactionCoins) semiConstruct() error {
	transaction.Services.ServiceName = "transactions-coins"
	return transaction.Services.semiConstruct(transaction.Transaction.Transaction.BillDetails.Transaction_id)
}

const ErrRequestFaildDueTotalMoneyFroze string = "Request To  sent %.3f from makes  (total  -  frozen ) %.3f\n"

/*
*

	CreateTransaction -  create  a  transaction
	@Returns  error
*/
func (t *TransactionCoins) CreateTransaction() error {
	transaction := &t.Transaction.Transaction
	Services := &t.Services
	logger := t.Services.LoggerService
	err := Services.valid()
	if err != nil {
		return err
	}

	transaction.BillDetails.Bill.From.Address = Services.WalletService.GetPub()
	bill := transaction.BillDetails.Bill
	transactionDetails := fmt.Sprintf("%s From %v to %v  For  %f", transaction.BillDetails.Transaction_id, bill.From.Address, bill.To.Address, transaction.Amount)

	logger.Log(fmt.Sprintf("Start creating  Transaction %s", transactionDetails))
	//issue  time to transaction
	transaction.BillDetails.Created_at = time.Now().Unix()

	amount := transaction.Amount
	//check if acount has the amount

	//? i need  a service  to find  the balance  of the sender

	total, err := Services.BalanceServiceInstance.findAndLock(amount) // ensure  the  balance  will not  change
	//wallet will have  balance -frozenMoney coins  that
	//trnsction  is  not yet  added  in chain
	//Be  optimist and  froze  the  coins
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	//find  balance
	if total < 0 {
		message := fmt.Sprintf(ErrRequestFaildDueTotalMoneyFroze, amount, total)
		// not  valid  trancation the total
		Services.WalletService.UnFreeze(amount)
		logger.Error(message)
		return errors.New(message)

	}
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s", transactionDetails))

	return nil
}

/*
setSign  -  set the signiture
@Param  siginiture  string
*/
func (t *TransactionCoins) setSign(signature []byte) {
	t.Transaction.Signiture = signature
}

/*
*

	GetAmount  -  get amount
	@Param  siginiture  string
*/
func (t TransactionCoins) GetAmount() float64 {
	return t.Transaction.Transaction.Amount
}

/*
*

	GetTransaction  -  get the transaction
	@Param  siginiture  string
*/
func (t TransactionCoins) GetTransaction() ([]byte, error) {
	//*The  signature  must not  be in document IMPORTANT
	return t.Services.jsonStringfy(t.Transaction.Transaction)
}

/*
*

	VerifySignature  -  verify the transaction
	@Param  siginiture  string
*/
func (t TransactionCoins) VerifySignature() error {
	logger := t.Services.LoggerService
	logger.Log(fmt.Sprintf("Start  Verifying  Transaction_id :%s ", t.Transaction.Transaction.BillDetails.Transaction_id))
	//define  verify
	verify := func(transaction []byte) error {
		PublicKey := t.Transaction.Transaction.BillDetails.Bill.From.Address
		hashed := sha256.Sum256(transaction)
		logger.Log(fmt.Sprintf("%v", PublicKey))

		//err := rsa.VerifyPKCS1v15(&PublicKey, crypto.SHA256, hashed[:], t.GetSigniture())
		err := t.Services.verifyMethod(&PublicKey, crypto.SHA256, hashed[:], t.GetSigniture())
		if err != nil {
			return err
		}
		return nil
	}
	data, err := t.GetTransaction()
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort  Verifying  Transaction_id :%s ", t.Transaction.Transaction.BillDetails.Transaction_id))
		return err
	}
	err = verify(data)
	if err != nil {
		logger.Log(fmt.Sprintf("Commit  Verifying  Transaction_id :%s (invalid)", t.Transaction.Transaction.BillDetails.Transaction_id))
		return err
	}
	logger.Log(fmt.Sprintf("Commit  Verifying  Transaction_id :%s (valid)", t.Transaction.Transaction.BillDetails.Transaction_id))
	return nil
}
func (t TransactionCoins) GetSigniture() []byte {
	return t.Transaction.Signiture
}
func (t TransactionCoins) GetInterface() interface{} {
	return t.Transaction
}

type TransactionMsg struct {
	Transaction entitys.TransactionMsgEntityRoot
	Services    TransactionsStandard
}

/*
*

	Construct - create  transactionMsg  service
	@Returns  error
*/
func (transaction *TransactionMsg) Construct() error {
	id, err := transaction.Services.Construct()
	if err != nil {
		return err
	}
	transaction.Transaction.Transaction.BillDetails.Transaction_id = id
	createdMessage := fmt.Sprintf("Created  service : %s \n", transaction.Services.
		ServiceName)
	LoggerService := transaction.Services.LoggerService
	LoggerService.Log(createdMessage)
	return nil
}
func (transaction *TransactionMsg) semiConstruct() error {
	transaction.Services.ServiceName = "transactions-msg"
	return transaction.Services.semiConstruct(transaction.Transaction.Transaction.BillDetails.Transaction_id)
}

/*
*

	CreateTransaction -  create		 a  transaction
	@Returns  error
*/
func (t *TransactionMsg) CreateTransaction() error {
	transaction := &t.Transaction.Transaction
	Services := &t.Services
	logger := t.Services.LoggerService
	err := Services.valid()
	if err != nil {
		return err
	}

	transaction.BillDetails.Bill.From.Address = Services.WalletService.GetPub()
	bill := transaction.BillDetails.Bill
	transactionDetails := fmt.Sprintf("%s From %v to %v  For  %s", transaction.BillDetails.Transaction_id, bill.From.Address, bill.To.Address, transaction.Msg)

	logger.Log(fmt.Sprintf("Start creating  Transaction %s", transactionDetails))
	//issue  time to transaction
	transaction.BillDetails.Created_at = time.Now().Unix()
	logger.Log(fmt.Sprintf("Commit creating  Transaction %s", transactionDetails))
	return nil
}

/*
*

	setSign  -  set the signiture
	@Param  siginiture  string
*/
func (t *TransactionMsg) setSign(signature []byte) {
	t.Transaction.Signiture = signature
}

/*
*

	GetTransaction  -  get the transaction
	@Param  siginiture  string
*/
func (t TransactionMsg) GetTransaction() ([]byte, error) {
	//*The  signature  must not  be in document IMPORTANT
	return t.Services.jsonStringfy(t.Transaction.Transaction)
}

/*
*

	GetAmount  -  get amount
*/
func (t TransactionMsg) GetAmount() float64 {
	return float64(0)
}

/*
*

	VerifySignature  -  verify the transaction
	@Param  siginiture  string
*/
func (t TransactionMsg) VerifySignature() error {
	//define  verify
	verify := func(transaction []byte) error {
		PublicKey := t.Transaction.Transaction.BillDetails.Bill.From.Address
		hashed := sha256.Sum256(transaction)

		//err := rsa.VerifyPKCS1v15(&PublicKey, crypto.SHA256, hashed[:], t.GetSigniture())
		err := t.Services.verifyMethod(&PublicKey, crypto.SHA256, hashed[:], t.GetSigniture())
		if err != nil {
			return err
		}
		return nil
	}
	data, err := t.GetTransaction()
	if err != nil {
		return err
	}
	return verify(data)
}
func (t TransactionMsg) GetSigniture() []byte {
	return t.Transaction.Signiture
}
func (t TransactionMsg) GetInterface() interface{} {
	return t.Transaction
}

//Mock service

type MockTransactionService struct {
	ConstructErr          error
	CreateTransactionErr  error
	VerifyErr             error
	SpyCallParamSign      []byte
	Transaction           []byte
	ErrorGet              error
	Signiture             []byte
	Amount                float64
	Interface             interface{}
	CallConstruct         int
	CallCreateTransaction int
	CallVerifySigniture   int
	CallSetSign           int
	CallGetTransaction    int
	CallGetAmount         int
	CallGetSigniture      int
	CallGetInterface      int
}

func (service *MockTransactionService) Construct() error {
	service.CallConstruct++
	return service.ConstructErr
}
func (service *MockTransactionService) semiConstruct() error {
	return nil
}

func (service *MockTransactionService) CreateTransaction() error {
	service.CallCreateTransaction++
	return service.CreateTransactionErr
}
func (service *MockTransactionService) VerifySignature() error {
	service.CallVerifySigniture++
	return service.VerifyErr
}
func (service *MockTransactionService) setSign(sign []byte) {
	service.CallSetSign++
	service.SpyCallParamSign = sign
}
func (service *MockTransactionService) GetTransaction() ([]byte, error) {
	service.CallGetTransaction++
	return service.Transaction, service.ErrorGet
}
func (service *MockTransactionService) GetSigniture() []byte {
	service.CallGetSigniture++
	return service.Signiture
}
func (service *MockTransactionService) GetAmount() float64 {
	service.CallGetAmount++
	return service.Amount
}
func (service *MockTransactionService) GetInterface() interface{} {
	service.CallGetInterface++
	return service.Interface
}
