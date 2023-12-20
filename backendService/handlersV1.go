package main 

import (
    "net/http"
    "fmt"
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
