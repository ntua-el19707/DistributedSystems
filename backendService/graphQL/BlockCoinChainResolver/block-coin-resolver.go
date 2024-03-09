package BlockCoinChainResolver

import (
	"fmt"
	"services"

	"github.com/graphql-go/graphql"
)

// define BlockModel
func DefineBlock(name string) *graphql.Object {
	blockModel := graphql.NewObject(graphql.ObjectConfig{
		Name:        name,
		Description: "Model for the Block Standards",
		Fields: graphql.Fields{
			"index": &graphql.Field{
				Type:        graphql.Int,
				Description: "Block Index in chain",
			},
			"created_at": &graphql.Field{
				Type:        graphql.Int,
				Description: "The time that block is created",
			},
			"validator": &graphql.Field{
				Type:        graphql.Int,
				Description: "Validator of block in index",
			},
			"capacity": &graphql.Field{
				Type:        graphql.Int,
				Description: "Capacity of block",
			},
			"current_hash": &graphql.Field{
				Type:        graphql.String,
				Description: "Current hash of block",
			},
			"parrent_hash": &graphql.Field{
				Type:        graphql.String,
				Description: "Parent hash of block",
			},
		},
	})
	return blockModel
}
func defineTransactionModel(name string) *graphql.Object {
	transactionCoinRowType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TransactionCoinRow" + name,
		Fields: graphql.Fields{
			"From": &graphql.Field{
				Type: graphql.Int,
			},
			"To": &graphql.Field{
				Type: graphql.Int,
			},
			"Coins": &graphql.Field{
				Type: graphql.Float,
			},
			"Nonce": &graphql.Field{
				Type: graphql.Int,
			},
			"Reason": &graphql.Field{
				Type: graphql.String,
			},
			"Time": &graphql.Field{
				Type: graphql.Int,
			},
			"TransactionId": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	return transactionCoinRowType
}
func DefineBlockCoin(name string) *graphql.Object {
	blockModelName := fmt.Sprintf("%s%s", name, "blockDetails")
	blockModel := DefineBlock(blockModelName)
	transactionModelName := fmt.Sprintf("%s%s", name, "transactionModel")
	transactionModel := defineTransactionModel(transactionModelName)
	model := graphql.NewObject(graphql.ObjectConfig{
		Name:        name,
		Description: "model for block coin model ",
		Fields: graphql.Fields{
			"block": &graphql.Field{
				Type:        blockModel,
				Description: "block details",
			},
			"transactions": &graphql.Field{
				Type:        graphql.NewList(transactionModel),
				Description: "transaction list ",
			},
		},
	})
	return model
}
func DefineChainCoinResolver() *graphql.Field {
	model := DefineBlockCoin("BlockChainCoinResolver")
	field := &graphql.Field{
		Name:        "BlockChainCoinField",
		Type:        graphql.NewList(model),
		Description: "Get the vlock coin chain ",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			list := services.FindBalanceService.GetChain()
			return list, nil
		},
	}
	return field
}
