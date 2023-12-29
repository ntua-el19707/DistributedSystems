package main 

import (
    "net/http"
    "fmt"
     "crypto/rsa"
    "services"

 )
//-- HEALTH --
/**
    healthV1 - test healthiness of process 
    @Param  w http.ResponseWriter
    @Param  r * http.Request 
*/
func  healthV1(w  http.ResponseWriter  ,  r * http.Request ){
    switch r.Method {
        case http.MethodGet :
          jsonBuilder(w ,http.StatusOK ,  struct {}{})
        default:
        //methods not  implemented
        message :=   fmt.Sprintf(httpErrorResponseNotImplemented , r.Method ,  r.URL.Path)
        jsonErrorBuilder(w ,  http.StatusMethodNotAllowed ,  message)
    }
}

/**
    TransferMoneyV1 - test healthiness of process 
    @Param  w http.ResponseWriter
    @Param  r * http.Request 
*/
func  TransferMoneyV1(w  http.ResponseWriter  ,  r * http.Request ){
    switch r.Method {
        case http.MethodGet :
             to :=  rsa.PublicKey{}
            list ,  err := services.TransactionManagerInstanceService.TransferMoney(&to , &to , float64(10))
            if  err != nil {
                   jsonErrorBuilder(w ,  http.StatusBadRequest ,  err.Error())
                return 
            }
            jsonBuilder(w ,http.StatusOK ,  list ) 
        default:
        //methods not  implemented
        message :=   fmt.Sprintf(httpErrorResponseNotImplemented , r.Method ,  r.URL.Path)
        jsonErrorBuilder(w ,  http.StatusMethodNotAllowed ,  message)
    }
}

/**
    SendMsgV1 - test healthiness of process 
    @Param  w http.ResponseWriter
    @Param  r * http.Request 
*/
func  SendMsgV1(w  http.ResponseWriter  ,  r * http.Request ){
    switch r.Method {
        case http.MethodGet :
             to :=  rsa.PublicKey{}
            list ,  err := services.TransactionManagerInstanceService.SendMessage(&to , &to , &to, "Paries")
            if  err != nil {
                   jsonErrorBuilder(w ,  http.StatusBadRequest ,  err.Error())
                return 
            }
            jsonBuilder(w ,http.StatusOK ,  list ) 
        default:
        //methods not  implemented
        message :=   fmt.Sprintf(httpErrorResponseNotImplemented , r.Method ,  r.URL.Path)
        jsonErrorBuilder(w ,  http.StatusMethodNotAllowed ,  message)
    }
}