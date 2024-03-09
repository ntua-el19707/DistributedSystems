package TransactionMsgResolver

import (
	"crypto/rsa"
	"services"

	"github.com/graphql-go/graphql"
)

// define  model
func GetModel(name string) *graphql.Object {
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
	model := graphql.NewObject(graphql.ObjectConfig{
		Name: name + "List",
		Fields: graphql.Fields{
			"transactions": &graphql.Field{
				Type:        graphql.NewList(modelRow),
				Description: "transaction msg list",
			},
		},
	})
	return model

}

func GetInboxField() *graphql.Field {
	//get model
	model := GetModel("TransactionMsgInbox")
	field := &graphql.Field{
		Name:        "Inbox",
		Type:        model,
		Description: "Message  that te reciever is  this node",
		Args: graphql.FieldConfigArgument{
			"From": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Who  send the message",
			},
			"TimeAfter": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Transactions  After",
			},
			"TimeBefore": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Transactions  After",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var keys []rsa.PublicKey
			fromIndex, ok := p.Args["From"].(int)
			if ok {
				key, err := services.SystemInfoService.Who(fromIndex)
				if err != nil {
					return nil, err
				}
				keys = append(keys, key)
			}
			keys = append(keys, services.WalletService.GetPub())
			var times []int64
			timeAfter, ok := p.Args["TimeAfter"].(int)

			if ok {
				times = append(times, int64(timeAfter))
			}
			timeBefore, ok := p.Args["TimeBefore"].(int)
			if ok {
				if len(times) == 0 {
					times = append(times, 0)
				}
				times = append(times, int64(timeBefore))
			}
			err, list := services.InboxService.Receive(keys, times)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"transactions": list}, nil
		},
	}
	return field
}
func GetSendField() *graphql.Field {
	//get model
	model := GetModel("TransactionSendMsg")
	field := &graphql.Field{
		Name:        "Send",
		Type:        model,
		Description: "Message  that this node send",
		Args: graphql.FieldConfigArgument{
			"To": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Who is  the reciever the message",
			},
			"TimeAfter": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Transactions  After",
			},
			"TimeBefore": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Transactions  After",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			keys := []rsa.PublicKey{services.WalletService.GetPub()}
			toIndex, ok := p.Args["To"].(int)
			if ok {
				key, err := services.SystemInfoService.Who(toIndex)
				if err != nil {
					return nil, err
				}
				keys = append(keys, key)
			}
			var times []int64
			timeAfter, ok := p.Args["TimeAfter"].(int)

			if ok {
				times = append(times, int64(timeAfter))
			}
			timeBefore, ok := p.Args["TimeBefore"].(int)
			if ok {
				if len(times) == 0 {
					times = append(times, 0)
				}
				times = append(times, int64(timeBefore))
			}
			err, list := services.InboxService.Send(keys, times)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"transactions": list}, nil
		},
	}
	return field
}
func GetAllField() *graphql.Field {
	//get model
	model := GetModel("transactionMsgAll")
	field := &graphql.Field{
		Name:        "AllMessages",
		Type:        model,
		Description: "all messages in  blockchain",
		Args: graphql.FieldConfigArgument{
			"TimeAfter": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Transactions  After",
			},
			"TimeBefore": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Transactions  After",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var times []int64
			timeAfter, ok := p.Args["TimeAfter"].(int)

			if ok {
				times = append(times, int64(timeAfter))
			}
			timeBefore, ok := p.Args["TimeBefore"].(int)
			if ok {
				if len(times) == 0 {
					times = append(times, 0)
				}
				times = append(times, int64(timeBefore))
			}
			err, list := services.InboxService.All(times)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"transactions": list}, nil
		},
	}
	return field
}
func GetField() *graphql.Field {
	fieldInbox := GetInboxField()
	fieldSend := GetSendField()
	fieldAll := GetAllField()
	model := graphql.NewObject(graphql.ObjectConfig{
		Name:        "TransactionMessageModel",
		Description: "Transaction Message Model",
		Fields: graphql.Fields{
			"inbox": &graphql.Field{
				Type:        fieldInbox.Type,
				Description: "Get messages received by this node",
				Args:        fieldInbox.Args,
				Resolve:     fieldInbox.Resolve,
			},
			"send": &graphql.Field{
				Type:        fieldSend.Type,
				Description: "Get messages sent by this node",
				Args:        fieldSend.Args,
				Resolve:     fieldSend.Resolve,
			},
			"allMessages": &graphql.Field{
				Type:        fieldAll.Type,
				Description: "Get all messages in the blockchain",
				Args:        fieldAll.Args,
				Resolve:     fieldAll.Resolve,
			},
		},
	})
	field := &graphql.Field{
		Type:        model,
		Description: "Transaction Message Model",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	}
	return field
}
