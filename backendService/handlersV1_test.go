package  main 
import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const codeErrorMessageTesting = "Expected  status  code %d  but  get %d"
const contentTypeErrorMessageTesting ="Expected  content type to be   %s  but  get %s"
//TEST for healthV1

func  TestForHealthV1Get(t * testing.T){

    request := httptest.NewRequest(http.MethodGet ,  "/" , nil )
    w :=  httptest.NewRecorder()

    healthV1(w, request )
    expectedCode := http.StatusOK
    if  w.Code  != expectedCode {
        t.Errorf(codeErrorMessageTesting , expectedCode , w.Code   )

    } 
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf(contentTypeErrorMessageTesting,  expectedContentType ,contentType)

   }
   }
func  TestForHealthV1Post(t * testing.T){

    request := httptest.NewRequest(http.MethodPost ,  "/" , nil )
    w :=  httptest.NewRecorder()

    healthV1(w, request )
    expectedCode := http.StatusMethodNotAllowed
    if  w.Code  != expectedCode {
        t.Errorf(codeErrorMessageTesting , expectedCode , w.Code   )

    } 
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf(contentTypeErrorMessageTesting,  expectedContentType ,contentType)

   }
   }
func  TestForHealthV1Put(t * testing.T){

    request := httptest.NewRequest(http.MethodPut ,  "/" , nil )
    w :=  httptest.NewRecorder()

    healthV1(w, request )
    expectedCode := http.StatusMethodNotAllowed
    if  w.Code  != expectedCode {
        t.Errorf(codeErrorMessageTesting , expectedCode , w.Code   )

    } 
    // prepare  to check if its is  json  response
    contentType := w.Header().Get("Content-Type")
	const expectedContentType = "application/json"
    if contentType !=  expectedContentType {
       t.Errorf(contentTypeErrorMessageTesting,  expectedContentType ,contentType)

   }
}
