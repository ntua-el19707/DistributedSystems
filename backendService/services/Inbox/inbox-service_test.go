package Inbox

import (
	"Logger"
	"SystemInfo"
	"WalletAndTransactions"
	"crypto/rand"
	"crypto/rsa"
	"entitys"
	"fmt"
	rand2 "math/rand"
	"sort"
	"testing"
	"time"
)

type listTrMsg []entitys.TransactionMsg

func (l listTrMsg) Len() int {
	return len(l)
}
func (l listTrMsg) Less(i, j int) bool {
	return l[i].BillDetails.Created_at > l[j].BillDetails.Created_at
}
func (l listTrMsg) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
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
func stringGenerator(size int) string {
	const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, size)
	for i, _ := range id {
		id[i] = allChars[rand2.Intn(len(allChars))]
	}
	return string(id)
}
func createMsgTransactions(n, between int, keys []rsa.PublicKey) []entitys.TransactionMsg {
	list := make([]entitys.TransactionMsg, n)
	for i := 0; i < n; i++ {
		keyFrom := rand2.Intn(len(keys))
		var keyTo int
		for {
			keyTo = rand2.Intn(len(keys))
			if keyTo != keyFrom {
				break
			}
		}
		message := stringGenerator(rand2.Intn(10) + 5) // 5<=x<15
		timeShift := int64(rand2.Intn(between))        //to create transaction random in 1 hour period
		var transaction entitys.TransactionMsg
		transaction.Msg = message
		var billingInfo entitys.BillingInfo
		billingInfo.From.Address = keys[keyFrom]
		billingInfo.To.Address = keys[keyTo]
		var details entitys.TransactionDetails
		details.Bill = billingInfo
		details.Created_at = time.Now().Unix() + timeShift
		details.Transaction_id = stringGenerator(8)
		details.Nonce = i
		transaction.BillDetails = details
		list[i] = transaction
	}
	return list
}
func createPublicKey(n int, t *testing.T) []rsa.PublicKey {
	var publicKeys []rsa.PublicKey

	for i := 0; i < n; i++ {
		privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			t.Errorf("failed to generate RSA key pair: %v", err)
		}
		publicKeys = append(publicKeys, privateKey.PublicKey)
	}
	return publicKeys
}
func TestMapper(t *testing.T) {
	fmt.Printf("* Test Mapper  \n")
	const prefix = "----"
	keys := createPublicKey(3, t)
	itShouldMap0 := func(prefixOld string) {
		var inbox Inbox
		var transactions []entitys.TransactionMsg
		sysInfo := &SystemInfo.MockSystemInfoService{}
		inbox.Map(transactions, sysInfo)
		callExpector[int](0, len(inbox), t, prefixOld, "size  of  inbox")
		fmt.Printf("%s it should  Map 0 transactions\n", prefixOld)

	}
	itShouldMap10 := func(prefixOld string) {
		var inbox Inbox
		size := 10
		indexId := 2
		list := createMsgTransactions(size, 60*60, keys)
		sysInfo := &SystemInfo.MockSystemInfoService{}
		sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
		inbox.Map(list, sysInfo)
		callExpector[int](size, len(inbox), t, prefixOld, "size  of  inbox")
		callExpector[int](size*2, sysInfo.CallNodeDetails, t, prefixOld, "call node details ")

		for i := 0; i < size; i++ {

			callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
			callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
			callExpector(list[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
			callExpector(list[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
			callExpector(list[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
			callExpector(list[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
		}

		fmt.Printf("%s it should  Map 10 transactions\n", prefixOld)

	}

	itShouldMap0(prefix)
	itShouldMap10(prefix)
}
func TestInbox(t *testing.T) {
	fmt.Printf("* Test For Inbox  Service\n")
	const prefix string = "----"
	keys := createPublicKey(3, t)
	list := createMsgTransactions(10, 60*60, keys)

	TestForCreateService := func(prefixOld string) {
		fmt.Printf("%s Test For Creating Inbox  Service\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		itShouldCreate := func(prefixOld string) {
			logger := Logger.MockLogger{}
			sysInfo := SystemInfo.MockSystemInfoService{}
			blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
			providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
			service := InboxImpl{Providers: providers}
			err := service.Construct()
			expectorNoErr(t, prefixOld, err)

			fmt.Printf("%s it should  create service\n", prefixOld)
		}
		TestFailToCreate := func(prefixOld string) {
			fmt.Printf("%s Test For Fail to create Inbox  Service\n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldNotCreateNoBlockChainService := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				errMsg := fmt.Sprintf(errNoProvided, "Block Chain Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				fmt.Printf("%s it should  fail  to create service 'no block chain service'\n", prefixOld)

			}
			itShouldNotCreateNoSystemInfoService := func(prefixOld string) {
				logger := Logger.MockLogger{}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				providers := InboxProviders{LoggerService: &logger, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				errMsg := fmt.Sprintf(errNoProvided, "System Info Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				fmt.Printf("%s it should  fail  to create service 'no system info service'\n", prefixOld)

			}
			itShouldNotCreateNoBlockChainService(prefixNew)
			itShouldNotCreateNoSystemInfoService(prefixNew)
		}
		itShouldCreate(prefixNew)
		TestFailToCreate(prefixNew)

	}
	TestSend := func(prefixOld string) {
		tested := "Send"
		fmt.Printf("%s Test For Inbox  Service  %s \n", prefixOld, tested)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestSucceed := func(prefixOld string) {
			fmt.Printf("%s Test For Success Inbox  Service %s \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldGetInboxNoTimeParmasOneKey := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 2
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				err, inbox := service.Send(keyList, []int64{})
				expectorNoErr(t, prefixOld, err)

				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")

				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}

				fmt.Printf("%s it should get the %s  messages  from 1 key no time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParmsTwoKey := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				err, inbox := service.Send(keyList, []int64{})

				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be two")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[rsa.PublicKey](keyList[1], blockChainMsg.CallGetTransactionsWith[0].Keys[1], t, prefixOld, "key  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key to key 2 no time  params \n", prefixOld, tested)
			}
			itShouldGetInbox1TimeParms1Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				times := []int64{time.Now().Unix()}
				err, inbox := service.Send(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the send  messages  from 1 key ,1  time  params \n", prefixOld)
			}
			itShouldGetInbox2TimeParms1Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 3
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Send(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInbox2TimeParms2Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Send(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be two")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[rsa.PublicKey](keyList[1], blockChainMsg.CallGetTransactionsWith[0].Keys[1], t, prefixOld, "key  at index 1 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from key 1 to key 2 ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParmasOneKey(prefixNew)
			itShouldGetInboxNoTimeParmsTwoKey(prefixNew)
			itShouldGetInbox1TimeParms1Key(prefixNew)
			itShouldGetInbox2TimeParms1Key(prefixNew)
			itShouldGetInbox2TimeParms2Key(prefixNew)
		}
		TestFail := func(prefixOld string) {

			fmt.Printf("%s Test For fail Inbox  Service %s  \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldFailInvalidServiceNoSynfo := func(prefixOld string) {
				logger := Logger.MockLogger{}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Send(keyList, times)
				errMsg := fmt.Sprintf(errNoProvided, "System Info Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no sysInfo'\n", prefixOld, tested)
			}
			itShouldFailInvalidServiceNoBlockChain := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo}
				service := InboxImpl{Providers: providers}
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Send(keyList, times)
				errMsg := fmt.Sprintf(errNoProvided, "Block Chain Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no block chain'\n", prefixOld, tested)

			}
			itShouldFailNoKeys := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Send(keyList, times)
				callExpector[string](errAtleast1Rsa, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to not give rsa  Keys \n", prefixOld, tested)
			}
			itShouldFailMorethan2keys := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1], keys[2]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Send(keyList, times)
				callExpector[string](errMax2Rsa, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than 2 rsa keys \n", prefixOld, tested)
			}
			itShouldFailMorethan2times := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24)), time.Now().Unix()}
				err, inbox := service.Send(keyList, times)
				callExpector[string](errMax2Times, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than  2 times \n", prefixOld, tested)

			}

			itShouldFailInvalidServiceNoSynfo(prefixNew)
			itShouldFailInvalidServiceNoBlockChain(prefixNew)
			itShouldFailNoKeys(prefixNew)
			itShouldFailMorethan2keys(prefixNew)
			itShouldFailMorethan2times(prefixNew)

		}
		TestSucceed(prefixNew)
		TestFail(prefixNew)
	}
	TestAll := func(prefixOld string) {
		tested := "All"
		fmt.Printf("%s Test For Inbox  Service  %s \n", prefixOld, tested)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestSucceed := func(prefixOld string) {
			fmt.Printf("%s Test For Success Inbox  Service %s \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldGetInboxNoTimeParms := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				err, inbox := service.All([]int64{})

				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be zero")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages no time  params \n", prefixOld, tested)
			}
			itShouldGetInbox1TimeParams := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				times := []int64{time.Now().Unix()}
				err, inbox := service.All(times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be 0")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages 1  time  params \n", prefixOld, tested)
			}
			itShouldGetInbox2TimeParams := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 3
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.All(times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =true")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParms(prefixNew)
			itShouldGetInbox1TimeParams(prefixNew)
			itShouldGetInbox2TimeParams(prefixNew)
		}
		TestFail := func(prefixOld string) {

			fmt.Printf("%s Test For fail Inbox  Service %s  \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldFailInvalidServiceNoSynfo := func(prefixOld string) {
				logger := Logger.MockLogger{}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.All(times)
				errMsg := fmt.Sprintf(errNoProvided, "System Info Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no sysInfo'\n", prefixOld, tested)
			}
			itShouldFailInvalidServiceNoBlockChain := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo}
				service := InboxImpl{Providers: providers}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.All(times)
				errMsg := fmt.Sprintf(errNoProvided, "Block Chain Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no block chain'\n", prefixOld, tested)

			}
			itShouldFailMorethan2times := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24)), time.Now().Unix()}
				err, inbox := service.All(times)
				callExpector[string](errMax2Times, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than  2 times \n", prefixOld, tested)

			}

			itShouldFailInvalidServiceNoSynfo(prefixNew)
			itShouldFailInvalidServiceNoBlockChain(prefixNew)
			itShouldFailMorethan2times(prefixNew)
		}
		TestSucceed(prefixNew)
		TestFail(prefixNew)
	}
	TestRecieve := func(prefixOld string) {
		tested := "Receive"
		fmt.Printf("%s Test For Inbox  Service  %s \n", prefixOld, tested)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestSucceed := func(prefixOld string) {
			fmt.Printf("%s Test For Success Inbox  Service %s \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldGetInboxNoTimeParmasOneKey := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 2
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				err, inbox := service.Receive(keyList, []int64{})
				expectorNoErr(t, prefixOld, err)

				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")

				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}

				fmt.Printf("%s it should get the %s  messages  from 1 key no time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParmsTwoKey := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				err, inbox := service.Receive(keyList, []int64{})

				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be two")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[rsa.PublicKey](keyList[1], blockChainMsg.CallGetTransactionsWith[0].Keys[1], t, prefixOld, "key  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key to key 2 no time  params \n", prefixOld, tested)
			}
			itShouldGetInbox1TimeParms1Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				times := []int64{time.Now().Unix()}
				err, inbox := service.Receive(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key ,1  time  params \n", prefixOld, tested)
			}
			itShouldGetInbox2TimeParms1Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 3
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Receive(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInbox2TimeParms2Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Receive(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =false")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be two")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[rsa.PublicKey](keyList[1], blockChainMsg.CallGetTransactionsWith[0].Keys[1], t, prefixOld, "key  at index 1 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from key 1 to key 2 ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParmasOneKey(prefixNew)
			itShouldGetInboxNoTimeParmsTwoKey(prefixNew)
			itShouldGetInbox1TimeParms1Key(prefixNew)
			itShouldGetInbox2TimeParms1Key(prefixNew)
			itShouldGetInbox2TimeParms2Key(prefixNew)
		}
		TestFail := func(prefixOld string) {

			fmt.Printf("%s Test For fail Inbox  Service %s  \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldFailInvalidServiceNoSynfo := func(prefixOld string) {
				logger := Logger.MockLogger{}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Receive(keyList, times)
				errMsg := fmt.Sprintf(errNoProvided, "System Info Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no sysInfo'\n", prefixOld, tested)
			}
			itShouldFailInvalidServiceNoBlockChain := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo}
				service := InboxImpl{Providers: providers}
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Receive(keyList, times)
				errMsg := fmt.Sprintf(errNoProvided, "Block Chain Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no block chain'\n", prefixOld, tested)

			}
			itShouldFailNoKeys := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Receive(keyList, times)
				callExpector[string](errAtleast1Rsa, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to not give rsa  Keys \n", prefixOld, tested)
			}
			itShouldFailMorethan2keys := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1], keys[2]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.Receive(keyList, times)
				callExpector[string](errMax2Rsa, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than 2 rsa keys \n", prefixOld, tested)
			}
			itShouldFailMorethan2times := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24)), time.Now().Unix()}
				err, inbox := service.Receive(keyList, times)
				callExpector[string](errMax2Times, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than  2 times \n", prefixOld, tested)

			}

			itShouldFailInvalidServiceNoSynfo(prefixNew)
			itShouldFailInvalidServiceNoBlockChain(prefixNew)
			itShouldFailNoKeys(prefixNew)
			itShouldFailMorethan2keys(prefixNew)
			itShouldFailMorethan2times(prefixNew)

		}
		TestSucceed(prefixNew)
		TestFail(prefixNew)
	}
	TestRecieveAndSend := func(prefixOld string) {
		tested := "Receive"
		fmt.Printf("%s Test For Inbox  Service  %s \n", prefixOld, tested)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestSucceed := func(prefixOld string) {
			fmt.Printf("%s Test For Success Inbox  Service %s \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldGetInboxNoTimeParmasOneKey := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 2
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				err, inbox := service.SendAndReceived(keyList, []int64{})
				expectorNoErr(t, prefixOld, err)

				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =true")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")

				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}

				fmt.Printf("%s it should get the %s  messages  from 1 key no time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParmsTwoKey := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				err, inbox := service.SendAndReceived(keyList, []int64{})

				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =true")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be two")
				callExpector[int](0, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 0")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[rsa.PublicKey](keyList[1], blockChainMsg.CallGetTransactionsWith[0].Keys[1], t, prefixOld, "key  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key to key 2 no time  params \n", prefixOld, tested)
			}
			itShouldGetInbox1TimeParms1Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				times := []int64{time.Now().Unix()}
				err, inbox := service.SendAndReceived(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =true")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the send  %s from 1 key ,1  time  params \n", prefixOld, tested)
			}
			itShouldGetInbox2TimeParms1Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 3
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.SendAndReceived(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =true")
				callExpector[int](1, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be one")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from 1 key ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInbox2TimeParms2Key := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.SendAndReceived(keyList, times)
				expectorNoErr(t, prefixOld, err)
				callExpector[int](len(list), len(inbox), t, prefixOld, "inbox size")
				callExpector[int](1, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions")
				callExpector[bool](false, blockChainMsg.CallGetTransactionsWith[0].From, t, prefixOld, "from mmust be =false")
				callExpector[bool](true, blockChainMsg.CallGetTransactionsWith[0].TwoWay, t, prefixOld, "TwoWay mmust be =true")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Keys), t, prefixOld, "Keys mmust be two")
				callExpector[int](2, len(blockChainMsg.CallGetTransactionsWith[0].Times), t, prefixOld, "times mmust be 1")
				callExpector[rsa.PublicKey](keyList[0], blockChainMsg.CallGetTransactionsWith[0].Keys[0], t, prefixOld, "key  at index 0 ")
				callExpector[rsa.PublicKey](keyList[1], blockChainMsg.CallGetTransactionsWith[0].Keys[1], t, prefixOld, "key  at index 1 ")
				callExpector[int64](times[0], blockChainMsg.CallGetTransactionsWith[0].Times[0], t, prefixOld, "times  at index 0 ")
				callExpector[int64](times[1], blockChainMsg.CallGetTransactionsWith[0].Times[1], t, prefixOld, "times  at index 1 ")
				copied := make([]entitys.TransactionMsg, len(list))
				copy(copied, list)
				sortedList := listTrMsg(copied)
				sort.Sort(sortedList)
				for i := 0; i < len(inbox); i++ {
					callExpector(indexId, inbox[i].From, t, prefixOld, fmt.Sprintf("From  at index %d", i))
					callExpector(indexId, inbox[i].To, t, prefixOld, fmt.Sprintf("To  at index %d", i))
					callExpector(sortedList[i].BillDetails.Transaction_id, inbox[i].TransactionId, t, prefixOld, fmt.Sprintf("Transaction Id at index %d", i))
					callExpector(sortedList[i].BillDetails.Nonce, inbox[i].Nonce, t, prefixOld, fmt.Sprintf("Nonce  at index %d", i))
					callExpector(sortedList[i].Msg, inbox[i].Msg, t, prefixOld, fmt.Sprintf("Transaction Msg at index %d", i))
					callExpector(sortedList[i].BillDetails.Created_at, inbox[i].Time, t, prefixOld, fmt.Sprintf("Time  at index %d", i))
				}
				fmt.Printf("%s it should get the %s  messages  from key 1 to key 2 ,2 time  params \n", prefixOld, tested)
			}
			itShouldGetInboxNoTimeParmasOneKey(prefixNew)
			itShouldGetInboxNoTimeParmsTwoKey(prefixNew)
			itShouldGetInbox1TimeParms1Key(prefixNew)
			itShouldGetInbox2TimeParms1Key(prefixNew)
			itShouldGetInbox2TimeParms2Key(prefixNew)
		}
		TestFail := func(prefixOld string) {

			fmt.Printf("%s Test For fail Inbox  Service %s  \n", prefixOld, tested)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			itShouldFailInvalidServiceNoSynfo := func(prefixOld string) {
				logger := Logger.MockLogger{}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.SendAndReceived(keyList, times)
				errMsg := fmt.Sprintf(errNoProvided, "System Info Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no sysInfo'\n", prefixOld, tested)
			}
			itShouldFailInvalidServiceNoBlockChain := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo}
				service := InboxImpl{Providers: providers}
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.SendAndReceived(keyList, times)
				errMsg := fmt.Sprintf(errNoProvided, "Block Chain Service")
				callExpector[string](errMsg, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to 'no block chain'\n", prefixOld, tested)

			}
			itShouldFailNoKeys := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.SendAndReceived(keyList, times)
				callExpector[string](errAtleast1Rsa, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to not give rsa  Keys \n", prefixOld, tested)
			}
			itShouldFailMorethan2keys := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1], keys[2]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24))}
				err, inbox := service.SendAndReceived(keyList, times)
				callExpector[string](errMax2Rsa, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than 2 rsa keys \n", prefixOld, tested)
			}
			itShouldFailMorethan2times := func(prefixOld string) {
				logger := Logger.MockLogger{}
				sysInfo := SystemInfo.MockSystemInfoService{}
				indexId := 1
				sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: indexId}, Total: 1}
				blockChainMsg := &WalletAndTransactions.MockBlockChainMsg{}
				blockChainMsg.Transactions = list
				providers := InboxProviders{LoggerService: &logger, SystemInfoService: &sysInfo, BlockChainService: blockChainMsg}
				service := InboxImpl{Providers: providers}
				err := service.Construct()
				expectorNoErr(t, prefixOld, err)
				keyList := []rsa.PublicKey{keys[0], keys[1]}
				times := []int64{time.Now().Unix(), time.Now().Unix() + int64(rand2.Intn(60*60*24)), time.Now().Unix()}
				err, inbox := service.SendAndReceived(keyList, times)
				callExpector[string](errMax2Times, err.Error(), t, prefixOld, "error")
				callExpector[int](0, len(inbox), t, prefixOld, "size of inbox")
				callExpector[int](0, blockChainMsg.CallGetTransactions, t, prefixOld, "call get transactions ")
				callExpector[int](0, sysInfo.CallNodeDetails, t, prefixOld, "size of inbox")
				fmt.Printf("%s it should fail  get the %s  messages due  to give more than  2 times \n", prefixOld, tested)

			}

			itShouldFailInvalidServiceNoSynfo(prefixNew)
			itShouldFailInvalidServiceNoBlockChain(prefixNew)
			itShouldFailNoKeys(prefixNew)
			itShouldFailMorethan2keys(prefixNew)
			itShouldFailMorethan2times(prefixNew)

		}
		TestSucceed(prefixNew)
		TestFail(prefixNew)
	}
	TestForCreateService(prefix)
	TestSend(prefix)
	TestAll(prefix)
	TestRecieve(prefix)
	TestRecieveAndSend(prefix)
}
