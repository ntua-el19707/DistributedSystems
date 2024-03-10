package main

import (
	"Inbox"
	"WalletAndTransactions"
	"crypto/rsa"
	"encoding/json"
	"entitys"
	"fmt"
	"log"
	"net/http"
	"services"
)

//-- HEALTH --
/**
  healthV1 - test healthiness of process
  @Param  w http.ResponseWriter
  @Param  r * http.Request
*/
func healthV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonBuilder(w, http.StatusOK, struct{}{})
	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

//-- node details --
/**
  nodeDetailsV1 - test healthiness of process
  @Param  w http.ResponseWriter
  @Param  r * http.Request
*/
func nodeDetailsV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		type rsp struct {
			Client entitys.ClientInfo `json:"client"`
			Total  int                `json:"total"`
		}
		var response rsp
		response.Client, response.Total = services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)
	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

//-- node details --
/**
clientsV1 - test healthiness of process
 @Param  w http.ResponseWriter
 @Param  r * http.Request
*/
func clientsV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		type rsp struct {
			Clients []entitys.ClientInfo `json:"clients"`
			Total   int                  `json:"total"`
		}
		var response rsp
		response.Clients, response.Total = services.SystemInfoService.ClientList(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)
	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

//-- Register --
/**
  registerV1 - registerNode
  @Param  w http.ResponseWriter
  @Param  r * http.Request
*/
func registerV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		d := json.NewDecoder(r.Body)
		var payload entitys.ClientRequestBody
		err := d.Decode(&payload)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		err = services.SystemInfoService.AddWorker(payload)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		if services.SystemInfoService.IsFull() {
			err = services.SystemInfoService.BroadCast(sFm, sFc)
			if err != nil {
				jsonErrorBuilder(w, http.StatusInternalServerError, err.Error())
				log.Fatal(err.Error())
			}

			err, _, _ = services.SystemInfoService.Consume()
			if err != nil {
				jsonErrorBuilder(w, http.StatusInternalServerError, err.Error())
				log.Fatal(err.Error())
			}
			services.SetUp()

		}
		jsonBuilder(w, http.StatusOK, struct{}{})
	default:
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

/*
*

	TransferMoneyV1 - test healthiness of process
	@Param  w http.ResponseWriter
	@Param  r * http.Request
*/
func TransferMoneyV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		me, total := services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		code, err, body := ParseRequest[TranferMoney](r, me.IndexId, total)
		if err != nil {
			jsonErrorBuilder(w, code, err.Error())
			return
		}
		to, err := services.SystemInfoService.Who(body.To)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		list, err := services.TransactionManagerService.TransferMoney(to, body.Amount)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}

		err = services.RabbitMqS.PublishTractioncoinSet(list)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonBuilder(w, http.StatusOK, list)
	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

/*
*

	sendAndReceicedV1 - get all  send  and  recieved
	@Param  w http.ResponseWriter
	@Param  r * http.Request
*/
func sendAndReceicedV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		keys := []rsa.PublicKey{services.WalletService.GetPub()}
		var times []int64
		err, list := services.InboxService.SendAndReceived(keys, times)

		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonBuilder(w, http.StatusOK, list)
	default:
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

/*
*

	allMsgV1 - test healthiness of process
	@Param  w http.ResponseWriter
	@Param  r * http.Request
*/
func allMsgV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var times []int64
		err, list := services.InboxService.All(times)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		type rsp struct {
			Transactions Inbox.Inbox        `json:"transactions"`
			NodeDetails  entitys.ClientInfo `json:"nodeDetails"`
		}
		var response rsp
		response.Transactions = list
		response.NodeDetails, _ = services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)
	default:
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

/*
*

	stakeController  -  acontroller  to set stake  change  from initial
	@Param w  http.ResponseWriter
	@Param r  * http.Request
*/
func stakeController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var stakeR SetStakeRequest
		errcode, err := stakeR.Parse(r)
		if err != nil {
			jsonErrorBuilder(w, errcode, err.Error())
			return
		}
		services.FindBalanceService.SetStake(stakeR.Stake)
		jsonBuilder(w, http.StatusOK, struct{}{})
	default:
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)

	}
}

/*
*

	inboxV1 - test healthiness of process
	@Param  w http.ResponseWriter
	@Param  r * http.Request
*/
func inboxV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		keys := []rsa.PublicKey{services.WalletService.GetPub()}
		var times []int64
		err, list := services.InboxService.Receive(keys, times)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		type rsp struct {
			Transactions Inbox.Inbox        `json:"transactions"`
			NodeDetails  entitys.ClientInfo `json:"nodeDetails"`
		}
		var response rsp
		response.Transactions = list
		response.NodeDetails, _ = services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)
	default:
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}

/*
*

	SendMsgV1 - test healthiness of process
	@Param  w http.ResponseWriter
	@Param  r * http.Request
*/
func SendMsgV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		keys := []rsa.PublicKey{services.WalletService.GetPub()}
		var times []int64
		err, list := services.InboxService.Send(keys, times)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		type rsp struct {
			Transactions Inbox.Inbox        `json:"transactions"`
			NodeDetails  entitys.ClientInfo `json:"nodeDetails"`
		}
		var response rsp
		response.Transactions = list
		response.NodeDetails, _ = services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)

	case http.MethodPost:
		me, total := services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		code, err, body := ParseRequest[TranferMsg](r, me.IndexId, total)
		if err != nil {
			jsonErrorBuilder(w, code, err.Error())
			return
		}
		to, err := services.SystemInfoService.Who(body.To)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		list, err := services.TransactionManagerService.SendMessage(to, body.Msg)
		if err != nil {
			jsonErrorBuilder(w, http.StatusBadRequest, err.Error())
			return
		}
		err = services.RabbitMqS.PublishTractionMsgSet(list)
		jsonBuilder(w, http.StatusOK, struct{}{})
	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}
func TransactionsV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		transactions := services.FindBalanceService.GetTransactions([]rsa.PublicKey{services.WalletService.GetPub()}, []int64{})
		type rsp struct {
			Transactions WalletAndTransactions.TransactionListCoin `json:"transactions"`
			NodeDetails  entitys.ClientInfo                        `json:"nodeDetails"`
		}
		var response rsp
		response.Transactions = transactions
		response.NodeDetails, _ = services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)

	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}
func TransactionsAllV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		transactions := services.FindBalanceService.GetTransactions([]rsa.PublicKey{}, []int64{})
		type rsp struct {
			Transactions WalletAndTransactions.TransactionListCoin `json:"transactions"`
			NodeDetails  entitys.ClientInfo                        `json:"nodeDetails"`
		}
		var response rsp
		response.Transactions = transactions
		response.NodeDetails, _ = services.SystemInfoService.NodeDetails(services.WalletService.GetPub())
		jsonBuilder(w, http.StatusOK, response)

	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}
func balanceV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		amount := services.FindBalanceService.FindBalance(services.WalletService.GetPub())
		/*balance, err := services.BlockChainCoinsService.FindBalance()
		if err != nil {
			jsonErrorBuilder(w, http.StatusInternalServerError, err.Error())
			return
		}*/
		type rsp struct {
			Balance float64 `json:"availableBalance"`
		}
		jsonBuilder(w, http.StatusOK, rsp{Balance: amount})

	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}
