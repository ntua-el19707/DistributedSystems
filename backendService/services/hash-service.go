package services
import (	
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"errors"
	"time"
)
type  hashService  interface  
{
	Service
	Hash(previous ,  currentHash string ) ( error ,string)
	Valid(previous ,  currentHash ,  expected string ) error 
	Seed(currentHash string ) (int64 ,error)
	ParrentOFall() string
	instantHash(seed  int64 ) string  
}
type  hashIpmpl struct  {
	loggerService LogerService
} 

func  (service  *  hashIpmpl) construct () error   {
	if  service.loggerService == nil {
		service.loggerService = &Logger{ServiceName:"HasherService"}
		err :=  service.loggerService.construct()
		if err != nil {
			return  err
		}
	}
	service.loggerService.Log("Service  created")
	return nil 
}
func  (service hashIpmpl ) Hash(previous ,  currentHash string ) ( error ,string) {
	service.loggerService.Log("Start  hashing next  block")
	xor ,err := performXorString(previous , currentHash , service.loggerService) 
	if err != nil {
		return err ,""
	}
	times, err :=  perfomAndAndGetSum(previous , currentHash , service.loggerService) 
	if err != nil {
		return err ,""
	}
	hashed :=  hasher(xor ,times , service.loggerService)
	service.loggerService.Log("Commit  hashing next  block")
	return  nil , hashed
}
func  (service hashIpmpl )  Valid(previous ,  currentHash ,  expected string ) error {
	err ,hash  := service.Hash(previous , currentHash)
	if  err  != nil {
		return err 
	}
	if hash ==  expected {
		return  nil 
	}
	return errors.New("Hashes  do not match")
}
func  (service hashIpmpl )  Seed(hash  string ) (int64 ,error){
	var seed int64
	bytes := [] byte (hash)
	if len(bytes) >= 8 {
		seed = int64(binary.LittleEndian.Uint64(bytes[:8]))
	} else {
		return 0 , errors.New("Error: Not enough bytes to generate a seed.")
		
	} 
	return seed ,nil
}
func  (service hashIpmpl ) ParrentOFall()  (string){
	var  parent string
	for i:=0 ;  i<64 ;  i++ {
		parent += "1"
	} 
	return parent
}
func  (service hashIpmpl )  instantHash(seed  int64 ) string  {
	timeInt := time.Now().Unix() + seed
	hash :=   hasher(fmt.Sprint(timeInt)  ,  1 ,  service.loggerService)
	return  hash 
}
//-- USE - FULL

func  performXorString(str1 ,  str2 string  ,  loggerService LogerService ) ( string ,error) {
	loggerService.Log(fmt.Sprintf("Start performing XOR %s ,  %s " , str1 ,  str2))
	if  len(str1) !=  len(str2) {
		errmsg := fmt.Sprintf("Abbort performing XOR %s ,  %s string  does not  have  the same size " , str1 ,  str2)
		loggerService.Error(errmsg)
		return  "" , errors.New(errmsg)
	}
	result :=  make([]  byte  , len(str1 ))

	for  i:=0 ; i<len(str1) ; i++ {
		result[i] = str1[i] ^ str2[i]
	} 
	strResult :=  string(result)
	loggerService.Log(fmt.Sprintf("Commit performing XOR %s ,  %s " , str1 ,  str2))
	return strResult ,  nil 
}

func perfomAndAndGetSum(str1 ,  str2 string  ,  loggerService LogerService ) ( int ,error) {
	loggerService.Log(fmt.Sprintf("Start performing and  to get sum %s ,  %s " , str1 ,  str2))
	if  len(str1) !=  len(str2) {
		errmsg := fmt.Sprintf("Abbort performing And  to get sum %s ,  %s string  does not  have  the same size " , str1 ,  str2)
		loggerService.Error(errmsg)
		return   0 , errors.New(errmsg)
	}
	result :=  make([]  byte  , len(str1 ))
	for  i:=0 ; i<len(str1) ; i++ {
		result[i] = str1[i] & str2[i]
	} 
	strResult :=  string(result)
	sum := 1 // for at lest on iteration unlikly to be 1 
	for  _, charcter := range strResult {
		sum += int(charcter)
	}
	loggerService.Log(fmt.Sprintf("Commit performing and  and get sum %s ,  %s " , str1 ,  str2))
	return  sum , nil 

}
func  hasher(str string  ,  times int   ,logger LogerService) string {
	logger.Log(fmt.Sprintf("Start  loop hasshing for %d" ,  times))
	hash := sha256.New()
	hash.Write([]byte(str))
	for i:= 0  ; i <times-1 ; i++ {
		bytes :=  hash.Sum(nil)
		hash = sha256.New()
		hash.Write(bytes)
	}
	logger.Log(fmt.Sprintf("Commit  loop hasshing for %d" ,  times)) 
	return  hex.EncodeToString(hash.Sum(nil))
}

// -- Mock Hasher -- 
type mockHasher struct {
	hashvalue string 
	hashFailed  bool 
	invalid  bool
	seed int64
	seedFailed  bool 
	instantHashValue string 
	callHash int 
	callParentOfAll int 
	callValid int 
	callSeed int 
	callInstand int 


} 
func (m * mockHasher)construct() error  {
	
	return nil
}

func (m * mockHasher) Hash(previous, currentHash string) (error, string) {

	m.callHash++ 
	if  m.hashFailed {
		return  errors.New("has  faield ") , "" 
	}
	return nil, m.hashvalue 

}

func (m *  mockHasher) Valid(previous, currentHash, expected string) error {
	m.callValid++
	if  m.invalid {
		return  errors.New("invalid  block  ")
	}
	return nil
}


func (m * mockHasher) Seed(currentHash string) (int64, error) {
	m.callSeed++
	if  m.seedFailed {
		return  0 , errors.New("seed failed ") 
	}
	return m.seed , nil
}


func (m *mockHasher) ParrentOFall() string {
	m.callParentOfAll++ 
	return "1" 
}


func (m * mockHasher) instantHash(seed int64) string {
	m.callInstand++
	return m.instantHashValue
}