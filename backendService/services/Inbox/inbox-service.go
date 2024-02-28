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
	"sort"
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

func (inbox Inbox) Len() int {
	return len(inbox)
}
func (inbox Inbox) Less(i, j int) bool {
	return inbox[i].Time > inbox[j].Time
}
func (inbox Inbox) Swap(i, j int) {
	inbox[i], inbox[j] = inbox[j], inbox[i]
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
const errAtleast1Rsa = "Need at least 1 rsa.PublicKey"
const errMax2Rsa = "Max 2 rsa.PublicKey"
const errMax2Times = "Max 2 times  'unix' - int64"

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
	if len(times) > 2 {
		return errors.New(errMax2Times), nil
	}
	keys := make([]rsa.PublicKey, 0)
	transactions := service.Providers.BlockChainService.GetTransactions(true, false, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	sort.Sort(inbox)
	return nil, inbox
}
func (service *InboxImpl) Send(keys []rsa.PublicKey, times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	if len(keys) == 0 {
		return errors.New(errAtleast1Rsa), nil
	}
	if len(keys) > 2 {
		return errors.New(errMax2Rsa), nil
	}
	if len(times) > 2 {
		return errors.New(errMax2Times), nil
	}
	transactions := service.Providers.BlockChainService.GetTransactions(true, false, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	sort.Sort(inbox)
	return nil, inbox
}
func (service *InboxImpl) Receive(keys []rsa.PublicKey, times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	if len(keys) == 0 {
		return errors.New(errAtleast1Rsa), nil
	}
	if len(keys) > 2 {
		return errors.New(errMax2Rsa), nil
	}
	if len(times) > 2 {
		return errors.New(errMax2Times), nil
	}
	transactions := service.Providers.BlockChainService.GetTransactions(false, false, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	sort.Sort(inbox)
	return nil, inbox
}
func (service *InboxImpl) SendAndReceived(keys []rsa.PublicKey, times []int64) (error, Inbox) {
	err := service.Providers.valid()
	if err != nil {
		return err, nil
	}
	if len(keys) == 0 {
		return errors.New(errAtleast1Rsa), nil
	}
	if len(keys) > 2 {
		return errors.New(errMax2Rsa), nil
	}
	if len(times) > 2 {
		return errors.New(errMax2Times), nil
	}
	transactions := service.Providers.BlockChainService.GetTransactions(false, true, keys, times)
	var inbox Inbox
	inbox.Map(transactions, service.Providers.SystemInfoService)
	sort.Sort(inbox)
	return nil, inbox
}
