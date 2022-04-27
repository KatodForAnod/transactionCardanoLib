package cardanocli

type BaseTransactionParams struct {
	PaymentAddr string //1
	PolicyID    string //2
	ID          string //3
}

type TransactionParams struct {
	TxHash         string
	Txix           string
	Funds          string
	Fee            string
	Output         string
	Receiver       string
	ReceiverOutput string
}
