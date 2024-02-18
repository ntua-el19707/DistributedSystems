package WalletAndTransactions

import (
	"FindBalance"
	"Hasher"
	"Logger"
	"Lottery"
	"RabbitMqService"
	"crypto/rsa"
	"entitys"
	"errors"
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
	TestMsgImpl := func(prefixOld string) {
		fmt.Printf("%s  Test For  Message   Impl\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		IntegrationTest := func(prefixOld string) {
			fmt.Printf("%s Integration   Test For  Message  Impl\n", prefixOld)
			prefixNew := fmt.Sprintf("%s-integation-%s", prefixOld, prefix)

			TestGenesis := func(prefixOld string) {
				fmt.Printf("%s Integration   Test For  Message  Impl  Genesis\n", prefixOld)
				prefixNew := fmt.Sprintf("%s-integation-%s", prefixOld, prefix)
				itShouldCreate := func(prefixOld string) {
					wallet := WalletStructV1Implementation{}
					err := wallet.Construct()
					expectorNoErr(t, prefixOld, err)
					hash := Hasher.HashImpl{}
					err = hash.Construct()
					expectorNoErr(t, prefixOld, err)
					providers := BlockServiceProviders{WalletServiceInstance: &wallet, HashService: &hash}
					service := BlockChainMsgImpl{Services: providers}
					err = service.Construct()
					expectorNoErr(t, prefixOld, err)
					callExpector[int](0, len(service.Chain), t, prefixOld, "size  of chain before")
					err = service.Genesis()
					expectorNoErr(t, prefixOld, err)
					callExpector[int](1, len(service.Chain), t, prefixOld, "size  of chain ")
					callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "size  of chain[0].Transactions ")

					fmt.Printf("%s it should  genesis \n", prefixOld)
				}
				itShouldCreate(prefixNew)
			}
			TestInsertTransactions := func(prefixOld string) {
				fmt.Printf("%s Integration   Test For  Insert Message  Transactions \n", prefixOld)
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

				lage := RabbitMqService.MockRabbitMqImpl{}
				walletGen := func(n int) []WalletStructV1Implementation {
					list := make([]WalletStructV1Implementation, n)
					for i := 0; i < n; i++ {
						wallet := WalletStructV1Implementation{}
						err := wallet.Construct()
						expectorNoErr(t, prefixOld, err)
						list[i] = wallet
					}

					return list
				}
				createTransaction := func(wallet1, wallet2, walletSign *WalletStructV1Implementation, msg string) entitys.TransactionMessageSet {
					amount := float64(len(msg))
					tax := 0.03

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
						err = walletSign.Sign(&transaction)
						expectorNoErr(t, prefix, err)

						return transaction
					}
					var tSet entitys.TransactionCoinSet
					tSet.Tax = TransactionMake(rsa.PublicKey{}, taxBCC, "Tax").Transaction
					tSet.Transfer = TransactionMake(rsa.PublicKey{}, transfer, "Transfer").Transaction
					receiver := entitys.Client{Address: wallet2.GetPub()}
					pair := entitys.BillingInfo{To: receiver}
					standard := TransactionsStandard{ServiceName: "transaction-msg-service", BalanceServiceInstance: findBalance, WalletService: wallet1}
					transactionInfo := entitys.TransactionMsg{BillDetails: entitys.TransactionDetails{Bill: pair}, Msg: msg}
					transaction := TransactionMsg{Transaction: entitys.TransactionMsgEntityRoot{Transaction: transactionInfo}, Services: standard}
					err := transaction.Construct()
					expectorNoErr(t, prefix, err)
					err = transaction.CreateTransaction()
					expectorNoErr(t, prefix, err)
					err = walletSign.Sign(&transaction)
					expectorNoErr(t, prefix, err)
					var set entitys.TransactionMessageSet
					set.TransactionMessage = transaction.Transaction
					set.TransactionCoin = tSet
					wallet1.UnFreeze(wallet1.GetFreeze())
					return set
				}
				hash := Hasher.HashImpl{}
				err := hash.Construct()
				expectorNoErr(t, prefixOld, err)
				wallets := walletGen(3)
				zeroMock := func() {
					lage.CallPublishBlockMsg = 0
					lage.CallConsumeBlockMsg = 0
					lage.CallPublishTransactionCoinSet = 0
					lage.CallPublishBlock = 0
					lage.CallConsumeBlock = 0
					lage.ErrPublishBlockMsg = nil
					lage.ErrPublishTransactionCoin = nil
					lage.ErrPublishBlockCoin = nil
				}
				lottery := Lottery.MockLottery{}
				providers := BlockServiceProviders{WalletServiceInstance: &wallets[2], HashService: &hash, RabbitMqService: &lage, LotteryService: &lottery}
				var service BlockChainMsgImpl

				TestSucceed := func(prefixOld string) {
					fmt.Printf("%s Integration   Test For  Insert Message  Transactions Success \n", prefixOld)
					prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
					itShouldInsertOneAndNotSentTheTCoin := func(prefixOld string) {
						tr := createTransaction(&wallets[1], &wallets[2], &wallets[1], "hello world")

						defer zeroMock()

						service = BlockChainMsgImpl{Services: providers}
						err = service.Construct()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](0, len(service.Chain), t, prefixOld, "chain  size ")
						err = service.Genesis()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](1, len(service.Chain), t, prefixOld, "chain  size ")
						callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						service.InsertTransaction(tr)
						callExpector[int](1, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						callExpector[int](0, lage.CallPublishTransactionCoinSet, t, prefixOld, "rabbit mq calll publish transaction ")

						callExpector[entitys.TransactionMsg](tr.TransactionMessage.Transaction, service.Chain[0].Transactions[0], t, prefixOld, "transaction at index 0")

						fmt.Printf("%s it should  insert one  and  not  publish  coin  tr  \n", prefixOld)
					}

					itShouldInsertOneAndSentTheTCoin := func(prefixOld string) {
						tr := createTransaction(&wallets[2], &wallets[1], &wallets[2], "hello world")

						defer zeroMock()

						service = BlockChainMsgImpl{Services: providers}
						err = service.Construct()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](0, len(service.Chain), t, prefixOld, "chain  size ")
						err = service.Genesis()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](1, len(service.Chain), t, prefixOld, "chain  size ")
						callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						service.InsertTransaction(tr)
						callExpector[int](1, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						callExpector[int](1, lage.CallPublishTransactionCoinSet, t, prefixOld, "rabbit mq calll publish transaction ")

						callExpector[entitys.TransactionMsg](tr.TransactionMessage.Transaction, service.Chain[0].Transactions[0], t, prefixOld, "transaction at index 0")

						tr.TransactionCoin.Transfer.Transaction.BillDetails.Bill.To.Address = wallets[2].GetPub()
						var transactionTransfer entitys.TransactionCoinEntityRoot
						transactionTransfer.Transaction = tr.TransactionCoin.Transfer.Transaction
						transferTransactionService := TransactionCoins{Transaction: transactionTransfer}
						err = transferTransactionService.semiConstruct()
						expectorNoErr(t, prefixOld, err)
						err = wallets[2].Sign(&transferTransactionService)
						expectorNoErr(t, prefixOld, err)
						oldSign := string(tr.TransactionCoin.Transfer.Signiture)
						tr.TransactionCoin.Transfer = transferTransactionService.Transaction
						newSign := string(tr.TransactionCoin.Transfer.Signiture)
						if oldSign == newSign {
							err = errors.New("sign does  not change")
						}
						expectorNoErr(t, prefixOld, err)

						callExpector(tr.TransactionCoin.Transfer.Transaction, lage.TransactionSetCoin.Transfer.Transaction, t, prefixOld, "transaction  coin to publish ")

						callExpector(tr.TransactionCoin.Tax.Transaction, lage.TransactionSetCoin.Tax.Transaction, t, prefixOld, "transaction  coin to publish ")
						callExpector(string(tr.TransactionCoin.Transfer.Signiture), string(lage.TransactionSetCoin.Transfer.Signiture), t, prefixOld, "transaction  coin to publish ")

						callExpector(string(tr.TransactionCoin.Tax.Signiture), string(lage.TransactionSetCoin.Tax.Signiture), t, prefixOld, "transaction  coin to publish ")
						fmt.Printf("%s it should  insert one   publish  coin  tr  \n", prefixOld)
					}
					itShouldInsert6AndNotMine := func(prefixOld string) {

						tr1 := createTransaction(&wallets[1], &wallets[2], &wallets[1], "apples")
						tr2 := createTransaction(&wallets[1], &wallets[1], &wallets[1], "peach")

						tr3 := createTransaction(&wallets[0], &wallets[2], &wallets[0], "orange")
						tr4 := createTransaction(&wallets[2], &wallets[0], &wallets[2], "lemons")
						tr5 := createTransaction(&wallets[2], &wallets[0], &wallets[2], "watermellons")
						tr6 := createTransaction(&wallets[0], &wallets[2], &wallets[0], "straberries")
						defer zeroMock()
						lottery.SpinRsp = wallets[1].GetPub()

						service = BlockChainMsgImpl{Services: providers}
						err = service.Construct()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](0, len(service.Chain), t, prefixOld, "chain  size ")
						err = service.Genesis()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](1, len(service.Chain), t, prefixOld, "chain  size ")
						callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						service.InsertTransaction(tr1)
						service.InsertTransaction(tr2)
						service.InsertTransaction(tr3)
						service.InsertTransaction(tr4)
						service.InsertTransaction(tr5)
						block := entitys.BlockMessage{}
						err = block.MineBlock(lottery.SpinRsp, service.Chain[0].BlockEntity, &Logger.MockLogger{}, &hash)
						expectorNoErr(t, prefixOld, err)
						lage.BlockMsgRsp = entitys.BlockMessageMessageRabbitMq{BlockMsg: block}
						service.InsertTransaction(tr6)
						list1 := service.Chain[0].Transactions
						list2 := service.Chain[1].Transactions
						callExpector[entitys.TransactionMsg](tr1.TransactionMessage.Transaction, list1[0], t, prefixOld, "trSet1")
						callExpector[entitys.TransactionMsg](tr2.TransactionMessage.Transaction, list1[1], t, prefixOld, "trSet2")
						callExpector[entitys.TransactionMsg](tr3.TransactionMessage.Transaction, list1[2], t, prefixOld, "trSet3")
						callExpector[entitys.TransactionMsg](tr4.TransactionMessage.Transaction, list1[3], t, prefixOld, "trSet4")
						callExpector[entitys.TransactionMsg](tr5.TransactionMessage.Transaction, list1[4], t, prefixOld, "trSet5")
						callExpector[entitys.TransactionMsg](tr6.TransactionMessage.Transaction, list2[0], t, prefixOld, "trSet6")
						callExpector[int](2, lage.CallPublishTransactionCoinSet, t, prefixOld, "rabbit mq calll publish transaction ")
						callExpector[int](1, lage.CallConsumeBlockMsg, t, prefixOld, "rabbit mq calll consume  block  ")
						callExpector[int](0, lage.CallPublishBlockMsg, t, prefixOld, "rabbit mq calll publish block  ")
						fmt.Printf("%s it should  insert 6 and  not  mine  \n", prefixOld)
					}
					itShouldInsert6AndMine := func(prefixOld string) {

						tr1 := createTransaction(&wallets[1], &wallets[2], &wallets[1], "apples")
						tr2 := createTransaction(&wallets[1], &wallets[1], &wallets[1], "peach")

						tr3 := createTransaction(&wallets[0], &wallets[2], &wallets[0], "orange")
						tr4 := createTransaction(&wallets[2], &wallets[0], &wallets[2], "lemons")
						tr5 := createTransaction(&wallets[2], &wallets[0], &wallets[2], "watermellons")
						tr6 := createTransaction(&wallets[0], &wallets[2], &wallets[0], "straberries")
						defer zeroMock()
						lottery.SpinRsp = wallets[2].GetPub()

						service = BlockChainMsgImpl{Services: providers}
						err = service.Construct()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](0, len(service.Chain), t, prefixOld, "chain  size ")
						err = service.Genesis()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](1, len(service.Chain), t, prefixOld, "chain  size ")
						callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						service.InsertTransaction(tr1)
						service.InsertTransaction(tr2)
						service.InsertTransaction(tr3)
						service.InsertTransaction(tr4)
						service.InsertTransaction(tr5)
						block := entitys.BlockMessage{}
						err = block.MineBlock(lottery.SpinRsp, service.Chain[0].BlockEntity, &Logger.MockLogger{}, &hash)
						expectorNoErr(t, prefixOld, err)
						lage.BlockMsgRsp = entitys.BlockMessageMessageRabbitMq{BlockMsg: block}
						service.InsertTransaction(tr6)
						list1 := service.Chain[0].Transactions
						list2 := service.Chain[1].Transactions
						callExpector[entitys.TransactionMsg](tr1.TransactionMessage.Transaction, list1[0], t, prefixOld, "trSet1")
						callExpector[entitys.TransactionMsg](tr2.TransactionMessage.Transaction, list1[1], t, prefixOld, "trSet2")
						callExpector[entitys.TransactionMsg](tr3.TransactionMessage.Transaction, list1[2], t, prefixOld, "trSet3")
						callExpector[entitys.TransactionMsg](tr4.TransactionMessage.Transaction, list1[3], t, prefixOld, "trSet4")
						callExpector[entitys.TransactionMsg](tr5.TransactionMessage.Transaction, list1[4], t, prefixOld, "trSet5")
						callExpector[entitys.TransactionMsg](tr6.TransactionMessage.Transaction, list2[0], t, prefixOld, "trSet6")
						callExpector[int](2, lage.CallPublishTransactionCoinSet, t, prefixOld, "rabbit mq calll publish transaction ")
						callExpector[int](1, lage.CallConsumeBlockMsg, t, prefixOld, "rabbit mq calll consume  block  ")
						callExpector[int](1, lage.CallPublishBlockMsg, t, prefixOld, "rabbit mq calll publish block  ")
						callExpector(lage.BlockMsgRsp.BlockMsg.BlockEntity, lage.BlockMsg.BlockMsg.BlockEntity, t, prefixOld, "block tha publish hase been  called with")
						fmt.Printf("%s it should  insert 6 and  not  mine  \n", prefixOld)
					}
					itShouldInsertOneAndNotSentTheTCoin(prefixNew)
					itShouldInsertOneAndSentTheTCoin(prefixNew)
					itShouldInsert6AndNotMine(prefixNew)
					itShouldInsert6AndMine(prefixNew)
				}
				TestFail := func(prefixOld string) {

					fmt.Printf("%s Integration   Test For  Failures \n", prefixOld)
					prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
					itShouldFailToInsertAnInvalid := func(prefixOld string) {

						tr1 := createTransaction(&wallets[1], &wallets[2], &wallets[2], "apples")
						defer zeroMock()
						lottery.SpinRsp = wallets[1].GetPub()

						service = BlockChainMsgImpl{Services: providers}
						err = service.Construct()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](0, len(service.Chain), t, prefixOld, "chain  size ")
						err = service.Genesis()
						expectorNoErr(t, prefixOld, err)
						callExpector[int](1, len(service.Chain), t, prefixOld, "chain  size ")
						callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
						service.InsertTransaction(tr1)
						callExpector[int](0, len(service.Chain[0].Transactions), t, prefixOld, "transactions  size ")
					}
					itShouldFailToInsertAnInvalid(prefixNew)
				}
				TestSucceed(prefixNew)
				TestFail(prefixNew)
			}
			TestGenesis(prefixNew)
			TestInsertTransactions(prefixNew)
		}
		IntegrationTest(prefixNew)
	}
	TestMsgImpl(prefix)
	TestCoinsImpl(prefix)

}
