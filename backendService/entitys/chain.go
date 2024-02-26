package entitys

import (
	"Hasher"
	"Logger"
	"crypto/rsa"
	"errors"
	"fmt"
)

const ErrTransactionListIsNorFullYet string = "Tranasaction List of size  %d  in not yet Full  has %d"

// Block Chain  [] blockCoinEntity
type BlockChainCoins []BlockCoinEntity
type BlockChainMessage []BlockMessage

/*
*

	Genesis -  Genesis a  new Chain
*/
func (chain *BlockChainCoins) ChainGenesis(logger Logger.LoggerService, hasher Hasher.HashService, Validator rsa.PublicKey, shiftTime int64, capicity, workers int, perNode float64) {
	logger.Log("Start  creating a  new  chain -- GENESIS --  ")
	parent := hasher.ParrentOFall()
	current := hasher.InstantHash(shiftTime)
	var empty BlockChainCoins
	*chain = empty
	genessisBlock := &BlockCoinEntity{}
	genessisBlock.Genesis(Validator, parent, current, capicity, workers, perNode, logger)
	*chain = append(*chain, *genessisBlock)
	logger.Log("Commit  creating a  new  chain -- GENESIS --  ")
}

func (chain *BlockChainCoins) InsertNewBlock(logger Logger.LoggerService, hasher Hasher.HashService, blockDetails BlockCoinEntity) error {
	logger.Log("Start insert a new block in chain")
	size := len(*chain)
	if size != 0 {
		latest := (*chain)[size-1]
		capicity := latest.BlockEntity.Capicity
		hasTransactions := len(latest.Transactions)
		if capicity != hasTransactions {
			errmsg := fmt.Sprintf(ErrTransactionListIsNorFullYet, capicity, hasTransactions)
			logger.Error(fmt.Sprintf("Abbort: %s", errmsg))
			return errors.New(errmsg)
		}
		/*
			expectedIndex := len(*chain)
			bIndex := blockDetails.BlockEntity.Index
			if expectedIndex != bIndex {
				const errmsgTemplate string = "expecting the next block  would  have  index %d  but  had  %d "
				errmsg := fmt.Sprintf(errmsgTemplate, expectedIndex, bIndex)
				logger.Error(fmt.Sprintf("Abbort: %s", errmsg))
				return errors.New(errmsg)
			}
		*/
		logger.Log("Start validation of block")
		err := blockDetails.BlockEntity.ValidateBlock(logger, hasher.Valid, (*chain)[len(*chain)-1].BlockEntity)
		if err != nil {
			errmsg := err.Error()
			logger.Error(fmt.Sprintf("Abbort: Failed validation  due to %s", errmsg))
			return err
		}
		logger.Log("Commit validation of block")
	}
	*chain = append(*chain, blockDetails)
	logger.Log("Commit insert a new block in chain")
	return nil
}
func (chain *BlockChainCoins) InsertTransactions(transactions []TransactionCoins) {
	(*chain)[len(*chain)-1].Transactions = append((*chain)[len(*chain)-1].Transactions, transactions...)
}
func (chain *BlockChainCoins) FindBalance(key rsa.PublicKey) float64 {
	sum := float64(0)
	total := len(*chain)
	collector := make(chan float64, total)
	for i := 0; i < total; i++ {
		(*chain)[i].FindLocaleBalanceOf(key, collector)
	}
	for i := 0; i < total; i++ {
		blockSum := <-collector
		sum += blockSum
	}
	return sum
}

func (chain *BlockChainCoins) GetTransactions(from, twoWay bool, keys []rsa.PublicKey, times []int64) []TransactionCoins {
	list := make([]TransactionCoins, 0)
	total := len(*chain)
	collector := make(chan []TransactionCoins, total)
	executor := func(Block BlockCoinEntity) {
		subList := Block.GetTransactions(from, twoWay, keys, times)
		collector <- subList
	}
	for i := 0; i < total; i++ {
		//	subList := (*chain)[i].GetTransactions(from, twoWay, keys, times)
		//	list = append(list, subList...)
		go executor((*chain)[i])
	}
	for i := 0; i < total; i++ {
		subList := <-collector
		list = append(list, subList...)

	}
	return list
}

/*
*

	Genesis -  Genesis a  new Chain
*/
func (chain *BlockChainMessage) ChainGenesis(logger Logger.LoggerService, hasher Hasher.HashService, Validator rsa.PublicKey, shiftTime int64, capicity int) {
	logger.Log("Start  creating a  new  chain -- GENESIS --  ")
	parent := hasher.ParrentOFall()
	current := hasher.InstantHash(shiftTime)
	var empty BlockChainMessage
	*chain = empty
	genessisBlock := &BlockMessage{}
	genessisBlock.Genesis(Validator, parent, current, capicity, logger)
	*chain = append(*chain, *genessisBlock)
	logger.Log("Commit  creating a  new  chain -- GENESIS --  ")
}
func (chain *BlockChainMessage) InsertNewBlock(logger Logger.LoggerService, hasher Hasher.HashService, blockDetails BlockMessage) error {
	logger.Log("Start insert a new block in chain")
	size := len(*chain)
	if size != 0 {
		latest := (*chain)[size-1]
		capicity := latest.BlockEntity.Capicity
		hasTransactions := len(latest.Transactions)
		if capicity != hasTransactions {
			errmsg := fmt.Sprintf(ErrTransactionListIsNorFullYet, capicity, hasTransactions)
			logger.Error(fmt.Sprintf("Abbort: %s", errmsg))
			return errors.New(errmsg)
		}
		logger.Log("Start validation of block")
		err := blockDetails.BlockEntity.ValidateBlock(logger, hasher.Valid, (*chain)[len(*chain)-1].BlockEntity)
		if err != nil {
			errmsg := err.Error()
			logger.Error(fmt.Sprintf("Abbort: Failed validation  due to %s", errmsg))
			return err
		}
		logger.Log("Commit validation of block")
	}
	*chain = append(*chain, blockDetails)
	logger.Log("Commit insert a new block in chain")
	return nil
}
func (chain *BlockChainMessage) GetTransactions(from, twoWay bool, keys []rsa.PublicKey, times []int64) []TransactionMsg {

	list := make([]TransactionMsg, 0)
	total := len(*chain)
	collector := make(chan []TransactionMsg, total)
	executor := func(Block BlockMessage) {
		subList := Block.GetTransactions(from, twoWay, keys, times)
		collector <- subList
	}
	for i := 0; i < total; i++ {
		//	subList := (*chain)[i].GetTransactions(from, twoWay, keys, times)
		//	list = append(list, subList...)
		go executor((*chain)[i])
	}
	for i := 0; i < total; i++ {
		subList := <-collector
		list = append(list, subList...)

	}

	return list
}
func (chain *BlockChainMessage) InsertTransactions(transaction TransactionMsg) {
	(*chain)[len(*chain)-1].InsertTransaction(transaction)
}
