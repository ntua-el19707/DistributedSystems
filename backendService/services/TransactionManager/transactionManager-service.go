package TransactionManager

import (
	"FindBalance"
	"Logger"
	"Service"
	"WalletAndTransactions"
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
	"log"
)

const TransactionManagerServiceName = "TransactionManagerService"

// TransactionMangerService  interface  Rules
type TransactionManagerService interface {
	Service.Service // it  is  indeed a  service
	TransferMoney(to rsa.PublicKey, ammount float64) (entitys.TransactionCoinSet, error)
	SendMessage(to rsa.PublicKey, msg string) (entitys.TransactionMessageSet, error)
	unValidService() error
}

// TransactionManger  struct

type TransactionManager struct {
	WalletServiceInstance      *WalletAndTransactions.WalletStructV1Implementation
	LoggerServiceInstance      Logger.LoggerService
	FindBalanceServiceInstance FindBalance.BalanceService
}

const abbortTemplate = "Abbort due: %s"
const errNoWalletProvider = "The is  now wallet-service  instance"
const errNoFindBalanceProvider = "The is  now findBalance-service  instance"

// Construct
func (transactionManager *TransactionManager) Construct() error {
	// if  logger service  nil  =>  create  autmaticly
	if transactionManager.LoggerServiceInstance == nil {
		transactionManager.LoggerServiceInstance = &Logger.Logger{ServiceName: TransactionManagerServiceName}
		err := transactionManager.LoggerServiceInstance.Construct()
		if err != nil {
			return err
		}
	}

	if transactionManager.WalletServiceInstance == nil {
		//error  message
		errMsg := errNoWalletProvider
		transactionManager.LoggerServiceInstance.Error(errMsg)
		return errors.New(errMsg)
	}
	if transactionManager.FindBalanceServiceInstance == nil {
		//error  message
		errMsg := errNoFindBalanceProvider
		transactionManager.LoggerServiceInstance.Error(errMsg)
		return errors.New(errMsg)
	}
	transactionManager.LoggerServiceInstance.Log("Service  created")
	return nil
}

// Transfer  Money
func (transactionManager TransactionManager) TransferMoney(to rsa.PublicKey, amount float64) (entitys.TransactionCoinSet, error) {
	log.Println("hello")
	err := transactionManager.unValidService()
	var zeroSet entitys.TransactionCoinSet
	if err != nil {
		return zeroSet, err
	}
	type transactionErrorPair struct {
		Service WalletAndTransactions.TransactionService
		Err     error
	}
	channel := make(chan transactionErrorPair, 2)
	createTransaction := func(money float64, Receiver rsa.PublicKey, reason string, channel chan transactionErrorPair) {
		// Create  Transaction
		receiver := entitys.Client{Address: Receiver}
		pair := entitys.BillingInfo{To: receiver}
		standard := WalletAndTransactions.TransactionsStandard{ServiceName: "transaction-coin-service", BalanceServiceInstance: transactionManager.FindBalanceServiceInstance, WalletService: transactionManager.WalletServiceInstance}
		transactionInfo := entitys.TransactionCoins{BillDetails: entitys.TransactionDetails{Bill: pair}, Amount: money, Reason: reason}
		transaction := WalletAndTransactions.TransactionCoins{Transaction: entitys.TransactionCoinEntityRoot{Transaction: transactionInfo}, Services: standard}
		var rsp transactionErrorPair
		err := transaction.Construct()
		if err != nil {
			// Failed  to  create  a transaction   service
			rsp = transactionErrorPair{Service: nil, Err: err}
			channel <- rsp
			return
		}
		err = transaction.CreateTransaction()
		if err != nil {
			// Failed  to  create  a transaction
			rsp = transactionErrorPair{Service: nil, Err: err}
			channel <- rsp
			return
		}

		//Sign Transation
		err = transactionManager.WalletServiceInstance.Sign(&transaction)
		if err != nil {
			//* Note  that unfreeze  throw err  wheb coins< 0 and frozen - coins <0 t
			//s due to sucesffuly frozen  the  coins  in transationService  if  this fail
			transactionManager.WalletServiceInstance.UnFreeze(amount)
		}
		rsp = transactionErrorPair{Service: &transaction, Err: err}
		channel <- rsp
	}

	//calculate  money
	const taxConst float64 = 3.0
	howMuch, tax := Tax(amount, taxConst, transactionManager.LoggerServiceInstance)
	//? IF  i  make  the  lock  of  the  money(Full instead  of  partial ) here in  block  state  will  impove  the  next  routines
	// create  transactions Parallel
	const transfer = "Transfer"
	const fee = "fee"
	go createTransaction(howMuch, to, transfer, channel)
	go createTransaction(tax, rsa.PublicKey{}, fee, channel)
	transactions := make([]WalletAndTransactions.TransactionService, 0)
	var transactionSet entitys.TransactionCoinSet
	var errTransaction error
	for i := 0; i < 2; i++ {
		tpair := <-channel

		if tpair.Err == nil {

			entity, ok := tpair.Service.GetInterface().(entitys.TransactionCoinEntityRoot)
			if !ok {
				errTransaction = errors.New("Could Not Cast transaction to 'entitys.TransactionCoinEntityRoot'")
			} else {
				switch entity.Transaction.Reason {
				case transfer:
					transactionSet.Transfer = entity
					transactions = append(transactions, tpair.Service)
				case fee:
					transactionSet.Tax = entity
					transactions = append(transactions, tpair.Service)
				default:
					errTransaction = errors.New("UnKnown  Reason fo transactionCoin ")
				}
			}
		} else {
			errTransaction = tpair.Err
		}
	}
	if len(transactions) != 2 {
		for _, t := range transactions {
			//? Fall Back one  or more  transactions  failed  for those who succeed cancel the Frozen money
			depend := t.GetAmount()

			transactionManager.WalletServiceInstance.UnFreeze(depend)
		}
	}

	return transactionSet, errTransaction

}

func (transactionManager TransactionManager) unValidService() error {
	if transactionManager.LoggerServiceInstance == nil {
		return errors.New("Service  has  no  logger  service")
	}
	if transactionManager.WalletServiceInstance == nil {
		return errors.New("Service  has  no  wallet  service")
	}
	return nil

}
func (transactionManager TransactionManager) SendMessage(to rsa.PublicKey, msg string) (entitys.TransactionMessageSet, error) {
	var transactions entitys.TransactionMessageSet

	type transactionErrorPair struct {
		Service WalletAndTransactions.TransactionService
		Err     error
	}

	msgExec := func(msg string, tpair chan transactionErrorPair) {
		//create msg transactio

		receiver := entitys.Client{Address: to}
		pair := entitys.BillingInfo{To: receiver}
		standard := WalletAndTransactions.TransactionsStandard{ServiceName: "transaction-message-service", BalanceServiceInstance: transactionManager.FindBalanceServiceInstance, WalletService: transactionManager.WalletServiceInstance}
		transactionInfo := entitys.TransactionMsg{BillDetails: entitys.TransactionDetails{Bill: pair}, Msg: msg}

		t := &WalletAndTransactions.TransactionMsg{Services: standard, Transaction: entitys.TransactionMsgEntityRoot{Transaction: transactionInfo}}
		err := t.Construct()
		if err != nil {
			tpair <- transactionErrorPair{Service: nil, Err: err}
			return
		}
		err = t.CreateTransaction()
		if err != nil {
			tpair <- transactionErrorPair{Service: nil, Err: err}
			return
		}
		err = transactionManager.WalletServiceInstance.Sign(t)
		tpair <- transactionErrorPair{Service: t, Err: err}

	}
	type transactionsErrorPairForCoin struct {
		TransactionsCoin entitys.TransactionCoinSet
		Err              error
	}
	enforcerParalell := func(pay float64, tpair chan transactionsErrorPairForCoin) {
		//Create  Ope Tranasactions For  the Block-Chain-service  to stamp the  'TO'
		transactionSet, err := transactionManager.TransferMoney(rsa.PublicKey{}, pay)
		tpair <- transactionsErrorPairForCoin{TransactionsCoin: transactionSet, Err: err}
	}

	pay := float64(len(msg))
	singleTransaction := make(chan transactionErrorPair, 1)
	multipleTransaction := make(chan transactionsErrorPairForCoin, 1)
	go enforcerParalell(pay, multipleTransaction)
	go msgExec(msg, singleTransaction)
	var errInTransactions error
	for i := 0; i < 2; i++ {
		select {
		case single := <-singleTransaction:
			if single.Err != nil {
				errInTransactions = single.Err
			} else {
				//cast transaction to message
				entity, ok := single.Service.GetInterface().(entitys.TransactionMsgEntityRoot)
				if !ok {
					errInTransactions = errors.New("Could Not Cast transaction to 'entitys.TransactionMsgEntityRoot'")
				} else {
					transactions.TransactionMessage = entity
				}
			}
		case multiple := <-multipleTransaction:
			if multiple.Err != nil {
				errInTransactions = multiple.Err
			} else {
				transactions.TransactionCoin = multiple.TransactionsCoin
			}
		}
	}
	if errInTransactions != nil {
		return transactions, errInTransactions
	}
	return transactions, nil
}

// Helpfull  Functionality
/**
  Tax -  calculates  tax and  pay
  @Param amount float64
  @Param tax float64  (%)
  @logger  Logger.LoggerService
*/
func Tax(amount, tax float64, logger Logger.LoggerService) (float64, float64) {
	logger.Log(fmt.Sprintf("Start  Calculate  Tax for %.3f at  %.3f   %s", amount, tax, "%"))
	tax = tax / 100 // regulate(%) => 1  ex.  3.5% =>  0.035
	pay := amount * tax
	amount = amount - pay
	logger.Log(fmt.Sprintf("Commit  Calculate  Tax  %.3f coins value     , %.3f coins tax   ", amount, pay))
	return amount, pay

}
