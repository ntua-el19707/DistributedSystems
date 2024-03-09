package TransactionCoinResolver

import (
	"crypto/rsa"
	"services"

	"github.com/graphql-go/graphql"
)

//Define Model

func defineModel(name string) *graphql.Object {
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

	// Define TransactionCoinRsp object type
	transactionCoinRspType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TransactionCoinRsp" + name,
		Fields: graphql.Fields{
			"Transactions": &graphql.Field{
				Type:        graphql.NewList(transactionCoinRowType),
				Description: "List of transactions",
			},
			/*"NodeDetails": &graphql.Field{
				Type:        graphql.Int,
				Description: "Node details",
			},*/
		},
	})

	return transactionCoinRspType
}

func NodeTransactions() *graphql.Field {
	model := defineModel("nodetransactions")
	transactionsField := &graphql.Field{
		Type:        model,
		Description: "Get all transactions",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			transactions := services.FindBalanceService.GetTransactions([]rsa.PublicKey{services.WalletService.GetPub()}, []int64{})
			return map[string]interface{}{
				"Transactions": transactions,
			}, nil
		},
	}
	return transactionsField
}

//Defifine  resolver

func Query() *graphql.Field {
	model := defineModel("allTransactions")
	transactionsField := &graphql.Field{
		Type:        model,
		Description: "Get all transactions",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			transactions := services.FindBalanceService.GetTransactions([]rsa.PublicKey{}, []int64{})
			return map[string]interface{}{
				"Transactions": transactions,
			}, nil
		},
	}

	return transactionsField
}
