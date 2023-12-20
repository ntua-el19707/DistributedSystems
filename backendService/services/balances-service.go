package services ;

import (
	"crypto/rsa"
	"fmt"
)


 
type  BalanceService  interface {
	Service
	findBalance(sender rsa.PublicKey ) ( float64 ,  error)
}

type  balanceImplementation  struct {
	
	
}
func (balance * balanceImplementation ) construct() error {
	LogerServiceBalance.Log("service  created ")
	
	return nil 
}
func (balance * balanceImplementation )	findBalance(sender rsa.PublicKey ) ( float64 ,  error) {
	const lookingTemplate = "Look for  \n%v\n balance"
	lookingMessage := fmt.Sprintf(lookingTemplate  ,sender)  
	LogerServiceBalance.Log(lookingMessage)
	return 0 , nil
}
// Mocks
//mockFindBalance
type mockFindBalance struct {
	amount float64 
	err  error 
	findBalanceCalledTimes int

} 
func (balance *  mockFindBalance ) construct() error {	
	return nil 
}
func (balance * mockFindBalance )	findBalance(sender rsa.PublicKey ) ( float64 ,  error) {
	balance.findBalanceCalledTimes++
	return balance.amount  , balance.err
}