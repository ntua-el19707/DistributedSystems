package entitys

import "crypto/rsa"

type TransactionCoinSet struct {
	Tax      TransactionCoinEntityRoot `json:"tax"`
	Transfer TransactionCoinEntityRoot `json:"transfer"`
}

type TransactionMessageSet struct {
	TransactionCoin    TransactionCoinSet       `json:"transaction_money"`
	TransactionMessage TransactionMsgEntityRoot `json:"transaction_message"`
}
type BlockCoinMessageRabbitMq struct {
	BlockCoin BlockCoinEntity `json:"block"`
	//TODO ?  may sign the  block Signiture []byte          `json:"signiture"`
}
type BlockMessageMessageRabbitMq struct {
	BlockMsg BlockMessage `json:"block"`
	//TODO ?  may sign the  block Signiture []byte          `json:"signiture"`
}
type ClientInfo struct {
	Id        string `json:"nodeId"`
	IndexId   int    `json:"indexId"`
	Uri       string `json:"uri"`
	UriPublic string `json:"uriPublic"`
}
type RabbitMqSystemInfoPack struct {
	Clients         []ClientRequestBody `json:"Clients"`
	ExpectedWorkers int                 `json:"expectedWorkers"`
	ScaleFactorMsg  float64             `json:"ScaleFactorMsg"`
	ScaleFactorCoin float64             `json:"ScaleFactorCoin"`
}
type ClientRequestBody struct {
	Client    ClientInfo    `json:"clientInfo"`
	PublicKey rsa.PublicKey `json:"key"`
}
