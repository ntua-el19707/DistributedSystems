package Hasher

import (
	"Logger"
	"fmt"
	"testing"
)

func hashServiceCreator() (HashService, *Logger.MockLogger, *HashImpl, error) {
	mockLogger := &Logger.MockLogger{}

	service := &HashImpl{loggerService: mockLogger}
	err := service.Construct()
	return service, mockLogger, service, err
}
func TestHashService(t *testing.T) {
	fmt.Println("Test  For Hash Service")
	const prefix string = "----"
	// -- Testiing Coin Implemetation --
	func(t *testing.T, prefixOld string) {
		_, _, _, err := hashServiceCreator()
		if err != nil {
			t.Errorf("Expected to get no err  but  got  %v", err)
		}
		fmt.Printf("%s it  should create a  service for hashing\n", prefixOld)
	}(t, prefix)
	func(t *testing.T, prefixOld string) {
		service, _, _, _ := hashServiceCreator()
		err, _ := service.Hash("Giannaki", "ikannaiG")

		if err != nil {
			t.Errorf("Expected  to get no  err  but gout %v", err)
		}

		fmt.Printf("%s it  should hash \n", prefixOld)

	}(t, prefix)
	validation := func(t *testing.T, prefixOld string) {
		fmt.Printf("%s Test  For  Validation\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		func(t *testing.T, prefixOld string) {
			service, _, _, _ := hashServiceCreator()

			err := service.Valid("Giannaki", "ikannaiG", "df02bb2e80b4ac7fff122a1dd43931108d1ff2344c36cd7fe4788ee04c35221d")
			if err != nil {
				t.Errorf("Expected  no  err  but  gou %v", err)
			}
			fmt.Printf("%s it  should be a valid  hash\n", prefixOld)
		}(t, prefixNew)
		func(t *testing.T, prefixOld string) {
			service, _, _, _ := hashServiceCreator()

			err := service.Valid("Giannaki", "ikannaiG", "df02bb2e80b4ac7fff122a1dd43931108d1ff2344c36cd7fe4788ee04c65221d")
			if err.Error() != "Hashes  do not match" {
				t.Errorf("Expected  err  'Hashes  do  not match' but  got %v", err)
			}
			fmt.Printf("%s it  should  no be  valid Hash \n", prefixOld)
		}(t, prefixNew)

	}
	validation(t, prefix)
	func(t *testing.T, prefixOld string) {
		service, _, _, _ := hashServiceCreator()
		parent := service.ParrentOFall()
		if parent != "1111111111111111111111111111111111111111111111111111111111111111" {
			t.Errorf("parent  should  be  1111111111111111111111111111111111111111111111111111111111111111  but  got %s", parent)
		}

		fmt.Printf("%s it  should be 11111...11 (64)\n", prefixOld)

	}(t, prefix)

	func(t *testing.T, prefixOld string) {
		service, _, _, _ := hashServiceCreator()
		seed := int64(-1000)
		hash := service.InstantHash(seed)
		if len(hash) != 64 {
			t.Errorf("it  should  reaturn a eandom string  with len 64  vut  got %s  with %d  len ", hash, len(hash))
		}

		fmt.Printf("%s it  should be a string  of  (64)\n", prefixOld)

	}(t, prefix)
	func(t *testing.T, prefixOld string) {
		service, _, _, _ := hashServiceCreator()
		err, _ := service.Hash("abc", "de")
		if err.Error() != "Abbort performing XOR abc ,  de string  does not  have  the same size " {
			t.Errorf("it  should get this err Abbort performing XOR abc ,  de string  does not  have  the same size   but  got  %s ", err.Error())
		}
		fmt.Printf("%s it should  fail to create  hash\n", prefixOld)
	}(t, prefix)

}

func TestByteOperations(t *testing.T) {
	fmt.Printf("Test Byte  Operations")
	const prefix string = "----"
	XOR := func(t *testing.T, prefixOld string) {
		fmt.Printf("%sXOR Test Cases\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

		tester := func(t *testing.T, prefixOld, str1, str2, expected string) {
			mockLogger := &Logger.MockLogger{}
			actual, _ := performXorString(str1, str2, mockLogger)
			if expected != actual {
				t.Errorf("Expected  to get %s but  got  %s\n", expected, actual)
			}
			fmt.Printf("%sit  should  get %s  when %s %s\n", prefixOld, expected, str1, str2)
		}
		Failed := func(t *testing.T, prefixOld, str1, str2 string) {
			mockLogger := &Logger.MockLogger{}
			_, err := performXorString(str1, str2, mockLogger)
			errmsg := fmt.Sprintf("Abbort performing XOR %s ,  %s string  does not  have  the same size ", str1, str2)

			if err.Error() != errmsg {
				t.Errorf("Expected  to get %s but  got %v", errmsg, err)
			}
			fmt.Printf("%sit  should fail  get err  %s  when %s %s\n", prefixOld, errmsg, str1, str2)
		}
		tester(t, prefixNew, "B", "a", "#")
		tester(t, prefixNew, "B2", "a}", "#O")
		tester(t, prefixNew, "B22", "a}s", "#OA")
		tester(t, prefixNew, "B22_", "a}s>", "#OAa")
		Failed(t, prefixNew, "B", "")
		Failed(t, prefixNew, "B", "aa")
		Failed(t, prefixNew, "df", "d")
	}
	AND := func(t *testing.T, prefixOld string) {
		fmt.Printf("%sAND Test Cases\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		tester := func(t *testing.T, prefixOld, str1, str2 string, expected int) {
			mockLogger := &Logger.MockLogger{}
			actual, _ := perfomAndAndGetSum(str1, str2, mockLogger)
			if expected != actual {
				t.Errorf("Expected  to get %d but  got  %d\n", expected, actual)
			}
			fmt.Printf("%sit  should  get %d  when %s %s\n", prefixOld, expected, str1, str2)
		}
		Failed := func(t *testing.T, prefixOld, str1, str2 string) {
			mockLogger := &Logger.MockLogger{}
			sum, err := perfomAndAndGetSum(str1, str2, mockLogger)
			errmsg := fmt.Sprintf("Abbort performing And  to get sum %s ,  %s string  does not  have  the same size ", str1, str2)

			if err.Error() != errmsg {
				t.Errorf("Expected  to get %s but  got %v", errmsg, err)
			}
			if sum != 0 {
				t.Errorf("Expected  to get 0 but  got %d", sum)

			}
			fmt.Printf("%sit  should fail  get err  %s  when %s %s\n", prefixOld, errmsg, str1, str2)
		}
		tester(t, prefixNew, "B", "a", 64+1)                // My and sum return +1
		tester(t, prefixNew, "BT", "ab", 64+64+1)           // My and sum return + 1
		tester(t, prefixNew, "BT", "ab", 64+64+1)           // My and sum return + 1
		tester(t, prefixNew, "BTE", "abE", 64+64+69+1)      // My and sum return + 1
		tester(t, prefixNew, "BTE?", "abE:", 64+64+69+58+1) // My and sum return + 1

		Failed(t, prefixNew, "B", "")
		Failed(t, prefixNew, "B", "aa")
		Failed(t, prefixNew, "df", "d")
	}
	XOR(t, prefix)
	AND(t, prefix)
}
