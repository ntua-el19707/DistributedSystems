package  services 
import (
	"math/rand"
	"time"
	"errors"
	"fmt"
)

type  GeneratorService interface {
	Service 
	generateId(size int ) string 
	setSeed()
}
type  generatorImplementation struct {
    ServiceName string
	CharSet string 
	seededRand * rand.Rand
	logger LogerService 
}
func ( generator *  generatorImplementation ) construct () error {
	var zero LogerService
	if generator.logger == zero {
		generator.logger =  &Logger{ServiceName:generator.ServiceName}
	}
	err :=  construct(generator.logger)
	if err != nil {
		return err
	}
	generator.setSeed()
	if  generator.CharSet == "" {
		const template = "The  CharSet  has  not been  set" 
		return  errors.New(generator.logger.Sprintf(template))
	}
	generator.logger.Log("Service Created")


 return nil 
} 
func ( generator *  generatorImplementation ) generateId(size int ) string{
	generator.logger.Log("About  to  create random id")
	id := make([] byte , size)
	for i :=  range id {
		id[i] = generator.CharSet[generator.seededRand.Intn(len(generator.CharSet))]
	}
	const template = "Created  random id:%s\n"
	generator.logger.Log(fmt.Sprintf(template ,  string(id)))
	return string(id)
} 
func ( generator *  generatorImplementation )setSeed() {
	generator.logger.Log("About to set seed")
	generator.seededRand  =  rand.New(rand.NewSource(time.Now().UnixNano()))
	generator.logger.Log("Seed has been set")
} 

// -- MOCKS -- 
type mockGenerator struct {
	response string 
	timesCallgenerateId int 
	timesCallSeed int 
}
func ( generator *  mockGenerator ) construct () error {
 return nil 
} 
func ( generator *  mockGenerator) generateId(size int ) string{
	generator.timesCallgenerateId++
	return  generator.response
} 
func ( generator *  mockGenerator ) setSeed() {
	generator.timesCallSeed++
} 
