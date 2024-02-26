package Inbox

import (
	"Logger"
	"Service"
	"SystemInfo"
	"WalletAndTransactions"
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
)

type InboxRecord struct {
	From          int    `json:"From"`
	To            int    `json:"To"`
	Msg           string `json:"Msg"`
	Nonce         int    `json:"Nonce"`
	Time          int64  `json:"SendTime"`
	TransactionId string `json:"TransactionId"`
}
type Inbox []InboxRecord

// --  Map =>  TrnsactionMsg =>  InboxRecord
func (inbox *Inbox) Map(t []entitys.TransactionMsg, SystemInfoService SystemInfo.SystemInfoService) {
	*inbox = make([]InboxRecord, len(t))
	for i, transaction := range t {
		billDetails := transaction.BillDetails
		nodeFrom, _ := SystemInfoService.NodeDetails(billDetails.Bill.From.Address)
		nodeTo, _ := SystemInfoService.NodeDetails(billDetails.Bill.To.Address)
		row := InboxRecord{Nonce: transaction.BillDetails.Nonce, From: nodeFrom.IndexId, To: nodeTo.IndexId, Msg: transaction.Msg, Time: billDetails.Created_at, TransactionId: billDetails.Transaction_id}
		(*inbox)[i] = row
	}

}
func (inbox *Inbox) Sort() {
	if len(*inbox) < 2 {
		return
	}
	quickSort(*inbox, 0, len(*inbox)-1)
}

// Helper function for QuickSort algorithm
func quickSort(inbox []InboxRecord, low, high int) {
	if low < high {
		pi := partition(inbox, low, high)
		quickSort(inbox, low, pi-1)
		quickSort(inbox, pi+1, high)
	}
}

// Helper function to partition the array for QuickSort
func partition(inbox []InboxRecord, low, high int) int {
	pivot := inbox[high].Time
	i := low - 1
	for j := low; j < high; j++ {
		if inbox[j].Time > pivot {
			i++
			inbox[i], inbox[j] = inbox[j], inbox[i]
		}
	}
	inbox[i+1], inbox[high] = inbox[high], inbox[i+1]
	return i + 1
}

type InboxProviders struct {
	LoggerService     Logger.LoggerService
	BlockChainService WalletAndTransactions.BlockChainMsgService
	SystemInfoService SystemInfo.SystemInfoService
}

func (p *InboxProviders) Construct() error {
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: "inbox-service"}
		err := p.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	return p.valid()

}

const errNoProvided = "There is  no provider for '%s'"

func (p *InboxProviders) valid() error {
	if p.BlockChainService == nil {
		errMsg := fmt.Sprintf(errNoProvided, "Block Chain Service")
		return errors.New(errMsg)
	}
	if p.SystemInfoService == nil {
		errMsg := fmt.Sprintf(errNoProvided, "System Info Service")
		return errors.New(errMsg)
	}
	return nil
}

type InboxService interface {
	Service.Service
	Send(keys []rsa.PublicKey, times []int64) (error, Inbox)
	All(times []int64) (error, Inbox)
	Receive(keys []rsa.PublicKey, times []int64) (error, Inbox)
	SendAndReceived(keys []rsa.PublicKey, times []int64) (error, Inbox)
}

type InboxImpl struct {
	Providers InboxProviders
}

func (service *InboxImpl) Construct() error {
	return service.Providers.Construct()
}
func (service *InboxImpl) All(times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	keys := make([]rsa.PublicKey, 0)
	transactions := service.Providers.BlockChainService.GetTransactions(true, false, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	inbox.Sort()
	return nil, inbox
}
func (service *InboxImpl) Send(keys []rsa.PublicKey, times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	if len(keys) == 0 {
		return errors.New("need  at leat one  key"), nil
	}
	transactions := service.Providers.BlockChainService.GetTransactions(true, false, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	inbox.Sort()
	return nil, inbox
}
func (service *InboxImpl) Receive(keys []rsa.PublicKey, times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	if len(keys) == 0 {
		return errors.New("need  at leat one  key"), nil
	}
	transactions := service.Providers.BlockChainService.GetTransactions(false, false, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	inbox.Sort()
	return nil, inbox
}
func (service *InboxImpl) SendAndReceived(keys []rsa.PublicKey, times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	if len(keys) == 0 {
		return errors.New("need  at leat one  key"), nil
	}
	transactions := service.Providers.BlockChainService.GetTransactions(false, true, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	inbox.Sort()
	return nil, inbox
}
