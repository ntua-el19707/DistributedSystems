package services

import (
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
	"time"
)

type blockChainCoins []blockCoin

func (chain *blockChainCoins) chainGenesis(providers blockServiceProviders, validator rsa.PublicKey) {
	logger := providers.loggerService
	hasher := providers.hashService
	logger.Log("Start  creating a  new  chain -- GENESIS --  ")
	parent := hasher.ParrentOFall()
	current := hasher.instantHash(-1000)
	var empty blockChainCoins
	*chain = empty
	genessisBlock := &blockCoin{}
	err := genessisBlock.genesis(validator, parent, current)
	if err != nil {
		logger.Error("Abbort chain  creation  ")
		return
	}
	*chain = append(*chain, *genessisBlock)
	fmt.Println(chain)
	logger.Log("Commit  creating a  new  chain -- GENESIS --  ")
}
func (chain *blockChainCoins) insertNewBlock(logger LogerService, hasher hashService, blockDetails blockCoin) error {
	logger.Log("Start insert  a  new  block in chain")
	expectedIndex := len(*chain)
	if expectedIndex != blockDetails.b.index {
		const errmsgTemplate string = "expecting the next block  would  have  index %d  but  had  %d "
		errmsg := fmt.Sprintf(errmsgTemplate, expectedIndex, blockDetails.b.index)
		logger.Error(fmt.Sprintf("Abbort: %s", errmsg))
		return errors.New(errmsg)
	}
	logger.Log("Start  validation of  block ")
	err := blockDetails.b.validateBlock(logger, hasher, (*chain)[len(*chain)-1].b)
	if err != nil {
		errmsg := err.Error()
		logger.Error(fmt.Sprintf("Abbort: Failed valiadtion  due to %s", errmsg))
		return errors.New(errmsg)
	}
	logger.Log("Commit  validation of  block ")

	*chain = append(*chain, blockDetails)
	logger.Log("Commit  insert a  new  block in chain ")
	return nil
}
func (chain blockChainCoins) insertTransaction(logger LogerService, walletService WalletService, transactions []entitys.TransactionCoinEntityRoot) blockChainCoins {
	logger.Log("Start insert  a  Transactions")
	for _, t := range transactions {
		trancationService := TransactionCoins{Transaction: t, services: TransactionsStandard{loggerService: logger}}
		err := trancationService.VerifySignature()
		if err == nil {
			chain[len(chain)-1].Transactions = append(chain[len(chain)-1].Transactions, t.Transaction)
			if *walletService.GetPub() == t.Transaction.BillDetails.Bill.From.Address {
				walletService.UnFreeze(t.Transaction.Amount)
			}
		}

	}

	logger.Log("Commit   new  transaction ")
	return chain
}
func (chain blockChainCoins) findBalance(key rsa.PublicKey) float64 {
	sumChannel := make(chan float64, len(chain))
	for _, block := range chain {
		go block.findLocaleBalanceOf(key, sumChannel)
	}
	var sum float64
	for i := 0; i < len(chain); i++ {
		sum += <-sumChannel
	}
	return sum
}

type block struct {
	index       int
	createdAt   int64
	validator   rsa.PublicKey
	capicity    int
	currentHash string
	purrentHash string
}

func (b *block) genesis(validator rsa.PublicKey, parrent, current string) error {
	b.index = 0                     //first  block
	b.createdAt = time.Now().Unix() //creation  time  stamp
	b.validator = validator
	b.capicity = 5
	b.purrentHash = parrent
	b.currentHash = current //later
	return nil
}
func (b *block) mineBlock(index int, validator rsa.PublicKey, parrent, current string) {
	b.index = index                 //first  block
	b.createdAt = time.Now().Unix() //creation  time  stamp
	b.validator = validator
	b.capicity = 5
	b.purrentHash = parrent
	b.currentHash = current //later

}
func (b *block) validateBlock(logger LogerService, hasher hashService, previous block) error {
	logger.Log(fmt.Sprintf("Start validation  Process  for  block %s to connect  from %s ", b.currentHash, previous.currentHash))
	err := hasher.Valid(previous.purrentHash, previous.currentHash, b.currentHash)
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort validation  Process  for  block %s to connect  from %s  Failed due to %s", b.currentHash, previous.currentHash, err.Error()))
		return err
	}
	if previous.currentHash != b.purrentHash {
		logger.Error(fmt.Sprintf("Abbort validation  Process  for  block %s to connect  from %s  Failed due toParent  hash does  not match it previous  currentHash", b.currentHash, previous.currentHash))
		return errors.New("Parent  hash does  not match it previous  currentHash")
	}
	if previous.index+1 != b.index {
		logger.Error(fmt.Sprintf("Abbort validation  Process  for  block %s to connect  from %s  Failed due to  has  not  correct indexing", b.currentHash, previous.currentHash))
		return errors.New("has  not  correct indexing")
	}
	logger.Log(fmt.Sprintf("Commit validation  Process  for  block %s to connect  from %s ", b.currentHash, previous.currentHash))
	return nil
}

type blockCoin struct {
	b            block
	Transactions []entitys.TransactionCoins
}

func (b *blockCoin) genesis(validator rsa.PublicKey, parrent, current string) error {
	err := b.b.genesis(validator, parrent, current)
	if err != nil {
		return err
	}
	const workers int = 5
	b.Transactions = make([]entitys.TransactionCoins, 2)
	bill := entitys.BillingInfo{To: entitys.Client{Address: validator}}
	BillDetails := entitys.TransactionDetails{Bill: bill, Created_at: time.Now().Unix()}
	initial := entitys.TransactionCoins{Amount: float64(workers * 500), Reason: "BootStrap", BillDetails: BillDetails}
	b.Transactions[0] = initial
	b.Transactions[1] = initial
	return nil

}

func (b blockCoin) findLocaleBalanceOf(key rsa.PublicKey, sumNotify chan float64) {
	var sum float64
	for _, t := range b.Transactions {
		if t.BillDetails.Bill.From.Address == key {
			sum -= t.Amount
		}
		if t.BillDetails.Bill.To.Address == key {
			sum += t.Amount
		}
	}
	sumNotify <- sum

}

type blockChainBlock interface {
	genesis(validator rsa.PublicKey, parrent, current string)
	validateBlock() error
}
type blockChainService interface {
	Service
	genesis() error
	FindBalance() (float64, error)
	Mine() error
	InsertTransaction(t []entitys.TransactionCoinEntityRoot)
}
type blockChainCoinsImpl struct {
	chain    blockChainCoins
	services blockServiceProviders
}
type blockServiceProviders struct {
	loggerService LogerService
	walletService WalletService // i will  wand  my rsa for  the brodcasted  block
	hashService   hashService
}

const blockChainServiceName = "blockChainService"

func (p *blockServiceProviders) construct() error {
	var err error
	if p.loggerService == nil {
		p.loggerService = &Logger{ServiceName: blockChainServiceName}
		err = p.loggerService.construct()
		if err != nil {
			return err
		}
	}
	return p.valid()

}
func (p *blockServiceProviders) valid() error {

	if p.loggerService == nil {
		const errmsg = "The are  is no loggerService"
		return errors.New(errmsg)
	}

	if p.walletService == nil {
		const errmsg = "The are  is no walletService"
		return errors.New(errmsg)
	}
	if p.hashService == nil {
		const errmsg = "The are  is no hashService"
		return errors.New(errmsg)
	}
	return nil

}

func (service *blockChainCoinsImpl) construct() error {
	err := service.services.construct()
	if err != nil {
		return err
	}
	logger := service.services.loggerService
	logger.Log("Service  Created")
	return nil
}
func (service *blockChainCoinsImpl) genesis() error {
	err := service.services.valid()
	if err != nil {
		return err
	}

	service.chain.chainGenesis(service.services, *service.services.walletService.GetPub())
	return nil
}
func (service *blockChainCoinsImpl) FindBalance() (float64, error) {
	err := service.services.valid()
	if err != nil {
		return 0, err
	}

	b := service.chain.findBalance(*service.services.walletService.GetPub())
	return b, nil
}

func (service *blockChainCoinsImpl) InsertTransaction(t []entitys.TransactionCoinEntityRoot) {
	// TODO Check cappicity  after pushing  it (Transactions  will be  transfer  with  Rabbit  MQ =>  there  will be Quee the  same for all nodes )
	//?  if failed  see

	service.chain.insertTransaction(service.services.loggerService, service.services.walletService, t)

}
func (service *blockChainCoinsImpl) Mine() error {

	latest := service.chain[len(service.chain)-1].b
	err, nextHash := service.services.hashService.Hash(latest.purrentHash, latest.currentHash)
	if err != nil {
		return err
	}
	block := blockCoin{b: block{}, Transactions: make([]entitys.TransactionCoins, 0)}

	block.b.mineBlock(len(service.chain), *service.services.walletService.GetPub(), latest.currentHash, nextHash)
	err = service.chain.insertNewBlock(service.services.loggerService, service.services.hashService, block)
	if err != nil {
		return err
	}
	return nil

}
func (service *blockChainCoinsImpl) insertCoinBlockMine(block blockCoin) error {

	err := service.chain.insertNewBlock(service.services.loggerService, service.services.hashService, block)
	if err != nil {
		return err
	}
	return nil

}
