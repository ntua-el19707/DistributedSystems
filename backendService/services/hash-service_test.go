package services

import  (
	"fmt"
    "testing"
)
func  hashServiceCreator()(   hashService , * mockLogger  , * hashIpmpl , error  ){
	mockLogger :=  &mockLogger{}

	service :=  &hashIpmpl {loggerService:mockLogger }
	err := service.construct()
	return  service ,  mockLogger  , service ,err
}
// -- Testiing Coin Implemetation -- 
func  TestCreateSeviceHash(t * testing.T){
	_,_,_ , err :=   hashServiceCreator()
	if err != nil {
		t.Errorf("Expected to get no err  but  got  %v" ,err)
	}
	fmt.Println("created  service for hashing")
}
func  TestSuccedToCreateHash(t * testing.T){
	service , _ ,_ ,_  :=   hashServiceCreator()
	err, _ := service.Hash("Giannaki" ,  "ikannaiG")

	if  err !=  nil {
		t.Errorf("Expected  to get no  err  but gout %v" , err)
	}

	fmt.Println("It  should hash  ")
	
}
func TestValidHash(t * testing.T){
	service , _ ,_ ,_  :=   hashServiceCreator()
	
	err := service.Valid("Giannaki" ,  "ikannaiG","df02bb2e80b4ac7fff122a1dd43931108d1ff2344c36cd7fe4788ee04c35221d")
	if  err != nil {
		t.Errorf("Expected  no  err  but  gou %v" , err)
	}
	fmt.Println("it  should be a valid  hash ")
}
func TestHashParent(t * testing.T){
	service , _ ,_ ,_  :=   hashServiceCreator()
	parent :=  service.ParrentOFall()
	if  parent != "1111111111111111111111111111111111111111111111111111111111111111"{
		t.Errorf("parent  should  be  1111111111111111111111111111111111111111111111111111111111111111  but  got %s" , parent)
	}
	
	fmt.Println("it  should be 11111...11 (64)")

}

func   TestHashInstant(t * testing.T){
	service , _ ,_ ,_  :=   hashServiceCreator()
	seed := int64(-1000)
	hash := service.instantHash(seed)
	if len(hash) != 64 {
		t.Errorf("it  should  reaturn a eandom string  with len 64  vut  got %s  with %d  len " , hash , len(hash))
	}

	
	fmt.Println("it  should be a string  of  (64)" )

}
func  TestFailedToCreateHash(t * testing.T){
	service , _ ,_ ,_  :=   hashServiceCreator()
	err , _ :=  service.Hash("abc" , "de")
	if err.Error() !=  "Abbort performing XOR abc ,  de string  does not  have  the same size " {
		t.Errorf("it  should get this err Abbort performing XOR abc ,  de string  does not  have  the same size   but  got  %s " ,err.Error())
	}
	fmt.Println("it should  fail to create  hash")
}
func  TestInvalidHash(t * testing.T){
	service , _ ,_ ,_  :=   hashServiceCreator()
	
	err := service.Valid("Giannaki" ,  "ikannaiG","df02bb2e80b4ac7fff122a1dd43931108d1ff2344c36cd7fe4788ee04c65221d")
	if  err.Error() != "Hashes  do not match" {
		t.Errorf("Expected  no  err  but  gou %v" , err)
	}
	fmt.Println("it  should be  invalid  hash ")
}