package Generator

import (
	"Logger"
	"errors"
	"fmt"
	"testing"
)

func BuildTheWorldGeneratorSevice(CharSet string) (GeneratorService, *Logger.MockLogger, error) {
	Logger := &Logger.MockLogger{}
	var generatorService GeneratorService = &GeneratorImplementation{ServiceName: "TestingGeneratorService", CharSet: CharSet, Logger: Logger}
	err := generatorService.Construct()
	if err != nil {
		return generatorService, Logger, errors.New(fmt.Sprintf("It should create generator  service  instead get Error:%s", err.Error()))
	}

	return generatorService, Logger, err
}

func TestGenerator(t *testing.T) {
	const prefix string = "----"
	fmt.Println("* Test  For  generator-service")
	TestCreationService := func(prefixOld string) {
		fmt.Printf("%s  Test For  construct\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestCreateGeneratorService := func(prefixOld string) {
			_, Logger, err := BuildTheWorldGeneratorSevice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
			if err != nil {
				t.Errorf("It should create generator  service  instead get Error:%s", err.Error())
			}

			if len(Logger.Logs) != 3 {
				t.Errorf("It should create generator  service  instead get Error: Logger should have  3 messages  but  got %d", len(Logger.Logs))
			}
			fmt.Printf("%s it should create generator  service\n", prefixOld)
		}
		TestCreateGeneratorService(prefixNew)

	}
	TestGenerateId := func(prefixOld string) {
		fmt.Printf("%s  Test For  GenerateId\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestCreateGeneratorServiceGenerateaaaa := func(prefixOld string) {
			generatorService, Logger, _ := BuildTheWorldGeneratorSevice("a")
			id := generatorService.GenerateId(4)
			if id != "aaaa" {
				t.Errorf("Expected  id:aaaa , got %s ", id)
			}
			if len(Logger.Logs) != 5 {
				t.Errorf("It should create generator  service  instead get Error: Logger should have  5 messages  but  got %d", len(Logger.Logs))
			}

			fmt.Printf("%s it should return aaaa\n", prefixOld)
		}
		TestCreateGeneratorServiceGenerate15sizeid := func(prefixOld string) {
			generatorService, _, _ := BuildTheWorldGeneratorSevice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
			id := generatorService.GenerateId(15)
			if len(id) != 15 {
				t.Errorf("Expected  len(id) = 15 , got  " + id)
			}

			fmt.Printf("%s it should get id with size 15\n", prefixOld)
		}
		TestCreateGeneratorServiceGenerateShouldebeinCharset := func(prefixOld string) {
			generatorService, _, _ := BuildTheWorldGeneratorSevice("a1B#")
			id := generatorService.GenerateId(1000)
			for i, _ := range id {
				if string(id[i]) == "a" {
				} else if string(id[i]) == "1" {
				} else if string(id[i]) == "B" {
				} else if string(id[i]) == "#" {
				} else {
					t.Errorf("id should not have  charcted %s  ", string(id[i]))
				}

			}
			fmt.Printf("%s it should get id within the charset\n", prefixOld)

		}
		TestCreateGeneratorServiceGenerateaaaa(prefixNew)
		TestCreateGeneratorServiceGenerate15sizeid(prefixNew)
		TestCreateGeneratorServiceGenerateShouldebeinCharset(prefixNew)

	}
	TestSeed := func(prefixOld string) {
		fmt.Printf("%s  Test For  setSeed\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestSetSeed := func(prefixOld string) {
			// Create an instance of generatorImplementation
			generator := &GeneratorImplementation{
				Logger: &Logger.MockLogger{},
			}

			// Set the seed
			generator.SetSeed()
			if generator.seededRand == nil {
				t.Error("Expected seededRand to be set, but it is nil")
			}
			mockLogger, ok := generator.Logger.(*Logger.MockLogger)
			if !ok {
				t.Fatal("Failed to assert logger type")
			}
			if len(mockLogger.Logs) != 2 {
				t.Errorf("Expected to  get 2  messages  but  got  %d", len(mockLogger.Logs))
			}
			expectedLogMessages := []string{"About to set seed", "Seed has been set"}

			for i, expected := range expectedLogMessages {
				if mockLogger.Logs[i] != expected {
					t.Errorf("Expected log message: %s, but got: %s", expected, mockLogger.Logs[i])
				}
			}
			fmt.Printf("%s it should set Seed\n", prefixOld)

		}
		TestSetSeed(prefixNew)

	}
	TestCreationService(prefix)
	TestGenerateId(prefix)
	TestSeed(prefix)
}
