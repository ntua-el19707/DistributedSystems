package NodeDetails

import (
	"services"

	"github.com/graphql-go/graphql"
)

// Define NodeDetailsModel
func NodeDetailsModel(name string) (*graphql.Object, *graphql.Object) {
	nodeModel := graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodeId": &graphql.Field{
				Type:        graphql.String,
				Description: "node string Id ",
			},
			"indexId": &graphql.Field{
				Type:        graphql.Int,
				Description: "node Index Id",
			},
			"uri": &graphql.Field{
				Type:        graphql.String,
				Description: "private  uri of  node accesible only by node",
			},
			"uriPublic": &graphql.Field{
				Type:        graphql.String,
				Description: "public uri of node  accseble  by public ",
			},
		},
	})
	nodeExtensiveModel := graphql.NewObject(graphql.ObjectConfig{
		Name: name + "self",
		Fields: graphql.Fields{
			"client": &graphql.Field{
				Type:        nodeModel,
				Description: "Node details",
			},
			"total": &graphql.Field{
				Type:        graphql.Int,
				Description: "total workers",
			},
		},
	})

	return nodeModel, nodeExtensiveModel
}

func GetNodeSelfField() *graphql.Field {

	_, model := NodeDetailsModel("nodeDetailsSelf")

	field := &graphql.Field{
		Type:        model,
		Description: "self  node details",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			node, total := services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
			return map[string]interface{}{
				"client": node,
				"total":  total,
			}, nil
		},
	}
	return field

}

func GetNodeClientsField() *graphql.Field {
	model, _ := NodeDetailsModel("nodeDetailsClient")
	field := &graphql.Field{
		Type:        graphql.NewList(model),
		Description: "clients list",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			clients, _ := services.SystemInfoService.ClientList(services.WalletService.GetPub())
			return clients, nil
		},
	}
	return field

}
func GetAllField() *graphql.Field {
	model, _ := NodeDetailsModel("nodeDetailsAll")
	field := &graphql.Field{
		Type:        graphql.NewList(model),
		Description: "all node list",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			nodes := services.SystemInfoService.Nodes()
			return nodes, nil
		},
	}
	return field
}
