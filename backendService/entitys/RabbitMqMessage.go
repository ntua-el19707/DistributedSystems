package entitys

type TransactionCoinSet struct {
	Tax      TransactionCoinEntityRoot `json:"tax"`
	Transfer TransactionCoinEntityRoot `json:"transfer"`
}

type TransactionMessageSet struct {
	TransactionCoin    TransactionCoinSet       `json:"transaction_money"`
	TransactionMessage TransactionMsgEntityRoot `json:"transaction_message"`
}
type BlockCoinMessage struct {
	Index int `json:"index"`
}
