package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"services"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func Fall(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
func GetEnviroments() (int, int, bool, string, string, string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could  not  load  enviroment  varibales  due to %v", err)
	}
	serverPort, err := strconv.Atoi(os.Getenv("serverPort"))
	Fall(err)
	var workers int
	coordinatorS := os.Getenv("coordinator")
	var coordinator bool
	if strings.ToUpper(coordinatorS) == "TRUE" {
		coordinator = true
	}

	if coordinator {
		workers, err = strconv.Atoi(os.Getenv("workers"))
		Fall(err)
		if workers <= 0 {
			//cordinator is a worker
			err = errors.New("Must  have  at least 1  worker")
			Fall(err)
		}
	}
	host := os.Getenv("hostCoordinator")
	me := os.Getenv("myNetwork")

	portC, err := strconv.Atoi(os.Getenv("coordinatorPort"))
	Fall(err)
	host = fmt.Sprintf("http://%s:%d", host, portC)
	me = fmt.Sprintf("http://%s:%d", me, serverPort)
	node := os.Getenv("nodeId")
	rabbitMqUri := os.Getenv("rabbitMQ")
	publicUri := os.Getenv("publicUri")
	return serverPort, workers, coordinator, host, me, node, rabbitMqUri, publicUri

}
func GetScaleFactors() (float64, float64) {
	scaleFactorMsg, err := strconv.ParseFloat(os.Getenv("scaleFactorMsg"), 64)
	Fall(err)
	scaleFactorCoin, err := strconv.ParseFloat(os.Getenv("scaleFactorCoin"), 64)
	Fall(err)
	return scaleFactorMsg, scaleFactorCoin

}
func GetCapicitys() (int, int) {
	capicityBlockMsg, err := strconv.Atoi(os.Getenv("CapicityBlockMsg"))
	Fall(err)
	capicityBlockCoin, err := strconv.Atoi(os.Getenv("CapicityBlockCoin"))
	Fall(err)
	if capicityBlockCoin%2 == 1 {
		err = errors.New("Cappicity for  block Coin must be  even ")
	}
	Fall(err)
	return capicityBlockMsg, capicityBlockCoin
}

var sFm, sFc float64

/*
*

	main - function  of  the project START
*/
func main() {
	port, workers, coordinator, curi, muri, id, rabbitMqUri, publicUri := GetEnviroments()
	var perNode float64
	var capicityMsg, capicityCoin int
	if coordinator {
		sFm, sFc = GetScaleFactors()
		capicityMsg, capicityCoin = GetCapicitys()
		var err error
		perNode, err = strconv.ParseFloat(os.Getenv("perNode"), 64)
		Fall(err)
	}
	serverPort := fmt.Sprintf(":%d", port)
	services.BootOrDie(id, curi, muri, rabbitMqUri, publicUri, coordinator, workers, capicityMsg, capicityCoin, sFm, sFc, perNode)
	setUpServer(serverPort, coordinator)
}

/*
*

	setUpServer() - set the server up
	@Param  serverPort string
*/
func setUpServer(serverPort string, c bool) {
	server := &http.ServeMux{}

	//*  SET THE  ROUTER FUNCTION
	setUpMainRouter(server, c)
	log.Printf("Server  is  Listening  on  Port %s...\n", serverPort)
	err := http.ListenAndServe(serverPort, server)
	if err != nil {
		log.Fatalf("Server  has  fallen due  to  %v ", err)
	}

}
