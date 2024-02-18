package main

import (
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
			err = services.SystemInfoService.BroadCast(2.0, 0.5)
			if err != nil {
				jsonErrorBuilder(w, http.StatusInternalServerError, err.Error())
				log.Fatal(err.Error())
			}

			err = services.SystemInfoService.Consume()
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
		jsonBuilder(w, http.StatusOK, list)
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
		jsonBuilder(w, http.StatusOK, list)

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

func balanceV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		/*balance, err := services.BlockChainCoinsService.FindBalance()
		if err != nil {
			jsonErrorBuilder(w, http.StatusInternalServerError, err.Error())
			return
		}
		type rsp struct {
			Balance float64 `json:"available balance"`
		}*/
		jsonBuilder(w, http.StatusOK, struct{}{})

	default:
		//methods not  implemented
		message := fmt.Sprintf(httpErrorResponseNotImplemented, r.Method, r.URL.Path)
		jsonErrorBuilder(w, http.StatusMethodNotAllowed, message)
	}
}
