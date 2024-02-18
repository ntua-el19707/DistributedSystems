package TransactionManager

//Transaction Manager Test
import (
	"FindBalance"
	"Logger"
	"WalletAndTransactions"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"testing"
)

// This  is  will be  integration - with unit providers
func TestTransactionManager(t *testing.T) {
	const prefix = "----"
	fmt.Println("* Test  For transactionManager-service")
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
	testImplementationTransactionManger := func(prefixOld string) {
		fmt.Printf("%s Test  For TransactionManager implenetation\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		serviceGen := func() (TransactionManagerService, *Logger.MockLogger, *WalletAndTransactions.MockWallet, *FindBalance.MockFindBalance, *TransactionManager, error) {
			// *  Build  the  world 'Providers'
			wallet := &WalletAndTransactions.MockWallet{}
			logger := &Logger.MockLogger{}
			findBallance := &FindBalance.MockFindBalance{}
			service := &TransactionManager{WalletServiceInstance: wallet, FindBalanceServiceInstance: findBallance, LoggerServiceInstance: logger}
			err := service.Construct()
			return service, logger, wallet, findBallance, service, err
		}
		TestCreateService := func(prefixOld string) {
			fmt.Printf("%s Test  For TransactionManager create service\n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldCreate := func(prefixOld string) {
				_, _, _, _, _, err := serviceGen()
				if err != nil {
					t.Errorf("expected no err  but  got %v", err)
				}
				fmt.Printf("%s it should  create  transaction manager service \n", prefixOld)
			}
			itShouldFail := func(prefixOld string) {
				fmt.Printf("%s Test  for Fail  to create service\n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

				wallet := &WalletAndTransactions.MockWallet{}
				logger := &Logger.MockLogger{}
				findBallance := &FindBalance.MockFindBalance{}
				itShouldFailNoWallet := func(prefixOld string) {
					errMsg := errNoWalletProvider
					service := &TransactionManager{FindBalanceServiceInstance: findBallance, LoggerServiceInstance: logger}
					err := service.Construct()
					if err.Error() != errMsg {
						t.Errorf("expected to get %s error  but  got %s", errMsg, err.Error())
					}
					fmt.Printf("%s it should  fail to create  transaction manager service 'no wallet service  provider' \n", prefixOld)
				}
				itShouldFailNoFindBalance := func(prefixOld string) {
					errMsg := errNoFindBalanceProvider
					service := &TransactionManager{WalletServiceInstance: wallet, LoggerServiceInstance: logger}
					err := service.Construct()
					if err.Error() != errMsg {
						t.Errorf("expected to get %s error  but  got %s", errMsg, err.Error())
					}
					fmt.Printf("%s it should  fail to create  transaction manager service 'no findBalance service  provider' \n", prefixOld)
				}
				itShouldFailNoWallet(prefixNew)
				itShouldFailNoFindBalance(prefixNew)
			}
			itShouldCreate(prefixNew)
			itShouldFail(prefixNew)
		}
		TestTransferMoney := func(prefixOld string) {
			fmt.Printf("%s Test  for Transfer coins \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			keys, _ := keyGen(3)
			service, _, wallet, findBalance, _, err := serviceGen()
			if err != nil {
				t.Errorf("expected no err  but  got %v", err)
			}
			wallet.MockPublicKey = keys[0]
			findBalance.Amount = float64(15) // 15bcc
			itShouldCreateTrnactionForTransferMoney := func(prefixOld string) {
				transactionSet, err := service.TransferMoney(keys[1], float64(10))

				if err != nil {
					t.Errorf("expected no err  but  got %v", err)
				}
				tax := 0.03 * float64(10)
				transfer := 0.97 * float64(10)
				taxActual := transactionSet.Tax.Transaction.Amount
				transferActual := transactionSet.Transfer.Transaction.Amount
				sender1 := transactionSet.Tax.Transaction.BillDetails.Bill.From.Address
				sender2 := transactionSet.Transfer.Transaction.BillDetails.Bill.From.Address
				Recever1 := transactionSet.Tax.Transaction.BillDetails.Bill.To.Address
				Recever2 := transactionSet.Transfer.Transaction.BillDetails.Bill.To.Address
				reason1 := transactionSet.Tax.Transaction.Reason
				reason2 := transactionSet.Transfer.Transaction.Reason
				if sender1 != sender2 {
					t.Errorf("expected senders  to be thw same but got %v , %v ", sender1, sender2)
				}
				if sender1 != keys[0] {
					t.Errorf("expected senders  to be %v but got %v  ", keys[0], sender1)
				}
				zero := rsa.PublicKey{}
				if Recever1 != zero {
					t.Errorf("expected receiver - tax  to be %v but got %v  ", rsa.PublicKey{}, Recever1)
				}
				if Recever2 != keys[1] {
					t.Errorf("expected receiver  transfer  to be %v but got %v  ", keys[1], Recever2)
				}
				if tax != taxActual {
					t.Errorf("expected tax to be  %.3f , but  got  %.3f", tax, taxActual)
				}
				if transfer != transferActual {
					t.Errorf("expected transfer to be  %.3f , but  got  %.3f", transfer, transferActual)
				}
				if reason1 != "fee" {
					t.Errorf("expected reason to be  %s , but  got  %s", "fee", reason1)

				}
				if reason2 != "Transfer" {
					t.Errorf("expected reason to be  %s , but  got  %s", "Transfer", reason2)

				}
				expectorCaller := func(exp, act int, caller string) {
					if exp != act {
						t.Errorf("Expected to call %d %s  but call %d", exp, caller, act)
					}
				}
				expectorCaller(2, wallet.CounterSign, "walletService.Sign")
				expectorCaller(0, wallet.CounterUnFreeze, "walletService.UnFreeze")
				fmt.Printf("%s it should be able  to transfer 10 bcc \n", prefixOld)

			}
			itShouldFail := func(prefixOld string) {

				fmt.Printf("%s Test  for Fails  \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				itShouldNotNotEnoughMoneyInBalance := func(prefixOld string) {
					wallet := &WalletAndTransactions.MockWallet{}
					logger := &Logger.MockLogger{}
					findBallance := &FindBalance.MockFindBalance{}
					findBallance.Amount = float64(15)
					wallet.Frozen = float64(16)
					service := &TransactionManager{WalletServiceInstance: wallet, FindBalanceServiceInstance: findBallance, LoggerServiceInstance: logger}
					err := service.Construct()
					_, err = service.TransferMoney(keys[1], float64(10))
					fmt.Println(err)
					if err == nil {
						t.Errorf("expected  err   but  got nothing")
					}
					if wallet.CounterUnFreeze != 2 {
						t.Errorf("expected  to call %d call UnFreeze but  %d times", 2, wallet.CounterUnFreeze)
					}
					fmt.Printf("%s it should be not  able  to transfer 10 bcc \n", prefixOld)

				}
				itShouldNotNotNegativeMoneyInBalance := func(prefixOld string) {
					wallet := &WalletAndTransactions.WalletStructV1Implementation{}
					wallet.Construct()
					logger := &Logger.MockLogger{}
					findBallance := &FindBalance.MockFindBalance{}
					findBallance.Amount = float64(15)
					service := &TransactionManager{WalletServiceInstance: wallet, FindBalanceServiceInstance: findBallance, LoggerServiceInstance: logger}
					err := service.Construct()
					_, err = service.TransferMoney(keys[1], float64(-10))
					fmt.Println(err)
					if err == nil {
						t.Errorf("expected  err   but  got nothing")
					}
					if wallet.GetFreeze() != 0 {
						t.Errorf("Freeze  money should  be  0 not  %.3f ", wallet.GetFreeze())
					}
					fmt.Printf("%s it should be not  able  to transfer -10 bcc \n", prefixOld)
				}
				itShouldNotNotSign := func(prefixOld string) {
					wallet := &WalletAndTransactions.MockWallet{}
					logger := &Logger.MockLogger{}
					findBallance := &FindBalance.MockFindBalance{}
					findBallance.Amount = float64(15)
					wallet.ErrorSignWallet = errors.New("Failed to  sign")
					service := &TransactionManager{WalletServiceInstance: wallet, FindBalanceServiceInstance: findBallance, LoggerServiceInstance: logger}
					err := service.Construct()
					_, err = service.TransferMoney(keys[1], float64(10))
					if err.Error() != wallet.ErrorSignWallet.Error() {
						t.Errorf("expected  err %s   but  got %s", wallet.ErrorSignWallet.Error(), err.Error())
					}
					if wallet.CounterUnFreeze != 2 {
						t.Errorf("expected  to call UnFreeze %d  but %d", 2, wallet.CounterUnFreeze)

					}
					fmt.Printf("%s it should be not  able  to transfer 10 bcc \n", prefixOld)
				}
				itShouldNotNotEnoughMoneyInBalance(prefixNew)
				itShouldNotNotNegativeMoneyInBalance(prefixNew)
				itShouldNotNotSign(prefixNew)
			}
			itShouldCreateTrnactionForTransferMoney(prefixNew)
			itShouldFail(prefixNew)
		}
		TestSendMessage := func(prefixOld string) {
		}
		TestCreateService(prefixNew)
		TestTransferMoney(prefixNew)
		TestSendMessage(prefixNew)
	}
	testImplementationTransactionManger(prefix)
}
