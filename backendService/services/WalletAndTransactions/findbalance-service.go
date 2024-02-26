package WalletAndTransactions

import (
	"Logger"
	"Service"
	"SystemInfo"
	"crypto/rsa"
	"entitys"
	"fmt"
)

type TransactionCoinRow struct {
	From          int     `json:"From"`
	To            int     `json:"To"`
	Coins         float64 `json:"Coins"`
	Reason        string  `json:"Reason"`
	Time          int64   `json:"SendTime"`
	TransactionId string  `json:"TransactionId"`
}

type TransactionListCoin []TransactionCoinRow

func (transactions *TransactionListCoin) Map(t []entitys.TransactionCoins, SystemInfoService SystemInfo.SystemInfoService) {
	*transactions = make(TransactionListCoin, len(t))
	var zero rsa.PublicKey
	for i, transaction := range t {
		billDetails := transaction.BillDetails
		from := billDetails.Bill.From.Address
		var nodeFrom entitys.ClientInfo
		if zero == from {
			nodeFrom.IndexId = -1
		} else {
			nodeFrom, _ = SystemInfoService.NodeDetails(from)
		}
		nodeTo, _ := SystemInfoService.NodeDetails(billDetails.Bill.To.Address)
		row := TransactionCoinRow{From: nodeFrom.IndexId, To: nodeTo.IndexId, Coins: transaction.Amount, Reason: transaction.Reason, Time: billDetails.Created_at, TransactionId: billDetails.Transaction_id}
		(*transactions)[i] = row
	}
}
func (list *TransactionListCoin) Sort() {
	if len(*list) < 2 {
		return
	}
	quickSort(*list, 0, len(*list)-1)
}

// Helper function for QuickSort algorithm
func quickSort(list TransactionListCoin, low, high int) {
	if low < high {
		pi := partition(list, low, high)
		quickSort(list, low, pi-1)
		quickSort(list, pi+1, high)
	}
}

// Helper function to partition the array for QuickSort
func partition(list TransactionListCoin, low, high int) int {
	pivot := list[high].Time
	i := low - 1
	for j := low; j < high; j++ {
		if list[j].Time > pivot {
			i++
			list[i], list[j] = list[j], list[i]
		}
	}
	list[i+1], list[high] = list[high], list[i+1]
	return i + 1
}

type BalanceService interface {
	Service.Service
	FindBalance(sender rsa.PublicKey) float64
	findAndLock(amount float64) (float64, error)
	GetTransactions(keys []rsa.PublicKey, times []int64) TransactionListCoin
}

type BalanceImplementation struct {
	LoggerService     Logger.LoggerService
	BlockChainService BlockChainCoinsService

	SystemInfoService SystemInfo.SystemInfoService
}

func (balance *BalanceImplementation) Construct() error {
	balance.LoggerService = &Logger.Logger{ServiceName: "balance-service"}
	err := balance.LoggerService.Construct()
	if err != nil {
		return err
	}
	balance.LoggerService.Log("service  created ")

	return nil
}
func (balance *BalanceImplementation) FindBalance(sender rsa.PublicKey) float64 {
	const lookingTemplate = "Look for  \n%v\n balance"
	lookingMessage := fmt.Sprintf(lookingTemplate, sender)
	balance.LoggerService.Log(fmt.Sprintf("Start %s", lookingMessage))
	amount := balance.BlockChainService.FindBalance(sender)
	balance.LoggerService.Log(fmt.Sprintf("Commit %s", lookingMessage))
	return amount
}
func (balance *BalanceImplementation) findAndLock(amount float64) (float64, error) {
	const lookingTemplate = "Look for my  balance And Lock"
	lookingMessage := fmt.Sprintf(lookingTemplate)
	balance.LoggerService.Log(fmt.Sprintf("Start %s", lookingMessage))
	amount, err := balance.BlockChainService.findAndLock(amount)
	balance.LoggerService.Log(fmt.Sprintf("Commit %s", lookingMessage))
	return amount, err
}
func (balance *BalanceImplementation) GetTransactions(keys []rsa.PublicKey, times []int64) TransactionListCoin {
	var list TransactionListCoin
	list.Map(balance.BlockChainService.GetTransactions(false, true, keys, times), balance.SystemInfoService)
	list.Sort()
	return list

}

// Mocks
// MockFindBalance
type MockFindBalance struct {
	Amount                 float64
	Err                    error
	FindBalanceCalledTimes int
	Locked                 bool
	LockedCall             int
	UnlockedCall           int
}

func (balance *MockFindBalance) Construct() error {
	return nil
}
func (balance *MockFindBalance) FindBalance(sender rsa.PublicKey) (float64, error) {
	balance.FindBalanceCalledTimes++
	return balance.Amount, balance.Err
}
func (balance *MockFindBalance) LockBalance() {
	balance.Locked = true
	balance.LockedCall++
}
func (balance *MockFindBalance) UnLockBalance() {
	balance.Locked = false
	balance.UnlockedCall++
}
