package utils

import (
	"ppov/MetaData"
	"sync"
)

type TransactionPool struct {
	txsPool  map[int][][]byte
	txsSize  int
	amout    int
	capacity int
	writeNum int
	readNum  int
	lock	sync.Mutex
}

func (tp *TransactionPool) Init(txsSize int, capacity int) {
	var k int = 0
	if txsSize <= 0 || capacity <= 0 {
		tp.txsSize = 0
		tp.capacity = 0
		k = 1
	} else {
		tp.lock.Lock()
		tp.txsPool = make(map[int][][]byte)
		if tp.txsSize > tp.capacity {
			tp.txsSize = capacity
		} else {
			tp.txsSize = txsSize
		}
		tp.capacity = capacity
		if capacity%txsSize > 0 {
			k = capacity/txsSize + 1
		} else {
			k = capacity / txsSize
		}
		for i := 0; i < k; i++ {
			txs := make([][]byte, 0, tp.txsSize)
			tp.txsPool[i] = txs
		}
		tp.lock.Unlock()
	}
}

func (tp *TransactionPool) PushbackTransaction(header MetaData.TransactionHeader, transactionInterface MetaData.TransactionInterface) int {
	tp.lock.Lock()
	if tp.amout >= tp.capacity {
		tp.lock.Unlock()
		return -1
	}
	for {
		if len(tp.txsPool[tp.writeNum]) >= tp.txsSize {
			tp.writeNum = (tp.writeNum + 1) % len(tp.txsPool)
		} else {
			break
		}
	}
	tx := MetaData.EncodeTransaction(header, transactionInterface)
	tp.txsPool[tp.writeNum] = append(tp.txsPool[tp.writeNum], tx)
	tp.amout++
	tp.lock.Unlock()
	return 1
}

func (tp *TransactionPool) PushbackTransactionFromTxByte(tx []byte) int {
	tp.lock.Lock()
	if tp.amout >= tp.capacity {
		tp.lock.Unlock()
		return -1
	}
	for {
		if len(tp.txsPool[tp.writeNum]) >= tp.txsSize {
			tp.writeNum = (tp.writeNum + 1) % len(tp.txsPool)
		} else {
			break
		}
	}
	tp.txsPool[tp.writeNum] = append(tp.txsPool[tp.writeNum], tx)
	tp.amout++
	tp.lock.Unlock()
	return 1
}

func (tp *TransactionPool) GetCurrentTxsList() (txs [][]byte) {
	tp.lock.Lock()
	txs = tp.txsPool[tp.readNum]
	tp.readNum = (tp.readNum + 1) % len(tp.txsPool)
	tp.lock.Unlock()
	return
}

func (tp *TransactionPool) GetCurrentTxsListDelete() (txs [][]byte) {
	tp.lock.Lock()
	txs = tp.txsPool[tp.readNum]
	tp.amout-=len(txs)
	tp.txsPool[tp.readNum] = (tp.txsPool[tp.readNum])[:0]
	tp.readNum = (tp.readNum + 1) % len(tp.txsPool)
	tp.lock.Unlock()
	return
}

func (tp *TransactionPool) IsFull() bool {
	if tp.amout >= tp.capacity {
		return true
	} else {
		return false
	}
}
