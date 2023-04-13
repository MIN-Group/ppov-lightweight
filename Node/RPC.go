package Node

import (
	"ppov/KeyManager"
	"ppov/Message"
	"ppov/MetaData"
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

type PPoV struct {
	Node *Node
	Port int
}

func (node *Node)StartRPC() error{
	p := PPoV{node, node.config.ServicePort}
	var ppov *PPoV = &p

	//注册rpc服务
	rpc.Register(ppov)

	//获取tcpaddr
	port := p.Node.config.ServicePort
	tcpaddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
		return err
	}

	//监听端口
	tcplisten, err2 := net.ListenTCP("tcp4", tcpaddr);
	if err2 != nil {
		log.Fatal(err2)
		return err2
	}

	for {
		conn, err3 := tcplisten.Accept();
		if err3 != nil {
			log.Println(err3)
			continue
		}
		//使用goroutine单独处理rpc连接请求
		//这里使用jsonrpc进行处理
		go jsonrpc.ServeConn(conn);
	}

	return nil
}


//注意大写
type BasicMessage struct {
	Code     int 			`json:"code"`
	Message  string 		`json:"message"`
	Result   interface{} 	`json:"result"`
}

type Request struct {

}

/*
 *1.获取当前区块高度信息
 */
type GetHeightResult struct {
	Height int `json:"height"`
}

type GetHeightMessage struct {
	Code     int 			`json:"code"`
	Message  string 		`json:"message"`
	Result   interface{} 	`json:"result"`
}

func (p *PPoV) GetCurrentHeight(request *Request,basic *BasicMessage) error{
	height := p.Node.mongo.QueryHeight()

	basic.Code = 0
	basic.Message = "SUCCESS"

	var result = GetHeightResult{Height: height}
	basic.Result = result.Height

	return nil
}

/*
 *2.获取指定高度的区块组
 */
type GetBlockByHeightRequest struct {
	Height int `json:"height"`
}
func (p *PPoV) GetBlockGroupByHeight (request *GetBlockByHeightRequest,basic *BasicMessage) error{
	height := request.Height
	result := p.Node.mongo.GetBlockFromDatabase(height)

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = result

	return nil
}

/*
*3.获取某个区块组的所有区块
*/
type GetBlocksInGroupRequest struct {
	Height int  `json:"height"`
}
func (p *PPoV) GetBlocksInGroup( request GetBlocksInGroupRequest, basic *BasicMessage) error{
	height := request.Height
	blocks := p.Node.mongo.GetBlockFromDatabase(height).Blocks

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = blocks
	return nil
}

/*
 *4.获取某个具体的区块信息
 */
type GetBlockRequest struct {
	Height int  `json:"height"`
	Num    int  `json:"num"`
}
func (p *PPoV) GetBlock ( request *GetBlockRequest, basic *BasicMessage) error{
	height := request.Height
	num := request.Num
	block := p.Node.mongo.GetBlockFromDatabase(height).Blocks[num]

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = block

	return nil
}

/*
*5.获取某个范围内所有的区块组
*/
type GetBlockRangeRequest struct {
	From int  `json:"from"`
	To    int  `json:"to"`
}
func (p *PPoV) GetBlockRange( request *GetBlockRangeRequest, basic *BasicMessage) error{
	from := request.From
	to := request.To

	if from < to {
		basic.Code = 403
		basic.Message = "TO LESS THAN FROM"
		basic.Result = nil
		return nil
	}

	blockGroups := make([]MetaData.BlockGroup, from-to+1)

	for i:=from ; i <= to ; i++{
		blockGroups = append(blockGroups, p.Node.mongo.GetBlockFromDatabase(from))
	}
	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = blockGroups

	return nil
}

/*
*6.获取所有节点信息
 */
type GetNodeMessage struct {
	ID                uint64 `json:"id"`
	IP                string `json:"ip"`
	Port              int    `json:"port"`
	Pubkey            string `json:"pubkey"`
	Status            int    `json:"status"`
	IsWorker          int    `json:"isWorker"`
	IsVoter           int    `json:"isVoter"`
	IsWorkerCandidate int    `json:"isWorkerCandidate"`
	IsDutyWorker      int    `json:"isDutyWorker"`
}

func (p *PPoV) GetNodeList( request *Request, basic *BasicMessage) error{
	nodes := make(map[uint64]GetNodeMessage)

	for pubkey,id := range p.Node.accountManager.WorkerSet {
		var node GetNodeMessage
		node.ID = id
		node.Pubkey = pubkey
		node.IsWorker = 1
		node.IP = p.Node.network.NodeList[id].IP
		node.Port = p.Node.network.NodeList[id].PORT
		node.Status = p.Node.network.NodeList[id].Status
		if pubkey == p.Node.accountManager.WorkerNumberSet[p.Node.true_dutyWorkerNum] {
			node.IsDutyWorker = 1
		}

		nodes[id] = node
	}

	for pubkey,id := range p.Node.accountManager.VoterSet {
		node,ok := nodes[id]
		if ok{
			node.IsVoter = 1
			nodes[id] = node
		} else{
			node.ID = id
			node.Pubkey = pubkey
			node.IsWorker = 1
			node.IP = p.Node.network.NodeList[id].IP
			node.Port = p.Node.network.NodeList[id].PORT
			node.Status = p.Node.network.NodeList[id].Status

			nodes[id] = node
		}
	}

	for pubkey,id := range p.Node.accountManager.WorkerCandidateSet {
		node,ok := nodes[id]
		if ok{
			node.IsWorkerCandidate = 1
			nodes[id] = node
		} else{
			node.ID = id
			node.Pubkey = pubkey
			node.IsWorkerCandidate = 1
			node.IP = p.Node.network.NodeList[id].IP
			node.Port = p.Node.network.NodeList[id].PORT
			node.Status = p.Node.network.NodeList[id].Status

			nodes[id] = node
		}
	}

	result := make([]GetNodeMessage, len(nodes))

	index := 0
	for _, val := range nodes{
		result[index] = val
		index++
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = result

	return nil
}

/*
*7.根据公钥获取节点信息
 */
type GetNodeByPubkeyRequest struct {
	Pubkey string `json:"pubkey"`
}
func (p *PPoV) GetNodeByPubkey( request *GetNodeByPubkeyRequest, basic *BasicMessage) error{
	pubkey := request.Pubkey
	_, ok := p.Node.accountManager.VoterSet[pubkey]
	_, ok1 := p.Node.accountManager.WorkerSet[pubkey]
	_, ok2 := p.Node.accountManager.WorkerCandidateSet[pubkey]

	if !ok && !ok1 && !ok2{
		basic.Code = 402
		basic.Message = "NO NODE WITH THIS PUBKEY"
		basic.Result = nil
		return nil
	}
	var node  GetNodeMessage
	if ok{
		id := p.Node.accountManager.VoterSet[pubkey]
		node.IP = p.Node.network.NodeList[id].IP
		node.Port = p.Node.network.NodeList[id].PORT
		node.Status = p.Node.network.NodeList[id].Status
		node.ID = id
		node.Pubkey = pubkey
		node.IsVoter = 1
	}
	if ok1{
		id := p.Node.accountManager.VoterSet[pubkey]
		node.IP = p.Node.network.NodeList[id].IP
		node.Port = p.Node.network.NodeList[id].PORT
		node.Status = p.Node.network.NodeList[id].Status
		node.ID = id
		node.Pubkey = pubkey
		node.IsWorker = 1
	}
	if ok2{
		id := p.Node.accountManager.VoterSet[pubkey]
		node.IP = p.Node.network.NodeList[id].IP
		node.Port = p.Node.network.NodeList[id].PORT
		node.Status = p.Node.network.NodeList[id].Status
		node.ID = id
		node.Pubkey = pubkey
		node.IsWorkerCandidate = 1
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = node

	return nil
}

/*
*8.获取投票节点列表
 */
func (p *PPoV) GetVoterList(request *Request, basic *BasicMessage) error{
	nodes := make([]string, len(p.Node.accountManager.VoterSet))

	index := 0
	for key, _ := range p.Node.accountManager.VoterSet{
		nodes[index] = key
		index++
	}

	result := map[string][]string{
		"list" : nodes,
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = result

	return nil
}

/*
*9.获取记账节点列表
 */
func (p *PPoV) GetWorkerList(request *Request, basic *BasicMessage) error{
	nodes := make([]string, len(p.Node.accountManager.WorkerSet))

	index := 0
	for key, _ := range p.Node.accountManager.WorkerSet{
		nodes[index] = key
		index++
	}

	result := map[string][]string{
		"list" : nodes,
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = result

	return nil
}

/*
*10.获取候选记账节点列表
 */
func (p *PPoV) GetWorkerCandidateList(request *Request, basic *BasicMessage) error{
	nodes := make([]string, len(p.Node.accountManager.WorkerCandidateSet))

	index := 0
	for key, _ := range p.Node.accountManager.WorkerCandidateSet{
		nodes[index] = key
		index++
	}

	result := map[string][]string{
		"list" : nodes,
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = result

	return nil
}

/*
*11.获取当前轮值记账节点
 */
func (p *PPoV) GetCurrentDutyWorker(request *Request, basic *BasicMessage) error{
	pubkey := p.Node.accountManager.WorkerNumberSet[p.Node.true_dutyWorkerNum]

	result := map[string]string{
		"pubkey":pubkey,
	}
	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = result

	return nil
}

/*
*12.上传一条记录
 */
type PostRecordRequest struct {
	Key    string  `json:"key"`
	Value  string  `json:"value"`
	Type   string  `json:"type"`
}

func(p *PPoV) PostRecord(request *PostRecordRequest, basic *BasicMessage) error{
	if request.Key == "" || request.Value == ""{
		basic.Code = 401
		basic.Message = "PARAMETER WRONG"
		basic.Result = nil

		return nil
	}
	if request.Type == ""{
		request.Type = "default"
	}
	var transaction MetaData.Record
	transaction.Key = request.Key
	transaction.Value = request.Value
	transaction.Sender = p.Node.config.MyPubkey
	transaction.Timestamp =  strconv.FormatInt(time.Now().UTC().UnixNano(),10)
	transaction.Type = request.Type
	transaction.Command = MetaData.ADD
	hash := KeyManager.GetHash(transaction.ToByteArray())
	transaction.Sig, _ = p.Node.keymanager.Sign(hash)


	var transactionHeader MetaData.TransactionHeader
	transactionHeader.TXType = MetaData.Records

	item := MetaData.EncodeTransaction(transactionHeader, &transaction)
	p.Node.txPool.PushbackTransactionFromTxByte(item)
	txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(item))
	p.Node.SignResultCache.Store(txhash,1)

	basic.Code = 0
	basic.Message = "SUCCESS"
	result := map[string]string{
		"hash":txhash,
	}

	basic.Result = result

	return nil
}

/*
*13.查询一条记录
 */
type GetRecordRequest struct {
	Key    string  `json:"key"`
	Type   string  `json:"type"`
}
func(p *PPoV) GetRecord(request *GetRecordRequest, basic *BasicMessage) error{
	myType := request.Type
	key := request.Key
	record := p.Node.mongo.GetResultFromDatabase("Record","key",myType,"key",key)
	basic.Code = 0
	basic.Message = "SUCCESS"
	types := ""
	keys := ""
	value := ""
	if record["type"] != nil{
		types = record["type"].(string)
	}
	if record["key"] != nil{
		keys = record["key"].(string)
	}
	if record["value"] != nil{
		value = record["value"].(string)
	}

	result := MetaData.KVRecord{Type: types, Key: keys, Value: value}
	basic.Result = result

	return nil
}

/*
*14.根据哈希查询一条交易
 */
type GetTransactionByHashRequest struct {
	Hash string `json:"hash"`
}
type GetTransactionByHashResponse struct {
	Tx map[string]interface{}
	ExecutionResult bool
}
func (p *PPoV) GetTransactionByHash(request *GetTransactionByHashRequest, basic *BasicMessage) error{
	hash := request.Hash
	result := make(map[string]interface{})

	hb ,_:= base64.StdEncoding.DecodeString(hash)
	t := p.Node.mongo.GetBlockByTxHashFromDatabase(hb)
	result["height"] = t.Height
	executionResult := false
	if _,ok := t.ExecutionResult[hash]; ok == true{
		executionResult = t.ExecutionResult[hash]
	}

	for j,b := range t.Blocks{
		for i,h := range b.TransactionsHash{
			if string(h) == string(hb) {
				result["blockNums"] = j
				_,transaction := MetaData.DecodeTransaction(b.Transactions[i])
				result["data"] = transaction
			}
		}
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = GetTransactionByHashResponse{Tx: result, ExecutionResult: executionResult}
	return nil
}

/*
*15.获取指定区组块内的交易信息
 */
type GetTransactionsInBlockGroupRequest struct {
	Height int `json:"height"`
}
type GetTransactionsInBlockGroupResponse struct {
	txs []interface{}
	ExecutionResult map[string]bool
}

func (p *PPoV) GetTransactionsInBlockGroup(request *GetTransactionsInBlockGroupRequest, basic *BasicMessage) error{
	height := request.Height
	bg := p.Node.mongo.GetBlockFromDatabase(height)
	result := make([]interface{}, len(bg.Blocks))
	executionResult := bg.ExecutionResult

	for i,b := range bg.Blocks{
		btxs := make([]interface{}, len(b.Transactions))
		for j,t := range b.Transactions{
			_, tx := MetaData.DecodeTransaction(t)
			btxs[j] = tx
		}
		result[i] = btxs
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = GetTransactionsInBlockGroupResponse{
		ExecutionResult: executionResult,
		txs: result,
	}
	return nil
}

/*
*16.获取指定区组块内指定高度的交易信息
 */
type GetTransactionsInBlockRequest struct {
	Height int `json:"height"`
	BlockNums int `json:"blockNums"`
}
type GetTransactionsInBlockResponse struct {
	Txs []interface{}
	ExecutionResult map[string]bool
}
func (p *PPoV) GetTransactionsInBlock(request *GetTransactionsInBlockRequest, basic *BasicMessage) error{
	height := request.Height
	nums := request.BlockNums
	bg := p.Node.mongo.GetBlockFromDatabase(height)
	ertmp := bg.ExecutionResult

	executionResult := make(map[string]bool)
	if len(bg.Blocks[nums].Transactions) < nums {
		basic.Code = 403
		basic.Message = "OUT OF RANGE"
		basic.Result = nil
	}

	result := make([]interface{}, len(bg.Blocks[nums].Transactions))
	for i,tx := range bg.Blocks[nums].Transactions{
		_, t := MetaData.DecodeTransaction(tx)
		result[i] = t
		executionResult[base64.StdEncoding.EncodeToString(bg.Blocks[nums].TransactionsHash[i])] = ertmp[base64.StdEncoding.EncodeToString(bg.Blocks[nums].TransactionsHash[i])]
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	basic.Result = GetTransactionsInBlockResponse{Txs: result, ExecutionResult: executionResult}
	return nil
}

/*
*17.创建账户
 */
type CreatAccountMsg struct {
	Address    string     `json:"address"`
	Pubkey     string 	  `json:"pubkey"`
	Timestamp  string     `json:"timestamp"`
	Sig        string     `json:"sig"`
}

func(p *PPoV) CreatAccount(request CreatAccountMsg, basic *BasicMessage) error{
	if request.Address == "" || request.Pubkey == "" ||
		request.Timestamp == "" || request.Sig == "" {
		basic.Code = 401
		basic.Message = "PARAMETER WRONG"
		basic.Result = nil

		return nil
	}

	_, existed := p.Node.BalanceTable.Load(request.Address)
	if existed {
		basic.Code = 402
		basic.Message = "Address existed"
		basic.Result = nil

		return nil
	}

	var transaction MetaData.CreatAccount
	transaction.Address = request.Address
	transaction.Pubkey = request.Pubkey
	transaction.Timestamp = request.Timestamp
	transaction.Sig = request.Sig

	var transactionHeader MetaData.TransactionHeader
	transactionHeader.TXType = MetaData.CreatACCOUNT

	item := MetaData.EncodeTransaction(transactionHeader, &transaction)

	txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(item))
	p.Node.SignResultCache.Store(txhash,1)

	if p.Node.txPool.IsFull() == true{
		basic.Code = 500
		basic.Message = "Tx Pool is full"
		basic.Result = nil

		return nil
	}

	//验证公钥和钱包地址是否对应
	temp := request
	temp.Sig = ""
	temp2, _ := json.Marshal(temp)
	hash := KeyManager.GetHash(temp2)
	ok, err := p.Node.keymanager.Verify(hash,request.Sig,request.Pubkey)
	if err != nil || !ok {
		basic.Code = 403
		basic.Message = "key error"
		basic.Result = nil

		return nil
	}

	ok, err = p.Node.keymanager.VerifyAddressWithPubkey(request.Pubkey,request.Address)
	if err != nil || !ok {
		basic.Code = 403
		basic.Message = "address error"
		basic.Result = nil

		return nil
	}

	timestamp,err :=strconv.Atoi(transaction.Timestamp)
	if err != nil{
		return nil
	}
	recordNode := uint32(timestamp % len(p.Node.accountManager.WorkerNumberSet))
	//fmt.Println(p.Port)
	if p.Node.accountManager.WorkerNumberSet[recordNode] == p.Node.config.MyPubkey{

		p.Node.txPool.PushbackTransactionFromTxByte(item)
	}else{
		p.Node.txPool.PushbackTransactionFromTxByte(item)
		if false{
			var blockmsg Message.BlockMsg //消息体
			blockmsg.Data = item

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = p.Node.network.MyNodeInfo.ID
			msgheader.Receiver = p.Node.accountManager.WorkerSet[p.Node.accountManager.WorkerNumberSet[recordNode]]
			msgheader.Pubkey = p.Node.config.MyPubkey
			msgheader.MsgType = Message.TransactionMsg
			p.Node.SendMessage(msgheader, &blockmsg)
		}
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	result := txhash
	basic.Result = result

	return nil
}

/*
*18.转账
 */
type TransferMoneyMsg struct {
	From       string     `json:"from"`
	To     	   string 	  `json:"to"`
	Pubkey	   string 	  `json:"pubkey"`
	Amount     int 	      `json:"amount"`
	Timestamp  string     `json:"timestamp"`
	Sig        string     `json:"sig"`
}

func(p *PPoV) TransferMoney(request TransferMoneyMsg, basic *BasicMessage) error{
	if request.From == "" || request.To == "" || request.Amount <= 0 || request.Pubkey == ""||
		request.Timestamp == "" || request.Sig == "" {
		basic.Code = 401
		basic.Message = "PARAMETER WRONG"
		basic.Result = nil

		return nil
	}

	if p.Node.txPool.IsFull() == true{
		basic.Code = 500
		basic.Message = "Tx Pool is full"
		basic.Result = nil

		return nil
	}

	value, existed := p.Node.BalanceTable.Load(request.From)
	_, existed2 := p.Node.BalanceTable.Load(request.To)
	if !existed || !existed2 {
		basic.Code = 402
		basic.Message = "Address not existed"
		basic.Result = nil

		return nil
	}

	if request.Amount > value.(int) {
		basic.Code = 403
		basic.Message = "Insufficient Balance"
		basic.Result = nil

		return nil
	}

	//验证公钥和钱包地址是否对应
	temp := request
	temp.Sig = ""
	temp2, _ := json.Marshal(temp)

	ok, err := p.Node.keymanager.Verify(KeyManager.GetHash(temp2),request.Sig,request.Pubkey)
	if err != nil || !ok {
		basic.Code = 404
		basic.Message = "key error"
		basic.Result = nil

		return nil
	}
	ok, err = p.Node.keymanager.VerifyAddressWithPubkey(request.Pubkey,request.From)
	if err != nil || !ok {
		basic.Code = 403
		basic.Message = "address error"
		basic.Result = nil

		return nil
	}

	var transaction MetaData.TransferMoney
	transaction.From = request.From
	transaction.To = request.To
	transaction.Pubkey = request.Pubkey
	transaction.Amount = request.Amount
	transaction.Timestamp = request.Timestamp
	transaction.Sig = request.Sig

	var transactionHeader MetaData.TransactionHeader
	transactionHeader.TXType = MetaData.TransferMONEY

	item := MetaData.EncodeTransaction(transactionHeader, &transaction)
	txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(item))

	p.Node.SignResultCache.Store(txhash,1)

	timestamp,err :=strconv.Atoi(transaction.Timestamp)
	if err != nil{
		return nil
	}
	recordNode := uint32(timestamp % len(p.Node.accountManager.WorkerNumberSet))
	if p.Node.accountManager.WorkerNumberSet[recordNode] == p.Node.config.MyPubkey{
		p.Node.txPool.PushbackTransactionFromTxByte(item)
	}else{
		p.Node.txPool.PushbackTransactionFromTxByte(item)
		if false{
			var blockmsg Message.BlockMsg //消息体
			blockmsg.Data = item

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = p.Node.network.MyNodeInfo.ID
			msgheader.Receiver = p.Node.accountManager.WorkerSet[p.Node.accountManager.WorkerNumberSet[recordNode]]
			msgheader.Pubkey = p.Node.config.MyPubkey
			msgheader.MsgType = Message.TransactionMsg
			p.Node.SendMessage(msgheader, &blockmsg)
		}
	}

	basic.Code = 0
	basic.Message = "SUCCESS"

	basic.Result = txhash

	return nil
}

/*
*18.查询余额
 */
type GetBalanceMsg struct {
	Address       string     `json:"address"`
}

func(p *PPoV) GetBalance(request *GetBalanceMsg, basic *BasicMessage) error{
	if request.Address == "" {
		basic.Code = 401
		basic.Message = "PARAMETER WRONG"
		basic.Result = nil

		return nil
	}

	value, existed := p.Node.BalanceTable.Load(request.Address)
	if !existed  {
		basic.Code = 402
		basic.Message = "Address not existed"
		basic.Result = nil

		return nil
	}

	basic.Code = 200
	basic.Message = "SUCCESS"
	val := strconv.Itoa(value.(int))
	result := map[string]string{
		request.Address:val,
	}

	basic.Result = result

	return nil
}