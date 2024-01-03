package entitys 
import (
	"crypto/rsa"
	"encoding/json"
)
type  Client  struct  {
	Address rsa.PublicKey `json:"address"`
} 
type  BillingInfo  struct {
	From  Client `json:"from"`
	To  Client `json:"to"`
} 
type  TransactionDetails struct {
	Bill   BillingInfo `json:"bill"`
	Nonce  int   `json:"nonce"`
	Transaction_id string `json:"transaction_id"` 
	Created_at  int64 `json:"created_at"`
}

type  TransactionCoins  struct{
	BillDetails TransactionDetails `json:"transactionDetails"`
	Amount  float64	`json:"amount"`
	Reason  string 	`json:"reason"`
}  
type  TransactionMsg struct {
	BillDetails TransactionDetails `json:"transactionDetails"`
	Msg  string 	`json:"Msg"`
}

//* i will  create a specific method  for marshaling  in order  to define the marshal 

func  (t TransactionCoins ) JsonStringfy()  ([] byte , error) {
	data ,  err := json.Marshal(t) 
	if  err != nil {
		return nil ,err
	}
	return data, nil 
}  

func  (t TransactionMsg ) JsonStringfy()  ([] byte , error) {
	data ,  err := json.Marshal(t) 
	if  err != nil {
		return nil ,err
	}
	return data, nil 
}  


type  TransactionRecord  interface {
	 JsonStringfy()  ([] byte , error)
} 

type  TransactionTransport struct  {
	Record TransactionRecord  `json:"Transaction"`
	Signiture  []  byte  `json:"Signiture"`
}

func  JsonStringfy( i  TransactionRecord )([] byte , error){
	return  i.JsonStringfy()
}

