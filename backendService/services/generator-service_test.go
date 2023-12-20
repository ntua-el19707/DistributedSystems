package  services  
import (
	"fmt"
    "testing"
	"errors"
)
func  BuildTheWorldGeneratorSevice(CharSet string)( GeneratorService ,* mockLogger , error) {
    Logger := &mockLogger{}
	var generatorService GeneratorService  =  &generatorImplementation{ServiceName:"TestingGeneratorService" ,CharSet:CharSet , logger:Logger }
	err := generatorService.construct()
	if err != nil {
		return generatorService ,Logger , errors.New(fmt.Sprintf("It should create generator  service  instead get Error:%s" ,  err.Error()))
	}
	
	
	return generatorService ,Logger,  err
}
func TestCreateGeneratorService(t  * testing.T ){
	_ , Logger ,  err :=   BuildTheWorldGeneratorSevice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if err != nil {
		t.Errorf("It should create generator  service  instead get Error:%s" ,  err.Error())
	}

	if len(Logger.logs) != 3 {
		t.Errorf("It should create generator  service  instead get Error: Logger should have  3 messages  but  got %d" ,len(Logger.logs)  )	
	}
	fmt.Println("It should create generator  service")
}
func TestCreateGeneratorServiceGenerateaaaa(t  * testing.T ){
    generatorService , Logger, _ :=   BuildTheWorldGeneratorSevice("a")
	id := generatorService.generateId(4)
	if id != "aaaa" {
		t.Errorf("Expected  is:aaaa , got  "+ id )
	}
	if len(Logger.logs) != 5  {
		t.Errorf("It should create generator  service  instead get Error: Logger should have  5 messages  but  got %d" ,len(Logger.logs)  )	
	}

	fmt.Println("It should return aaaa")
}
func TestCreateGeneratorServiceGenerate15sizeid(t  * testing.T ){
    generatorService ,_,  _ :=   BuildTheWorldGeneratorSevice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	id := generatorService.generateId(15)
	if len(id) != 15  {
		t.Errorf("Expected  len(id) = 15 , got  "+ id )
	}

	fmt.Println("It should get id with size 15")
}
func  TestCreateGeneratorServiceGenerateShouldebeinCharset(t * testing.T){
    generatorService ,_,  _ :=   BuildTheWorldGeneratorSevice("a1B#")
	id := generatorService.generateId(1000)
	for i,_:=  range id {
		if  string(id[i]) == "a"  {
		}else  if  string(id[i]) == "1"  {
		}else  if  string(id[i]) == "B"  {
		}else  if  string(id[i]) == "#"  {
		}else {
			t.Errorf("id should not have  charcted %s  " , string( id[i]))
		}

	}
	fmt.Println("It should get id within the charset")


}


func TestSetSeed(t *testing.T) {
	// Create an instance of generatorImplementation
	generator := &generatorImplementation{
		logger: &mockLogger{},
	}

	// Set the seed
	generator.setSeed()
	if generator.seededRand == nil {
		t.Error("Expected seededRand to be set, but it is nil")
	}
	mockLogger, ok := generator.logger.(*mockLogger)
	if !ok {
		t.Fatal("Failed to assert logger type")
	}
	if len(mockLogger.logs) != 2  {
		t.Errorf("Expected to  get 2  messages  but  got  %d" , len(mockLogger.logs))
	}
	expectedLogMessages := []string{"About to set seed", "Seed has been set"}

	for i, expected := range expectedLogMessages {
		if mockLogger.logs[i] != expected {
			t.Errorf("Expected log message: %s, but got: %s", expected, mockLogger.logs[i])
		}
	}
	fmt.Println("It should set Seed")


}