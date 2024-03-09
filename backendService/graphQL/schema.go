package graphQL

import (
	"BlockCoinChainResolver"
	"Logger"
	"NodeDetails"
	TransactionCoinResolver "TransactionCoinsResolver"
	"TransactionMsgResolver"
	"balanceResolver"
	"blockChainMsgResolver"
	"log"

	"github.com/graphql-go/graphql"
)

func SetUp() *graphql.Schema {

	logger := Logger.Logger{ServiceName: "set-up-graph-ql"}
	err := logger.Construct()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Log("Start  loading  findBalance  Resolver")
	findBalanceField := balanceResolver.GetField()
	logger.Log("Commit  loading  findBalance  Resolver")

	logger.Log("Start  loading  blockChaincoin  Resolver")
	blockCoinChainField := BlockCoinChainResolver.DefineChainCoinResolver()
	logger.Log("Commit  blockchain coin resolver  Resolver")
	logger.Log("Start  loading  blockChainMsg  Resolver")
	blockMsgChainField := blockChainMsgResolver.DefineChainMsgResolver()
	logger.Log("Commit  blockchainMsg resolver  Resolver")

	logger.Log("Start  loading  Transaction Coin  Resolver")
	tranactionCoinField := TransactionCoinResolver.Query()
	logger.Log("Commit  loading  Transaction Coin  Resolver")

	logger.Log("Start  loading  Transaction Coin Node Resolver")
	tranactionCoinNodeField := TransactionCoinResolver.NodeTransactions()

	logger.Log("commit  loading  Transaction Coin Node Resolver")
	logger.Log("Start  loading  Node Info self Resolver")
	nodeFieldSelf := NodeDetails.GetNodeSelfField()
	logger.Log("Commit  loading  Node Info  Resolver")
	logger.Log("Start  loading  Node Info clients Resolver")
	nodeFieldClients := NodeDetails.GetNodeClientsField()
	logger.Log("Commit  loading  Node Info clients Resolver")
	logger.Log("Start  loading  Node Info all Resolver")
	nodeFieldAll := NodeDetails.GetAllField()
	logger.Log("Commit  loading  Node Info all Resolver")

	logger.Log("Start  loading  Transaction Msg Inbox  Resolver")
	transactionMsgInboxResolver := TransactionMsgResolver.GetInboxField()
	logger.Log("Commit  loading  Transaction Msg Resolver Inbox")
	logger.Log("Start  loading  Transaction Msg Send  Resolver")
	transactionMsgSendResolver := TransactionMsgResolver.GetSendField()
	logger.Log("Commit  loading  Transaction Msg Resolver Send")
	logger.Log("Start  loading  Transaction Msg all  Resolver")
	transactionMsgAllResolver := TransactionMsgResolver.GetAllField()
	logger.Log("Commit  loading  Transaction Msg all Resolver ")

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"getTransactionsCoins": tranactionCoinField,
			"nodeTransactions":     tranactionCoinNodeField,
			"allTransactionMsg":    transactionMsgAllResolver,
			"inbox":                transactionMsgInboxResolver,
			"send":                 transactionMsgSendResolver,
			"self":                 nodeFieldSelf,
			"clients":              nodeFieldClients,
			"allNodes":             nodeFieldAll,
			"blockChainCoins":      blockCoinChainField,
			"blockChainMsg":        blockMsgChainField,
			"balance":              findBalanceField,
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	fall(err, &logger)
	return &schema
}

func fall(err error, logger Logger.LoggerService) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}
