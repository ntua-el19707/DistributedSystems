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
	Chain    entitys.BlockChainCoins
	Workers  []rsa.PublicKey
	Services BlockServiceProviders
	mu       sync.Mutex
}
type BlockServiceProviders struct {
	LoggerService         Logger.LoggerService
	WalletServiceInstance *WalletStructV1Implementation // i will  wand  my rsa for  the brodcasted  block
	HashService           Hasher.HashService
	LotteryService        Lottery.LotteryService
	RabbitMqService       RabbitMqService.RabbitMqService
}

const blockChainServiceName = "Block-Chain-Service"

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
func (service *BlockChainCoinsImpl) FindBalance() (float64, error) {
	err := service.Services.valid()
	if err != nil {
		return 0, err
	}

	return 0, nil
}

const ErrMsgNotValidator = ""

func (service *BlockChainCoinsImpl) InsertTransaction(t entitys.TransactionCoinSet) error {
	logger := service.Services.LoggerService
	logger.Log(fmt.Sprintf("start insert  transactions SET %s  and %s", t.Tax.Transaction.BillDetails.Transaction_id, t.Transfer.Transaction.BillDetails.Transaction_id))
	service.mu.Lock()
	defer service.mu.Unlock()
	logger.Log("Lock  BlockChain ")
	lastBlock := service.Chain[len(service.Chain)-1]
	validator := lastBlock.BlockEntity.Validator

	transactions := make([]entitys.TransactionCoins, 2)
	transactions[0] = t.Tax.Transaction
	transactions[1] = t.Transfer.Transaction
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
	if errHappen != nil {
		return errHappen
	}

	if lastBlock.BlockEntity.Capicity <= len(lastBlock.Transactions) {
		//Mine Block
		logger.Log("Start  Mine")
		stake := Stake.StakeCoinBlockChain{
			Block:   lastBlock,
			Workers: service.Workers}
		err := stake.Construct()
		if err != nil {
			logger.Fatal(err.Error())
		}
		service.Services.LotteryService.LoadStakeService(&stake)
		luckyOne, err := service.Services.LotteryService.Spin(1.5)
		if err != nil {
			logger.Fatal(err.Error())
		}
		if service.Services.WalletServiceInstance.GetPub() == luckyOne {
			//Win And Miner
			block := entitys.BlockCoinEntity{}
			err := block.MineBlock(luckyOne, lastBlock.BlockEntity, service.Services.LoggerService, service.Services.HashService)
			if err != nil {
				logger.Fatal(err.Error())
			}
			broadCastBlock := entitys.BlockCoinMessageRabbitMq{
				BlockCoin: block,
			}
			service.Services.RabbitMqService.PublishBlock(broadCastBlock)

		}
		block := service.Services.RabbitMqService.ConsumeNextBlock()
		validator = block.BlockCoin.BlockEntity.Validator
		//NOW The re is  litle  a chance to fail Mine only if  internal error if err =>  commit  harakiri
		err = service.Chain.InsertNewBlock(service.Services.LoggerService, service.Services.HashService, block.BlockCoin)
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Log("Commit  Mine")
	}
	//Re Sign  and  send
	From := transactions[0].BillDetails.Bill.From.Address
	//Stamp Validatori
	transactions[0].BillDetails.Bill.To.Address = validator
	service.Chain.InsertTransactions(transactions)
	if service.Services.WalletServiceInstance.GetPub() == From {
		amount := transactions[0].Amount + transactions[1].Amount
		fmt.Println(amount)
		err := service.Services.WalletServiceInstance.UnFreeze(amount)
		if err != nil {
			logger.Fatal(err.Error())
		}

	}
	logger.Log("UnLock  BlockChain")
	logger.Log(fmt.Sprintf("commit insert  transactions SET %s  and %s", t.Tax.Transaction.BillDetails.Transaction_id, t.Transfer.Transaction.BillDetails.Transaction_id))
	return nil
}
func (service *BlockChainCoinsImpl) Mine() error {
	/*
		latest := service.chain[len(service.chain)-1].b
		err, nextHash := service.Services.hashService.Hash(latest.purrentHash, latest.currentHash)
		if err != nil {
			return err
		}
		block := blockCoin{b: block{}, Transactions: make([]entitys.TransactionCoins, 0)}

		block.b.mineBlock(len(service.chain), *service.Services.walletService.GetPub(), latest.currentHash, nextHash)
		err = service.chain.insertNewBlock(service.Services.loggerService, service.Services.hashService, block)
		if err != nil {
			return err
		}*/
	return nil

}
func (service *BlockChainCoinsImpl) InsertCoinBlockMine(block entitys.BlockCoinEntity) error {

	err := service.Chain.InsertNewBlock(service.Services.LoggerService, service.Services.HashService, block)
	if err != nil {
		return err
	}
	return nil

}
