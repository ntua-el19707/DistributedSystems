package FindBalance

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
	mu            sync.Mutex
	LoggerService Logger.LoggerService
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
