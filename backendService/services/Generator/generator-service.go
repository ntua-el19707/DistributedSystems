package Generator

import (
	"Logger"
	"Service"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type GeneratorService interface {
	Service.Service
	GenerateId(size int) string
	SetSeed()
}
type GeneratorImplementation struct {
	ServiceName string
	CharSet     string
	seededRand  *rand.Rand
	Logger      Logger.LoggerService
}

func (generator *GeneratorImplementation) Construct() error {
	var zero Logger.LoggerService
	if generator.Logger == zero {
		generator.Logger = &Logger.Logger{ServiceName: generator.ServiceName}
		err := generator.Logger.Construct()
		if err != nil {
			return err
		}
	}
	generator.SetSeed()
	if generator.CharSet == "" {
		const template = "The  CharSet  has  not been  set"
		generator.Logger.Error(template)
		return errors.New(generator.Logger.Sprintf(template))
	}
	generator.Logger.Log("Service Created")

	return nil
}
func (generator *GeneratorImplementation) GenerateId(size int) string {
	generator.Logger.Log("About  to  create random id")
	id := make([]byte, size)
	for i := range id {
		id[i] = generator.CharSet[generator.seededRand.Intn(len(generator.CharSet))]
	}
	const template = "Created  random id:%s\n"
	generator.Logger.Log(fmt.Sprintf(template, string(id)))
	return string(id)
}
func (generator *GeneratorImplementation) SetSeed() {
	generator.Logger.Log("About to set seed")
	generator.seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	generator.Logger.Log("Seed has been set")
}

// -- MOCKS --
type MockGenerator struct {
	Response            string
	TimesCallGenerateId int
	TimesCallSeed       int
}

func (generator *MockGenerator) Construct() error {
	return nil
}
func (generator *MockGenerator) GenerateId(size int) string {
	generator.TimesCallGenerateId++
	return generator.Response
}
func (generator *MockGenerator) SetSeed() {
	generator.TimesCallSeed++
}
