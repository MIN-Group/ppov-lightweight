package Node

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"ppov/KeyManager"
	"ppov/MetaData"
	"ppov/Network"
	"ppov/utils"
	"runtime"
	"sync"
	"time"
)

//交易排序需要
type TimePair struct {
	Key   int
	Value float64
}
type TimePairList []TimePair

func (t TimePairList) Len() int {
	return len(t)
}
func (t TimePairList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TimePairList) Less(i, j int) bool {
	return t[i].Value < t[j].Value
}

type TxPair struct {
	TransactionInterface MetaData.TransactionInterface
	Tx []byte
	Result bool
}
type TransTxPair struct {
	Id  int
	TransactionInterface MetaData.TransactionInterface
	Tx []byte
	Result bool
}

func (node *Node) UpdateIdTransformationVaribles(transactionInterface MetaData.TransactionInterface) {
	if transaction, ok := transactionInterface.(*MetaData.IdentityTransformation); ok {
		switch transaction.Type {
		case "ApplyForVoter":
			_, ok := node.accountManager.VoterSet[transaction.Pubkey]
			if !ok {
				node.accountManager.VoterSet[transaction.Pubkey] = transaction.NodeID
			} else {
				fmt.Println("申请成为投票节点失败，已经是投票节点")
			}
			_, ok = node.network.NodeList[transaction.NodeID]
			if !ok {
				var nodelist Network.NodeInfo
				nodelist.IP = transaction.IPAddr
				nodelist.PORT = transaction.Port
				nodelist.ID = transaction.NodeID
				node.network.NodeList[transaction.NodeID] = nodelist
			}
		case "ApplyForWorkerCandidate":
			_, ok := node.accountManager.WorkerCandidateSet[transaction.Pubkey]
			if !ok {
				node.accountManager.WorkerCandidateSet[transaction.Pubkey] = transaction.NodeID
			} else {
				fmt.Println("申请成为候选记账节点失败，已经是候选记账节点")
			}
			_, ok = node.network.NodeList[transaction.NodeID]
			if !ok {
				var nodelist Network.NodeInfo
				nodelist.IP = transaction.IPAddr
				nodelist.PORT = transaction.Port
				nodelist.ID = transaction.NodeID
				node.network.NodeList[transaction.NodeID] = nodelist
			}
		case "QuitVoter":
			delete(node.accountManager.VoterSet, transaction.Pubkey)
			delete(node.network.NodeList, transaction.NodeID)
			fmt.Println("退出投票节点成功")
		case "QuitWorkerCandidate":
			delete(node.accountManager.WorkerCandidateSet, transaction.Pubkey)
			delete(node.network.NodeList, transaction.NodeID)
			fmt.Println("退出候选记账节点成功")
		}
	}
}

func (node *Node) UpdateRecordVaribles(transactionInterface MetaData.TransactionInterface) {
	if transaction, ok := transactionInterface.(*MetaData.Record); ok {
		if transaction.Command == MetaData.ADD {
			record := node.mongo.GetResultFromDatabase("Record", "key", transaction.Key, "type", transaction.Type)
			_, ok := record["key"]
			if !ok {
				node.mongo.SaveRecordToDatabase("Record", *transaction)
			} else {
				record["value"] = transaction.Value
				node.mongo.UpdateRecordToDatabase("Record", record)
			}
			fmt.Println("TYPE", transaction.Type, "KEY:", transaction.Key, "VALUE", transaction.Value, "记录成功")
		}
	}
}
/*
func (node *Node) UpdateCreatAccountTx(transactionInterface MetaData.TransactionInterface, mmap map[string]bool, mutex *sync.RWMutex) bool {
	if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
		mutex.RLock()
		_, existed := mmap[transaction.Address]
		mutex.RUnlock()
		if !existed {
			_, existed := node.BalanceTable.Load(transaction.Address)
			if !existed {
				//验证签名

				temp := CreatAccountMsg{}
				temp.Pubkey = transaction.Pubkey
				temp.Address = transaction.Address
				temp.Timestamp = transaction.Timestamp
				temp.Sig = ""
				temp2, _ := json.Marshal(temp)

				ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2), transaction.Sig, transaction.Pubkey)
				if err != nil || !ok {
					fmt.Println("sign 4")
					return false
				}

				//验证公钥和钱包地址匹配情况

				ok, err = node.keymanager.VerifyAddressWithPubkey(transaction.Pubkey, transaction.Address)
				if err != nil || !ok {
					fmt.Println("address 2")
					return false
				}

				//正常情况
				mutex.Lock()
				node.BalanceTable.Store(transaction.Address, 100)
				mmap[transaction.Address] = true
				mutex.Unlock()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}*/

func (node *Node) UpdateCreatAccountTx(transactionInterface MetaData.TransactionInterface, mmap sync.Map, tx []byte) bool {
	if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
		//mutex.Lock()
		_, existed := mmap.Load(transaction.Address)
		//mutex.Unlock()
		if !existed {
			_, existed := node.BalanceTable.Load(transaction.Address)
			if !existed {
				//start := time.Now()
				////验证签名
				txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(tx))
				value , existed0 := node.SignResultCache.Load(txhash)
				if !existed0 {
					temp := CreatAccountMsg{}
					temp.Pubkey = transaction.Pubkey
					temp.Address = transaction.Address
					temp.Timestamp = transaction.Timestamp
					temp.Sig = ""
					temp2, _ := json.Marshal(temp)

					ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2),transaction.Sig,transaction.Pubkey)
					if err != nil || !ok {
						return false
					}

					//验证公钥和钱包地址匹配情况
					ok, err = node.keymanager.VerifyAddressWithPubkey(transaction.Pubkey,transaction.Address)
					if err != nil || !ok {
						return false
					}
				}else{
					if value == 0{
						return false
					}

				}

				//fmt.Println(time.Since(start))

				//正常情况
				_, existed := mmap.Load(transaction.Address)
				if !existed{
					node.BalanceTable.Store(transaction.Address, 1000)
				}
				//mutex.Lock()
				mmap.LoadOrStore(transaction.Address,true)
				//mutex.Unlock()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

/*
func (node *Node) UpdateCreatAccountTx(transactionInterface MetaData.TransactionInterface, mmap map[string]bool, mutex *sync.RWMutex) bool {
	if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
		mutex.RLock()
		_, existed := mmap[transaction.Address]
		mutex.RUnlock()
		if !existed {
			mutex.RLock()
			_, existed := node.BalanceTable.Load(transaction.Address)
			mutex.RUnlock()
			if !existed {
				//验证签名

				temp := CreatAccountMsg{}
				temp.Pubkey = transaction.Pubkey
				temp.Address = transaction.Address
				temp.Timestamp = transaction.Timestamp
				temp.Sig = ""
				temp2, _ := json.Marshal(temp)

				ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2), transaction.Sig, transaction.Pubkey)
				if err != nil || !ok {
					fmt.Println("sign 4")
					return false
				}

				//验证公钥和钱包地址匹配情况

				ok, err = node.keymanager.VerifyAddressWithPubkey(transaction.Pubkey, transaction.Address)
				if err != nil || !ok {
					fmt.Println("address 2")
					return false
				}

				//正常情况
				mutex.Lock()
				node.BalanceTable.Store(transaction.Address, 100)
				mmap[transaction.Address] = true
				mutex.Unlock()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}*/


func (node *Node) UpdateCreatAccountTxSerilized(transactionInterface MetaData.TransactionInterface, mmap map[string]bool) bool {
	if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
		//mutex.Lock()
		_, existed := mmap[transaction.Address]
		//mutex.Unlock()
		if !existed {
			_, existed := node.BalanceTable.Load(transaction.Address)
			if !existed {

				//正常情况
				node.BalanceTable.Store(transaction.Address, 1000)
				//mutex.Lock()
				mmap[transaction.Address] = true
				//mutex.Unlock()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

/*
func (node *Node) UpdateCreatAccountTx(transactionInterface MetaData.TransactionInterface, mmap map[string]bool, tx []byte) bool {
	if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
		//mutex.Lock()
		_, existed := mmap[transaction.Address]
		//mutex.Unlock()
		if !existed {
			_, existed := node.BalanceTable.Load(transaction.Address)
			if !existed {
				//start := time.Now()
				////验证签名
				node.TotalNum++

				txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(tx))
				value , existed0 := node.SignResultCache.Load(txhash)
				if !existed0 {
					temp := CreatAccountMsg{}
					temp.Pubkey = transaction.Pubkey
					temp.Address = transaction.Address
					temp.Timestamp = transaction.Timestamp
					temp.Sig = ""
					temp2, _ := json.Marshal(temp)

					ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2),transaction.Sig,transaction.Pubkey)
					if err != nil || !ok {
						return false
					}

					//验证公钥和钱包地址匹配情况
					ok, err = node.keymanager.VerifyAddressWithPubkey(transaction.Pubkey,transaction.Address)
					if err != nil || !ok {
						return false
					}
				}else{
					if value == 0{
						return false
					}else{
						node.CacheNum++
					}

				}

				//fmt.Println(time.Since(start))

				//正常情况
				node.BalanceTable.Store(transaction.Address, 1000)
				//mutex.Lock()
				mmap[transaction.Address] = true
				//mutex.Unlock()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}
*/

func (node *Node) UpdateTransferMoneyTx(transactionInterface MetaData.TransactionInterface, mmap sync.Map) bool {
	if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
		_, ok := mmap.Load(transaction.From)
		if !ok {
			balance1, existed1 := node.BalanceTable.Load(transaction.From)
			balance2, existed2 := node.BalanceTable.Load(transaction.To)
			if existed1 && existed2 && transaction.Amount > 0 && balance1.(int) >= transaction.Amount {
				//正常情况
				node.BalanceTable.Store(transaction.From, balance1.(int)-transaction.Amount)
				node.BalanceTable.Store(transaction.To, balance2.(int)+transaction.Amount)
				mmap.LoadOrStore(transaction.From, true)
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

func (node *Node) UpdateTransferMoneyTxWithVerify(transactionInterface MetaData.TransactionInterface, mmap sync.Map, tx []byte) bool {
	if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
		_, ok := mmap.Load(transaction.From)
		if !ok {
			txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(tx))
			value , existed0 := node.SignResultCache.Load(txhash)
			if !existed0 {
				var temp MetaData.TransferMoney
				temp = *transaction
				temp.Sig = ""
				temp2, _ := json.Marshal(temp)
				ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2), transaction.Sig, transaction.Pubkey)
				if err != nil || !ok {
					return false
				}

				//验证公钥和钱包地址匹配情况
				ok, err = node.keymanager.VerifyAddressWithPubkey(transaction.Pubkey, transaction.From)
				if err != nil || !ok {
					return false
				}
			}else{
				if value == 0{
					return false
				}

			}

			balance1, existed1 := node.BalanceTable.Load(transaction.From)
			balance2, existed2 := node.BalanceTable.Load(transaction.To)
			if existed1 && existed2 && transaction.Amount > 0 && balance1.(int) >= transaction.Amount {
				//正常情况
				node.BalanceTable.Store(transaction.From, balance1.(int)-transaction.Amount)
				node.BalanceTable.Store(transaction.To, balance2.(int)+transaction.Amount)
				mmap.LoadOrStore(transaction.From, true)
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

func (node *Node) UpdateTransferMoneyTxSerilized(transactionInterface MetaData.TransactionInterface, mmap map[string]bool) bool {
	if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
		_, ok := mmap[transaction.From]
		if !ok {
			balance1, existed1 := node.BalanceTable.Load(transaction.From)
			balance2, existed2 := node.BalanceTable.Load(transaction.To)
			if existed1 && existed2 && transaction.Amount > 0 && balance1.(int) >= transaction.Amount {
				//正常情况
				node.BalanceTable.Store(transaction.From, balance1.(int)-transaction.Amount)
				node.BalanceTable.Store(transaction.To, balance2.(int)+transaction.Amount)
				mmap[transaction.From] = true
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return false
}

func (node *Node) UpdateGenesisVaribles(transactionInterface MetaData.TransactionInterface) {
	if genesisTransaction, ok := transactionInterface.(*MetaData.GenesisTransaction); ok {
		node.config.WorkerNum = genesisTransaction.WorkerNum
		node.config.VotedNum = genesisTransaction.VotedNum
		node.config.BlockGroupPerCycle = genesisTransaction.BlockGroupPerCycle
		node.config.Tcut = genesisTransaction.Tcut
		node.accountManager.WorkerSet = genesisTransaction.WorkerPubList
		node.accountManager.WorkerCandidateSet = genesisTransaction.WorkerCandidatePubList
		node.accountManager.VoterSet = genesisTransaction.VoterPubList
		var index uint32 = 0
		for _, key := range genesisTransaction.WorkerSet {
			node.accountManager.WorkerNumberSet[index] = key
			node.accountManager.WorkerSetNumber[key] = index
			index = index + 1
		}
		index = 0
		for _, key1 := range genesisTransaction.VoterSet {
			node.accountManager.VoterNumberSet[index] = key1
			node.accountManager.VoterSetNumber[key1] = index
			index = index + 1
		}
		for key2, _ := range genesisTransaction.WorkerCandidatePubList {
			node.accountManager.WorkerCandidateList = append(node.accountManager.WorkerCandidateList, key2)
		}

	}
}

func (node *Node) UpdateVaribles(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		//去重
		var mmap  sync.Map

		//交易排序
		var tempTimePairList TimePairList
		for k, v := range bg.Blocks {
			pair := TimePair{
				Key:   k,
				Value: v.Timestamp,
			}
			tempTimePairList = append(tempTimePairList, pair)
		}
		//sort.Sort(tempTimePairList)

		if bg.ExecutionResult == nil {
			bg.ExecutionResult = make(map[string]bool)
		}


		createTxs := []TxPair{}
		tranTxConflict := []int{}
		tranTxNotConflict := []int{}
		tranTxs := []TransTxPair{}
		sub := 0
		conflictMap := make(map[string]int)

		deVoteResult := utils.DecompressToIntSlice(bg.VoteResult)
		for _, v := range tempTimePairList {
			//test
			//if node.dutyWorkerNumber == node.GetMyWorkerNumber() {
			//	fmt.Println(v.Key)
			//}
			if deVoteResult[v.Key] != 1 {
				continue
			}

			for _, eachTransaction := range bg.Blocks[v.Key].Transactions {
				transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
				switch transactionHeader.TXType {
				case MetaData.IdTransformation:
					node.UpdateIdTransformationVaribles(transactionInterface)
				case MetaData.Records:
					node.UpdateRecordVaribles(transactionInterface)
				case MetaData.CreatACCOUNT:
					createTxs = append(createTxs, TxPair{transactionInterface, eachTransaction,false})
				case MetaData.TransferMONEY:
					if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
						if conflictMap[transaction.From] != 1 && conflictMap[transaction.To] != 1{
							tranTxNotConflict = append(tranTxNotConflict, sub)
						}else{
							tranTxConflict = append(tranTxConflict,  sub)
						}
						tranTxs = append(tranTxs, TransTxPair{sub, transactionInterface, eachTransaction, false})
						sub++
						conflictMap[transaction.From] = 1
						conflictMap[transaction.To] = 1
					}

					/*
					res := node.UpdateTransferMoneyTxWithVerify(transactionInterface, mmap,eachTransaction)
					if bg.ExecutionResult == nil {
						bg.ExecutionResult = make(map[string]bool)
					}
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res*/

				}
			}


			node.TxsAmount += uint64(len(bg.Blocks[v.Key].Transactions))
			node.TxsPeriodAmount += uint64(len(bg.Blocks[v.Key].Transactions))
		}
		//create account
		if len(createTxs) > 0{
			worker := func(jobs <-chan TxPair, results chan<- TxPair) {
				for job := range jobs{
					res := node.UpdateCreatAccountTx(job.TransactionInterface, mmap, job.Tx)
					job.Result = res
					results <- job
				}
			}
			jobs := make(chan TxPair, len(createTxs))
			results := make(chan TxPair, len(createTxs))
			for w := 0 ; w < runtime.NumCPU() ; w++{
				go worker(jobs, results)
			}

			for j := 0 ; j< len(createTxs); j++{
				jobs <- TxPair{createTxs[j].TransactionInterface, createTxs[j].Tx, false }
			}
			close(jobs)
			for r := 0; r < len(createTxs) ; r++{
				result := <- results
				bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(result.Tx))] = result.Result
			}
		}

		//transfer money
		if len(tranTxs) > 0{

			transWorker := func(jobs <-chan TransTxPair, results chan<- TransTxPair) {
				for job := range jobs{
					res := false
					if transaction, ok := job.TransactionInterface.(*MetaData.TransferMoney); ok {
						txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(job.Tx))
						value, existed0 := node.SignResultCache.Load(txhash)
						if !existed0 {
							var temp MetaData.TransferMoney
							temp = *transaction
							temp.Sig = ""
							temp2, _ := json.Marshal(temp)
							ok, err := node.keymanager.Verify(KeyManager.GetHash(temp2), transaction.Sig, transaction.Pubkey)
							if err != nil || !ok {
								fmt.Println("key0 error")
								res = false
							}

							//验证公钥和钱包地址匹配情况
							ok, err = node.keymanager.VerifyAddressWithPubkey(transaction.Pubkey, transaction.From)
							if err != nil || !ok {
								fmt.Println("key error")
								res = false
							}else{
								res =true
							}

						}else{
							if value == 0 {
								fmt.Println("key2 error")
								res =  false
							}else{
								res = true
							}
						}
					}

					job.Result = res
					results <- job
				}
			}
			transJobs:= make(chan TransTxPair, len(tranTxs))
			transResults := make(chan TransTxPair, len(tranTxs))
			for w := 0 ; w < runtime.NumCPU() ; w++{
				go transWorker(transJobs, transResults)
			}

			for j := 0 ; j< len(tranTxs); j++{
				transJobs <-  TransTxPair{j,tranTxs[j].TransactionInterface, tranTxs[j].Tx, false }
			}
			close(transJobs)
			for r := 0; r < len(tranTxs) ; r++{
				result := <- transResults
				tranTxs[result.Id].Result = result.Result
			}

			transWorker1 := func(jobs <-chan TransTxPair, results chan<- TransTxPair) {
				for job := range jobs{
					res := false
					if tranTxs[job.Id].Result == false{
						res = false
					}else{
						res = node.UpdateTransferMoneyTx(tranTxs[job.Id].TransactionInterface, mmap)
					}
					job.Result = res
					results <- job
				}

			}
			transJobs1:= make(chan TransTxPair, len(tranTxNotConflict))
			transResults1 := make(chan TransTxPair, len(tranTxNotConflict))
			for w := 0 ; w < runtime.NumCPU() ; w++{
				go transWorker1(transJobs1, transResults1)
			}
			for _, j := range tranTxNotConflict{
				transJobs1 <-  TransTxPair{j,tranTxs[j].TransactionInterface, tranTxs[j].Tx, false }
			}
			close(transJobs1)
			for r := 0; r < len(tranTxNotConflict) ; r++{
				result := <- transResults1
				bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(tranTxs[result.Id].Tx))] = result.Result
			}

			for _, i := range tranTxConflict{

				if tranTxs[i].Result == false{
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(tranTxs[i].Tx))] = false
				}else{
					res := node.UpdateTransferMoneyTx(tranTxs[i].TransactionInterface, mmap)
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(tranTxs[i].Tx))] = res
				}
			}
		}


		node.BlockGroups.Store(bg.Height, *bg)
		node.StartTime = bg.Timestamp
	}
}
/*
func (node *Node) UpdateVaribles(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		//去重
		mmap := make(map[string]bool)

		//交易排序
		var tempTimePairList TimePairList
		for k, v := range bg.Blocks {
			pair := TimePair{
				Key:   k,
				Value: v.Timestamp,
			}
			tempTimePairList = append(tempTimePairList, pair)
		}
		//sort.Sort(tempTimePairList)

		if bg.ExecutionResult == nil {
			bg.ExecutionResult = make(map[string]bool)
		}

		node.CacheNum = 0
		node.TotalNum = 0
		for _, v := range tempTimePairList {
			//test
			//if node.dutyWorkerNumber == node.GetMyWorkerNumber() {
			//	fmt.Println(v.Key)
			//}
			if bg.VoteResult[v.Key] != 1 {
				continue
			}
			/*
			//并行执行交易
			mutex := sync.RWMutex{}
			length := len(bg.Blocks[v.Key].Transactions)
			wg := sync.WaitGroup{}
			wg.Add(length)
			for i := 0; i < length; i++ {
				go func(eachTransaction []byte) {
					defer wg.Done()
					transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
					switch transactionHeader.TXType {
					case MetaData.IdTransformation:
						node.UpdateIdTransformationVaribles(transactionInterface)
					case MetaData.Records:
						node.UpdateRecordVaribles(transactionInterface)

					case MetaData.CreatACCOUNT:
						res := node.UpdateCreatAccountTx(transactionInterface, mmap, &mutex)
						mutex.Lock()
						if bg.ExecutionResult == nil{
							bg.ExecutionResult = make(map[string]bool)
						}
						bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res
						mutex.Unlock()

					case MetaData.TransferMONEY:
						res := node.UpdateTransferMoneyTx(transactionInterface, mmap)
						mutex.Lock()
						if bg.ExecutionResult == nil {
							bg.ExecutionResult = make(map[string]bool)
						}
						bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res
						mutex.Unlock()
					}
				}(bg.Blocks[v.Key].Transactions[i])
			}
			wg.Wait()*/
/*
			for _, eachTransaction := range bg.Blocks[v.Key].Transactions {

				transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
				switch transactionHeader.TXType {
				case MetaData.IdTransformation:
					node.UpdateIdTransformationVaribles(transactionInterface)
				case MetaData.Records:
					node.UpdateRecordVaribles(transactionInterface)

				case MetaData.CreatACCOUNT:

					res := node.UpdateCreatAccountTx(transactionInterface, mmap, eachTransaction)
					if bg.ExecutionResult == nil{
						bg.ExecutionResult = make(map[string]bool)
					}
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res

				case MetaData.TransferMONEY:
					res := node.UpdateTransferMoneyTx(transactionInterface, mmap)
					if bg.ExecutionResult == nil {
						bg.ExecutionResult = make(map[string]bool)
					}
					bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] = res
				}
			}
			node.TxsAmount += uint64(len(bg.Blocks[v.Key].Transactions))
			node.TxsPeriodAmount += uint64(len(bg.Blocks[v.Key].Transactions))
		}
		fmt.Println("cache num", node.CacheNum, "total num", node.TotalNum)
		node.BlockGroups.Store(bg.Height, *bg)
		node.StartTime = bg.Timestamp
	}
}*/

func (node *Node) UpdateVariblesFromDisk(bg *MetaData.BlockGroup) {
	if bg.Height > 0 { //normal blockgroup
		node.dutyWorkerNumber = bg.NextDutyWorker
		node.StartTime = bg.Timestamp

		//去重
		mmap := make(map[string]bool)

		//交易排序
		var tempTimePairList TimePairList
		for k, v := range bg.Blocks {
			pair := TimePair{
				Key:   k,
				Value: v.Timestamp,
			}
			tempTimePairList = append(tempTimePairList, pair)
		}
		//sort.Sort(tempTimePairList)

		deVoteResult := utils.DecompressToIntSlice(bg.VoteResult)
		for _, v := range tempTimePairList {
			if deVoteResult[v.Key] != 1 {
				continue
			}
			for _, eachTransaction := range bg.Blocks[v.Key].Transactions {
				transactionHeader, transactionInterface := MetaData.DecodeTransaction(eachTransaction)
				switch transactionHeader.TXType {
				case MetaData.IdTransformation:
					node.UpdateIdTransformationVaribles(transactionInterface)
				case MetaData.CreatACCOUNT:
					if bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] {
						node.UpdateCreatAccountTxSerilized(transactionInterface, mmap)
					}
				case MetaData.TransferMONEY:
					if bg.ExecutionResult[base64.StdEncoding.EncodeToString(KeyManager.GetHash(eachTransaction))] {
						node.UpdateTransferMoneyTxSerilized(transactionInterface, mmap)
					}
				}
			}
		}
	}
}

func (node *Node) UpdateGenesisBlockVaribles(bg *MetaData.BlockGroup) {
	if bg.Height == 0 { //genesis blockgroup
		node.dutyWorkerNumber = 0
		node.StartTime = bg.Timestamp
		if bg.Blocks[0].Height == 0 {
			transactionHeader, transactionInterface := MetaData.DecodeTransaction(bg.Blocks[0].Transactions[0])
			if transactionHeader.TXType == MetaData.Genesis {
				node.UpdateGenesisVaribles(transactionInterface)
			}
		}
		_ = node.state
		node.state <- Normal
		time.Sleep(time.Second)
	} else {
		fmt.Println("更新变量错误")
	}
}
