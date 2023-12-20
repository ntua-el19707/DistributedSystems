package  main 

import (
    "log"
    "net/http"
    "os"
    "github.com/joho/godotenv"
)

/**
    main - function  of  the project START
*/
func  main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Could  not  load  enviroment  varibales  due to %v" ,err)
    }
    serverPort :=  os.Getenv("serverPort")
    setUpServer(serverPort)
}

/**
    setUpServer() - set the server up 
    @Param  serverPort string 
 */
 func  setUpServer(serverPort string){
     server := &http.ServeMux{}
     //*  SET THE  ROUTER FUNCTION 
     setUpMainRouter(server)
     log.Printf("Server  is  Listening  on  Port %s...\n" ,serverPort)
    err :=  http.ListenAndServe(serverPort , server )
     if err != nil {
        log.Fatalf("Server  has  fallen due  to  %v " ,err)
     }
    
 }
