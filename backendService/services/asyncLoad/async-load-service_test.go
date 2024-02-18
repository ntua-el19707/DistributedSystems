package asyncLoad

import (
	"FindBalance"
	"Hasher"
	"Logger"
	"Lottery"
	"RabbitMqService"
	"WalletAndTransactions"
	"crypto/rsa"
	"entitys"
	"fmt"
	"testing"
	"time"
)

func callExpector[T comparable](obj1, obj2 T, t *testing.T, prefix, what string) {
	if obj1 != obj2 {
		t.Errorf("%s  Expected  '%s' to get %v but %v ", prefix, what, obj1, obj2)
	}
}
func TestAsyncLoad(t *testing.T) {
	//* UseFull
	/*	expectorCaller := func(t *testing.T, prefix, caller string, expected, actual int) {
		if expected != actual {
			t.Errorf("%s  Expcted  func '%s' to be called  %d  times  but  call %d ", prefix, caller, expected, actual)
		}
	}*/
	expectorNoErr := func(t *testing.T, prefix string, err error) {
		if err != nil {
			t.Errorf("%s  Expected  no err  but got %v err ", prefix, err)
		}
	}
	walletGen := func(prefix string, n int) []WalletAndTransactions.WalletStructV1Implementation {
		wallets := make([]WalletAndTransactions.WalletStructV1Implementation, n)
		for i := 0; i < n; i++ {
			wallets[i] = WalletAndTransactions.WalletStructV1Implementation{}
			err := wallets[i].Construct()
			expectorNoErr(t, prefix, err)
		}
		return wallets
	}
	transactionSetGen := func(prefix string, wallet1, wallet2, signWallet *WalletAndTransactions.WalletStructV1Implementation, amount, tax float64) entitys.TransactionCoinSet {
		transfer := (1.0 - tax) * amount
		taxBCC := tax * amount
		findBalance := &FindBalance.MockFindBalance{}
		findBalance.Amount = amount
		TransactionMake := func(Receiver rsa.PublicKey, money float64, reason string) WalletAndTransactions.TransactionCoins {

			receiver := entitys.Client{Address: Receiver}
			pair := entitys.BillingInfo{To: receiver}
			standard := WalletAndTransactions.TransactionsStandard{ServiceName: "transaction-coin-service", BalanceServiceInstance: findBalance, WalletService: wallet1}
			transactionInfo := entitys.TransactionCoins{BillDetails: entitys.TransactionDetails{Bill: pair}, Amount: money, Reason: reason}
			transaction := WalletAndTransactions.TransactionCoins{Transaction: entitys.TransactionCoinEntityRoot{Transaction: transactionInfo}, Services: standard}
			err := transaction.Construct()
			expectorNoErr(t, prefix, err)
			err = transaction.CreateTransaction()
			expectorNoErr(t, prefix, err)
			err = signWallet.Sign(&transaction)
			expectorNoErr(t, prefix, err)

			return transaction
		}
		var tSet entitys.TransactionCoinSet
		tSet.Tax = TransactionMake(rsa.PublicKey{}, taxBCC, "Tax").Transaction
		tSet.Transfer = TransactionMake(wallet2.GetPub(), transfer, "Transfer").Transaction
		return tSet

	}

	const prefix string = "----"
	fmt.Println("*  Test  For  AsyncLoadService")
	TestImplemtation := func(prefixOld string) {
		fmt.Printf("%s Test  For  AsyncLoadImpl\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestCreateForCreateService := func(prefixOld string) {
			//TODO	TEST CASES  FOR creation  of   service
		}
		TestConsumeTransactions := func(prefixOld string) {
			fmt.Printf("%s Test  For  InsertTransaction 'Hybrid' testing Mocking & Integration\n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			//*Services  of world  3  wallets

			itShouldSucceed := func(prefixOld string) {
				//Create  MockRabbitMQ Control data flow
				fmt.Printf("%s Test  For  success InsertTransaction 'Hybrid' testing Mocking & Integration\n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				itShouldInsert2transaction := func(prefixOld string) {

					lage := RabbitMqService.MockRabbitMqImpl{Channel: make(chan entitys.TransactionCoinSet)}
					logger := Logger.MockLogger{}
					wallets := walletGen(prefixNew, 3)
					trSet1 := transactionSetGen(prefixOld, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
					err := wallets[0].UnFreeze(10)
					expectorNoErr(t, prefixOld, err)
					trSet2 := transactionSetGen(prefixOld, &wallets[0], &wallets[1], &wallets[0], 5, 0.03)
					hash := Hasher.HashImpl{}
					hash.Construct()

					providers := WalletAndTransactions.BlockServiceProviders{WalletServiceInstance: &wallets[2], HashService: &hash}
					service := WalletAndTransactions.BlockChainCoinsImpl{Services: providers}

					err = service.Construct()
					expectorNoErr(t, prefixOld, err)

					err = service.Genesis()
					expectorNoErr(t, prefixOld, err)

					providersAsync := AsyncLoadProviders{RabbitMqService: &lage, LoggerService: &logger, BlockCoinService: &service}
					serviceAsync := AsyncLoadImpl{Providers: providersAsync}

					err = serviceAsync.Construct()
					expectorNoErr(t, prefixOld, err)
					go serviceAsync.consumeTransactions()
					time.Sleep(time.Second)
					lage.Channel <- trSet1
					lage.Channel <- trSet2
					time.Sleep(time.Second)
					list := service.Chain[0].Transactions
					trSet1.Tax.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
					trSet2.Tax.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
					callExpector(list[2], trSet1.Tax.Transaction, t, prefixOld, "Trsansaction Tax 0.3")
					callExpector(list[3], trSet1.Transfer.Transaction, t, prefixOld, "Trsansaction Transafer 9.7")
					callExpector(list[4], trSet2.Tax.Transaction, t, prefixOld, "Trsansaction trSet2")
					callExpector(list[5], trSet2.Transfer.Transaction, t, prefixOld, "Trsansaction trSet2")
					fmt.Printf("%s it should  sucessfuly  insert  2  intergate transactions \n", prefixOld)
				}
				itShouldInsert2transaction(prefixNew)
				itShouldInsert3transactionAndMine := func(prefixOld string) {

					lage := RabbitMqService.MockRabbitMqImpl{Channel: make(chan entitys.TransactionCoinSet)}
					logger := Logger.MockLogger{}
					wallets := walletGen(prefixNew, 3)
					trSet1 := transactionSetGen(prefixOld, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
					err := wallets[0].UnFreeze(10)
					expectorNoErr(t, prefixOld, err)
					trSet2 := transactionSetGen(prefixOld, &wallets[0], &wallets[1], &wallets[0], 5, 0.03)
					err = wallets[0].UnFreeze(5)
					expectorNoErr(t, prefixOld, err)
					trSet3 := transactionSetGen(prefixOld, &wallets[0], &wallets[2], &wallets[0], 100, 0.03)
					err = wallets[0].UnFreeze(100)
					expectorNoErr(t, prefixOld, err)

					hash := Hasher.HashImpl{}
					hash.Construct()
					lottery := Lottery.MockLottery{}
					lottery.SpinRsp = wallets[1].GetPub()
					providers := WalletAndTransactions.BlockServiceProviders{RabbitMqService: &lage, WalletServiceInstance: &wallets[2], HashService: &hash, LotteryService: &lottery}
					service := WalletAndTransactions.BlockChainCoinsImpl{Services: providers}

					err = service.Construct()
					expectorNoErr(t, prefixOld, err)
					err = service.Genesis()
					expectorNoErr(t, prefixOld, err)

					block := entitys.BlockCoinEntity{}
					err = block.MineBlock(lottery.SpinRsp, service.Chain[0].BlockEntity, &Logger.MockLogger{}, &hash)
					expectorNoErr(t, prefixOld, err)
					lage.Blocks = make([]entitys.BlockCoinMessageRabbitMq, 2)
					lage.Blocks[0] = entitys.BlockCoinMessageRabbitMq{BlockCoin: block}
					providersAsync := AsyncLoadProviders{RabbitMqService: &lage, LoggerService: &logger, BlockCoinService: &service}
					serviceAsync := AsyncLoadImpl{Providers: providersAsync}

					err = serviceAsync.Construct()
					expectorNoErr(t, prefixOld, err)
					go serviceAsync.consumeTransactions()
					time.Sleep(time.Second)
					lage.Channel <- trSet1
					lage.Channel <- trSet2
					lage.Channel <- trSet3

					time.Sleep(time.Second * 6)
					list := service.Chain[0].Transactions
					list2 := service.Chain[1].Transactions
					trSet1.Tax.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
					trSet2.Tax.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
					trSet3.Tax.Transaction.BillDetails.Bill.To.Address = lottery.SpinRsp
					callExpector(list[2], trSet1.Tax.Transaction, t, prefixOld, "Trsansaction Tax 0.3")
					callExpector(list[3], trSet1.Transfer.Transaction, t, prefixOld, "Trsansaction Transafer 9.7")
					callExpector(list[4], trSet2.Tax.Transaction, t, prefixOld, "Trsansaction trSet2")
					callExpector(list[5], trSet2.Transfer.Transaction, t, prefixOld, "Trsansaction trSet2")
					callExpector(list2[0], trSet3.Tax.Transaction, t, prefixOld, "Trsansaction trSet3")
					callExpector(list2[1], trSet3.Transfer.Transaction, t, prefixOld, "Trsansaction trSet3")
					fmt.Printf("%s it should  sucessfuly  insert  2  intergate transactions \n", prefixOld)
				}
				itShouldInsert3transactionAndMine(prefixNew)
			}
			itShouldFail := func(prefixOld string) {
				fmt.Printf("%s Test  For  Fail InsertTransaction 'Hybrid' testing Mocking & Integration\n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

				lage := RabbitMqService.MockRabbitMqImpl{Channel: make(chan entitys.TransactionCoinSet)}
				logger := Logger.MockLogger{}
				wallets := walletGen(prefixNew, 3)
				trSet1 := transactionSetGen(prefixOld, &wallets[0], &wallets[1], &wallets[2], 10, 0.03)
				err := wallets[0].UnFreeze(10)
				expectorNoErr(t, prefixOld, err)
				trSet2 := transactionSetGen(prefixOld, &wallets[0], &wallets[1], &wallets[0], 5, 0.03)
				itShouldFailtsSet1transaction := func(prefixOld string) {

					hash := Hasher.HashImpl{}
					hash.Construct()

					providers := WalletAndTransactions.BlockServiceProviders{WalletServiceInstance: &wallets[2], HashService: &hash}
					service := WalletAndTransactions.BlockChainCoinsImpl{Services: providers}

					err = service.Construct()
					expectorNoErr(t, prefixOld, err)

					err = service.Genesis()
					expectorNoErr(t, prefixOld, err)

					providersAsync := AsyncLoadProviders{RabbitMqService: &lage, LoggerService: &logger, BlockCoinService: &service}
					serviceAsync := AsyncLoadImpl{Providers: providersAsync}

					err = serviceAsync.Construct()
					expectorNoErr(t, prefixOld, err)
					go serviceAsync.consumeTransactions()
					time.Sleep(time.Second)
					lage.Channel <- trSet1
					lage.Channel <- trSet2
					time.Sleep(time.Second)
					list := service.Chain[0].Transactions
					trSet1.Tax.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
					trSet2.Tax.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
					callExpector(len(service.Chain[0].Transactions), 4, t, prefixOld, "Only bootstrab and  trset2")

					callExpector(list[2], trSet2.Tax.Transaction, t, prefixOld, "trset2 tax")
					callExpector(list[3], trSet2.Transfer.Transaction, t, prefixOld, "Trsansaction trSet2")
					fmt.Printf("%s it should fail innsert  the tset1 but succed trset2  intergate transactions \n", prefixOld)
				}
				itShouldFailtsSet1transaction(prefixNew)

			}
			itShouldSucceed(prefixNew)
			itShouldFail(prefixNew)

		}
		TestCreateForCreateService(prefixNew)
		TestConsumeTransactions(prefixNew)
	}
	TestImplemtation(prefix)

}
