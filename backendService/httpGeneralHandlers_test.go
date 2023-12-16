package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//-- TEST FOR  defaultEmptyHttpController --
func  TestForDefaultEmptyHttpControllerStatusNotFoundGet(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/a", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusNotFound{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusNotFound ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  Resource '/a' is  not Found"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
func  TestForDefaultEmptyHttpControllerStatusNotFoundPost(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/asdoipjdid/aaskasaso", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusNotFound{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusNotFound ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  Resource '/asdoipjdid/aaskasaso' is  not Found"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}

func  TestForDefaultEmptyHttpControllerStatusNotFoundPOSTt(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusMethodNotAllowed{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusMethodNotAllowed ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  method POST is  not  ipplemented in route '/'"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
func  TestForDefaultEmptyHttpControllerStatusNotAllowedGet(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusMethodNotAllowed{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusMethodNotAllowed ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  method GET is  not  ipplemented in route '/'"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
func  TestForDefaultEmptyHttpControllerStatusNotAllowedPut(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusMethodNotAllowed{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusMethodNotAllowed ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  method PUT is  not  ipplemented in route '/'"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
func  TestForDefaultEmptyHttpControllerStatusNotAllowedPatch(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPatch, "/", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusMethodNotAllowed{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusMethodNotAllowed ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  method PATCH is  not  ipplemented in route '/'"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
func  TestForDefaultEmptyHttpControllerStatusNotAllowedDelete(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodDelete, "/", nil)
    
    defaultEmptyHttpController(w ,req)
    //check code response
    if w.Code !=  http.StatusMethodNotAllowed{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusMethodNotAllowed ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"The  method DELETE is  not  ipplemented in route '/'"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
//-- TEST FOR  jsonErrorBuilder --
func TestForjsonErrorBuilderBadRequest ( t * testing.T){
    //create  http response
    w :=  httptest.NewRecorder()
    code := int16(http.StatusBadRequest)
    message  :=  "bad request"
    jsonErrorBuilder(w ,code , message)
    if w.Code !=  http.StatusBadRequest{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusBadRequest ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Message":"bad request"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }
}
//-- TEST FOR JSON BUILDER --
func  TestJsonBuilderOK(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    code := int16(http.StatusOK)
    type  OkStruct struct  {
        Msg  string `Msg:"string"`
    }
    payload := OkStruct{Msg:"ok",}

    jsonBuilder(w , code ,payload)
    //check code response
    if w.Code !=  http.StatusOK{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusOK ,w.Code)
    }
    
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf("Expected  content type to be   %s  but  get %s" ,  expectedContentType ,contentType)
    }

    //check body 
    response :=  w.Body.String()
    const expectedResponse = `{"Msg":"ok"}`
    if  response != expectedResponse{
       t.Errorf("Expected  response to be   %s  but  get %s" ,  expectedResponse ,response)
    }

}
func  TestJsonBuilderStatusFound(t * testing.T){
    // create mock HTTP response  writer
    w :=  httptest.NewRecorder()
    code := int16(http.StatusFound)
    type  OkStruct struct  {
        Msg  string `Msg:"string"`
    }
    payload := OkStruct{Msg:"",}

    jsonBuilder(w , code ,payload)
    //check code response
    if w.Code !=  http.StatusFound{
       t.Errorf("Expected  status  code %d  but  get %d" ,  http.StatusOK ,w.Code)
    }
}

