package Node

import (
	"ppov/KeyManager"
	"ppov/MetaData"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func (node *Node) ValidateBlockHeader(b *MetaData.Block) bool {
	_, existed := node.accountManager.WorkerSet[b.Generator]
	if !existed {
		return false
	}
	if node.accountManager.WorkerNumberSet[b.BlockNum] != b.Generator {
		return false
	}
	return true
}

/*func (node *Node)ValidateBlock(b *MetaData.Block) bool{
	_,existed:=node.accountManager.WorkerSet[b.Generator]
	if !existed {
		return false
	}
	if(node.accountManager.WorkerNumberSet[b.BlockNum]!=b.Generator){
		return false
	}
	if(!node.ValidateTransactions(b.Transactions)){
		return false
	}
	return true
}*/

func (node *Node) ValidateTransactions(txs *([][]byte)) bool {
	return true
	start := time.Now()
	length := len(*txs)
	res := true
	wg := sync.WaitGroup{}
	wg.Add(length)
	sum := 0
	for i := 0; i < length; i++ {
		go func(tx []byte) {
			defer wg.Done()
			txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(tx))
			_, existed := node.SignResultCache.Load(txhash)
			if !existed {
				sum++
				header, transactionInterface := MetaData.DecodeTransaction(tx)
				switch header.TXType {
				case MetaData.Zero:
					if transaction, ok := transactionInterface.(*MetaData.ZeroTransaction); ok {
						if !node.ValidateZeroTransaction(transaction) {
							res = false
						}
					}
				case MetaData.Genesis:
					if transaction, ok := transactionInterface.(*MetaData.GenesisTransaction); ok {
						if !node.ValidateGenesisTransaction(transaction) {
							res = false
						}
					}
				case MetaData.CreatACCOUNT:
					if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
						if !node.ValidateCreatAccountTransaction(transaction) {
							res = false
						}
					}
				case MetaData.TransferMONEY:
					if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
						if !node.ValidateTransferMoneyTransaction(transaction) {
							res = false
						}
					}
				}
			}else{
			} //else if value.(int) == 0 {
			//	res = false
			//}
		}((*txs)[i])
	}
	wg.Wait()
	fmt.Println("validate ",time.Since(start),"-",sum,"-",length)
	return res
}

func (node *Node) ValidateZeroTransaction(tx *MetaData.ZeroTransaction) bool {
	return true
}

func (node *Node) ValidateGenesisTransaction(tx *MetaData.GenesisTransaction) bool {
	return true
}

func (node *Node) ValidateCreatAccountTransaction(tx *MetaData.CreatAccount) bool {
	//验证签名
	var temp MetaData.CreatAccount
	temp = *tx
	temp.Sig = ""
	temp2, _ := json.Marshal(temp)
	ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2),tx.Sig,tx.Pubkey)
	if err != nil || !ok {
		return false
	}

	//验证公钥和钱包地址匹配情况
	ok, err = node.keymanager.VerifyAddressWithPubkey(tx.Pubkey,tx.Address)
	if err != nil || !ok {
		return false
	}
	return true
}

func (node *Node) ValidateTransferMoneyTransaction(tx *MetaData.TransferMoney) bool {
	//验证签名
	var temp MetaData.TransferMoney
	temp = *tx
	temp.Sig = ""
	temp2, _ := json.Marshal(temp)
	ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2), tx.Sig, tx.Pubkey)
	if err != nil || !ok {
		return false
	}

	//验证公钥和钱包地址匹配情况
	ok, err = node.keymanager.VerifyAddressWithPubkey(tx.Pubkey, tx.From)
	if err != nil || !ok {
		return false
	}
	return true
}