package balanceResolver

import (
	"services"

	"github.com/graphql-go/graphql"
)

func defineBlanceModel() *graphql.Object {
	model := graphql.NewObject(graphql.ObjectConfig{
		Name:        "FindBalanceModel",
		Description: "A Model to  fo findBalance 'coins' of a node ",
		Fields: graphql.Fields{
			"availableBalance": &graphql.Field{
				Type:        graphql.Float,
				Description: "node  coins",
			},
		},
	})
	return model
}
func resolver(p graphql.ResolveParams) (interface{}, error) {
	Node, ok := p.Args["Node"].(int)
	key := services.WalletService.GetPub()
	var err error
	if ok {
		key, err = services.SystemInfoService.Who(Node)
		if err != nil {
			return nil, err
		}
	}
	amount := services.FindBalanceService.FindBalance(key)
	return map[string]interface{}{
		"availableBalance": amount,
	}, nil
}
func GetField() *graphql.Field {
	model := defineBlanceModel()
	field := &graphql.Field{
		Name:        "FindBalanceField",
		Description: "A  Field to Find Balance  of a node ",
		Type:        model,
		Args: graphql.FieldConfigArgument{
			"Node": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "The balance  of node  your looking",
			},
		},
		Resolve: resolver,
	}
	return field

}
