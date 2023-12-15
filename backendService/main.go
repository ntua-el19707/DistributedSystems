package  main 

import (
    "log"
    "net/http"
)

/**
    main - function  of  the project START
*/
func  main() {
    setUpServer(":3000")
}

/**
    setUpServer() - set the server up 
    @Param  serverPort string 
 */
 func  setUpServer(serverPort string){
     server := &http.ServeMux{}
     //*  SET THE  ROUTER FUNCTION 

     log.Printf("Server  is  Listening  on  Port %s...\n" ,serverPort)
    err :=  http.ListenAndServe(serverPort , server )
     if err != nil {
        log.Fatalf("Server  has  fallen due  to  %v " ,err)
     }
    
 }
