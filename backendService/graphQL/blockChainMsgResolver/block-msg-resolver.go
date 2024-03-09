package blockChainMsgResolver

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
	modelRow := graphql.NewObject(graphql.ObjectConfig{
		Name: name + "Row",
		Fields: graphql.Fields{
			"From": &graphql.Field{
				Type:        graphql.Int,
				Description: "From field  who send ",
			},
			"To": &graphql.Field{
				Type:        graphql.Int,
				Description: "To field who receved",
			},
			"Msg": &graphql.Field{
				Type:        graphql.String,
				Description: "message",
			},
			"Nonce": &graphql.Field{
				Type:        graphql.Int,
				Description: "nonce of transaction ",
			},
			"Time": &graphql.Field{
				Type:        graphql.Int,
				Description: "time  the  trasanction is created send  time",
			},
			"TransactionId": &graphql.Field{
				Type:        graphql.String,
				Description: "Transaction Id ",
			},
		},
	})

	return modelRow
}
func DefineBlockMsg(name string) *graphql.Object {
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
func DefineChainMsgResolver() *graphql.Field {
	model := DefineBlockMsg("BlockChainMsgResolver")
	field := &graphql.Field{
		Name:        "BlockChainMsgListField",
		Type:        graphql.NewList(model),
		Description: "Get the block chain msg",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			list, err := services.InboxService.GetBlockChain()
			return list, err
		},
	}
	return field
}
