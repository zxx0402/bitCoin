package someCompose

import "time"

type Transaction struct {
	Sender    string
	Receiver  string
	Num       int
	TimeStamp uint64
}

//创建交易信息

func NewTransaction(sender string, receiver string, num int) *Transaction {
	b := Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Num:       num,
		TimeStamp: uint64(time.Now().Unix()),
	}
	return &b
}

func Contains(slice []Transaction, item uint64) bool {
	for _, element := range slice {
		if element.TimeStamp == item {
			return true
		}
	}
	return false
}
func DeleteTransaction(LocalTrans []Transaction, BlockTrans []Transaction) []Transaction {
	var updatedTransactionPool []Transaction
	for _, tx := range LocalTrans {
		// 检查交易是否已被打包
		isPackaged := Contains(BlockTrans, tx.TimeStamp)

		// 如果交易未被打包，则保留在新的交易池中
		if !isPackaged {
			updatedTransactionPool = append(updatedTransactionPool, tx)
		}
	}
	return updatedTransactionPool
}
