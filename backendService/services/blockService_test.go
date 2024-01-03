package  services 
import  (
	"fmt"
    "testing"
)
func  blockCoinServiceCreator()(   blockChainService , * mockLogger , * mockWallet , * mockHasher , *blockChainCoinsImpl , error  ){
	mockLogger :=  &mockLogger{}
    wallet := &mockWallet{} 
	hash := &mockHasher{}

	service :=  &blockChainCoinsImpl{services:blockServiceProviders {loggerService:mockLogger , walletService:wallet , hashService:hash}}
	err := service.construct()
	return  service ,  mockLogger  , wallet ,hash, service ,err
}
// -- Testiing Coin Implemetation -- 
func  TestCreateSeviceBlockChain(t * testing.T){
	_,logger ,_,_ ,_,err  := blockCoinServiceCreator()
	if  err != nil {
		t.Errorf("Expected no  err  but  got %v" , err)
	}
	if  len(logger.logs) != 1{
		t.Errorf("Expected 1 msg  to Log   but  got %d" , len(logger.logs))
	}
	expectedMsg :=  "Service  Created"
	if  logger.logs[0] != expectedMsg { 
		t.Errorf("Expected msg to be %s     but  got %s" ,expectedMsg , logger.logs[0])
	}
	fmt.Println("it  should Create  block  service  from coinImpl ")
} 
func  TestCreateSeviceFailedBlockCoin(t * testing.T){
	logger :=  &mockLogger{}
 
	service :=  &blockChainCoinsImpl{services:blockServiceProviders {loggerService:logger ,}}
	err := service.construct()
	if  err == nil {
		t.Errorf("Expected   err  but  got nothing " )
	}
	if err.Error() != "The are  is no walletService" {
		t.Errorf("expected  error to be The are  is no walletService  but  got  %s" , err.Error())
	}

	fmt.Println("it  should not  Create  block  service  from coinImpl ")
} 
func  TestServiceBlockCoingenesis(t * testing.T){
	service , _ , _, hasher,_ ,_  := blockCoinServiceCreator()
	hasher.instantHashValue ="curent  hash"
	err :=  service.genesis()
	if  err != nil {
		t.Errorf("It  should  not  get  err but  got %v" ,err) 
	}
    if hasher.callInstand != 1 {
		t.Errorf("It  should call  hasher ibstand at leat 1 time  but call  %d" , hasher.callInstand) 
	} 
	if hasher.callParentOfAll != 1 {
		t.Errorf("It  should call  parentofAll  at leat 1 time  but call  %d" , hasher.callParentOfAll) 
	}
	
	fmt.Println("It should  genesis blockcoinservice")
}