package  main 

import (
    "net/http"
    "log"
    "fmt"
)

// message  that  will display to see  the availble  links  of server in log
const setUpRouteMessage = "The route :%s is now  available\n"
const connectRouterMessage = "The  router '%s' is  connected to '%s' router  " 

/**
    setUpMainRouter - set ups  the  router for  the api 
    @Param s * http.ServeMux 
*/
func setUpMainRouter(s * http.ServeMux){
    log.Printf(setUpRouteMessage, "/")
    s.HandleFunc("/" , defaultEmptyHttpController)
    //set api 
    api := http.NewServeMux()
    prefix := "/api"
    setRouterForApi(api , prefix)
    s.Handle(fmt.Sprintf("%s/" , prefix) ,  http.StripPrefix(prefix ,  api))
    log.Printf( connectRouterMessage, prefix ,  "/")
     
}
/**
    setUpRouterForApi - set ups  the  router for  the api 
    @Param s * http.ServeMux 
    @Param prefix string
*/
func  setRouterForApi(api *  http.ServeMux , prefix string ){
    log.Printf(setUpRouteMessage, prefix)
    api.HandleFunc("/" , defaultEmptyHttpController)
    //set api 
    v1 := http.NewServeMux()
    prefixV1 := "/v1"
    setRouterForV1(v1 , fmt.Sprintf("%s%s" ,prefix , prefixV1))
    api.Handle(fmt.Sprintf("%s/" , prefixV1) ,  http.StripPrefix(prefixV1 ,  v1))
    log.Printf( connectRouterMessage , prefix ,  prefixV1)


}
/**
    setUpRouterForV1- set up   the  router for  version 1
    @Param s * http.ServeMux 
    @Param prefixV1
*/
func setRouterForV1(v1 * http.ServeMux, prefixV1 string){

    log.Printf(setUpRouteMessage, prefixV1)
    v1.HandleFunc("/" , defaultEmptyHttpController)
    //Set up routes 
    //-- Health --
    log.Printf(setUpRouteMessage, fmt.Sprintf("%s/health" , prefixV1))
    v1.HandleFunc("/health" , healthV1) 
    
    log.Printf(setUpRouteMessage, fmt.Sprintf("%s/transfer" , prefixV1))
    v1.HandleFunc("/transfer"  , TransferMoneyV1)
        log.Printf(setUpRouteMessage, fmt.Sprintf("%s/send" , prefixV1))
    v1.HandleFunc("/send"  , SendMsgV1)
}
