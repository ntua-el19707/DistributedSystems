package  services 

import (
	"testing"
	"fmt"
	"crypto/rand"
	"crypto/rsa"
)

func  InstanceCreatorLottery() (lotteryService , * mockLogger , * mockHasher , * mockStake ,  error) {
 	
	mockLogger :=  &mockLogger{}
	hash := &mockHasher{}
	stake := &mockStake{}
	
	service := &lotteryImpl{services:lotteryProviders{logger:mockLogger ,  hasher:hash , stake:stake }}
	err :=  service.construct()
	return  service , mockLogger , hash , stake , err 
}   
func  TestCreateServiceLottery(t * testing.T) {
	_,_,_,_,err :=  InstanceCreatorLottery()
	if  err != nil{
		t.Errorf("Expected to get no err  but  got %v" ,err)
	} 
	fmt.Println("it should  create  service  for lottery")
 
} 

func  createMeRsaPublicKeys(n int  ) ([]  rsa.PublicKey ,error){
		var publicKeys [] rsa.PublicKey

	for i := 0; i < n; i++ {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
		}

		publicKeys = append(publicKeys,  privateKey.PublicKey)
	}

	return publicKeys, nil
}
func  TestSpinTest1(t * testing.T) {
	service,_,hash,stake,_ :=  InstanceCreatorLottery()
	keys ,_ := createMeRsaPublicKeys(5)
	stake.workers = keys 
	stake.distributedRoundUp  = make (map[rsa.PublicKey] int )
	for  _, k := range keys {
		stake.distributedRoundUp[k] = 20 
	}
	stake.total =  100 
	hash.seed  = 0 
	key , err := service.spin()
	if  err != nil {
			t.Errorf("Expected to get no err  but  got %v" ,err)
	} 
	if key != keys[0]{
		t.Errorf("spin 1 has  not  work corerctly" )
	}
	fmt.Println("it should  create  service  for lottery")
 
} 
func  TestSpinTest2(t * testing.T) {
	service,_,hash,stake,_ :=  InstanceCreatorLottery()
	keys ,_ := createMeRsaPublicKeys(5)
	stake.workers = keys 
	stake.distributedRoundUp  = make (map[rsa.PublicKey] int )
	for  _, k := range keys {
		stake.distributedRoundUp[k] = 20 
	}
	stake.total =  100 
	hash.seed  = 1 
	key , err := service.spin()
	if  err != nil {
			t.Errorf("Expected to get no err  but  got %v" ,err)
	} 
	if key != keys[3]{
		t.Errorf("spin 2 has  not  work corerctly" )
	}
	fmt.Println("it should  create  service  for lottery")
 
} 
//770 
func  TestSpinTest3(t * testing.T) {
	service,_,hash,stake,_ :=  InstanceCreatorLottery()
	keys ,_ := createMeRsaPublicKeys(5)
	stake.workers = keys 
	stake.distributedRoundUp  = make (map[rsa.PublicKey] int )
	for  _, k := range keys {
		stake.distributedRoundUp[k] = 20 
	}
	stake.total =  100 
	hash.seed  = 770 
	key , err := service.spin()
	if  err != nil {
			t.Errorf("Expected to get no err  but  got %v" ,err)
	} 
	if key != keys[2]{
		t.Errorf("spin 3 has  not  work corerctly")
	}
	fmt.Println("it should  create  service  for lottery")
 
} 