package Lottery

import (
	"Hasher"
	"Logger"
	"Stake"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"testing"
)

func InstanceCreatorLottery() (LotteryService, *Logger.MockLogger, *Hasher.MockHasher, *Stake.MockStake, error) {

	mockLogger := &Logger.MockLogger{}
	hash := &Hasher.MockHasher{}
	stake := &Stake.MockStake{}

	service := &LotteryImpl{Services: LotteryProviders{LoggerService: mockLogger, HasherService: hash, StakeService: stake}}
	err := service.Construct()
	return service, mockLogger, hash, stake, err
}

func createMeRsaPublicKeys(n int) ([]rsa.PublicKey, error) {
	var publicKeys []rsa.PublicKey

	for i := 0; i < n; i++ {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
		}

		publicKeys = append(publicKeys, privateKey.PublicKey)
	}

	return publicKeys, nil
}
func TestSpinService(t *testing.T) {
	const prefix string = "----"
	fmt.Println("*  Test For  LotteryService")
	TestLotteryImpl := func(prefixOld string) {
		fmt.Printf("%s  Test For  LotteryImpl\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestCreationService := func(prefixOld string) {
			fmt.Printf("%s  Test For  Creation  Lottery  Service\n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestCreateServiceLottery := func(prefixOld string) {
				_, _, _, _, err := InstanceCreatorLottery()
				if err != nil {
					t.Errorf("Expected to get no err  but  got %v", err)
				}
				fmt.Printf("%s it should  create  service  for lottery\n", prefixOld)
			}
			TestCreateServiceLottery(prefixNew)
		}
		TestSpins := func(prefixOld string) {
			fmt.Printf("%s  Test Spins \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestSpin3 := func(prefixOld string) {
				service, _, hash, stake, _ := InstanceCreatorLottery()
				keys, _ := createMeRsaPublicKeys(5)
				stake.Workers = keys
				stake.DistributedRoundUp = make(map[rsa.PublicKey]int)
				for _, k := range keys {
					stake.DistributedRoundUp[k] = 20
				}
				stake.Total = 100
				hash.SeedVal = 770
				key, err := service.Spin(0.0)
				if err != nil {
					t.Errorf("Expected to get no err  but  got %v", err)
				}
				if key != keys[2] {
					t.Errorf("spin 3 has  not  work corerctly")
				}
				fmt.Printf("%s it should spin correctly 3rd  test \n", prefixOld)

			}
			TestSpin2 := func(prefixOld string) {
				service, _, hash, stake, _ := InstanceCreatorLottery()
				keys, _ := createMeRsaPublicKeys(5)
				stake.Workers = keys
				stake.DistributedRoundUp = make(map[rsa.PublicKey]int)
				for _, k := range keys {
					stake.DistributedRoundUp[k] = 20
				}
				stake.Total = 100
				hash.SeedVal = 1
				key, err := service.Spin(0.0)
				if err != nil {
					t.Errorf("Expected to get no err  but  got %v", err)
				}
				if key != keys[3] {
					t.Errorf("spin 2 has  not  work corerctly")
				}
				fmt.Printf("%s it should spin correctly 2ond  test \n", prefixOld)

			}
			TestSpin1 := func(prefixOld string) {
				service, _, hash, stake, _ := InstanceCreatorLottery()
				keys, _ := createMeRsaPublicKeys(5)
				stake.Workers = keys
				stake.DistributedRoundUp = make(map[rsa.PublicKey]int)
				for _, k := range keys {
					stake.DistributedRoundUp[k] = 20
				}
				stake.Total = 100
				hash.SeedVal = 0
				key, err := service.Spin(0.0)
				if err != nil {
					t.Errorf("Expected to get no err  but  got %v", err)
				}
				if key != keys[0] {
					t.Errorf("spin 2 has  not  work corerctly")
				}
				fmt.Printf("%s it should spin correctly 1st  test \n", prefixOld)

			}
			TestFail := func(prefixOld string) {
				mockLogger := &Logger.MockLogger{}
				hash := &Hasher.MockHasher{}
				service := &LotteryImpl{Services: LotteryProviders{LoggerService: mockLogger, HasherService: hash}}
				err := service.Construct()
				if err != nil {
					t.Errorf("Expected to get no err  but  got %v", err)
				}
				_, err = service.Spin(0.0)
				exp := errors.New("Stake Service Not  Loaded ")
				if err.Error() != exp.Error() {
					t.Errorf("Expected to get %s but  got %s", exp.Error(), err.Error())
				}
				fmt.Printf("%s it should fail to  spin no 'stake-service' \n", prefixOld)
			}
			TestSpin1(prefixNew)
			TestSpin2(prefixNew)
			TestSpin3(prefixNew)
			TestFail(prefixNew)
		}
		TestLoadProvider := func(prefixOld string) {
			mockLogger := &Logger.MockLogger{}
			hash := &Hasher.MockHasher{}
			stake := &Stake.MockStake{}
			service := &LotteryImpl{Services: LotteryProviders{LoggerService: mockLogger, HasherService: hash}}
			err := service.Construct()
			if err != nil {
				t.Errorf("Expected to get no err  but  got %v", err)
			}
			service.LoadStakeService(stake)
			if service.Services.StakeService != stake {
				t.Errorf("Expected to load stake service  but failed ")
			}
			fmt.Printf("%s  it should  load  stake-service\n", prefixOld)

		}
		TestCreationService(prefixNew)
		TestSpins(prefixNew)
		TestLoadProvider(prefixNew)
	}
	TestLotteryImpl(prefix)

}
