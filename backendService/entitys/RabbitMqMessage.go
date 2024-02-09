package entitys

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
