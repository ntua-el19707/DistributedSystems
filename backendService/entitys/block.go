package entitys

import (
	"Hasher"
	"Logger"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"
)

// Type  Block
type Block struct {
	Index       int           `json:"index"`
	CreatedAt   int64         `json:"created_at"`
	Validator   rsa.PublicKey `json:"validator"`
	Capicity    int           `json:"capicity"`
	CurrentHash string        `json:"current_hash"`
	ParentHash  string        `json:"parrent_hash"`
}

func (b *Block) Genesis(Validator rsa.PublicKey, Parent, Current string, Capicity int, logger Logger.LoggerService) {
	b.Index = 0                     //first  block
	b.CreatedAt = time.Now().Unix() //creation  time  stamp
	b.Validator = Validator
	b.Capicity = Capicity
	b.ParentHash = Parent
	b.CurrentHash = Current // later
	logger.Log(fmt.Sprintf("Created  Genesis Block  %s ", b.CurrentHash))
}
func (b *Block) MineBlock(validator rsa.PublicKey, previousBlock Block, logger Logger.LoggerService, hasher Hasher.HashService) error {
	logger.Log("Start Mining a  new  Block ")
	err, current := hasher.Hash(previousBlock.ParentHash, previousBlock.CurrentHash)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	b.Index = previousBlock.Index + 1
	b.CreatedAt = time.Now().Unix() //creation  time  stamp
	b.Validator = validator
	b.Capicity = previousBlock.Capicity
	b.ParentHash = previousBlock.CurrentHash
	b.CurrentHash = current
	logger.Log("Commit Mining a  new  Block ")
	return nil
}
func (b *Block) ValidateBlock(logger Logger.LoggerService, Valid func(string, string, string) error, previous Block) error {
	logger.Log(fmt.Sprintf("Start validation  Process  for  block %s to connect  from %s ", b.CurrentHash, previous.CurrentHash))
	if previous.CurrentHash != b.ParentHash {
		logger.Error(fmt.Sprintf("Abbort validation  Process  for  block %s to connect  from %s  Failed due toParent  hash does  not match it previous  currentHash", b.CurrentHash, previous.CurrentHash))
		return errors.New("Parent  hash does  not match it previous  currentHash")
	}
	if previous.Index+1 != b.Index {
		logger.Error(fmt.Sprintf("Abbort validation  Process  for  block %s to connect  from %s  Failed due to  has  not  correct indexing", b.CurrentHash, previous.CurrentHash))
		return errors.New("has  not  correct indexing")
	}
	err := Valid(previous.ParentHash, previous.CurrentHash, b.CurrentHash)
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort validation  Process  for  block %s to connect  from %s  Failed due to %s", b.CurrentHash, previous.CurrentHash, err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("Commit validation  Process  for  block %s to connect  from %s ", b.CurrentHash, previous.CurrentHash))
	return nil
}

// Type  BlockCoin
type BlockCoinEntity struct {
	BlockEntity  Block              `json:"block"`
	Transactions []TransactionCoins `json:"transactions"`
	workers      int
	perNode      float64
}

/*
*

	Genesis - Genesis
	@Param Validator rsa.PublicKey
	@Param Parent  string
	@Param Current string
	@Param logger  Logger.LoggerService
*/
func (block *BlockCoinEntity) Genesis(Validator rsa.PublicKey, Parent, Current string, logger Logger.LoggerService) {
	//Genesis The BlockEnity
	blockSize := 6
	block.BlockEntity.Genesis(Validator, Parent, Current, blockSize, logger)
	/*	if block.perNode == 0 {i
		block.perNode = 1000.0
	}*/
	workers := 5             //block.workers
	perNode := float64(1000) // block.perNode
	total := float64(workers) * perNode
	block.Transactions = make([]TransactionCoins, 2)
	bill := BillingInfo{To: Client{Address: Validator}}
	BillDetails := TransactionDetails{Bill: bill, Created_at: time.Now().Unix()}
	initial := TransactionCoins{Amount: total / 2, Reason: "BootStrap", BillDetails: BillDetails}
	block.Transactions[0] = initial
	block.Transactions[1] = initial
	logger.Log(fmt.Sprintf("Created  Genesis Block Coin   %s  with  The total of  %.6f and  to distribute %.6f ", block.BlockEntity.CurrentHash, total, perNode))
}

func (b *BlockCoinEntity) MineBlock(validator rsa.PublicKey, previousBlock Block, logger Logger.LoggerService, hasher Hasher.HashService) error {
	return b.BlockEntity.MineBlock(validator, previousBlock, logger, hasher)
} /*
*

	FindLocaleBalanceOf -  Find The Balanace  of key   in  a Block
	@Param key rsa.PublicKey
	@Param sumNotify chan float64
*/
func (blockCoin BlockCoinEntity) FindLocaleBalanceOf(key rsa.PublicKey, sumNotify chan float64) {
	var sum float64
	for _, t := range blockCoin.Transactions {
		if t.BillDetails.Bill.From.Address == key {
			sum -= t.Amount
		}
		if t.BillDetails.Bill.To.Address == key {
			sum += t.Amount
		}
	}
	sumNotify <- sum
}

type BlockMessage struct {
	BlockEntity  Block            `json:"block"`
	Transactions []TransactionMsg `json:"transactions"`
}

func (block *BlockMessage) Genesis(Validator rsa.PublicKey, Parent, Current string, logger Logger.LoggerService) {
	blockSize := 5
	block.BlockEntity.Genesis(Validator, Parent, Current, blockSize, logger)
	logger.Log(fmt.Sprintf("Created  Genesis Block Message  %s", block.BlockEntity.CurrentHash))
}
func (b *BlockMessage) MineBlock(validator rsa.PublicKey, previousBlock Block, logger Logger.LoggerService, hasher Hasher.HashService) error {
	return b.BlockEntity.MineBlock(validator, previousBlock, logger, hasher)
}
func equalPublicKeys(key1, key2 *rsa.PublicKey) bool {
	return key1.N.Cmp(key2.N) == 0 && key1.E == key2.E
}
func (b *BlockMessage) GetTransactions(from, twoWay bool, keys []rsa.PublicKey, times []int64) []TransactionMsg {
	var list []TransactionMsg
	add := func(t TransactionMsg) {
		if len(times) == 0 {
			list = append(list, t)
		} else if len(times) == 1 {
			created := t.BillDetails.Created_at
			if created >= times[0] {
				list = append(list, t)
			}
		} else if len(times) >= 2 {
			created := t.BillDetails.Created_at
			if created >= times[0] && created <= times[1] {
				list = append(list, t)
			}
		}
	}

	if len(keys) == 0 {
		return nil
	} else if len(keys) == 1 {
		if twoWay {
			key := keys[0]
			for _, t := range b.Transactions {
				from := t.BillDetails.Bill.From.Address
				to := t.BillDetails.Bill.To.Address
				if equalPublicKeys(&key, &from) || equalPublicKeys(&key, &to) {
					add(t)
				}
			}
		} else if from {
			key := keys[0]
			for _, t := range b.Transactions {
				from := t.BillDetails.Bill.From.Address
				if equalPublicKeys(&key, &from) {
					add(t)
				}
			}
		} else {
			key := keys[0]
			for _, t := range b.Transactions {
				to := t.BillDetails.Bill.To.Address
				if equalPublicKeys(&key, &to) {
					add(t)
				}
			}
		}
	} else {
		fromKey := keys[0]
		toKey := keys[1]
		for _, t := range b.Transactions {
			from := t.BillDetails.Bill.From.Address
			to := t.BillDetails.Bill.To.Address

			if equalPublicKeys(&fromKey, &from) && equalPublicKeys(&toKey, &to) {
				add(t)
			} else if equalPublicKeys(&fromKey, &to) && equalPublicKeys(&toKey, &from) && twoWay {
				add(t)
			}
		}

	}
	return list

}

func (b *BlockMessage) InsertTransaction(t TransactionMsg) {
	b.Transactions = append(b.Transactions, t)
}
