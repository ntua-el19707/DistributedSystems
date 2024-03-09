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
	var zeroTransaction entitys.TransactionMsg
	for i, transaction := range t {
		if transaction == zeroTransaction {
			break
		}
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

type BlockDto struct {
	Index       int    `json:"index"`
	CreatedAt   int64  `json:"created_at"`
	Validator   int    `json:"validator"`
	Capicity    int    `json:"capacity"`
	CurrentHash string `json:"current_hash"`
	ParentHash  string `json:"parrent_hash"`
}

func (b *BlockDto) Map(block entitys.Block, SystemInfoService SystemInfo.SystemInfoService) {
	nodeValidator, _ := SystemInfoService.NodeDetails(block.Validator)
	b.Index = block.Index
	b.CreatedAt = block.CreatedAt
	b.Validator = nodeValidator.IndexId
	b.Capicity = block.Capicity
	b.CurrentHash = block.CurrentHash
	b.ParentHash = block.ParentHash

}

type BlockMsgDto struct {
	Block        BlockDto `json:"block"`
	Transactions Inbox    `json:"transactions"`
}

func (b *BlockMsgDto) Map(block entitys.BlockMessage, SystemInfoService SystemInfo.SystemInfoService) {
	b.Block.Map(block.BlockEntity, SystemInfoService)
	var transactions Inbox
	transactions.Map(block.Transactions, SystemInfoService)
	b.Transactions = transactions

}

type ChainMsgDTO []BlockMsgDto

func (c *ChainMsgDTO) Map(chain entitys.BlockChainMessage, SystemInfoService SystemInfo.SystemInfoService) {
	size := len(chain)
	*c = make([]BlockMsgDto, size)
	for i := 0; i < size; i++ {
		(*c)[i].Map(chain[i], SystemInfoService)
	}
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
	GetBlockChain() (ChainMsgDTO, error)
}

type InboxImpl struct {
	Providers InboxProviders
}

func (service *InboxImpl) Construct() error {
	return service.Providers.Construct()
}
func (service *InboxImpl) GetBlockChain() (ChainMsgDTO, error) {
	err := service.Providers.valid()
	if err != nil {
		return nil, err
	}
	var chain ChainMsgDTO
	chain.Map(service.Providers.BlockChainService.RetriveChain(), service.Providers.SystemInfoService)
	return chain, nil
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
