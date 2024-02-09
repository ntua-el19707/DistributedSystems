package entitys

import (
	"Hasher"
	"Logger"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"testing"
)

// Test cases   for  Block
const prefix = "----"

func TestBlock(t *testing.T) {
	expectorNoErr := func(t *testing.T, prefix string, err error) {
		if err != nil {
			t.Errorf("%s  Expected  no err  but got %v err ", prefix, err)
		}
	}
	fmt.Println("Test  For  Block ")
	TestGenesisBlock := func(t *testing.T, prefix string) {
		block := Block{}
		validator := rsa.PublicKey{}
		logger := &Logger.MockLogger{}
		block.Genesis(validator, "1111", "2222", 6, logger)
		if block.Index != 0 || block.ParentHash != "1111" || block.CurrentHash != "2222" || block.Validator != validator {
			t.Errorf("The  block Index %d ,  Parent %s  , Current %s  and validator  %v but got  %d_%s_%s_%v", 0, "1111", "2222", validator, block.Index, block.ParentHash, block.CurrentHash, block.Validator)
		}
		fmt.Printf("%sIt should  genesis  a general block\n", prefix)
	}
	TestValidBlock := func(t *testing.T, prefixOld string) {
		fmt.Printf("%sTest  For  Valid  Block\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		//Test for  validation  block
		func(t *testing.T, prefixOld string) {
			block := Block{}
			validator := rsa.PublicKey{}
			logger := &Logger.MockLogger{}
			block.Genesis(validator, "1111", "2222", 6, logger)

			block2 := Block{}
			block2.Genesis(validator, "2222", "3333", 6, logger)
			block2.Index = 1
			valid := func(string, string, string) error { return nil }
			err := block2.ValidateBlock(logger, valid, block)
			if err != nil {
				t.Errorf("It  should  not  get err  but  got %s", err.Error())
			}
			fmt.Printf("%sIt should validate a  valid block\n", prefixOld)
		}(t, prefixNew)
		func(t *testing.T, prefixOld string) {

			block := Block{}
			validator := rsa.PublicKey{}
			logger := &Logger.MockLogger{}
			block.Genesis(validator, "1111", "2222", 6, logger)

			block2 := Block{}
			block2.Genesis(validator, "2222", "3333", 6, logger)
			block2.Index = 1
			const errmsg string = "invalid puzzle"
			valid := func(string, string, string) error { return errors.New(errmsg) }
			err := block2.ValidateBlock(logger, valid, block)
			if err == nil {
				t.Errorf("It  should get err but  got  nothing ")
			}
			if err.Error() != errmsg {
				t.Errorf("It  should  not  get invalid puzzle  err  but  got %s", err.Error())
			}
			fmt.Printf("%sIt should validate an  invalid block 'puzzle  error '\n", prefixOld)
		}(t, prefixNew)
		func(t *testing.T, prefixOld string) {

			block := Block{}
			validator := rsa.PublicKey{}
			logger := &Logger.MockLogger{}
			block.Genesis(validator, "1111", "2222", 6, logger)

			block2 := Block{}
			block2.Genesis(validator, "2222", "3333", 6, logger)
			block2.Index = 3
			const errmsg string = "has  not  correct indexing"
			valid := func(string, string, string) error { return nil }
			err := block2.ValidateBlock(logger, valid, block)
			if err == nil {
				t.Errorf("It  should get err but  got  nothing ")
			}
			if err.Error() != errmsg {
				t.Errorf("It  should  not  get %s  err  but  got %s", errmsg, err.Error())
			}
			fmt.Printf("%sIt should validate an  invalid block '%s'\n", prefixOld, errmsg)
		}(t, prefixNew)
		func(t *testing.T, prefixOld string) {

			block := Block{}
			validator := rsa.PublicKey{}
			logger := &Logger.MockLogger{}
			block.Genesis(validator, "1111", "2222", 6, logger)

			block2 := Block{}
			block2.Genesis(validator, "1111", "3333", 6, logger)
			block2.Index = 2
			const errmsg string = "Parent  hash does  not match it previous  currentHash"
			valid := func(string, string, string) error { return nil }
			err := block2.ValidateBlock(logger, valid, block)
			if err == nil {
				t.Errorf("It  should get err but  got  nothing ")
			}
			if err.Error() != errmsg {
				t.Errorf("It  should  not  get %s  err  but  got %s", errmsg, err.Error())
			}
			fmt.Printf("%sIt should validate an  invalid block '%s'\n", prefixOld, errmsg)

		}(t, prefixNew)
	}
	TestMineBlock := func(t *testing.T, prefixOld string) {

		block := Block{}
		block.Index = 1
		validator := rsa.PublicKey{}
		blockP := Block{}
		block.CurrentHash = "1111"
		hash := Hasher.MockHasher{}
		hash.Hashvalue = "2222"

		err := blockP.MineBlock(validator, block, &Logger.MockLogger{}, &hash)

		expectorNoErr(t, prefixOld, err)
		if blockP.Index != 2 || blockP.ParentHash != "1111" || blockP.CurrentHash != "2222" || blockP.Validator != validator {
			t.Errorf("The  block Index %d ,  Parent %s  , Current %s  and validator  %v but got  %d_%s_%s_%v", 2, "1111", "2222", validator, blockP.Index, blockP.ParentHash, blockP.CurrentHash, blockP.Validator)
		}
		fmt.Printf("%sIt should Mine  a  block\n", prefixOld)
	}
	TestGenesisBlock(t, prefix)
	TestValidBlock(t, prefix)
	TestMineBlock(t, prefix)
}
func TestBlockCoin(t *testing.T) {
	fmt.Println("Test  fot  Block Coin ")
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
	TestGenesisBlockCoin := func(t *testing.T, prefixOld string) {
		blockCoin := BlockCoinEntity{}
		blockCoin.workers = 5
		blockCoin.perNode = 1000.0
		keys, _ := keyGen(1)
		validator := keys[0]
		logger := &Logger.MockLogger{}
		blockCoin.Genesis(validator, "1111", "2222", logger)
		block := blockCoin.BlockEntity
		if block.Index != 0 || block.ParentHash != "1111" || block.CurrentHash != "2222" || block.Validator != validator {
			t.Errorf("The  block Index %d ,  Parent %s  , Current %s  and validator  %v but got  %d_%s_%s_%v", 0, "1111", "2222", validator, block.Index, block.ParentHash, block.CurrentHash, block.Validator)
		}
		transactions := blockCoin.Transactions
		if len(transactions) != 2 {
			t.Errorf("Genesis Block  Should  Have  2  transactions but  got %d", len(transactions))
		}
		bill1 := transactions[0].BillDetails.Bill.To.Address
		bill2 := transactions[1].BillDetails.Bill.To.Address
		if bill1 != validator || bill2 != validator {
			t.Errorf("Genesis Block  Transactions  should  Have  To %v but have  %v %v ", validator, bill1, bill2)
		}
		amount1 := transactions[0].Amount
		amount2 := transactions[1].Amount

		if amount1+amount2 != float64(blockCoin.workers)*blockCoin.perNode {
			t.Errorf("Genesis should give   %.6f  but give %.6f   ", float64(blockCoin.workers)*blockCoin.perNode, amount1+amount2)

		}
		fmt.Printf("%sIt should  genesis  a general block coin\n", prefixOld)
	}
	TansactionMaker := func() ([]rsa.PublicKey, []TransactionCoins) {
		keys, _ := keyGen(5)
		//A  +500 , B -250  C  + 750 , D  -1000
		BillDetails1 := TransactionDetails{Bill: BillingInfo{From: Client{Address: keys[3]}, To: Client{Address: keys[0]}}}
		BillDetails2 := TransactionDetails{Bill: BillingInfo{From: Client{Address: keys[3]}, To: Client{Address: keys[2]}}}
		BillDetails3 := TransactionDetails{Bill: BillingInfo{From: Client{Address: keys[1]}, To: Client{Address: keys[3]}}}
		BillDetails4 := TransactionDetails{Bill: BillingInfo{From: Client{Address: keys[2]}, To: Client{Address: keys[3]}}}
		transactions := make([]TransactionCoins, 5)
		transactions[0] = TransactionCoins{BillDetails: BillDetails1, Amount: 250}
		transactions[1] = TransactionCoins{BillDetails: BillDetails3, Amount: 250}
		transactions[2] = TransactionCoins{BillDetails: BillDetails1, Amount: 250}
		transactions[3] = TransactionCoins{BillDetails: BillDetails2, Amount: 1000}
		transactions[4] = TransactionCoins{BillDetails: BillDetails4, Amount: 250}
		return keys, transactions
	}
	TestBalances := func(t *testing.T, prefixOld string) {
		fmt.Printf("%sTest Find Locale Balance\n", prefixOld)
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		keys, transactions := TansactionMaker()
		blockCoin := BlockCoinEntity{}
		blockCoin.Transactions = transactions
		lookFor := func(t *testing.T, key rsa.PublicKey, expected float64, prefixOld string) {
			sum := make(chan float64)
			go blockCoin.FindLocaleBalanceOf(key, sum)

			balance := <-sum
			if balance != expected {
				t.Errorf("Expected to got %.6f But  got  %.6f", expected, balance)
			}
			fmt.Printf("%s It should  have %.6f\n", prefixOld, expected)
		}
		lookFor(t, keys[0], 500.0, prefixNew)
		lookFor(t, keys[1], -250.0, prefixNew)
		lookFor(t, keys[2], 750.0, prefixNew)
		lookFor(t, keys[3], -1000.0, prefixNew)
		lookFor(t, keys[4], 0, prefixNew)

	}

	TestGenesisBlockCoin(t, prefix)
	TestBalances(t, prefix)
}
