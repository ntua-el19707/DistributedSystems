package WalletAndTransactions

import (
	"Generator"
	"Logger"
	"crypto"
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
	"testing"
)

func buildTransactionMsgService() (TransactionService, *Logger.MockLogger, *Generator.MockGenerator, *mockFindBalance, *TransactionMsg, *MockWallet, error) {

	mockLogger := &Logger.MockLogger{}
	mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
	mockFindBalance := &mockFindBalance{}
	wallet := &MockWallet{}
	transactionStandard := TransactionsStandard{ServiceName: "transactions-service", WalletService: wallet, LoggerService: mockLogger, BalanceServiceInstance: mockFindBalance, GeneratorService: mockGenerator}
	transaction := &TransactionMsg{Services: transactionStandard}
	transactionService := transaction
	err := transactionService.Construct()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return transactionService, mockLogger, mockGenerator, mockFindBalance, transaction, wallet, nil
}
func buildTransactionService() (TransactionService, *Logger.MockLogger, *Generator.MockGenerator, *mockFindBalance, *TransactionCoins, *MockWallet, error) {

	mockLogger := &Logger.MockLogger{}
	mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
	mockFindBalance := &mockFindBalance{}
	wallet := &MockWallet{}
	transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: mockLogger, BalanceServiceInstance: mockFindBalance, GeneratorService: mockGenerator}
	transaction := &TransactionCoins{Services: transactionStandard}
	transactionService := transaction
	err := transactionService.Construct()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return transactionService, mockLogger, mockGenerator, mockFindBalance, transaction, wallet, nil

}
func TestCreateServiceTransactionMsg(t *testing.T) {
	_, logger, _, _, _, _, err := buildTransactionMsgService()
	if err != nil {
		t.Errorf("Expceted  no err  but  got %v", err)
	}
	if len(logger.Logs) != 3 {
		t.Errorf("Expceted  to log 3  msg   but  got %d", len(logger.Logs))
	}
	fmt.Println("created Transaction Msg Service")
}

func TestCreateServiceTransactionMsgCreateValidTransaction(t *testing.T) {
	tr, _, _, _, _, _, err := buildTransactionMsgService()
	err = tr.CreateTransaction()
	if err != nil {
		t.Errorf("Expceted  no err  but  got %v", err)
	}
	fmt.Println("create transaction Msg Service")
}
func TestCreateServiceTransaction(t *testing.T) {
	_, logger, generator, _, _, _, err := buildTransactionService()
	if err != nil {
		t.Errorf("Failed  to create  trnsaction service  due  to %s", err.Error())
	}
	if generator.TimesCallGenerateId != 1 {
		t.Errorf("Generate  id function  should  be  call  1 time  but get called %d", generator.TimesCallGenerateId)
	}
	if len(logger.Logs) != 3 {
		t.Errorf("Logger logs  should  have  3  message but  have  %d messages ", len(logger.Logs))
	}
	msg := fmt.Sprintf("Created  service : %s \n", "transaction-service")
	if logger.Logs[2] != msg {
		t.Errorf("message  should  be  %s  but  got  %s  ", msg, logger.Logs[2])
	}

	fmt.Println("Create Transaction Service")
}
func TestCreateTransactionInvalid(t *testing.T) {
	service, _, _, findBalance, transaction, wallet, _ := buildTransactionService()
	findBalance.amount = 10
	transaction.Transaction.Transaction.Amount = 15
	wallet.Frozen = 15
	errmsg := fmt.Sprintf("Request To  sent  %.3f  from  %.3f  balance Failed  due to total Money Froze(for wallet  ) %.3f\n", transaction.Transaction.Transaction.Amount, findBalance.amount, wallet.Frozen)
	err := service.CreateTransaction()
	if err == nil {
		t.Errorf("It should be  invalid")
	}
	if err.Error() != errmsg {
		t.Errorf("It should get err: %s  but  got %s", errmsg, err.Error())
	}
	if findBalance.locked {
		t.Errorf("It should call first  lock and  then unlock ")

	}
	if findBalance.lockedCall != 1 {
		t.Errorf("It should locked balance once  but lock %d ", findBalance.lockedCall)

	}
	if findBalance.unlockedCall != 1 {
		t.Errorf("It should unlocked balance once  but unlock %d ", findBalance.unlockedCall)

	}
	if wallet.CounterFreeze != 1 {
		t.Errorf("It should freeze money  once  but freeze %d ", wallet.CounterFreeze)

	}
	if wallet.CounterUnFreeze != 1 {
		t.Errorf("It should un freeze money  once  but un  freeze %d ", wallet.CounterUnFreeze)

	}
	if wallet.CounterGetFreeze != 1 {
		t.Errorf("It should get freeze money  once  but get  freeze %d ", wallet.CounterGetFreeze)

	}

	fmt.Println("Should  not  create  invalid transaction")

}
func TestCreateTransactionvalid(t *testing.T) {
	service, logger, _, findBalance, transaction, wallet, _ := buildTransactionService()
	findBalance.amount = 100
	transaction.Transaction.Transaction.Amount = 15
	wallet.Frozen = transaction.Transaction.Transaction.Amount
	logger.Logs = make([]string, 0)
	err := service.CreateTransaction()
	if err != nil {
		t.Errorf("It should be  valid not  get  err  %s", err.Error())
	}
	if len(logger.Logs) != 2 {
		t.Errorf("Logger  should  receive 2  messages  but intead  got  %d ", len(logger.Logs))
	}
	if findBalance.findBalanceCalledTimes != 1 {
		t.Errorf("findbalance should  be  called  once but intead  called  %d ", findBalance.findBalanceCalledTimes)
	}
	if findBalance.locked {
		t.Errorf("It should call first  lock and  then unlock ")

	}
	if findBalance.lockedCall != 1 {
		t.Errorf("It should locked balance once  but lock %d ", findBalance.lockedCall)

	}
	if findBalance.unlockedCall != 1 {
		t.Errorf("It should unlocked balance once  but unlock %d ", findBalance.unlockedCall)

	}
	if wallet.CounterFreeze != 1 {
		t.Errorf("It should freeze money  once  but freeze %d ", wallet.CounterFreeze)

	}
	if wallet.CounterUnFreeze != 0 {
		t.Errorf("It should not un freeze money  but did un   freeze %d ", wallet.CounterUnFreeze)

	}
	if wallet.CounterGetFreeze != 1 {
		t.Errorf("It should get freeze money  once  but get  freeze %d ", wallet.CounterGetFreeze)

	}
	fmt.Println("Should    create  valid transaction")

}

func TestTransactions(t *testing.T) {
	const prefix string = "----"
	fmt.Println("* Test  cases for  TransactionService")
	TestTransactionCoinsImpl := func(prefixOld string) {
		fmt.Printf("%s transactions-coins implementation\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestCreateService := func(prefixOld string) {

			fmt.Printf("%s Create Service \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestCreate := func(prefixOld string) {
				mockLogger := &Logger.MockLogger{}
				mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
				mockFindBalance := &mockFindBalance{}
				wallet := &MockWallet{}
				transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: mockLogger, BalanceServiceInstance: mockFindBalance, GeneratorService: mockGenerator}
				transaction := &TransactionCoins{Services: transactionStandard}
				transactionService := transaction
				err := transactionService.Construct()
				if err != nil {
					t.Errorf("expected to got  no err  but  got %v ", err)
				}
				fmt.Printf("%s it  should  create a transaction-service\n", prefixOld)
			}
			TestFail := func(prefixOld string) {
				fmt.Printf("%s fail  to create  service \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestFailDueToWallet := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
					mockFindBalance := &mockFindBalance{}
					transactionStandard := TransactionsStandard{ServiceName: "transaction-service", LoggerService: mockLogger, BalanceServiceInstance: mockFindBalance, GeneratorService: mockGenerator}
					transaction := &TransactionCoins{Services: transactionStandard}
					transactionService := transaction
					err := transactionService.Construct()
					errWallet := "transaction standard has  no  WalletService "
					errMsg := fmt.Sprintf("Abbort   Services not  valid %s", errWallet)

					if err.Error() != errMsg {
						t.Errorf("expected to got err %s but  got %v ", errMsg, err)
					}
					fmt.Printf("%s it  should fail to create a transaction-service due to 'no wallet-service'\n", prefixOld)
				}
				TestFailDueToBalance := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
					wallet := &MockWallet{}
					transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: mockLogger, GeneratorService: mockGenerator}
					transaction := &TransactionCoins{Services: transactionStandard}
					transactionService := transaction
					err := transactionService.Construct()
					errBalance := "transaction standard has  no  BalanceServiceInstance "
					errMsg := fmt.Sprintf("Abbort   Services not  valid %s", errBalance)

					if err.Error() != errMsg {
						t.Errorf("expected to got err %s but  got %v ", errMsg, err)
					}
					fmt.Printf("%s it  should fail to create a transaction-service due to 'no balance-service'\n", prefixOld)
				}
				TestFailDueToWallet(prefixNew)
				TestFailDueToBalance(prefixNew)
			}
			TestCreate(prefixNew)
			TestFail(prefixNew)
		}
		TestCreateTransaction := func(prefixOld string) {

			fmt.Printf("%s Create Transaction \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestShouldCreate := func(prefixOld string) {
				logger := &Logger.MockLogger{}
				mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
				findBalance := &mockFindBalance{}
				wallet := &MockWallet{}
				transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: logger, BalanceServiceInstance: findBalance, GeneratorService: mockGenerator}
				transaction := &TransactionCoins{Services: transactionStandard}
				transactionService := transaction
				err := transactionService.Construct()
				if err != nil {
					t.Errorf("expected to got  no err  but  got %v ", err)
				}
				findBalance.amount = 100
				transaction.Transaction.Transaction.Amount = 15
				wallet.Frozen = transaction.Transaction.Transaction.Amount
				logger.Logs = make([]string, 0)
				err = transactionService.CreateTransaction()
				if err != nil {
					t.Errorf("It should be  valid not  get  err  %s", err.Error())
				}
				if len(logger.Logs) != 2 {
					t.Errorf("Logger  should  receive 2  messages  but intead  got  %d ", len(logger.Logs))
				}
				if findBalance.findBalanceCalledTimes != 1 {
					t.Errorf("findbalance should  be  called  once but intead  called  %d ", findBalance.findBalanceCalledTimes)
				}
				if findBalance.locked {
					t.Errorf("It should call first  lock and  then unlock ")

				}
				if findBalance.lockedCall != 1 {
					t.Errorf("It should locked balance once  but lock %d ", findBalance.lockedCall)

				}
				if findBalance.unlockedCall != 1 {
					t.Errorf("It should unlocked balance once  but unlock %d ", findBalance.unlockedCall)

				}
				if wallet.CounterFreeze != 1 {
					t.Errorf("It should freeze money  once  but freeze %d ", wallet.CounterFreeze)

				}
				if wallet.CounterUnFreeze != 0 {
					t.Errorf("It should not un freeze money  but did un   freeze %d ", wallet.CounterUnFreeze)

				}
				if wallet.CounterGetFreeze != 1 {
					t.Errorf("It should get freeze money  once  but get  freeze %d ", wallet.CounterGetFreeze)

				}
				fmt.Printf("%s it  should create a  valid transaction\n", prefixOld)
			}
			TestShouldFail := func(prefixOld string) {
				fmt.Printf("%s Fail to create Transaction  \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestShouldFailNotEnough := func(prefixOld string) {
					logger := &Logger.MockLogger{}
					mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
					findBalance := &mockFindBalance{}
					wallet := &MockWallet{}
					transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: logger, BalanceServiceInstance: findBalance, GeneratorService: mockGenerator}
					transaction := &TransactionCoins{Services: transactionStandard}
					transactionService := transaction
					err := transactionService.Construct()
					if err != nil {
						t.Errorf("expected to got  no err  but  got %v ", err)
					}

					wallet.Frozen = transaction.Transaction.Transaction.Amount
					logger.Logs = make([]string, 0)
					findBalance.amount = 10
					transaction.Transaction.Transaction.Amount = 15
					wallet.Frozen = 15
					errmsg := fmt.Sprintf("Request To  sent  %.3f  from  %.3f  balance Failed  due to total Money Froze(for wallet  ) %.3f\n", transaction.Transaction.Transaction.Amount, findBalance.amount, wallet.Frozen)
					err = transactionService.CreateTransaction()
					if err.Error() != errmsg {
						t.Errorf("expected to get %s but got  %s", errmsg, err.Error())
					}
					if len(logger.Logs) != 1 {
						t.Errorf("Logger  should  receive 1  messages  but intead  got  %d ", len(logger.Logs))
					}
					if findBalance.findBalanceCalledTimes != 1 {
						t.Errorf("findbalance should  be  called  once but intead  called  %d ", findBalance.findBalanceCalledTimes)
					}
					if findBalance.locked {
						t.Errorf("It should call first  lock and  then unlock ")

					}
					if findBalance.lockedCall != 1 {
						t.Errorf("It should locked balance once  but lock %d ", findBalance.lockedCall)

					}
					if findBalance.unlockedCall != 1 {
						t.Errorf("It should unlocked balance once  but unlock %d ", findBalance.unlockedCall)

					}
					if wallet.CounterFreeze != 1 {
						t.Errorf("It should freeze money  once  but freeze %d ", wallet.CounterFreeze)

					}
					if wallet.CounterUnFreeze != 1 {
						t.Errorf("It should  un freeze money  but did un   freeze %d ", wallet.CounterUnFreeze)

					}
					if wallet.CounterGetFreeze != 1 {
						t.Errorf("It should get freeze money  once  but get  freeze %d ", wallet.CounterGetFreeze)

					}
					fmt.Printf("%s it  should not  create an  invalid transaction not enough coins \n", prefixOld)
				}
				TestShouldFailFailedToGetBalance := func(prefixOld string) {
					logger := &Logger.MockLogger{}
					mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
					findBalance := &mockFindBalance{}
					wallet := &MockWallet{}
					transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: logger, BalanceServiceInstance: findBalance, GeneratorService: mockGenerator}
					transaction := &TransactionCoins{Services: transactionStandard}
					transactionService := transaction
					err := transactionService.Construct()
					if err != nil {
						t.Errorf("expected to got  no err  but  got %v ", err)
					}

					wallet.Frozen = transaction.Transaction.Transaction.Amount
					logger.Logs = make([]string, 0)
					findBalance.err = errors.New("Failed to get balance")
					transaction.Transaction.Transaction.Amount = 15
					findBalance.amount = 100
					transaction.Transaction.Transaction.Amount = 15
					wallet.Frozen = 15
					errmsg := findBalance.err.Error()
					err = transactionService.CreateTransaction()
					if err.Error() != errmsg {
						t.Errorf("expected to get %s but got  %s", errmsg, err.Error())
					}
					if len(logger.Logs) != 1 {
						t.Errorf("Logger  should  receive 1  messages  but intead  got  %d ", len(logger.Logs))
					}
					if findBalance.findBalanceCalledTimes != 1 {
						t.Errorf("findbalance should  be  called  once but intead  called  %d ", findBalance.findBalanceCalledTimes)
					}
					if findBalance.locked {
						t.Errorf("It should call first  lock and  then unlock ")

					}
					if findBalance.lockedCall != 1 {
						t.Errorf("It should locked balance once  but lock %d ", findBalance.lockedCall)

					}
					if findBalance.unlockedCall != 1 {
						t.Errorf("It should unlocked balance once  but unlock %d ", findBalance.unlockedCall)

					}
					if wallet.CounterFreeze != 1 {
						t.Errorf("It should  freeze money  once  but freeze %d ", wallet.CounterFreeze)

					}
					if wallet.CounterUnFreeze != 1 {
						t.Errorf("It should   un freeze money  but did un   freeze %d ", wallet.CounterUnFreeze)

					}
					if wallet.CounterGetFreeze != 0 {
						t.Errorf("It should not  get freeze money  once  but get  freeze %d ", wallet.CounterGetFreeze)

					}
					fmt.Printf("%s it  should not  create an  invalid transaction failed to get balance \n", prefixOld)
				}
				TestShouldFailNotEnough(prefixNew)
				TestShouldFailFailedToGetBalance(prefixNew)
			}
			TestShouldCreate(prefixNew)
			TestShouldFail(prefixNew)
		}
		TestVerifyMockSign := func(prefixOld string) {
			fmt.Printf("%s Test Mock verifications transactions\n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			logger := &Logger.MockLogger{}
			mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
			findBalance := &mockFindBalance{}
			wallet := &MockWallet{}
			transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: logger, BalanceServiceInstance: findBalance, GeneratorService: mockGenerator}
			transaction := &TransactionCoins{Services: transactionStandard}
			transactionService := transaction
			err := transactionService.Construct()
			if err != nil {
				t.Errorf("expected to got  no err  but  got %v ", err)
			}
			findBalance.amount = 100
			transaction.Transaction.Transaction.Amount = 15
			wallet.Frozen = transaction.Transaction.Transaction.Amount
			err = transactionService.CreateTransaction()
			if err != nil {
				t.Errorf("It should be  valid not  get  err  %s", err.Error())
			}
			logger.Logs = make([]string, 0)
			TestItShouldVerify := func(prefixOld string) {
				fn := func(*rsa.PublicKey, crypto.Hash, []byte, []byte) error {
					return nil
				}

				oldVerify := transactionService.Services.verifyMethod
				defer func() {
					transactionService.Services.verifyMethod = oldVerify
				}()
				transactionService.Services.verifyMethod = fn
				err := transactionService.VerifySignature()
				if err != nil {
					t.Errorf("expected to got  no err  but  got %v ", err)
				}
				fmt.Printf("%s it  should  verify  a  valid transactions \n", prefixOld)
			}
			TestFail := func(prefixOld string) {
				fmt.Printf("%s Test Fail to verify transactions \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestInvalidSgn := func(prefixOld string) {
					errExpected := errors.New("invalid signiture")
					fn := func(*rsa.PublicKey, crypto.Hash, []byte, []byte) error {
						return errExpected
					}
					oldVerify := transactionService.Services.verifyMethod
					defer func() {
						transactionService.Services.verifyMethod = oldVerify
					}()
					transactionService.Services.verifyMethod = fn
					err := transactionService.VerifySignature()
					if err.Error() != errExpected.Error() {
						t.Errorf("expected to got  err %s  but  got %s ", errExpected.Error(), err.Error())
					}
					fmt.Printf("%s it  should  fail to verify  an  invalid transaction 'invalid signiture' \n", prefixOld)
				}
				TestMarshalFail := func(prefixOld string) {
					errExpected := errors.New("invalid signiture")
					fn := func(*rsa.PublicKey, crypto.Hash, []byte, []byte) error {
						return nil
					}
					oldStringfy := entitys.JsonStringfy
					defer func() {
						transactionService.Services.jsonStringfy = oldStringfy
					}()
					transactionService.Services.jsonStringfy = func(v entitys.TransactionRecord) ([]byte, error) {
						return nil, errExpected
					}
					oldVerify := transactionService.Services.verifyMethod
					defer func() {
						transactionService.Services.verifyMethod = oldVerify
					}()
					transactionService.Services.verifyMethod = fn
					err := transactionService.VerifySignature()
					if err.Error() != errExpected.Error() {
						t.Errorf("expected to got  err %s  but  got %s ", errExpected.Error(), err.Error())
					}
					fmt.Printf("%s it  should  fail to verify trasnaction that fails jasonStrigfy\n", prefixOld)
				}
				TestInvalidSgn(prefixNew)
				TestMarshalFail(prefixNew)
			}
			TestItShouldVerify(prefixNew)
			TestFail(prefixNew)
		}
		TestSetSign := func(prefixOld string) {
			fmt.Printf("%s Test For  Set Sign \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldSetSign := func(prefixOld string) {

				logger := &Logger.MockLogger{}
				mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
				findBalance := &mockFindBalance{}
				wallet := &MockWallet{}
				transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: logger, BalanceServiceInstance: findBalance, GeneratorService: mockGenerator}
				transaction := &TransactionCoins{Services: transactionStandard}
				transactionService := transaction
				err := transactionService.Construct()
				if err != nil {
					t.Errorf("expected to got  no err  but  got %v ", err)
				}
				logger.Logs = make([]string, 0)
				data := []byte("signiture")
				transactionService.setSign(data)
				if string(transactionService.Transaction.Signiture) != string(data) {
					t.Errorf("expected to get  signiture %v but  got  %v  ", data, transactionService.Transaction.Signiture)

				}
				fmt.Printf("%s it  should  set  a signiture\n", prefixOld)

			}
			itShouldSetSign(prefixNew)
		}
		TestGetters := func(prefixOld string) {
			fmt.Printf("%s Test For  Getters  \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			logger := &Logger.MockLogger{}
			mockGenerator := &Generator.MockGenerator{Response: "aaaaabbbbb"}
			findBalance := &mockFindBalance{}
			wallet := &MockWallet{}
			transactionStandard := TransactionsStandard{ServiceName: "transaction-service", WalletService: wallet, LoggerService: logger, BalanceServiceInstance: findBalance, GeneratorService: mockGenerator}
			transaction := &TransactionCoins{Services: transactionStandard}
			transactionService := transaction
			err := transactionService.Construct()
			if err != nil {
				t.Errorf("expected to got  no err  but  got %v ", err)
			}
			findBalance.amount = 100
			transaction.Transaction.Transaction.Amount = 15
			wallet.Frozen = transaction.Transaction.Transaction.Amount
			err = transactionService.CreateTransaction()
			if err != nil {
				t.Errorf("It should be  valid not  get  err  %s", err.Error())
			}
			logger.Logs = make([]string, 0)
			itShouldGetSigniture := func(prefixOld string) {

				data := []byte("signiture")
				transactionService.Transaction.Signiture = data
				if string(transactionService.GetSigniture()) != string(data) {
					t.Errorf("expected to get  signiture %v but  got  %v  ", data, transactionService.GetSigniture())

				}
			}
			itShouldGetTransaction := func(prefixOld string) {
				fmt.Printf("%s Test For  Get Transaction   \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				oldStringfy := entitys.JsonStringfy
				itShouldGet := func(prefixOld string) {
					defer func() {
						transactionService.Services.jsonStringfy = oldStringfy
					}()
					tr := []byte("transaction")
					transactionService.Services.jsonStringfy = func(v entitys.TransactionRecord) ([]byte, error) {
						return tr, nil
					}
					trActual, err := transactionService.GetTransaction()
					if err != nil {
						t.Errorf("expected to got  no err  but  got %v ", err)
					}
					if string(tr) != string(trActual) {
						t.Errorf("expected to get %v  but got  %v ", tr, trActual)

					}
					fmt.Printf("%s it  should  get a  transaction \n", prefixOld)

				}
				//TODO fail test case
				itShouldGet(prefixNew)
			}
			itShouldGetAmount := func(prefixOld string) {
				expected := float64(15)
				actual := transactionService.GetAmount()
				if expected != actual {
					t.Errorf("Get  amount expceted %.3f but  got %.3f ", expected, actual)
				}
				fmt.Printf("%s it  should  get amount  %.3f \n", prefixOld, expected)

			}
			itShouldGetSigniture(prefixNew)
			itShouldGetTransaction(prefixNew)
			itShouldGetAmount(prefixNew)
		}
		TestCreateService(prefixNew)
		TestCreateTransaction(prefixNew)
		TestVerifyMockSign(prefixNew)
		TestSetSign(prefixNew)
		TestGetters(prefixNew)
	}
	TestTransactionCoinsImpl(prefix)
}
