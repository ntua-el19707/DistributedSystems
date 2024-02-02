package Stake

import (
	"Logger"
	"crypto/rand"
	"crypto/rsa"
	"entitys"
	"fmt"
	"testing"
)

/*
*

	create  stake  service
*/
func createStackService() (StakeService, *Logger.MockLogger, *StakeCoinBlockChain, error) {
	mockLogger := &Logger.MockLogger{}
	impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
	err := impl.Construct()
	return impl, mockLogger, impl, err
}

func TestServiceStake(t *testing.T) {
	const prefix string = "----"
	fmt.Println("* Test  For Stake  Service")
	keyGen := func(n int) ([]rsa.PublicKey, error) {
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
	TestForStakeCoins := func(prefixOld string) {
		fmt.Printf("%s  Test  For  Stake  Coin Impl\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		TestServiceConstruct := func(prefixOld string) {
			fmt.Printf("%s  Test  For  Service Creation \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestCreateStakeService := func(prefixOld string) {
				_, logger, _, err := createStackService()
				if err != nil {
					t.Errorf("Expcected  to not  get err  but  got %v", err)
				}
				const expected string = "Service  created"
				if len(logger.Logs) != 1 {
					t.Errorf("Expcected  to  get 1  log message   but  got %d", len(logger.Logs))
				}
				if logger.Logs[0] != expected {
					t.Errorf("Expcected  to not  get msg %s  but  got %s", expected, logger.Logs[0])
				}
				fmt.Printf("%s it  should  CreateStakeService\n", prefixOld)

			}
			TestFailCreateStakeService := func(prefixOld string) {

				mockLogger := &Logger.MockLogger{}
				impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
				impl.Block.BlockEntity.Capicity = 1
				err := impl.Construct()
				const expected string = "Block  is  not  full cappicity  is  1 but  has  0 "

				if err.Error() != expected {
					t.Errorf("Expcected  to not  get  err : %s  but  got %s", expected, err.Error())
				}
				fmt.Printf("%s it  should  failed   CreateStakeService not  full block\n", prefixOld)
			}
			TestCreateStakeService(prefixNew)
			TestFailCreateStakeService(prefixNew)
		}
		TestDistibutionPrivate := func(prefixOld string) {
			fmt.Printf("%s  Test  For  Service distibutionMap  \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			TestDistributionScale0 := func(prefixOld string) {
				mockLogger := &Logger.MockLogger{}
				impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
				impl.Workers, _ = keyGen(2)
				impl.Block.BlockEntity.Capicity = 5

				billDetails1 := entitys.TransactionDetails{}
				billDetails2 := entitys.TransactionDetails{}
				billDetails1.Bill.From.Address = impl.Workers[0]
				billDetails1.Bill.To.Address = impl.Workers[1]

				billDetails2.Bill.To.Address = impl.Workers[0]
				billDetails2.Bill.From.Address = impl.Workers[1]

				transactions := []entitys.TransactionCoins{entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails2, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails2, Amount: 20.0}}
				impl.Block.Transactions = transactions
				err := impl.Construct()
				if err != nil {
					t.Errorf("Expected no err  but  got %v", err)
				}
				dMap, total := impl.distributionOfStake(0)
				if dMap[impl.Workers[0]] != 60.0 || dMap[impl.Workers[1]] != 40.0 || total != 100.0 {
					t.Errorf("Expected worker0:%3f worker1:%3f , total %3f but  got  worker0:%3f worker1:%3f , total %3f ", 60.0, 40.0, 100.0, dMap[impl.Workers[0]], dMap[impl.Workers[1]], total)
				}
				fmt.Printf("%s it  should  create distibution of scale  0   60 ,  40 distibutionMap\n", prefixOld)
			}
			TestDistributionScale1AndHalf := func(prefixOld string) {
				mockLogger := &Logger.MockLogger{}
				impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
				impl.Workers, _ = keyGen(2)
				impl.Block.BlockEntity.Capicity = 5

				billDetails1 := entitys.TransactionDetails{}
				billDetails2 := entitys.TransactionDetails{}
				billDetails1.Bill.From.Address = impl.Workers[0]
				billDetails1.Bill.To.Address = impl.Workers[1]

				billDetails2.Bill.To.Address = impl.Workers[0]
				billDetails2.Bill.From.Address = impl.Workers[1]

				transactions := []entitys.TransactionCoins{entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails2, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails2, Amount: 20.0}}
				impl.Block.Transactions = transactions
				err := impl.Construct()
				if err != nil {
					t.Errorf("Expected no err  but  got %v", err)
				}
				dMap, total := impl.distributionOfStake(1.5)
				if dMap[impl.Workers[0]] != 120.0 || dMap[impl.Workers[1]] != 130.0 || total != 250.0 {
					t.Errorf("Expected worker0:%3f worker1:%3f , total %3f but  got  worker0:%3f worker1:%3f , total %3f ", 120.0, 130.0, 250.0, dMap[impl.Workers[0]], dMap[impl.Workers[1]], total)
				}
				fmt.Printf("%s it should  create distibution 1.5  48 , 52  \n", prefixOld)
			}
			TestDistributionScale0(prefixNew)
			TestDistributionScale1AndHalf(prefixNew)

		}
		TestDistibutionPublic := func(prefixOld string) {
			fmt.Printf("%s  Test  For  Service MapOfDistibutesRoundUp  \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)

			TestDistributionWeight1HalfIntMap := func(prefixOld string) {
				mockLogger := &Logger.MockLogger{}
				impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
				impl.Workers, _ = keyGen(2)
				impl.Block.BlockEntity.Capicity = 5

				billDetails1 := entitys.TransactionDetails{}
				billDetails2 := entitys.TransactionDetails{}
				billDetails1.Bill.From.Address = impl.Workers[0]
				billDetails1.Bill.To.Address = impl.Workers[1]

				billDetails2.Bill.To.Address = impl.Workers[0]
				billDetails2.Bill.From.Address = impl.Workers[1]

				transactions := []entitys.TransactionCoins{entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails1, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails2, Amount: 20.0}, entitys.TransactionCoins{BillDetails: billDetails2, Amount: 20.0}}
				impl.Block.Transactions = transactions
				err := impl.Construct()
				if err != nil {
					t.Errorf("Expected no err  but  got %v", err)
				}
				dMap, total := impl.MapOfDistibutesRoundUp(1.5)
				if dMap[impl.Workers[0]] != 48000 || dMap[impl.Workers[1]] != 52000 || total != 100000 {
					t.Errorf("Expected iworker0:%d worker1:%d , total %d but  got  worker0:%d worker1:%d , total %d ", 48000, 52000, 100000, dMap[impl.Workers[0]], dMap[impl.Workers[1]], total)
				}
				fmt.Printf("%s it  should  Create Distribution Map  Rouned  Up  weight 1.5  48000 , 52000\n", prefixOld)
			}
			TestDistributionWeight1HalfIntMap(prefixNew)
		}
		TestGetCurrentHash := func(prefixOld string) {
			fmt.Printf("%s  Test  For  Get Current Hash  \n", prefixOld)
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			GetCurrentHash := func(prefixOld string) {

				mockLogger := &Logger.MockLogger{}
				impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
				err := impl.Construct()
				if err != nil {
					t.Errorf("Expected to  not  get  err  but  got  %v", err)
				}
				impl.Block.BlockEntity.CurrentHash = "aa"
				actual := impl.GetCurrentHash()
				if actual != "aa" {
					t.Errorf("Expected to  not  get  'aa' but  got  '%s'", actual)
				}

				fmt.Printf("%s it should  get  correct  hash aa\n", prefixOld)

			}

			GetCurrentHash2 := func(prefixOld string) {

				mockLogger := &Logger.MockLogger{}
				impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
				err := impl.Construct()
				if err != nil {
					t.Errorf("Expected to  not  get  err  but  got  %v", err)
				}
				impl.Block.BlockEntity.CurrentHash = "bb"
				actual := impl.GetCurrentHash()
				if actual != "bb" {
					t.Errorf("Expected to  not  get  'bb' but  got  '%s'", actual)
				}

				fmt.Printf("%s it  should  get  correct  hash bb\n", prefixOld)

			}
			GetCurrentHash(prefixNew)
			GetCurrentHash2(prefixNew)
		}

		TestServiceConstruct(prefixNew)
		TestDistibutionPrivate(prefixNew)
		TestDistibutionPublic(prefixNew)
		TestGetCurrentHash(prefixNew)
		func(prefixOld string) {
			mockLogger := &Logger.MockLogger{}
			impl := &StakeCoinBlockChain{Services: StakeProviders{LoggerService: mockLogger}}
			impl.Workers, _ = keyGen(2)
			err := impl.Construct()
			if err != nil {
				t.Errorf("Expected to  not  get  err  but  got  %v", err)
			}
			equals := func(list2, list []rsa.PublicKey) bool {
				if len(list) != len(list2) {
					return false
				}
				comparePublicKeys := func(key1, key2 rsa.PublicKey) bool {
					return key1.N.Cmp(key2.N) == 0 && key1.E == key2.E
				}

				for i := 0; i < len(list); i++ {
					if comparePublicKeys(list[i], list2[i]) {
						return false
					}
				}
				return true
			}

			if equals(impl.Workers, impl.GetWorkers()) {
				t.Errorf("Expected to workers to be %v    but  got  %v", impl.Workers, impl.GetWorkers())
			}
			fmt.Printf("%s it should  get Workers\n", prefixOld)

		}(prefixNew)
	}
	TestForStakeCoins(prefix)
}
