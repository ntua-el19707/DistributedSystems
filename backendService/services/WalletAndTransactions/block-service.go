package WalletAndTransactions

import (
	"Hasher"
	"Logger"
	"Lottery"
	"RabbitMqService"
	"Service"
	"Stake"
	"crypto/rsa"
	"entitys"
	"errors"
	"fmt"
	"sync"
)

type BlockChainService interface {
	Service.Service
	Genesis() error
	FindBalance() (float64, error)
	Mine() error
	InsertTransaction(t []entitys.TransactionCoinEntityRoot)
}
type BlockChainCoinsImpl struct {
	Chain       entitys.BlockChainCoins
	Workers     []rsa.PublicKey
	Services    BlockServiceProviders
	mu          sync.Mutex
	ScaleFactor float64
}
type BlockServiceProviders struct {
	LoggerService         Logger.LoggerService
	WalletServiceInstance *WalletStructV1Implementation // i will  wand  my rsa for  the brodcasted  block
	HashService           Hasher.HashService
	LotteryService        Lottery.LotteryService
	RabbitMqService       RabbitMqService.RabbitMqService
}

const blockChainServiceName = "Block-Chain-Service"

func EqualPublicKeys(key1, key2 *rsa.PublicKey) bool {
	return key1.N.Cmp(key2.N) == 0 && key1.E == key2.E
}

func (p *BlockServiceProviders) Construct() error {
	var err error
	if p.LoggerService == nil {
		p.LoggerService = &Logger.Logger{ServiceName: blockChainServiceName}
		err = p.LoggerService.Construct()
		if err != nil {
			return err
		}
	}
	return p.valid()

}
func (p *BlockServiceProviders) valid() error {

	if p.LoggerService == nil {
		const errmsg = "The are  is no loggerService"
		return errors.New(errmsg)
	}

	if p.WalletServiceInstance == nil {
		const errmsg = "The are  is no walletService"
		return errors.New(errmsg)
	}
	if p.HashService == nil {
		const errmsg = "The are  is no hashService"
		return errors.New(errmsg)
	}
	return nil

}

func (service *BlockChainCoinsImpl) Construct() error {
	err := service.Services.Construct()
	if err != nil {
		return err
	}
	logger := service.Services.LoggerService
	logger.Log("Service  Created")
	return nil
}

func (service *BlockChainCoinsImpl) Genesis() error {
	err := service.Services.valid()
	if err != nil {
		return err
	}
	service.Chain.ChainGenesis(service.Services.LoggerService, service.Services.HashService, service.Services.WalletServiceInstance.GetPub(), 0)
	return nil
}

const ErrMsgNotValidator = ""

func (service *BlockChainCoinsImpl) FindBalance(key rsa.PublicKey) float64 {
	service.mu.Lock()
	defer service.mu.Unlock()
	amount := service.Chain.FindBalance(key)
	return amount

}
func (service *BlockChainCoinsImpl) findAndLock(amount float64) (float64, error) {

	service.mu.Lock()
	defer service.mu.Unlock()

	err := service.Services.WalletServiceInstance.Freeze(amount)
	if err != nil {
		errMsg := fmt.Sprintf("Could not Freeze  money due to %s ", err.Error())
		return 0, errors.New(errMsg)
	}
	frozen := service.Services.WalletServiceInstance.GetFreeze()

	total := service.Chain.FindBalance(service.Services.WalletServiceInstance.GetPub())

	return total - frozen, nil

}
func (service *BlockChainCoinsImpl) GetTransactions(from, twoWay bool, keys []rsa.PublicKey, times []int64) []entitys.TransactionCoins {
	logger := service.Services.LoggerService
	service.mu.Lock()
	defer func() {
		logger.Log("Unlock")
		service.mu.Unlock()
	}()
	logger.Log("Lock")

	list := service.Chain.GetTransactions(from, twoWay, keys, times)

	return list
}

func (service *BlockChainCoinsImpl) InsertTransaction(t entitys.TransactionCoinSet) error {
	logger := service.Services.LoggerService
	logger.Log(fmt.Sprintf("start insert  transactions SET %s  and %s", t.Tax.Transaction.BillDetails.Transaction_id, t.Transfer.Transaction.BillDetails.Transaction_id))
	service.mu.Lock()
	defer service.mu.Unlock()
	logger.Log("Lock  BlockChain ")
	processPublicKey := service.Services.WalletServiceInstance.GetPub()
	lastBlock := service.Chain[len(service.Chain)-1]
	validator := lastBlock.BlockEntity.Validator

	transactions := make([]entitys.TransactionCoins, 2)
	transactions[0] = t.Tax.Transaction
	transactions[1] = t.Transfer.Transaction
	err := verify(t)
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort due to %s", err.Error()))
		return err
	}

	if lastBlock.BlockEntity.Capicity <= len(lastBlock.Transactions) {
		//Mine Block
		logger.Log("Start  Mine Block Coins")

		stake := Stake.StakeCoinBlockChain{
			Block:   lastBlock,
			Workers: service.Workers}
		err := stake.Construct()
		if err != nil {
			logger.Fatal(err.Error())
		}
		service.Services.LotteryService.LoadStakeService(&stake)
		luckyOne, err := service.Services.LotteryService.Spin(service.ScaleFactor)
		if err != nil {
			logger.Fatal(err.Error())
		}
		if EqualPublicKeys(&processPublicKey, &luckyOne) {
			//Win And Miner
			block := entitys.BlockCoinEntity{}
			err := block.MineBlock(luckyOne, lastBlock.BlockEntity, service.Services.LoggerService, service.Services.HashService)
			if err != nil {
				logger.Fatal(err.Error())
			}
			broadCastBlock := entitys.BlockCoinMessageRabbitMq{
				BlockCoin: block,
			}
			err = service.Services.RabbitMqService.PublishBlockCoin(broadCastBlock)
			if err != nil {
				// harakiri
				logger.Fatal(err.Error())
				return err
			}

		}
		block := service.Services.RabbitMqService.ConsumeNextBlockCoin()
		validator = block.BlockCoin.BlockEntity.Validator
		//NOW The re is  litle  a chance to fail Mine only if  internal error if err =>  commit  harakiri
		err = service.Chain.InsertNewBlock(service.Services.LoggerService, service.Services.HashService, block.BlockCoin)
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Log("Commit  Mine Block Coins ")
	}
	//Re Sign  and  send
	From := transactions[0].BillDetails.Bill.From.Address
	//Stamp Validatori
	transactions[0].BillDetails.Bill.To.Address = validator
	service.Chain.InsertTransactions(transactions)
	if EqualPublicKeys(&processPublicKey, &From) {
		amount := transactions[0].Amount + transactions[1].Amount
		err := service.Services.WalletServiceInstance.UnFreeze(amount)
		if err != nil {
			logger.Fatal(err.Error())
		}

	}
	logger.Log("UnLock  BlockChain")
	logger.Log(fmt.Sprintf("commit insert  transactions SET %s  and %s", t.Tax.Transaction.BillDetails.Transaction_id, t.Transfer.Transaction.BillDetails.Transaction_id))
	return nil
}

type BlockChainMsgImpl struct {
	Chain       entitys.BlockChainMessage
	Workers     []rsa.PublicKey
	Services    BlockServiceProviders
	mu          sync.Mutex
	ScaleFactor float64
}

func (service *BlockChainMsgImpl) Construct() error {
	err := service.Services.Construct()
	if err != nil {
		return err
	}
	logger := service.Services.LoggerService
	logger.Log("Service  Created")
	return nil
}
func (service *BlockChainMsgImpl) GetTransactions(from, twoWay bool, keys []rsa.PublicKey, times []int64) []entitys.TransactionMsg {
	logger := service.Services.LoggerService
	service.mu.Lock()
	defer func() {
		logger.Log("Unlock")
		service.mu.Unlock()
	}()
	logger.Log("Lock")

	list := service.Chain.GetTransactions(from, twoWay, keys, times)

	return list
}
func (service *BlockChainMsgImpl) Genesis() error {
	err := service.Services.valid()
	if err != nil {
		return err
	}
	service.Chain.ChainGenesis(service.Services.LoggerService, service.Services.HashService, service.Services.WalletServiceInstance.GetPub(), 0)
	return nil
}

func (service *BlockChainMsgImpl) InsertTransaction(t entitys.TransactionMessageSet) error {

	logger := service.Services.LoggerService
	logger.Log(fmt.Sprintf("start insert  transaction Msg  %s", t.TransactionMessage.Transaction.BillDetails.Transaction_id))
	service.mu.Lock()
	defer func() {
		logger.Log("Unlock")
		service.mu.Unlock()
	}()
	logger.Log("Lock")
	processPublicKey := service.Services.WalletServiceInstance.GetPub()

	wallet := service.Services.WalletServiceInstance
	trMsg := t.TransactionMessage
	sender := trMsg.Transaction.BillDetails.Bill.From.Address
	// -- Verify  Purchase Transactions --
	err := verify(t.TransactionCoin)
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort due to %s", err.Error()))
		return err
	}
	// -- verify  Transaction --
	transactionService := TransactionMsg{Transaction: trMsg}
	err = transactionService.semiConstruct()
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort due to %s", err.Error()))
		return err
	}
	err = transactionService.VerifySignature()
	if err != nil {
		logger.Error(fmt.Sprintf("Abbort due to %s", err.Error()))
		return err
	}
	//if  full -Mine
	lastBlock := service.Chain[len(service.Chain)-1]
	validator := lastBlock.BlockEntity.Validator
	if lastBlock.BlockEntity.Capicity == len(lastBlock.Transactions) {
		// -- MINE --
		logger.Log("Start Mine")
		stake := Stake.StakeMesageBlockChain{
			Block:   lastBlock,
			Workers: service.Workers,
		}
		err := stake.Construct()
		if err != nil {
			// harakiri
			logger.Fatal(err.Error())
			return err
		}
		service.Services.LotteryService.LoadStakeService(&stake)
		luckyOne, err := service.Services.LotteryService.Spin(service.ScaleFactor)
		if err != nil {
			// harakiri
			logger.Fatal(err.Error())
			return err
		}
		if EqualPublicKeys(&processPublicKey, &luckyOne) {
			block := entitys.BlockMessage{}
			err := block.MineBlock(luckyOne, lastBlock.BlockEntity, service.Services.LoggerService, service.Services.HashService)
			if err != nil {
				// harakiri
				logger.Fatal(err.Error())
				return err
			}
			blockToBroadcast := entitys.BlockMessageMessageRabbitMq{
				BlockMsg: block,
			}
			err = service.Services.RabbitMqService.PublishBlockMsg(blockToBroadcast)

			if err != nil {
				// harakiri
				logger.Fatal(err.Error())
				return err
			}
		}
		recievedBlock := service.Services.RabbitMqService.ConsumeNextBlockMsg()
		validator = recievedBlock.BlockMsg.BlockEntity.Validator
		service.Chain.InsertNewBlock(service.Services.LoggerService, service.Services.HashService, recievedBlock.BlockMsg)
		logger.Log("Commit Mine")

	}
	service.Chain.InsertTransactions(trMsg.Transaction)

	if EqualPublicKeys(&processPublicKey, &sender) {
		//stamp to who go the money
		t.TransactionCoin.Transfer.Transaction.BillDetails.Bill.To.Address = validator
		var transactionTransfer entitys.TransactionCoinEntityRoot
		transactionTransfer.Transaction = t.TransactionCoin.Transfer.Transaction
		transferTransactionService := TransactionCoins{Transaction: transactionTransfer}
		transferTransactionService.semiConstruct()
		err = wallet.Sign(&transferTransactionService)
		if err != nil {
			// harakiri
			logger.Fatal(err.Error())
			return err
		}
		t.TransactionCoin.Transfer = transferTransactionService.Transaction
		//re sign
		err := service.Services.RabbitMqService.PublishTractioncoinSet(t.TransactionCoin)
		if err != nil {
			// harakiri
			logger.Fatal(err.Error())
			return err
		}

	}
	//add  transaction
	//add transfer dst  and  sign =. use rabbitMw to publish

	logger.Log(fmt.Sprintf("commit insert  transaction Msg  %s", t.TransactionMessage.Transaction.BillDetails.Transaction_id))
	return nil
}

/*
*

	verify - verify a tranction coin set
	@Param  t  entitys.TransactionCoinSet
	@Returns  error
*/
func verify(t entitys.TransactionCoinSet) error {
	taxTransactionService := TransactionCoins{Transaction: t.Tax}
	transferTransactionService := TransactionCoins{Transaction: t.Transfer}

	errChannel := make(chan error, 2)
	semiVerify := func(transactionService TransactionService, finish chan error) {
		err := transactionService.semiConstruct()
		if err != nil {
			finish <- err
			return
		}
		finish <- transactionService.VerifySignature()
	}
	go semiVerify(&taxTransactionService, errChannel)
	go semiVerify(&transferTransactionService, errChannel)
	var errHappen error
	for i := 0; i < 2; i++ {
		err := <-errChannel
		if err != nil {
			errHappen = err
		}
	}
	return errHappen
}
