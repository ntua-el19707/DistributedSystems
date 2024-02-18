package WalletAndTransactions

import (
	"Logger"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
)

func TestWalletService(t *testing.T) {
	const prefix string = "----"
	fmt.Println("* Test Cases  For  wallet-service")

	TestForV1StructImpl := func(prefixOld string) {
		fmt.Printf("%s Test  For  WalletStructV1Implementation\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		buildWalletService := func() (WalletService, *Logger.MockLogger, error) {
			mockLogger := &Logger.MockLogger{}

			walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
			err := walletService.Construct()
			if err != nil {
				return nil, nil, err
			}
			return walletService, mockLogger, nil

		}
		TestInstaseCreate := func(prefixOld string) {

			fmt.Printf("%s Test  For service  creation \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

			TestWalletCreateServiceVersionStruct1 := func(prefixOld string) {

				_, logger, err := buildWalletService()
				if err != nil {
					t.Errorf("%s", err.Error())
				}
				if len(logger.Logs) != 4 {
					t.Errorf("It  should  have  been 4  messages  to logs  but  got %d", len(logger.Logs))
				}
				fmt.Printf("%s it  should  create  a new  wallet service vesrion struct 1\n", prefixOld)

			}
			TestWalletCreateServiceVersionStruct1FailPrivateConstructor := func(prefixOld string) {

				mockLogger := &Logger.MockLogger{}
				walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
				errSet := errors.New("failed  to generate keys")
				const errorTemplate = "Abbort  error :%s"
				errorMessage := fmt.Sprintf(failedToGenerateKeys, errSet.Error())
				errMsg := fmt.Sprintf(errorTemplate, errorMessage)
				err := walletService.privateConstructor(func(io.Reader, int) (*rsa.PrivateKey, error) { return nil, errSet })

				if err.Error() != errMsg {
					t.Errorf("expected to got err : '%s' but got '%s' ", errSet.Error(), err.Error())
				}
				if len(mockLogger.ErrorList) != 2 {
					t.Errorf("It  should  have  been 2  error  messages  to logs  but  got %d", len(mockLogger.ErrorList))
				}
				if len(mockLogger.Logs) != 2 {
					t.Errorf("It  should  have  been 2  messages  to logs  but  got %d", len(mockLogger.Logs))
				}
				fmt.Printf("%s it  should fail  via constuct private constructor 'make sure that the construct return private constructor' \n", prefixOld)
			}
			TestWalletCreateServiceVersionStruct1(prefixNew)
			TestWalletCreateServiceVersionStruct1FailPrivateConstructor(prefixNew)
		}
		TestGenarateWallet := func(prefixOld string) {
			fmt.Printf("%s Test  For generate a wallet \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestShouldGenerate := func(prefixOld string) {
				service, _, _ := buildWalletService()
				err := service.GenerateWallet(2048, rsa.GenerateKey)
				if err != nil {
					t.Errorf("excpeted no err  but  got %v", err)
				}

				//mock Checks
				fmt.Printf("%s it should generate wallet \n", prefixOld)
			}
			TestShouldFail := func(prefixOld string) {
				fmt.Printf("%s Test  For Fail to generate  wallet  \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestShouldFailToGenerate := func(prefixOld string) {
					service, mockLogger, _ := buildWalletService()
					mockLogger.Logs = make([]string, 0)
					mockLogger.ErrorList = make([]string, 0)
					errSet := errors.New("failed  to generate keys")
					const errorTemplate = "Abbort  error :%s"
					errorMessage := fmt.Sprintf(failedToGenerateKeys, errSet.Error())
					errMsg := fmt.Sprintf(errorTemplate, errorMessage)
					err := service.GenerateWallet(2048, func(io.Reader, int) (*rsa.PrivateKey, error) { return nil, errSet })
					if err.Error() != errMsg {
						t.Errorf("expected to got err : '%s' but got '%s' ", errSet.Error(), err.Error())
					}
					if len(mockLogger.ErrorList) != 1 {
						t.Errorf("It  should  have  been 1  error  messages  to logs  but  got %d", len(mockLogger.ErrorList))

					}
					if len(mockLogger.Logs) != 1 {
						t.Errorf("It  should  have  been 1  messages  to logs  but  got %d", len(mockLogger.Logs))
					}
					fmt.Printf("%s it should fail  to generate wallet due to   ('rsa.GenerateKey')\n", prefixOld)

				}
				TestShouldFailToGenerate(prefixNew)

			}
			TestShouldGenerate(prefixNew)
			TestShouldFail(prefixNew)
		}
		TestSign := func(prefixOld string) {
			fmt.Printf("%s Test  For Sign \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestShouldSign := func(prefixOld string) {
				service, _, _ := buildWalletService()
				tr := &MockTransactionService{}
				err := service.Sign(tr)

				if err != nil {
					t.Errorf("excpeted no err  but  got %v", err)
				}
				if tr.CallGetTransaction != 1 {
					t.Errorf("excpeted to  call  get Transaction %d times but  call %d ", 1, tr.CallGetTransaction)
				}
				if tr.CallSetSign != 1 {
					t.Errorf("excpeted to  call  set Sign  %d times but  call %d ", 1, tr.CallSetSign)
				}
				//mock Checks
				fmt.Printf("%s it should sign  a transaction\n", prefixOld)
			}
			TestShouldFail := func(prefixOld string) {
				fmt.Printf("%s Test  For Fail to Sign \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestShouldFailToSign := func(prefixOld string) {
					service, _, _ := buildWalletService()
					tr := &MockTransactionService{}
					const errorTemplate = "Abbort  error :  Failed  to   getTransaction  due to %s"
					errSet := errors.New("Failed to Get transaction")
					tr.ErrorGet = errSet
					err := service.Sign(tr)
					errMsg := fmt.Sprintf(errorTemplate, errSet.Error())

					if err.Error() != errMsg {
						t.Errorf("excpeted to go err %s  but  got %s", errMsg, err.Error())
					}
					if tr.CallGetTransaction != 1 {
						t.Errorf("excpeted to  call  get Transaction %d times but  call %d ", 1, tr.CallGetTransaction)
					}
					if tr.CallSetSign != 0 {
						t.Errorf("excpeted to  call  set Sign  %d times but  call %d ", 0, tr.CallSetSign)
					}
					//mock Checks
					fmt.Printf("%s it should fail  to sign  a transaction fue to ('failed to get Transaction')\n", prefixOld)
				}
				TestShouldFailToSignMethod := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}

					service := &WalletStructV1Implementation{LoggerService: mockLogger}
					err := service.Construct()
					if err != nil {
						t.Errorf("Expected no err but  got %v", err)
					}
					tr := &MockTransactionService{}
					const errorTemplate = "Abbort  error :  Failed  to   SignTransaction  due to %s"
					errSet := errors.New("Failed to Sign")
					errRsaSign := func(random io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, msg []byte) ([]byte, error) {
						return []byte{}, errSet
					}
					defer func() {
						service.signMethod = rsa.SignPKCS1v15
					}()
					service.signMethod = errRsaSign
					err = service.Sign(tr)
					errMsg := fmt.Sprintf(errorTemplate, errSet.Error())

					if err.Error() != errMsg {
						t.Errorf("excpeted to go err %s  but  got %s", errMsg, err.Error())
					}
					if tr.CallGetTransaction != 1 {
						t.Errorf("excpeted to  call  get Transaction %d times but  call %d ", 1, tr.CallGetTransaction)
					}
					if tr.CallSetSign != 0 {
						t.Errorf("excpeted to  call  set Sign  %d times but  call %d ", 0, tr.CallSetSign)
					}
					//mock Checks
					fmt.Printf("%s it should fail  to sign  a transaction fue to ('sign method  should  fail for some reason ')\n", prefixOld)
				}
				TestShouldFailToSign(prefixNew)
				TestShouldFailToSignMethod(prefixNew)
			}
			TestShouldFail(prefixNew)
			TestShouldSign(prefixNew)
		}
		TestGetPub := func(prefixOld string) {

			keyGen := func(n int) ([]rsa.PublicKey, error) {
				var publicKeys []rsa.PublicKey

				for i := 0; i < n; i++ {
					privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
					if err != nil {
						return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
					}

					publicKeys = append(publicKeys, privateKey.PublicKey)

				}
				return publicKeys, nil
			}
			keys, _ := keyGen(1)
			mockLogger := &Logger.MockLogger{}
			walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
			walletService.PublicKey = keys[0]
			pubKey := walletService.GetPub()
			if pubKey != keys[0] {
				t.Errorf("the  pub  key it should  get %v  but  got %v  ", pubKey, keys[0])
			}

			fmt.Printf("%s it should get a  public key\n", prefixOld)
		}
		TestForFreezeUnFreezeGetFreeze := func(prefixOld string) {
			fmt.Printf("%s Test  For Freezing  UnFreezin and GetFrozen coins  \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestFreeze := func(prefixOld string) {
				fmt.Printf("%s Freeze Test \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestShouldFrozen := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
					coins := float64(+10)
					err := walletService.Freeze(coins)
					if err != nil {
						t.Errorf("excpected to get no err   but got %s  ", err.Error())
					}
					if walletService.frozen != coins {
						t.Errorf("excpected to get   %.3f but got  %.3f ", float64(10), coins)
					}
					fmt.Printf("%s it should to froze  %3f coins \n", prefixOld, coins)
				}
				TestShouldNotFrozen := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
					coins := float64(-10)
					err := walletService.Freeze(coins)
					expected := fmt.Sprintf(errCannotFreezeUnFreezeNegativeCoins, "freeze")

					if err.Error() != expected {
						t.Errorf("excpected to err %s  but got %s  ", expected, err.Error())
					}
					fmt.Printf("%s it should fail to froze  %3f coins \n", prefixOld, coins)
				}
				TestShouldFrozen(prefixNew)
				TestShouldNotFrozen(prefixNew)
			}
			TestUnFreeze := func(prefixOld string) {
				fmt.Printf("%s UnFreeze Test \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestShouldUnFrozen := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
					walletService.frozen = 11
					coins := float64(+10)
					err := walletService.UnFreeze(coins)
					if err != nil {
						t.Errorf("excpected to get no err   but got %s  ", err.Error())
					}
					if walletService.frozen != 1.0 {
						t.Errorf("excpected to get   %.3f but got  %.3f ", float64(1), coins)
					}
					fmt.Printf("%s it should to froze  %3f coins \n", prefixOld, coins)
				}
				TestShouldNotUnFrozenNegativeCoins := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
					walletService.frozen = 11
					coins := float64(-10)
					err := walletService.UnFreeze(coins)
					expected := fmt.Sprintf(errCannotFreezeUnFreezeNegativeCoins, "unfreeze")

					if err.Error() != expected {
						t.Errorf("excpected to err %s  but got %s  ", expected, err.Error())
					}
					fmt.Printf("%s it should fail to unfroze  %3f coins \n", prefixOld, coins)
				}
				TestShouldNotUnFrozenNotEnough := func(prefixOld string) {
					mockLogger := &Logger.MockLogger{}
					walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
					walletService.frozen = 9
					coins := float64(10)
					err := walletService.UnFreeze(coins)
					expected := errCannotUnFreezeAndGoNegativeCoins

					if err.Error() != expected {
						t.Errorf("excpected to err %s  but got %s  ", expected, err.Error())
					}
					fmt.Printf("%s it should fail to unfroze  %3f coins \n", prefixOld, coins)
				}
				TestShouldUnFrozen(prefixNew)
				TestShouldNotUnFrozenNegativeCoins(prefixNew)
				TestShouldNotUnFrozenNotEnough(prefixNew)
			}
			TestFreeze(prefixNew)
			TestUnFreeze(prefixNew)
			TestSenarioBlockStyle := func(prefixOld string) {
				fmt.Printf("%s Block Style \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				mockLogger := &Logger.MockLogger{}
				walletService := &WalletStructV1Implementation{LoggerService: mockLogger}

				TestShouldFroze := func(prefixOld string) {
					coins := float64(15)
					err := walletService.Freeze(coins)
					if err != nil {
						t.Errorf("excpected to get no err but  got %v ", err)
					}
					if walletService.frozen != coins {
						t.Errorf("excpected to froze  %.3f but froze %.3f ", coins, walletService.frozen)
					}
					fmt.Printf("%s it should froze  %3f coins \n", prefixOld, coins)
				}
				TestShouldUnFroze := func(prefixOld string) {
					coins := float64(3)
					err := walletService.UnFreeze(coins)
					if err != nil {
						t.Errorf("excpected to get no err but  got %v ", err)
					}
					if walletService.frozen != float64(12) {
						t.Errorf("excpevted to have   %.3f but have  %.3f frozen  ", float64(12), walletService.frozen)
					}
					fmt.Printf("%s it should un froze  %3f coins \n", prefixOld, coins)
				}
				TestShouldGetFrozen := func(prefixOld string) {
					coins := walletService.GetFreeze()
					if walletService.frozen != coins {
						t.Errorf("excpected to get   %.3f but got  %.3f ", float64(12), coins)
					}
					fmt.Printf("%s it should get  %3f coins \n", prefixOld, float64(12))
				}
				TestShouldFroze(prefixNew)
				TestShouldUnFroze(prefixNew)
				TestShouldGetFrozen(prefixNew)
			}
			TestSenarioParrallel := func(prefixOld string) {
				fmt.Printf("%s Parallel  froze  and  infroze  want to test mu.Lock \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				mockLogger := &Logger.MockLogger{}
				walletService := &WalletStructV1Implementation{LoggerService: mockLogger}
				TestStableS_NRoutines := func(prefixOld string, stable float64, timeFroze, timeUnFroze int) {
					const frozen float64 = 1000000
					walletService.frozen = frozen
					var wg sync.WaitGroup
					executeFreeze := func(times int, wg *sync.WaitGroup) {
						fr := func(wg *sync.WaitGroup) {
							defer wg.Done()
							walletService.Freeze(float64(1))
							//	time.Sleep(time.Millisecond * 10)
						}
						for i := 0; i < times; i++ {
							go fr(wg)
						}
					}
					executeUnFreeze := func(times int, wg *sync.WaitGroup) {
						unFr := func(wg *sync.WaitGroup) {
							defer wg.Done()
							walletService.UnFreeze(float64(1))
							//time.Sleep(time.Millisecond * 100)
						}
						for i := 0; i < times; i++ {
							go unFr(wg)
						}
					}
					wg.Add(timeFroze + timeUnFroze)
					go executeFreeze(timeFroze, &wg)
					go executeUnFreeze(timeUnFroze, &wg)
					wg.Wait()
					if walletService.GetFreeze() != frozen+stable {
						t.Errorf("excpected to get   %.3f but got  %.3f ", frozen+stable, walletService.GetFreeze())
					}
					fmt.Printf("%s it should freeze %d and  unFreeze %d  in Parallel  and leave stable  +%.3f coins  in frozen \n", prefixOld, timeFroze, timeUnFroze, stable)
				}
				TestStableS_NRoutines(prefixNew, 0, 500, 500)
				TestStableS_NRoutines(prefixNew, 100, 600, 500)
				TestStableS_NRoutines(prefixNew, 500, 6000, 5500)
				TestStableS_NRoutines(prefixNew, 10000, 60000, 50000)
				//TestStableS_NRoutines(prefixNew, 100000, 600000, 500000)
				//	TestCaseZero  :=
			}
			TestSenarioBlockStyle(prefixNew)
			TestSenarioParrallel(prefixNew)
		}
		TestInstaseCreate(prefixNew)
		TestGenarateWallet(prefixNew)
		TestSign(prefixNew)
		TestGetPub(prefixNew)
		TestForFreezeUnFreezeGetFreeze(prefixNew)
	}
	TestForV1StructImpl(prefix)

}

/**
func TestWalletSign(t *testing.T) {

	service, _, _ := buildWalletService()
	tservice, _, _, _, transaction, _, _ := BuildTransactionService()
	transaction.Transaction.Transaction.Amount = 10
	err := tservice.CreateTransaction()
	if err != nil {
		t.Error("failed to create  transaction" + err.Error())
	}
	err = service.sign(tservice)

	if err != nil {
		t.Error("Failed to sign  Error " + err.Error())
	}
	//  err = tservice.VerifySignature()
	// if err != nil {
	//t.Error("Failed to verify Error " + err.Error())
	//}

	fmt.Println("it  should  sign  transaction ")

}*/
