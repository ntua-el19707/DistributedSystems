package  main 

import (
    "net/http"
    "log"
    "fmt"
    "encoding/json"
)

//http standard  error  Messages 

const  httpErrorResposeNotFound = "The  Resource '%s' is  not Found"
const  httpErrorResponseNotImplemented = "The  method %s is  not  ipplemented in route '%s'"
const  httpServerInternalError = "SERVER INTERNAL ERROR:  %s\n"

/**
    defaultEmptyHttpController - set default  empty controllers for  route  path and 404 the child endpoints  
    @Param  w http.ResponseWriter
    @Param  r *  http.Request
*/
func defaultEmptyHttpController(w  http.ResponseWriter , r  *  http.Request ){
    url := r.URL.Path
    var message  string 
    if  url !=  "/" {
        message  =  fmt.Sprintf( httpErrorResposeNotFound ,  url)
        jsonErrorBuilder(w ,http.StatusNotFound , message)
        return 
    }
    message =   fmt.Sprintf(httpErrorResponseNotImplemented , r.Method ,  url)
    jsonErrorBuilder(w ,  http.StatusMethodNotAllowed ,  message)
}
/**
    jsonErrorBuilder - build  a standard error  response 
    @Param  w http.ResponseWriter
    @Param  code  int16 
    @Param message  string 
 */ 
func jsonErrorBuilder(w  http.ResponseWriter , code  int16 , message  string ){
    //check  if  internal  error and if => fall
    if  code > 499 {
        log.Fatalf(httpServerInternalError , message )  
    }
    //define  a  type  for response  
    type  ErrorStruct struct {
        Message  string  `json:errMsg`
    }
    //create  response  struct 
    messageStruct := ErrorStruct{
        Message: message,
    } 
    //build response
    jsonBuilder(w ,code, messageStruct) 
}

/** 
    jsonBuilder - build  json response 
    @Param w http.ResponseWriter
    @Param code  int16 
    @Param payload interface {}
*/
func  jsonBuilder(w  http.ResponseWriter ,  code  int16 ,  payload interface{}){
    //create  json  response
    data , err := json.Marshal(payload)
    if err != nil {
        log.Fatalf(httpServerInternalError , err.Error() )
        w.WriteHeader(http.StatusInternalServerError)
        return 
    }
    //set response
    w.Header().Add("Content-Type" ,"application/json")
    w.WriteHeader(int(code))
    w.Write(data)
}

