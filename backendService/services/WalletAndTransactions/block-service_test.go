package WalletAndTransactions

import (
	"FindBalance"
	"Hasher"
	"Logger"
	"Lottery"
	"RabbitMqService"
	"crypto/rsa"
	"entitys"
	"fmt"
	"testing"
)

func callExpector[T comparable](obj1, obj2 T, t *testing.T, prefix, what string) {
	if obj1 != obj2 {
		t.Errorf("%s  Expected  '%s' to get %v but %v ", prefix, what, obj1, obj2)
	}
}
func TestBlockChainService(t *testing.T) {
	expectorNoErr := func(t *testing.T, prefix string, err error) {
		if err != nil {
			t.Errorf("%s  Expected  no err  but got %v err ", prefix, err)
		}
	}
	walletGen := func(prefix string, n int) []WalletStructV1Implementation {
		wallets := make([]WalletStructV1Implementation, n)
		for i := 0; i < n; i++ {
			wallets[i] = WalletStructV1Implementation{}
			err := wallets[i].Construct()
			expectorNoErr(t, prefix, err)
		}
		return wallets
	}
	transactionSetGen := func(prefix string, wallet1, wallet2, signWallet *WalletStructV1Implementation, amount, tax float64) entitys.TransactionCoinSet {
		transfer := (1.0 - tax) * amount
		taxBCC := tax * amount
		findBalance := &FindBalance.MockFindBalance{}
		findBalance.Amount = amount
		TransactionMake := func(Receiver rsa.PublicKey, money float64, reason string) TransactionCoins {

			receiver := entitys.Client{Address: Receiver}
			pair := entitys.BillingInfo{To: receiver}
			standard := TransactionsStandard{ServiceName: "transaction-coin-service", BalanceServiceInstance: findBalance, WalletService: wallet1}
			transactionInfo := entitys.TransactionCoins{BillDetails: entitys.TransactionDetails{Bill: pair}, Amount: money, Reason: reason}
			transaction := TransactionCoins{Transaction: entitys.TransactionCoinEntityRoot{Transaction: transactionInfo}, Services: standard}
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
	giveMeCoinChain := func(prefixOld string, wallet *WalletStructV1Implementation) (*BlockChainCoinsImpl, *RabbitMqService.MockRabbitMqImpl, *Lottery.MockLottery, *Hasher.HashImpl) {

		lage := RabbitMqService.MockRabbitMqImpl{Channel: make(chan entitys.TransactionCoinSet)}
		hash := Hasher.HashImpl{}
		err := hash.Construct()
		expectorNoErr(t, prefixOld, err)
		lottery := Lottery.MockLottery{}
		providers := BlockServiceProviders{WalletServiceInstance: wallet, HashService: &hash, RabbitMqService: &lage, LotteryService: &lottery}
		service := BlockChainCoinsImpl{Services: providers}

		err = service.Construct()
		expectorNoErr(t, prefixOld, err)
		err = service.Genesis()
		expectorNoErr(t, prefixOld, err)
		return &service, &lage, &lottery, &hash
	}
	fmt.Println("* Test For  BlockChainService")
	const prefix string = "----"

	TestCoinsImpl := func(prefixOld string) {
		fmt.Printf("%s  Test For  Coins  Impl ", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		IntegrationTest := func(prefixOld string) {
			fmt.Printf("%s Integration   Test For  Coins  Impl ", prefixOld)
			prefixNew := fmt.Sprintf("%s-integation-%s", prefixOld, prefix)

			TestGenesis := func(prefixOld string) {
				wallet := WalletStructV1Implementation{}
				wallet.Construct()
				hash := Hasher.HashImpl{}
				hash.Construct()
				providers := BlockServiceProviders{WalletServiceInstance: &wallet, HashService: &hash}
				service := BlockChainCoinsImpl{Services: providers}
				service.Construct()
				service.Genesis()
				fmt.Println(service.Chain)
			}
			TestInsertTransactions := func(prefixOld string) {
				fmt.Printf("%s Test For Insert  Transaction Coins", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				TestSucceed := func(prefixOld string) {
					fmt.Printf("%s Test For succeed Insert  Transaction Coins", prefixOld)
					prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

					itShouldInsertAValidSet := func(prefixOld string) {
						wallets := walletGen(prefixNew, 3)
						tr1 := transactionSetGen(prefixNew, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
						service, lage, lottery, _ := giveMeCoinChain(prefixOld, &wallets[2])
						freezeBefore := wallets[2].GetFreeze()
						service.InsertTransaction(tr1)
						sizeActual := len(service.Chain[0].Transactions)
						callExpector[int](4, sizeActual, t, prefixOld, "lenght  of  Transaction  in chain index 0 ")
						taxTo := wallets[2].GetPub()
						tr1.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						freezeAfter := wallets[2].GetFreeze()
						wallets[0].UnFreeze(float64(10))
						list := service.Chain[0].Transactions
						callExpector[entitys.TransactionCoins](list[2], tr1.Tax.Transaction, t, prefixOld, "trSet  1 Tax")
						callExpector[entitys.TransactionCoins](list[3], tr1.Transfer.Transaction, t, prefixOld, "trSet  1 Transfer")
						callExpector[float64](freezeBefore, freezeAfter, t, prefixOld, "no change in frozen money")
						callExpector[int](0, lage.CallPublishBlock, t, prefixOld, "call publish block")
						callExpector[int](0, lage.CallConsumeBlock, t, prefixOld, "call consume block")
						callExpector[int](0, lottery.CallLoadStake, t, prefixOld, "call load stake service")
						callExpector[int](0, lottery.CallSpin, t, prefixOld, "call spin ")
						fmt.Printf("%s  it should  insert  a valid  set  of 10  coins\n", prefixOld)
					}
					itShouldInsertAValidSetToSenderChain := func(prefixOld string) {
						wallets := walletGen(prefixNew, 3)
						tr1 := transactionSetGen(prefixNew, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
						service, lage, lottery, _ := giveMeCoinChain(prefixOld, &wallets[0])
						freezeBefore := wallets[0].GetFreeze()
						service.InsertTransaction(tr1)
						sizeActual := len(service.Chain[0].Transactions)
						callExpector[int](4, sizeActual, t, prefixOld, "lenght  of  Transaction  in chain index 0 ")
						taxTo := wallets[0].GetPub()
						tr1.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						freezeAfter := wallets[0].GetFreeze()
						list := service.Chain[0].Transactions
						callExpector[entitys.TransactionCoins](list[2], tr1.Tax.Transaction, t, prefixOld, "trSet  1 Tax")
						callExpector[entitys.TransactionCoins](list[3], tr1.Transfer.Transaction, t, prefixOld, "trSet  1 Transfer")
						callExpector[float64](freezeBefore, freezeAfter+float64(10), t, prefixOld, "no change in frozen money")
						callExpector[int](0, lage.CallPublishBlock, t, prefixOld, "call publish block")
						callExpector[int](0, lage.CallConsumeBlock, t, prefixOld, "call consume block")
						callExpector[int](0, lottery.CallLoadStake, t, prefixOld, "call load stake service")
						callExpector[int](0, lottery.CallSpin, t, prefixOld, "call spin ")
						fmt.Printf("%s  it should  insert  a valid  set  of 10  coins Sender \n", prefixOld)
					}
					itShouldInsert2ValidTransactions := func(prefixOld string) {

						wallets := walletGen(prefixNew, 3)
						tr1 := transactionSetGen(prefixNew, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
						tr2 := transactionSetGen(prefixNew, &wallets[2], &wallets[1], &wallets[2], 5, 0.03)
						service, lage, lottery, _ := giveMeCoinChain(prefixOld, &wallets[0])
						freezeBefore := wallets[0].GetFreeze()
						freezeBefore2 := wallets[2].GetFreeze()
						service.InsertTransaction(tr1)
						service.InsertTransaction(tr2)
						sizeActual := len(service.Chain[0].Transactions)
						callExpector[int](6, sizeActual, t, prefixOld, "lenght  of  Transaction  in chain index 0 ")
						taxTo := wallets[0].GetPub()
						tr1.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						tr2.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						freezeAfter := wallets[0].GetFreeze()
						freezeAfter2 := wallets[2].GetFreeze()
						list := service.Chain[0].Transactions
						callExpector[entitys.TransactionCoins](list[2], tr1.Tax.Transaction, t, prefixOld, "trSet  1 Tax")
						callExpector[entitys.TransactionCoins](list[3], tr1.Transfer.Transaction, t, prefixOld, "trSet  1 Transfer")
						callExpector[entitys.TransactionCoins](list[4], tr2.Tax.Transaction, t, prefixOld, "trSet  2 Tax")
						callExpector[entitys.TransactionCoins](list[5], tr2.Transfer.Transaction, t, prefixOld, "trSet  2 Transfer")
						callExpector[float64](freezeBefore, freezeAfter+float64(10), t, prefixOld, "no change in frozen money")
						callExpector[float64](freezeBefore2, freezeAfter2, t, prefixOld, "no change in frozen money")
						callExpector[int](0, lage.CallPublishBlock, t, prefixOld, "call publish block")
						callExpector[int](0, lage.CallConsumeBlock, t, prefixOld, "call consume block")
						callExpector[int](0, lottery.CallLoadStake, t, prefixOld, "call load stake service")
						callExpector[int](0, lottery.CallSpin, t, prefixOld, "call spin ")
						fmt.Printf("%s  it should  insert 2 valid  set  of 10 + 5 coins\n", prefixOld)
					}
					itShouldInsertAndMineNotMiner := func(prefixOld string) {
						wallets := walletGen(prefixNew, 3)
						tr1 := transactionSetGen(prefixNew, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
						tr2 := transactionSetGen(prefixNew, &wallets[2], &wallets[1], &wallets[2], 5, 0.03)
						tr3 := transactionSetGen(prefixNew, &wallets[1], &wallets[0], &wallets[1], 3, 0.03)
						service, lage, lottery, hash := giveMeCoinChain(prefixOld, &wallets[0])
						block := entitys.BlockCoinEntity{}
						lottery.SpinRsp = wallets[1].GetPub()
						err := block.MineBlock(lottery.SpinRsp, service.Chain[0].BlockEntity, &Logger.MockLogger{}, hash)
						expectorNoErr(t, prefixOld, err)
						lage.Blocks = make([]entitys.BlockCoinMessageRabbitMq, 2)
						lage.Blocks[0] = entitys.BlockCoinMessageRabbitMq{BlockCoin: block}
						freezeBefore := wallets[0].GetFreeze()
						freezeBefore2 := wallets[2].GetFreeze()
						freezeBefore3 := wallets[1].GetFreeze()
						service.InsertTransaction(tr1)
						service.InsertTransaction(tr2)
						service.InsertTransaction(tr3)
						sizeActual := len(service.Chain[0].Transactions)
						callExpector[int](6, sizeActual, t, prefixOld, "lenght  of  Transaction  in chain index 0 ")
						taxTo := wallets[0].GetPub()
						tr1.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						tr2.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						tr3.Tax.Transaction.BillDetails.Bill.To.Address = lottery.SpinRsp
						freezeAfter := wallets[0].GetFreeze()
						freezeAfter2 := wallets[2].GetFreeze()
						freezeAfter3 := wallets[1].GetFreeze()
						list := service.Chain[0].Transactions
						list2 := service.Chain[1].Transactions
						callExpector[entitys.TransactionCoins](list[2], tr1.Tax.Transaction, t, prefixOld, "trSet  1 Tax")
						callExpector[entitys.TransactionCoins](list[3], tr1.Transfer.Transaction, t, prefixOld, "trSet  1 Transfer")
						callExpector[entitys.TransactionCoins](list[4], tr2.Tax.Transaction, t, prefixOld, "trSet  2 Tax")
						callExpector[entitys.TransactionCoins](list[5], tr2.Transfer.Transaction, t, prefixOld, "trSet  2 Transfer")
						callExpector[entitys.TransactionCoins](list2[0], tr3.Tax.Transaction, t, prefixOld, "trSet  3 Tax")
						callExpector[entitys.TransactionCoins](list2[1], tr3.Transfer.Transaction, t, prefixOld, "trSet  3 Transfer")
						callExpector[float64](freezeBefore, freezeAfter+float64(10), t, prefixOld, "no change in frozen money")
						callExpector[float64](freezeBefore2, freezeAfter2, t, prefixOld, "no change in frozen money")
						callExpector[float64](freezeBefore3, freezeAfter3, t, prefixOld, "no change in frozen money")
						callExpector[int](0, lage.CallPublishBlock, t, prefixOld, "call publish block")
						callExpector[int](1, lage.CallConsumeBlock, t, prefixOld, "call consume block")
						callExpector[int](1, lottery.CallLoadStake, t, prefixOld, "call load stake service")
						callExpector[int](1, lottery.CallSpin, t, prefixOld, "call spin ")

						fmt.Printf("%s  it should  insert 3 valid and  it is  no miner    set  of 10 + 5 coins\n", prefixOld)
					}
					itShouldInsertAndMineMiner := func(prefixOld string) {
						wallets := walletGen(prefixNew, 3)
						tr1 := transactionSetGen(prefixNew, &wallets[0], &wallets[1], &wallets[0], 10, 0.03)
						tr2 := transactionSetGen(prefixNew, &wallets[2], &wallets[1], &wallets[2], 5, 0.03)
						tr3 := transactionSetGen(prefixNew, &wallets[1], &wallets[0], &wallets[1], 3, 0.03)
						service, lage, lottery, hash := giveMeCoinChain(prefixOld, &wallets[0])
						block := entitys.BlockCoinEntity{}
						lottery.SpinRsp = wallets[0].GetPub()
						err := block.MineBlock(lottery.SpinRsp, service.Chain[0].BlockEntity, &Logger.MockLogger{}, hash)
						expectorNoErr(t, prefixOld, err)
						lage.Blocks = make([]entitys.BlockCoinMessageRabbitMq, 2)
						lage.Blocks[0] = entitys.BlockCoinMessageRabbitMq{BlockCoin: block}
						freezeBefore := wallets[0].GetFreeze()
						freezeBefore2 := wallets[2].GetFreeze()
						freezeBefore3 := wallets[1].GetFreeze()
						service.InsertTransaction(tr1)
						service.InsertTransaction(tr2)
						service.InsertTransaction(tr3)
						sizeActual := len(service.Chain[0].Transactions)
						callExpector[int](6, sizeActual, t, prefixOld, "lenght  of  Transaction  in chain index 0 ")
						taxTo := wallets[0].GetPub()
						tr1.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						tr2.Tax.Transaction.BillDetails.Bill.To.Address = taxTo
						tr3.Tax.Transaction.BillDetails.Bill.To.Address = lottery.SpinRsp
						freezeAfter := wallets[0].GetFreeze()
						freezeAfter2 := wallets[2].GetFreeze()
						freezeAfter3 := wallets[1].GetFreeze()
						list := service.Chain[0].Transactions
						list2 := service.Chain[1].Transactions
						callExpector[entitys.TransactionCoins](list[2], tr1.Tax.Transaction, t, prefixOld, "trSet  1 Tax")
						callExpector[entitys.TransactionCoins](list[3], tr1.Transfer.Transaction, t, prefixOld, "trSet  1 Transfer")
						callExpector[entitys.TransactionCoins](list[4], tr2.Tax.Transaction, t, prefixOld, "trSet  2 Tax")
						callExpector[entitys.TransactionCoins](list[5], tr2.Transfer.Transaction, t, prefixOld, "trSet  2 Transfer")
						callExpector[entitys.TransactionCoins](list2[0], tr3.Tax.Transaction, t, prefixOld, "trSet  3 Tax")
						callExpector[entitys.TransactionCoins](list2[1], tr3.Transfer.Transaction, t, prefixOld, "trSet  3 Transfer")
						callExpector[float64](freezeBefore, freezeAfter+float64(10), t, prefixOld, "no change in frozen money")
						callExpector[float64](freezeBefore2, freezeAfter2, t, prefixOld, "no change in frozen money")
						callExpector[float64](freezeBefore3, freezeAfter3, t, prefixOld, "no change in frozen money")
						callExpector[int](1, lage.CallPublishBlock, t, prefixOld, "call publish block")
						callExpector[int](1, lage.CallConsumeBlock, t, prefixOld, "call consume block")
						callExpector[int](1, lottery.CallLoadStake, t, prefixOld, "call load stake service")
						callExpector[int](1, lottery.CallSpin, t, prefixOld, "call spin ")

						fmt.Printf("%s  it should  insert 3 valid and  Mine   set  of 10 + 5 + 3 coins\n", prefixOld)
					}
					itShouldInsertAValidSet(prefixNew)
					itShouldInsertAValidSetToSenderChain(prefixNew)
					itShouldInsert2ValidTransactions(prefixNew)
					itShouldInsertAndMineNotMiner(prefixNew)
					itShouldInsertAndMineMiner(prefixNew)
				}
				TestSucceed(prefixNew)

			}
			TestGenesis(prefixNew)
			TestInsertTransactions(prefixNew)
		}

		IntegrationTest(prefixNew)
	}
	TestCoinsImpl(prefix)

}
