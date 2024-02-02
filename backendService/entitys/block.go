package entitys

import (
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

const Capicity int = 5

func (b *Block) Genesis(Validator rsa.PublicKey, Parent, Current string, logger Logger.LoggerService) {
	b.Index = 0                     //first  block
	b.CreatedAt = time.Now().Unix() //creation  time  stamp
	b.Validator = Validator
	b.Capicity = Capicity
	b.ParentHash = Parent
	b.CurrentHash = Current // later
	logger.Log(fmt.Sprintf("Created  Genesis Block  %s ", b.CurrentHash))
}
func (b *Block) MineBlock(index int, validator rsa.PublicKey, parrent, current string) {
	b.Index = index                 //first  block
	b.CreatedAt = time.Now().Unix() //creation  time  stamp
	b.Validator = validator
	b.Capicity = Capicity
	b.ParentHash = parrent
	b.CurrentHash = current //later

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
	block.BlockEntity.Genesis(Validator, Parent, Current, logger)
	if block.perNode == 0 {
		block.perNode = 1000.0
	}
	workers := block.workers
	perNode := block.perNode
	total := float64(workers) * perNode
	block.Transactions = make([]TransactionCoins, 2)
	bill := BillingInfo{To: Client{Address: Validator}}
	BillDetails := TransactionDetails{Bill: bill, Created_at: time.Now().Unix()}
	initial := TransactionCoins{Amount: total / 2, Reason: "BootStrap", BillDetails: BillDetails}
	block.Transactions[0] = initial
	block.Transactions[1] = initial
	logger.Log(fmt.Sprintf("Created  Genesis Block Coin   %s  with  The total of  %.6f and  to distribute %.6f ", block.BlockEntity.CurrentHash, total, perNode))
}

/*
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
