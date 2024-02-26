package asyncLoad

import (
	"Logger"
	"RabbitMqService"
	"WalletAndTransactions"
	"crypto/rand"
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
	rand2 "math/rand"
	"testing"
	"time"
)

func callExpector[T comparable](obj1, obj2 T, t *testing.T, prefix, what string) {
	if obj1 != obj2 {
		t.Errorf("%s  Expected  '%s' to get %v but %v ", prefix, what, obj1, obj2)
	}
}
func expectorNoErr(t *testing.T, prefix string, err error) {
	if err != nil {
		t.Errorf("%s  Expected  no err  but got %v err ", prefix, err)
	}
}
func createTransaction(from, to rsa.PublicKey, time int64, amount float64, reason string) entitys.TransactionCoinEntityRoot {
	var transaction entitys.TransactionCoinEntityRoot
	var billingInfo entitys.BillingInfo
	billingInfo.From.Address = from
	billingInfo.To.Address = to
	var details entitys.TransactionDetails
	details.Bill = billingInfo
	details.Created_at = time
	details.Transaction_id = stringGenerator(8)
	transaction.Transaction.BillDetails = details
	transaction.Transaction.Amount = amount
	transaction.Transaction.Reason = reason
	transaction.Signiture = []byte("a signiture")
	return transaction
}
func createTransactionCoinSet(from, to, validator rsa.PublicKey, amount float64, time int64) entitys.TransactionCoinSet {
	var transaction entitys.TransactionCoinSet
	transaction.Tax = createTransaction(from, validator, time, 0.03*amount, "Tax")
	transaction.Transfer = createTransaction(from, to, time, 0.97*amount, "Transfer")
	return transaction
}
func createPublicKey(n int, t *testing.T) []rsa.PublicKey {
	var publicKeys []rsa.PublicKey

	for i := 0; i < n; i++ {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Errorf("failed to generate RSA key pair: %v", err)
		}
		publicKeys = append(publicKeys, privateKey.PublicKey)
	}
	return publicKeys
}
func stringGenerator(size int) string {
	const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, size)
	for i, _ := range id {
		id[i] = allChars[rand2.Intn(len(allChars))]
	}
	return string(id)
}
func createMsgTransactions(n, between int, keys []rsa.PublicKey) []entitys.TransactionMessageSet {
	list := make([]entitys.TransactionMessageSet, n)
	for i := 0; i < n; i++ {
		keyFrom := rand2.Intn(len(keys))
		var keyTo int
		for {
			keyTo = rand2.Intn(len(keys))
			if keyTo != keyFrom {
				break
			}
		}
		key1 := rand2.Intn(len(keys))
		key3 := rand2.Intn(len(keys))
		message := stringGenerator(rand2.Intn(10) + 5) // 5<=x<15

		timeShift := int64(rand2.Intn(between)) //to create transaction random in 1 hour period
		coins := float64(len(message))          // 5.0 <x<10.0
		coinstr := createTransactionCoinSet(keys[keyFrom], keys[key1], keys[key3], coins, time.Now().Unix()+timeShift)
		var transaction entitys.TransactionMsgEntityRoot
		transaction.Signiture = []byte("a signiture")
		transaction.Transaction.Msg = message
		var billingInfo entitys.BillingInfo
		billingInfo.From.Address = keys[keyFrom]
		billingInfo.To.Address = keys[keyTo]
		var details entitys.TransactionDetails
		details.Bill = billingInfo
		details.Created_at = time.Now().Unix() + timeShift
		details.Transaction_id = stringGenerator(8)
		transaction.Transaction.BillDetails = details
		var tr entitys.TransactionMessageSet
		tr.TransactionCoin = coinstr
		tr.TransactionMessage = transaction
		list[i] = tr
	}
	return list
}
func TestAsyncLoad(t *testing.T) {
	const prefix = "----"
	fmt.Println("* Test Cases  For Async Load Service")
	keys := createPublicKey(6, t)
	listTrsanctions := []entitys.TransactionCoinSet{
		createTransactionCoinSet(keys[0], keys[1], keys[2], 7.0, time.Now().Unix()),
		createTransactionCoinSet(keys[1], keys[2], keys[2], 5.0, time.Now().Unix()+1000),
		createTransactionCoinSet(keys[0], keys[2], keys[1], 4.0, time.Now().Unix()+2000),
	}
	createService := func(prefixOld string) {
		fmt.Printf("%s Test Cases  For Creating Async Load Service\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		itShouldCreate := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)

			fmt.Printf("%s it should  create  \n", prefixOld)

		}
		itShouldFailNoRabbit := func(prefixOld string) {
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			errMsg := ErrNoRabbitMqProviders
			callExpector[string](errMsg, err.Error(), t, prefixOld, "error")

			fmt.Printf("%s it should fail to create 'no rabbitmqservice' \n", prefixOld)

		}
		itShouldFailNoBlockCoinService := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			errMsg := ErrNoBlockCoinService
			callExpector[string](errMsg, err.Error(), t, prefixOld, "error")

			fmt.Printf("%s it should fail to create 'no BlockChainCoinService' \n", prefixOld)

		}
		itShouldFailNoBlockMsgService := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			errMsg := ErrNoMsgService
			callExpector[string](errMsg, err.Error(), t, prefixOld, "error")

			fmt.Printf("%s it should fail to create 'no BlockChainMsgService' \n", prefixOld)

		}
		itShouldCreate(prefixNew)
		itShouldFailNoRabbit(prefixNew)
		itShouldFailNoBlockCoinService(prefixNew)
		itShouldFailNoBlockMsgService(prefixNew)
	}
	TestForConsumingTransaxtionCoinsSet := func(prefixOld string) {
		fmt.Printf("%s Test Cases  For Consuming Trasaction Coin set\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		itShouldConsumeOne := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)

			bunny.Channel = make(chan entitys.TransactionCoinSet)
			go service.consumeTransactionsCoins()
			time.Sleep(10 * time.Millisecond)
			bunny.Channel <- listTrsanctions[0]
			time.Sleep(10 * time.Millisecond)
			callExpector[int](1, blockChainCoins.CallInsertTransactions, t, prefixOld, "call insert transaction ")
			callExpector[int](1, bunny.CallConsumerTransactionsCoins, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetChannelTransactionCoin, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")
			callExpector[entitys.TransactionCoins](listTrsanctions[0].Tax.Transaction, blockChainCoins.CallInsertTransactionsWith[0].Tax.Transaction, t, prefixOld, "call with tranction set 1 Tax")
			callExpector[entitys.TransactionCoins](listTrsanctions[0].Transfer.Transaction, blockChainCoins.CallInsertTransactionsWith[0].Transfer.Transaction, t, prefixOld, "call with tranction set 1 Transfer")
			callExpector[string](string(listTrsanctions[0].Tax.Signiture), string(blockChainCoins.CallInsertTransactionsWith[0].Tax.Signiture), t, prefixOld, "call with tranction set 1 Tax singiture ")
			callExpector[string](string(listTrsanctions[0].Transfer.Signiture), string(blockChainCoins.CallInsertTransactionsWith[0].Transfer.Signiture), t, prefixOld, "call with tranction set 1 Transfer")
			fmt.Printf("%s it should consume  one  transaction set \n", prefixOld)

		}
		itShouldConsumeThree := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)

			bunny.Channel = make(chan entitys.TransactionCoinSet)
			go service.consumeTransactionsCoins()
			time.Sleep(10 * time.Millisecond)
			bunny.Channel <- listTrsanctions[0]
			bunny.Channel <- listTrsanctions[1]
			bunny.Channel <- listTrsanctions[2]
			time.Sleep(10 * time.Millisecond)
			callExpector[int](3, blockChainCoins.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsCoins, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetChannelTransactionCoin, t, prefixOld, "call rabbitmq get channel ")

			for i := 0; i < 3; i++ {
				callExpector[entitys.TransactionCoins](listTrsanctions[i].Tax.Transaction, blockChainCoins.CallInsertTransactionsWith[i].Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](listTrsanctions[i].Transfer.Transaction, blockChainCoins.CallInsertTransactionsWith[i].Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(listTrsanctions[i].Tax.Signiture), string(blockChainCoins.CallInsertTransactionsWith[i].Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(listTrsanctions[i].Transfer.Signiture), string(blockChainCoins.CallInsertTransactionsWith[i].Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
			}
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")
			fmt.Printf("%s it should consume  3 transaction set \n", prefixOld)

		}
		itShouldConsumeThreeFailInsertOne := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)

			bunny.Channel = make(chan entitys.TransactionCoinSet)
			go service.consumeTransactionsCoins()
			time.Sleep(10 * time.Millisecond)
			bunny.Channel <- listTrsanctions[0]
			time.Sleep(10 * time.Millisecond)
			errorInsert := errors.New("failed to insert transaction")
			blockChainCoins.ErrorInsertTransaction = errorInsert
			bunny.Channel <- listTrsanctions[1]
			time.Sleep(10 * time.Millisecond)

			blockChainCoins.ErrorInsertTransaction = nil
			bunny.Channel <- listTrsanctions[2]
			time.Sleep(10 * time.Millisecond)
			callExpector[int](3, blockChainCoins.CallInsertTransactions, t, prefixOld, "call insert transaction ")
			callExpector[int](1, len(logger.ErrorList), t, prefixOld, "error Messages")

			callExpector[int](1, bunny.CallConsumerTransactionsCoins, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetChannelTransactionCoin, t, prefixOld, "call rabbitmq get channel ")
			callExpector[string](errorInsert.Error(), logger.ErrorList[0], t, prefixOld, "error Message at index 0 ")

			for i := 0; i < 3; i++ {
				callExpector[entitys.TransactionCoins](listTrsanctions[i].Tax.Transaction, blockChainCoins.CallInsertTransactionsWith[i].Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](listTrsanctions[i].Transfer.Transaction, blockChainCoins.CallInsertTransactionsWith[i].Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(listTrsanctions[i].Tax.Signiture), string(blockChainCoins.CallInsertTransactionsWith[i].Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(listTrsanctions[i].Transfer.Signiture), string(blockChainCoins.CallInsertTransactionsWith[i].Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
			}
			fmt.Printf("%s it should consume  3 transaction set failed to insert  one  \n", prefixOld)

		}
		itShouldConsume100 := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			list := make([]entitys.TransactionCoinSet, 100)
			for i := 0; i < 100; i++ {
				key1 := rand2.Intn(len(keys))
				var key2 int
				for {
					key2 = rand2.Intn(len(keys))
					if key1 != key2 {
						break
					}
				}
				key3 := rand2.Intn(len(keys))
				timeShift := int64(rand2.Intn(60 * 60)) //to create transaction random in 1 hour period
				coins := (rand2.Float64() * 10) + 5     // 5.0 <x<10.0
				list[i] = createTransactionCoinSet(keys[key1], keys[key2], keys[key3], coins, time.Now().Unix()+timeShift)
			}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)

			bunny.Channel = make(chan entitys.TransactionCoinSet)
			go service.consumeTransactionsCoins()
			time.Sleep(10 * time.Millisecond)
			for i := 0; i < 100; i++ {
				bunny.Channel <- list[i]
			}
			time.Sleep(20 * time.Millisecond)
			callExpector[int](100, blockChainCoins.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsCoins, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetChannelTransactionCoin, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")

			for i := 0; i < 100; i++ {
				callExpector[entitys.TransactionCoins](list[i].Tax.Transaction, blockChainCoins.CallInsertTransactionsWith[i].Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](list[i].Transfer.Transaction, blockChainCoins.CallInsertTransactionsWith[i].Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(list[i].Tax.Signiture), string(blockChainCoins.CallInsertTransactionsWith[i].Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(list[i].Transfer.Signiture), string(blockChainCoins.CallInsertTransactionsWith[i].Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
			}
			fmt.Printf("%s it should consume  100 transaction set \n", prefixOld)

		}
		itShouldConsumeOne(prefixNew)
		itShouldConsumeThree(prefixNew)
		itShouldConsumeThreeFailInsertOne(prefixNew)
		itShouldConsume100(prefixNew)
	}
	TestForConsumingTransactionMsg := func(prefixOld string) {
		fmt.Printf("%s Test Cases  For Consuming Transaction Msg \n", prefixOld)

		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		list := createMsgTransactions(1000, 60*60, keys)
		itShouldConsumeOne := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)
			bunny.ChannelTransactionMsg = make(chan entitys.TransactionMessageSet)
			go service.consumeTransactionsMsg()
			time.Sleep(10 * time.Millisecond)
			bunny.ChannelTransactionMsg <- list[0]
			time.Sleep(10 * time.Millisecond)
			callExpector[int](1, blockChainMsg.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsMsg, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetTransactionMsgChannel, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")
			for i := 0; i < 1; i++ {
				TransactionCoins := list[i].TransactionCoin
				TrMsg := list[i].TransactionMessage
				ActualTrasactionCoin := blockChainMsg.CallInsertTransactionsWith[i].TransactionCoin
				ActualTranctionMsg := blockChainMsg.CallInsertTransactionsWith[i].TransactionMessage
				callExpector[entitys.TransactionMsg](TrMsg.Transaction, ActualTranctionMsg.Transaction, t, prefixOld, fmt.Sprintf("call with transaction set %d Tr message ", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Tax.Transaction, ActualTrasactionCoin.Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Transfer.Transaction, ActualTrasactionCoin.Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(TransactionCoins.Tax.Signiture), string(ActualTrasactionCoin.Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(TransactionCoins.Transfer.Signiture), string(ActualTrasactionCoin.Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))

				callExpector[string](string(TrMsg.Signiture), string(ActualTranctionMsg.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d transaction msg singiture ", i))
			}
			fmt.Printf("%s it should consume  1 transaction set \n", prefixOld)

		}
		itShouldConsumeTen := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)
			bunny.ChannelTransactionMsg = make(chan entitys.TransactionMessageSet)
			go service.consumeTransactionsMsg()
			time.Sleep(10 * time.Millisecond)
			for i := 0; i < 10; i++ {
				bunny.ChannelTransactionMsg <- list[i]
			}
			time.Sleep(10 * time.Millisecond)
			callExpector[int](10, blockChainMsg.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsMsg, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetTransactionMsgChannel, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")
			for i := 0; i < 10; i++ {
				TransactionCoins := list[i].TransactionCoin
				TrMsg := list[i].TransactionMessage
				ActualTrasactionCoin := blockChainMsg.CallInsertTransactionsWith[i].TransactionCoin
				ActualTranctionMsg := blockChainMsg.CallInsertTransactionsWith[i].TransactionMessage
				callExpector[entitys.TransactionMsg](TrMsg.Transaction, ActualTranctionMsg.Transaction, t, prefixOld, fmt.Sprintf("call with transaction set %d Tr message ", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Tax.Transaction, ActualTrasactionCoin.Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Transfer.Transaction, ActualTrasactionCoin.Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(TransactionCoins.Tax.Signiture), string(ActualTrasactionCoin.Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(TransactionCoins.Transfer.Signiture), string(ActualTrasactionCoin.Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))

				callExpector[string](string(TrMsg.Signiture), string(ActualTranctionMsg.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d transaction msg singiture ", i))
			}
			fmt.Printf("%s it should consume  10 transaction set \n", prefixOld)

		}
		itShouldConsume100 := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)
			bunny.ChannelTransactionMsg = make(chan entitys.TransactionMessageSet)
			go service.consumeTransactionsMsg()
			time.Sleep(10 * time.Millisecond)
			for i := 0; i < 100; i++ {
				bunny.ChannelTransactionMsg <- list[i]
			}
			time.Sleep(100 * time.Millisecond)
			callExpector[int](100, blockChainMsg.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsMsg, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetTransactionMsgChannel, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")
			for i := 0; i < 100; i++ {
				TransactionCoins := list[i].TransactionCoin
				TrMsg := list[i].TransactionMessage
				ActualTrasactionCoin := blockChainMsg.CallInsertTransactionsWith[i].TransactionCoin
				ActualTranctionMsg := blockChainMsg.CallInsertTransactionsWith[i].TransactionMessage
				callExpector[entitys.TransactionMsg](TrMsg.Transaction, ActualTranctionMsg.Transaction, t, prefixOld, fmt.Sprintf("call with transaction set %d Tr message ", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Tax.Transaction, ActualTrasactionCoin.Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Transfer.Transaction, ActualTrasactionCoin.Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(TransactionCoins.Tax.Signiture), string(ActualTrasactionCoin.Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(TransactionCoins.Transfer.Signiture), string(ActualTrasactionCoin.Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))

				callExpector[string](string(TrMsg.Signiture), string(ActualTranctionMsg.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d transaction msg singiture ", i))
			}
			fmt.Printf("%s it should consume  100 transaction set \n", prefixOld)
		}
		itShouldConsume1000 := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)
			bunny.ChannelTransactionMsg = make(chan entitys.TransactionMessageSet)
			go service.consumeTransactionsMsg()
			time.Sleep(10 * time.Millisecond)
			for i := 0; i < 1000; i++ {
				bunny.ChannelTransactionMsg <- list[i]
			}
			time.Sleep(100 * time.Millisecond)
			callExpector[int](1000, blockChainMsg.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsMsg, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetTransactionMsgChannel, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](0, len(logger.ErrorList), t, prefixOld, "error Messages")
			for i := 0; i < 1000; i++ {
				TransactionCoins := list[i].TransactionCoin
				TrMsg := list[i].TransactionMessage
				ActualTrasactionCoin := blockChainMsg.CallInsertTransactionsWith[i].TransactionCoin
				ActualTranctionMsg := blockChainMsg.CallInsertTransactionsWith[i].TransactionMessage
				callExpector[entitys.TransactionMsg](TrMsg.Transaction, ActualTranctionMsg.Transaction, t, prefixOld, fmt.Sprintf("call with transaction set %d Tr message ", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Tax.Transaction, ActualTrasactionCoin.Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Transfer.Transaction, ActualTrasactionCoin.Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(TransactionCoins.Tax.Signiture), string(ActualTrasactionCoin.Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(TransactionCoins.Transfer.Signiture), string(ActualTrasactionCoin.Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))

				callExpector[string](string(TrMsg.Signiture), string(ActualTranctionMsg.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d transaction msg singiture ", i))
			}
			fmt.Printf("%s it should consume  1000 transaction set \n", prefixOld)
		}
		itShouldConsume3failOne := func(prefixOld string) {
			bunny := &RabbitMqService.MockRabbitMqImpl{}
			blockChainCoins := &WalletAndTransactions.MockBlockChainCoins{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			logger := &Logger.MockLogger{}
			providers := AsyncLoadProviders{LoggerService: logger, RabbitMqService: bunny, BlockCoinService: blockChainCoins, BlockMsgService: blockChainMsg}
			service := AsyncLoadImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)
			bunny.ChannelTransactionMsg = make(chan entitys.TransactionMessageSet)
			go service.consumeTransactionsMsg()
			time.Sleep(10 * time.Millisecond)
			bunny.ChannelTransactionMsg <- list[0]
			time.Sleep(10 * time.Millisecond)
			errorInsert := errors.New("failed to insert transaction")
			blockChainMsg.ErrorInsertTransaction = errorInsert
			bunny.ChannelTransactionMsg <- list[1]
			time.Sleep(10 * time.Millisecond)
			blockChainMsg.ErrorInsertTransaction = nil
			bunny.ChannelTransactionMsg <- list[2]

			time.Sleep(10 * time.Millisecond)
			callExpector[int](3, blockChainMsg.CallInsertTransactions, t, prefixOld, "call insert transaction ")

			callExpector[int](1, bunny.CallConsumerTransactionsMsg, t, prefixOld, "call rabbitmq cobnsumer ")
			callExpector[int](1, bunny.CallGetTransactionMsgChannel, t, prefixOld, "call rabbitmq get channel ")
			callExpector[int](1, len(logger.ErrorList), t, prefixOld, "error Messages")
			callExpector[string](errorInsert.Error(), logger.ErrorList[0], t, prefixOld, "error Message at index 0 ")
			for i := 0; i < 3; i++ {
				TransactionCoins := list[i].TransactionCoin
				TrMsg := list[i].TransactionMessage
				ActualTrasactionCoin := blockChainMsg.CallInsertTransactionsWith[i].TransactionCoin
				ActualTranctionMsg := blockChainMsg.CallInsertTransactionsWith[i].TransactionMessage
				callExpector[entitys.TransactionMsg](TrMsg.Transaction, ActualTranctionMsg.Transaction, t, prefixOld, fmt.Sprintf("call with transaction set %d Tr message ", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Tax.Transaction, ActualTrasactionCoin.Tax.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Tax", i))
				callExpector[entitys.TransactionCoins](TransactionCoins.Transfer.Transaction, ActualTrasactionCoin.Transfer.Transaction, t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))
				callExpector[string](string(TransactionCoins.Tax.Signiture), string(ActualTrasactionCoin.Tax.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Tax singiture ", i))
				callExpector[string](string(TransactionCoins.Transfer.Signiture), string(ActualTrasactionCoin.Transfer.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d Transfer", i))

				callExpector[string](string(TrMsg.Signiture), string(ActualTranctionMsg.Signiture), t, prefixOld, fmt.Sprintf("call with tranction set %d transaction msg singiture ", i))
			}
			fmt.Printf("%s it should consume  3 transaction set fail one  \n", prefixOld)
		}
		itShouldConsumeOne(prefixNew)
		itShouldConsume3failOne(prefixNew)
		itShouldConsumeTen(prefixNew)
		itShouldConsume100(prefixNew)
		itShouldConsume1000(prefixNew)
	}
	createService(prefix)
	TestForConsumingTransaxtionCoinsSet(prefix)
	TestForConsumingTransactionMsg(prefix)

}
