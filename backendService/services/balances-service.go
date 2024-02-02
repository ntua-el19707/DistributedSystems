package services

import (
	"crypto/rsa"
	"fmt"
	"sync"
)

type BalanceService interface {
	Service
	findBalance(sender rsa.PublicKey) (float64, error)
	LockBalance()
	UnLockBalance()
}

type balanceImplementation struct {
	mu            sync.Mutex
	walletService WalletService
	loggerService LogerService
}

func (balance *balanceImplementation) construct() error {
	balance.loggerService = &Logger{ServiceName: balanceServiceName}
	err := balance.loggerService.construct()
	if err != nil {
		return err
	}
	balance.loggerService.Log("service  created ")

	return nil
}
func (balance *balanceImplementation) findBalance(sender rsa.PublicKey) (float64, error) {
	const lookingTemplate = "Look for  \n%v\n balance"
	lookingMessage := fmt.Sprintf(lookingTemplate, sender)
	balance.loggerService.Log(lookingMessage)
	b, err := BlockChainCoinsService.FindBalance()
	return b, err
}
func (balance *balanceImplementation) LockBalance() {
	balance.mu.Lock()
}
func (balance *balanceImplementation) UnLockBalance() {
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

func (balance *mockFindBalance) construct() error {
	return nil
}
func (balance *mockFindBalance) findBalance(sender rsa.PublicKey) (float64, error) {
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
