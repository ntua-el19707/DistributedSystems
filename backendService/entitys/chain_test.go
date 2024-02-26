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

func TestChain(t *testing.T) {
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
	fmt.Println("Test  Cases for  Chain")
	validators, _ := keyGen(2)
	const prefix string = "----"
	TestBlockCoin := func(prefixOld string) {
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		fmt.Printf("%s Test  Cases for  block  chain coin\n", prefixOld)
		//Genesis Test
		testGenesis := func(prefixOld string) {
			var chain BlockChainCoins

			logger := &Logger.MockLogger{}
			hasher := &Hasher.MockHasher{InstantHashValue: "123456"}
			chain.ChainGenesis(logger, hasher, validators[0], 0, 6, 5, 1000.0)
			if len(chain) != 1 {
				t.Errorf("chain size expecte to be %d  but  got  %d", 1, len(chain))
			}
			block := chain[0]
			blockDetails := block.BlockEntity
			if blockDetails.Validator != validators[0] {
				t.Errorf("block validator(miner) expected to be %v  but  got  %v", validators[0], blockDetails.Validator)
			}
			if blockDetails.Index != 0 {
				t.Errorf("block inedex  expected to be %d  but  got  %d", 0, blockDetails.Index)

			}
			if blockDetails.CurrentHash != "123456" {
				t.Errorf("current Hash  expected to be %s but  got  %s", "123456", blockDetails.CurrentHash)
			}
			if blockDetails.ParentHash != "1" {
				t.Errorf("parent Hash  expected to be %s but  got  %s", "1", blockDetails.ParentHash)
			}
			if hasher.CallInstand != 1 {
				t.Errorf("hasher service  call instand  expected to be called  1 but  got called  %d", hasher.CallInstand)
			}
			if hasher.CallParentOfAll != 1 {
				t.Errorf("hasher service  call parentOfAll  expected to be called  1 but  got called  %d", hasher.CallParentOfAll)
			}
			expectedLog := []string{"Start  creating a  new  chain -- GENESIS --  ",
				"Commit  creating a  new  chain -- GENESIS --  ",
			}
			if logger.Logs[0] != expectedLog[0] {
				t.Errorf("Expected log %s  and  got  %s", expectedLog[0], logger.Logs[0])
			}
			if logger.Logs[3] != expectedLog[1] {
				t.Errorf("Expected log %s  and  got  %s", expectedLog[1], logger.Logs[3])
			}
			fmt.Printf("%s it  should  genesis  \n", prefixOld)
		}
		testInsertBlock := func(prefixOld string) {
			var chain BlockChainCoins

			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			fmt.Printf("%s Test  Cases for  insert a block  chain coin\n", prefixOld)
			logger := &Logger.MockLogger{}
			hasher := &Hasher.MockHasher{InstantHashValue: "123456"}
			succeedInsert := func(prefixOld string) {
				chain.ChainGenesis(logger, hasher, validators[0], 0, 6, 5, 1000.0)
				chain[0].BlockEntity.Capicity = 2
				logger.Logs = make([]string, 0) // empty logger
				block := BlockCoinEntity{BlockEntity: Block{Index: 1,
					ParentHash: "123456", Validator: validators[1], CurrentHash: "654321",
				}}
				err := chain.InsertNewBlock(logger, hasher, block)
				if err != nil {
					t.Errorf("Expected to got no err but  got  %v", err)

				}
				if len(chain) != 2 {
					t.Errorf("chain size expecte to be %d  but  got  %d", 2, len(chain))
				}
				blockFromChain := chain[1]
				blockDetails := blockFromChain.BlockEntity
				if blockDetails.Validator != validators[1] {
					t.Errorf("block validator(miner) expected to be %v  but  got  %v", validators[0], blockDetails.Validator)
				}
				if blockDetails.Index != 1 {
					t.Errorf("block inedex  expected to be %d  but  got  %d", 1, blockDetails.Index)

				}
				if blockDetails.CurrentHash != "654321" {
					t.Errorf("current Hash  expected to be %s but  got  %s", "654321", blockDetails.CurrentHash)
				}
				if blockDetails.ParentHash != "123456" {
					t.Errorf("parent Hash  expected to be %s but  got  %s", "123456", blockDetails.ParentHash)
				}
				if hasher.CallValid != 1 {
					t.Errorf("hasher service  call valid  expected to be called  1 but  got called  %d", hasher.CallValid)
				}
				log1 := "Start insert a new block in chain"
				log2 := "Start validation of block"
				log3 := "Commit validation of block"
				log4 := "Commit insert a new block in chain"
				expectedLog := []string{log1, log2, log3, log4}
				offset := 0
				for i := 0; i < len(expectedLog)+2; i++ {
					if i == 2 || i == 3 {
						offset--
						continue
					}
					if logger.Logs[i] != expectedLog[i+offset] {
						t.Errorf("Expected log %s  and  got  %s", expectedLog[i+offset], logger.Logs[i])
					}
				}
				fmt.Printf("%s it  should  insert  a valid block   \n", prefixOld)
			}
			FailedInsertHash := func(prefixOld string) {
				chain.ChainGenesis(logger, hasher, validators[0], 0, 6, 5, 1000)
				chain[0].BlockEntity.Capicity = 2
				logger.Logs = make([]string, 0) // empty logger
				logger.ErrorList = make([]string, 0)
				hasher.Invalid = true
				hasher.CallValid = 0
				hasher.InvalidError = "invalid hash"
				genesisBlock := chain[0]
				if len(chain[0].Transactions) != 2 {
					t.Errorf("genesis  block should  have 2 trnsactions(bootstrab transactions )  if  logic change then change the testing and code ")
				}
				t1 := chain[0].Transactions[0]
				t2 := chain[0].Transactions[1]
				block := BlockCoinEntity{BlockEntity: Block{Index: 1,
					ParentHash: "123456", Validator: validators[1], CurrentHash: "654321",
				}}
				expectedErr := errors.New(hasher.InvalidError)
				err := chain.InsertNewBlock(logger, hasher, block)
				if err.Error() != expectedErr.Error() {
					t.Errorf("Expected to get  err %v  but  got  %v", expectedErr, err)

				}
				if len(chain) != 1 {
					t.Errorf("chain size expecte to be %d  but  got  %d", 1, len(chain))
				}
				blockFromChain := chain[0]
				//check integrity of  block
				if genesisBlock.BlockEntity != blockFromChain.BlockEntity {
					t.Errorf("The  previous  should  be  %v but  has  changed to %v ", genesisBlock.BlockEntity, blockFromChain.BlockEntity)
				}

				if genesisBlock.Transactions[0] != t1 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 0, genesisBlock.Transactions[0], t1)
				}
				if genesisBlock.Transactions[1] != t2 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 1, genesisBlock.Transactions[1], t2)
				}
				if hasher.CallValid != 1 {
					t.Errorf("hasher service  call valid  expected to be called  1 but  got called  %d", hasher.CallValid)
				}
				log1 := "Start insert a new block in chain"
				log2 := "Start validation of block"
				expectedLog := []string{log1, log2}
				offset := 0
				for i := 0; i < len(expectedLog)+2; i++ {
					if i == 2 || i == 3 {
						offset--
						continue
					}
					if logger.Logs[i] != expectedLog[i+offset] {
						t.Errorf("Expected log %s  and  got  %s", expectedLog[i+offset], logger.Logs[i])
					}
				}
				if len(logger.ErrorList) != 2 {
					t.Errorf("logger expceted to log  2  msg but  log %d", len(logger.ErrorList))

				}
				errMsg := fmt.Sprintf("Abbort: Failed validation  due to %s", err.Error())
				if logger.ErrorList[1] != errMsg {
					t.Errorf("Expected log %s  and  got  %s", logger.ErrorList[0], errMsg)
				}
				fmt.Printf("%s it  should  failed  (hash service Valid  not  valid ) insert  an invalid block   \n", prefixOld)
			}
			FailedInsertIndex := func(prefixOld string) {
				chain.ChainGenesis(logger, hasher, validators[0], 0, 6, 5, 1000)
				chain[0].BlockEntity.Capicity = 2
				logger.Logs = make([]string, 0) // empty logger
				logger.ErrorList = make([]string, 0)
				hasher.Invalid = false
				hasher.CallValid = 0
				genesisBlock := chain[0]
				if len(chain[0].Transactions) != 2 {
					t.Errorf("genesis  block should  have 2 trnsactions(bootstrab transactions )  if  logic change then change the testing and code ")
				}
				t1 := chain[0].Transactions[0]
				t2 := chain[0].Transactions[1]
				block := BlockCoinEntity{BlockEntity: Block{Index: 2,
					ParentHash: "123456", Validator: validators[1], CurrentHash: "654321",
				}}
				expected := errors.New("has  not  correct indexing")
				errmsg := expected.Error()
				err := chain.InsertNewBlock(logger, hasher, block)
				if err.Error() != errmsg {
					t.Errorf("Expected to get  err %v  but  got  %v", errmsg, err)

				}
				if len(chain) != 1 {
					t.Errorf("chain size expecte to be %d  but  got  %d", 1, len(chain))
				}
				blockFromChain := chain[0]
				//check integrity of  block
				if genesisBlock.BlockEntity != blockFromChain.BlockEntity {
					t.Errorf("The  previous  should  be  %v but  has  changed to %v ", genesisBlock.BlockEntity, blockFromChain.BlockEntity)
				}

				if genesisBlock.Transactions[0] != t1 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 0, genesisBlock.Transactions[0], t1)
				}
				if genesisBlock.Transactions[1] != t2 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 1, genesisBlock.Transactions[1], t2)
				}
				if hasher.CallValid != 0 {
					t.Errorf("hasher service  call valid  expected to be called  %d but  got called  %d", 0, hasher.CallValid)
				}
				log1 := "Start insert a new block in chain"
				log2 := "Start validation of block"
				expectedLog := []string{log1, log2}
				offset := 0
				for i := 0; i < len(expectedLog)+2; i++ {
					if i == 2 || i == 3 {
						offset--
						continue
					}
					if logger.Logs[i] != expectedLog[i+offset] {
						t.Errorf("Expected log %s  and  got  %s", expectedLog[i+offset], logger.Logs[i])
					}
				}
				if len(logger.ErrorList) != 2 {
					t.Errorf("logger expceted to log  2  msg but  log %d", len(logger.ErrorList))

				}
				errMsg := fmt.Sprintf("Abbort: Failed validation  due to %s", err.Error())
				if logger.ErrorList[1] != errMsg {
					t.Errorf("Expected log %s  and  got  %s", logger.ErrorList[1], errMsg)
				}
				fmt.Printf("%s it  should  failed  (index not  increased by 1 ) insert  an invalid block   \n", prefixOld)
			}
			FailedInsertParrentHash := func(prefixOld string) {
				chain.ChainGenesis(logger, hasher, validators[0], 0, 6, 5, 1000.0)
				chain[0].BlockEntity.Capicity = 2
				logger.Logs = make([]string, 0) // empty logger
				logger.ErrorList = make([]string, 0)
				hasher.Invalid = false
				hasher.CallValid = 0
				genesisBlock := chain[0]
				if len(chain[0].Transactions) != 2 {
					t.Errorf("genesis  block should  have 2 trnsactions(bootstrab transactions )  if  logic change then change the testing and code ")
				}
				t1 := chain[0].Transactions[0]
				t2 := chain[0].Transactions[1]
				block := BlockCoinEntity{BlockEntity: Block{Index: 1,
					ParentHash: "111111", Validator: validators[1], CurrentHash: "654321",
				}}
				expected := errors.New("Parent  hash does  not match it previous  currentHash")
				err := chain.InsertNewBlock(logger, hasher, block)
				if err.Error() != expected.Error() {
					t.Errorf("Expected to get  err %v  but  got  %v", expected.Error(), err)

				}
				if len(chain) != 1 {
					t.Errorf("chain size expecte to be %d  but  got  %d", 1, len(chain))
				}
				blockFromChain := chain[0]
				//check integrity of  block
				if genesisBlock.BlockEntity != blockFromChain.BlockEntity {
					t.Errorf("The  previous  should  be  %v but  has  changed to %v ", genesisBlock.BlockEntity, blockFromChain.BlockEntity)
				}

				if genesisBlock.Transactions[0] != t1 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 0, genesisBlock.Transactions[0], t1)
				}
				if genesisBlock.Transactions[1] != t2 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 1, genesisBlock.Transactions[1], t2)
				}
				if hasher.CallValid != 0 {
					t.Errorf("hasher service  call valid  expected to be called  %d but  got called  %d", 0, hasher.CallValid)
				}
				log1 := "Start insert a new block in chain"
				log2 := "Start validation of block"
				expectedLog := []string{log1, log2}
				offset := 0
				for i := 0; i < len(expectedLog)+2; i++ {
					if i == 2 || i == 3 {
						offset--
						continue
					}
					if logger.Logs[i] != expectedLog[i+offset] {
						t.Errorf("Expected log %s  and  got  %s", expectedLog[i+offset], logger.Logs[i])
					}
				}
				if len(logger.ErrorList) != 2 {
					t.Errorf("logger expceted to log 2  msg but  log %d", len(logger.ErrorList))

				}
				errMsg := fmt.Sprintf("Abbort: Failed validation  due to %s", err.Error())
				if logger.ErrorList[1] != errMsg {
					t.Errorf("Expected log %s  and  got  %s", logger.ErrorList[1], errMsg)
				}
				fmt.Printf("%s it  should  failed  (parrent  hash does  not  match  ) insert  an invalid block   \n", prefixOld)
			}
			FailedInsertCapicity := func(prefixOld string) {
				chain.ChainGenesis(logger, hasher, validators[0], 0, 6, 5, 1000)
				chain[0].BlockEntity.Capicity = 5
				logger.Logs = make([]string, 0) // empty logger
				logger.ErrorList = make([]string, 0)
				hasher.Invalid = false
				hasher.CallValid = 0
				genesisBlock := chain[0]
				if len(chain[0].Transactions) != 2 {
					t.Errorf("genesis  block should  have 2 trnsactions(bootstrab transactions )  if  logic change then change the testing and code ")
				}
				t1 := chain[0].Transactions[0]
				t2 := chain[0].Transactions[1]
				block := BlockCoinEntity{BlockEntity: Block{Index: 1,
					ParentHash: "111111", Validator: validators[1], CurrentHash: "654321",
				}}
				errmsg := fmt.Sprintf(ErrTransactionListIsNorFullYet, 5, 2)
				errMsg := fmt.Sprintf("Abbort: %s", errmsg)
				expected := errors.New(errmsg)
				err := chain.InsertNewBlock(logger, hasher, block)
				if err.Error() != expected.Error() {
					t.Errorf("Expected to get  err %v  but  got  %v", expected.Error(), err)

				}
				if len(chain) != 1 {
					t.Errorf("chain size expecte to be %d  but  got  %d", 1, len(chain))
				}
				blockFromChain := chain[0]
				//check integrity of  block
				if genesisBlock.BlockEntity != blockFromChain.BlockEntity {
					t.Errorf("The  previous  should  be  %v but  has  changed to %v ", genesisBlock.BlockEntity, blockFromChain.BlockEntity)
				}

				if genesisBlock.Transactions[0] != t1 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 0, genesisBlock.Transactions[0], t1)
				}
				if genesisBlock.Transactions[1] != t2 {
					t.Errorf("The  previous Transaction at index %d   should  be  %v but  has  changed to %v ", 1, genesisBlock.Transactions[1], t2)
				}
				if hasher.CallValid != 0 {
					t.Errorf("hasher service  call valid  expected to be called  %d but  got called  %d", 0, hasher.CallValid)
				}
				log1 := "Start insert a new block in chain"
				expectedLog := []string{log1}
				offset := 0
				for i := 0; i < len(expectedLog)+2; i++ {
					if i == 1 || i == 2 {
						offset--
						continue
					}
					if logger.Logs[i] != expectedLog[i+offset] {
						t.Errorf("Expected log %s  and  got  %s", expectedLog[i+offset], logger.Logs[i])
					}
				}
				if len(logger.ErrorList) != 1 {
					t.Errorf("logger expceted to log 1  msg but  log %d", len(logger.ErrorList))

				}
				if logger.ErrorList[0] != errMsg {
					t.Errorf("Expected log %s  and  got  %s", logger.ErrorList[0], errMsg)
				}
				fmt.Printf("%s it  should  failed  (transaction list is  not  full ) insert  an invalid block   \n", prefixOld)
			}
			succeedInsert(prefixNew)
			TestFails := func(prefixOld string) {
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				fmt.Printf("%s Test  Cases for  fail to insert  a block \n", prefixOld)
				FailedInsertHash(prefixNew)
				FailedInsertParrentHash(prefixNew)
				FailedInsertIndex(prefixNew)
				FailedInsertCapicity(prefixNew)
			}
			TestFails(prefixNew)
		}
		testGenesis(prefixNew)
		testInsertBlock(prefixNew)
	}
	TestBlockMsg := func(prefixOld string) {

		SenderToPair := func(from, to rsa.PublicKey) TransactionDetails {
			return TransactionDetails{Bill: BillingInfo{From: Client{Address: from}, To: Client{Address: to}}}
		}
		prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
		fmt.Printf("%s Test  Cases for block chain msg implemetation \n", prefixOld)
		TestGenesis := func(prefixOld string) {
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			fmt.Printf("%s Test  For  Genesis \n", prefixOld)
			itShouldGenesis := func(prefixOld string) {
				var chain BlockChainMessage

				logger := &Logger.MockLogger{}
				hasher := &Hasher.MockHasher{InstantHashValue: "123456"}
				chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
				callExpector[int](1, len(chain), t, prefixOld, "size  of  chain ")
				callExpector[int](1, hasher.CallParentOfAll, t, prefixOld, "call  hash.ParrnetOfall")
				callExpector[int](1, hasher.CallInstand, t, prefixOld, "call hash.InstantHash ")
				fmt.Printf("%s  it should  Genesis \n", prefixOld)

			}
			itShouldGenesis(prefixNew)

		}
		TestInsertBlock := func(prefixOld string) {

			AtoB := SenderToPair(validators[0], validators[1])
			BtoA := SenderToPair(validators[1], validators[0])
			list := make([]TransactionMsg, 5)
			list[0] = TransactionMsg{BillDetails: AtoB, Msg: "apples"}
			list[1] = TransactionMsg{BillDetails: BtoA, Msg: "bananas"}
			list[2] = TransactionMsg{BillDetails: AtoB, Msg: "oranges"}
			list[3] = TransactionMsg{BillDetails: BtoA, Msg: "lemon"}
			list[4] = TransactionMsg{BillDetails: AtoB, Msg: "peach"}
			logger := &Logger.MockLogger{}
			hasher := &Hasher.MockHasher{InstantHashValue: "123456"}

			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			fmt.Printf("%s Test  Cases for insert a block to chain \n", prefixOld)
			itShouldInsert := func(prefixOld string) {
				var chain BlockChainMessage
				chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
				chain[0].Transactions = list
				first := chain[0]
				candidateBlock := BlockMessage{}
				candidateBlock.BlockEntity.Index = 1
				candidateBlock.BlockEntity.ParentHash = "123456"
				err := chain.InsertNewBlock(logger, hasher, candidateBlock)
				expectorNoErr(t, err, prefixOld)
				callExpector[int](2, len(chain), t, prefixOld, "chain size")
				callExpector[Block](candidateBlock.BlockEntity, chain[1].BlockEntity, t, prefixOld, "item in  index 1")
				callExpector[Block](first.BlockEntity, chain[0].BlockEntity, t, prefixOld, "item in  index 0")
				callExpector[int](0, len(chain[1].Transactions), t, prefixOld, "Trnasaction list at index 1  size")
				fmt.Printf("%s  it should  insert a valid block \n", prefixOld)
			}
			itShouldFail := func(prefixOld string) {
				prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
				fmt.Printf("%s Test  Cases for fail to insert a block to chain \n", prefixOld)
				// -- test fail index  --
				itShouldFailDueToIndex := func(prefixOld string) {

					var chain BlockChainMessage
					chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
					chain[0].Transactions = list
					first := chain[0]
					candidateBlock := BlockMessage{}
					candidateBlock.BlockEntity.Index = 2
					candidateBlock.BlockEntity.ParentHash = "123456"
					err := chain.InsertNewBlock(logger, hasher, candidateBlock)
					expectedErr := errors.New("has  not  correct indexing")
					callExpector[string](expectedErr.Error(), err.Error(), t, prefixOld, "expected error")
					callExpector[int](1, len(chain), t, prefixOld, "chain size")
					callExpector[Block](first.BlockEntity, chain[0].BlockEntity, t, prefixOld, "item in  index 0")

					fmt.Printf("%s  it should  fail  insert an invalid block 'invalid index' \n", prefixOld)
				}
				itShouldFailDueToparrentHash := func(prefixOld string) {
					var chain BlockChainMessage
					chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
					chain[0].Transactions = list
					first := chain[0]
					candidateBlock := BlockMessage{}
					candidateBlock.BlockEntity.Index = 1
					candidateBlock.BlockEntity.ParentHash = "1234"
					err := chain.InsertNewBlock(logger, hasher, candidateBlock)
					expectedErr := errors.New("Parent  hash does  not match it previous  currentHash")
					callExpector[string](expectedErr.Error(), err.Error(), t, prefixOld, "expected error")
					callExpector[int](1, len(chain), t, prefixOld, "chain size")
					callExpector[Block](first.BlockEntity, chain[0].BlockEntity, t, prefixOld, "item in  index 0")
					fmt.Printf("%s  it should  fail  insert an invalid block 'invalid parent hash' \n", prefixOld)
				}
				itShouldFailDueToHashValid := func(prefixOld string) {
					hasher.Invalid = true
					hasher.CallValid = 0
					hasher.InvalidError = "invalid hash"
					var chain BlockChainMessage
					chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
					chain[0].Transactions = list
					first := chain[0]
					candidateBlock := BlockMessage{}
					candidateBlock.BlockEntity.Index = 1
					candidateBlock.BlockEntity.ParentHash = "123456"
					err := chain.InsertNewBlock(logger, hasher, candidateBlock)
					expectedErr := errors.New(hasher.InvalidError)
					callExpector[string](expectedErr.Error(), err.Error(), t, prefixOld, "expected error")
					callExpector[int](1, len(chain), t, prefixOld, "chain size")
					callExpector[Block](first.BlockEntity, chain[0].BlockEntity, t, prefixOld, "item in  index 0")
					callExpector[int](1, hasher.CallValid, t, prefixOld, "call  Hsh.valid")
					hasher.Invalid = false
					hasher.CallValid = 0
					hasher.InvalidError = ""
					fmt.Printf("%s  it should  fail  insert an invalid block 'invalid hash  compination' \n", prefixOld)
				}
				itShouldFailDueToCappicity := func(prefixOld string) {
					list := make([]TransactionMsg, 4)
					list[0] = TransactionMsg{BillDetails: AtoB, Msg: "apples"}
					list[1] = TransactionMsg{BillDetails: BtoA, Msg: "peach"}
					list[2] = TransactionMsg{BillDetails: BtoA, Msg: "bannanas"}
					list[3] = TransactionMsg{BillDetails: AtoB, Msg: "mellons"}
					var chain BlockChainMessage
					chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
					chain[0].Transactions = list
					first := chain[0]
					candidateBlock := BlockMessage{}
					candidateBlock.BlockEntity.Index = 1
					candidateBlock.BlockEntity.ParentHash = "123456"
					err := chain.InsertNewBlock(logger, hasher, candidateBlock)
					errmsg := fmt.Sprintf(ErrTransactionListIsNorFullYet, 5, 4)
					expectedErr := errors.New(errmsg)
					callExpector[string](expectedErr.Error(), err.Error(), t, prefixOld, "expected error")
					callExpector[int](1, len(chain), t, prefixOld, "chain size")
					callExpector[Block](first.BlockEntity, chain[0].BlockEntity, t, prefixOld, "item in  index 0")
					fmt.Printf("%s  it should  fail  insert an invalid block 'capicity error ' \n", prefixOld)
				}
				itShouldFailDueToIndex(prefixNew)
				itShouldFailDueToparrentHash(prefixNew)
				itShouldFailDueToHashValid(prefixNew)
				itShouldFailDueToCappicity(prefixNew)
			}

			itShouldInsert(prefixNew)
			itShouldFail(prefixNew)
		}
		TestInsertTransactions := func(prefixOld string) {
			AtoB := SenderToPair(validators[0], validators[1])
			prefixNew := fmt.Sprintf("%s%s", prefixOld, prefix)
			fmt.Printf("%s Test  for  insert transactions \n", prefixOld)
			logger := &Logger.MockLogger{}
			hasher := &Hasher.MockHasher{InstantHashValue: "123456"}
			itShouldInsert := func(prefixOld string) {
				var chain BlockChainMessage
				chain.ChainGenesis(logger, hasher, validators[0], 0, 5)
				transaction := TransactionMsg{BillDetails: AtoB, Msg: "apples"}
				chain.InsertTransactions(transaction)
				callExpector[int](1, len(chain), t, prefixOld, "chain size")
				callExpector[int](1, len(chain[0].Transactions), t, prefixOld, "tranaction list size")
				callExpector[TransactionMsg](transaction, chain[0].Transactions[0], t, prefixOld, "item in  index 0")
				fmt.Printf("%s  it should    insert a transaction  \n", prefixOld)

			}
			itShouldInsert(prefixNew)
		}
		TestGenesis(prefixNew)
		TestInsertBlock(prefixNew)
		TestInsertTransactions(prefixNew)

	}
	TestBlockCoin(prefix)
	TestBlockMsg(prefix)
}
