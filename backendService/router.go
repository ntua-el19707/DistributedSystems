package main

import (
	"encoding/json"
	"fmt"
	"graphQL"
	"log"
	"net/http"
	"services"

	"github.com/graphql-go/graphql"
)

type graphQlHandler struct {
	Schema *graphql.Schema
}

// ServeHTTP method to handle GraphQL requests
func (h *graphQlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	params := graphql.Params{
		Schema:        *h.Schema,
		RequestString: query,
	}
	result := graphql.Do(params)

	// Write the GraphQL response to the HTTP response writer
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// message  that  will display to see  the availble  links  of server in log
const setUpRouteMessage = "The route :%s is now  available\n"
const connectRouterMessage = "The  router '%s' is  connected to '%s' router  "

/*
*

	setUpMainRouter - set ups  the  router for  the api
	@Param s * http.ServeMux
*/
func setUpMainRouter(s *http.ServeMux, c bool) {
	log.Printf(setUpRouteMessage, "/")
	//s.HandleFunc("/", defaultEmptyHttpController)
	staticDir := "./staticServer/browser"
	fs := http.FileServer(http.Dir(staticDir))
	s.Handle("/", http.StripPrefix("/", fs))

	schema := graphQL.SetUp()
	s.Handle("/graphql", &graphQlHandler{Schema: schema})
	//set api
	api := http.NewServeMux()
	prefix := "/api"
	setRouterForApi(api, prefix, c)
	s.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, api))
	log.Printf(connectRouterMessage, prefix, "/")

}

/*
*

	setUpRouterForApi - set ups  the  router for  the api
	@Param s * http.ServeMux
	@Param prefix string
*/
func setRouterForApi(api *http.ServeMux, prefix string, c bool) {
	log.Printf(setUpRouteMessage, prefix)
	api.HandleFunc("/", defaultEmptyHttpController)
	//set api
	v1 := http.NewServeMux()
	prefixV1 := "/v1"
	setRouterForV1(v1, fmt.Sprintf("%s%s", prefix, prefixV1), c)
	api.Handle(fmt.Sprintf("%s/", prefixV1), http.StripPrefix(prefixV1, v1))
	log.Printf(connectRouterMessage, prefix, prefixV1)

}

/*
*

	setUpRouterForV1- set up   the  router for  version 1
	@Param s * http.ServeMux
	@Param prefixV1
*/
func setRouterForV1(v1 *http.ServeMux, prefixV1 string, c bool) {

	log.Printf(setUpRouteMessage, prefixV1)
	v1.HandleFunc("/", defaultEmptyHttpController)
	//Set up routes
	if c {
		log.Printf(setUpRouteMessage, fmt.Sprintf("%s/register", prefixV1))
		v1.HandleFunc("/register", SystemNotInitilize(registerV1))
	}
	//-- Health --
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/health", prefixV1))
	v1.HandleFunc("/health", healthV1)
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/stake", prefixV1))
	v1.HandleFunc("/stake", SystemInitilize(stakeController))

	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/transfer", prefixV1))
	v1.HandleFunc("/transfer", SystemInitilize(TransferMoneyV1))
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/send", prefixV1))
	v1.HandleFunc("/send", SystemInitilize(SendMsgV1))
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/NodeDetails", prefixV1))
	v1.HandleFunc("/NodeDetails", SystemInitilize(nodeDetailsV1))
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/NodeDetails/Clients", prefixV1))
	v1.HandleFunc("/NodeDetails/Clients", SystemInitilize(clientsV1))
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/inbox", prefixV1))
	v1.HandleFunc("/inbox", SystemInitilize(inboxV1))
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/allMsg", prefixV1))
	v1.HandleFunc("/allMsg", SystemInitilize(allMsgV1))

	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/inboxAndSent", prefixV1))
	v1.HandleFunc("/inboxAndSent", SystemInitilize(sendAndReceicedV1))
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/balance", prefixV1))
	v1.HandleFunc("/balance", balanceV1)
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/transactions", prefixV1))
	v1.HandleFunc("/transactions", TransactionsV1)
	log.Printf(setUpRouteMessage, fmt.Sprintf("%s/transactionsAll", prefixV1))
	v1.HandleFunc("/transactionsAll", TransactionsAllV1)
}

func SystemInitilize(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !services.SystemInfoService.IsOk() {
			jsonErrorBuilder(w, http.StatusBadRequest, "Workers  are not connected still waiting")
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
func SystemNotInitilize(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if services.SystemInfoService.IsOk() {
			jsonErrorBuilder(w, http.StatusBadRequest, "System is  initialised")
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
