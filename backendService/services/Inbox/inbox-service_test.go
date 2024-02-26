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
		fmt.Printf("%s Test For Inbox  Service Send \n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		itShouldGetTheSendInboxNoTimeParmasOneKey := func(prefixOld string) {
			logger := Logger.MockLogger{}
			sysInfo := SystemInfo.MockSystemInfoService{}
			sysInfo.NodeDetailsRsp = SystemInfo.NodeDetailsResponse{Info: entitys.ClientInfo{IndexId: 2}, Total: 1}
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

		}
		itShouldGetTheSendInboxNoTimeParmasOneKey(prefixNew)
	}
	TestForCreateService(prefix)
	TestSend(prefix)
}
