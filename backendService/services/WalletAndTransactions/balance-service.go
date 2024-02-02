package WalletAndTransactions

import (
	"Logger"
	"Service"
	"crypto/rsa"
	"fmt"
	"sync"
)

type BalanceService interface {
	Service.Service
	FindBalance(sender rsa.PublicKey) (float64, error)
	LockBalance()
	UnLockBalance()
}

type BalanceImplementation struct {
	mu                    sync.Mutex
	WalletServiceInstance WalletService
	LoggerService         Logger.LoggerService
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
func (balance *BalanceImplementation) FindBalance(sender rsa.PublicKey) (float64, error) {
	const lookingTemplate = "Look for  \n%v\n balance"
	lookingMessage := fmt.Sprintf(lookingTemplate, sender)
	balance.LoggerService.Log(lookingMessage)
	//	b, err := BlockChainCoinsService.FindBalance()
	return 0, nil
}
func (balance *BalanceImplementation) LockBalance() {
	balance.mu.Lock()
}
func (balance *BalanceImplementation) UnLockBalance() {
	balance.mu.Unlock()
}

// Mocks
// mockFindBalance
type mockFindBalance struct {
	amount                 float64
	err                    error
	findBalanceCalledTimes int
	locked                 bool
	lockedCall             int
	unlockedCall           int
}

func (balance *mockFindBalance) Construct() error {
	return nil
}
func (balance *mockFindBalance) FindBalance(sender rsa.PublicKey) (float64, error) {
	balance.findBalanceCalledTimes++
	return balance.amount, balance.err
}
func (balance *mockFindBalance) LockBalance() {
	balance.locked = true
	balance.lockedCall++
}
func (balance *mockFindBalance) UnLockBalance() {
	balance.locked = false
	balance.unlockedCall++
}
