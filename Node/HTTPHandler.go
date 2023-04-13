package Node

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"ppov/MetaData"
	"ppov/KeyManager"
	"ppov/Message"
)

type HTTPHandler struct {
	Node *Node
}
func (node *Node) StartHTTP(){
	handler := HTTPHandler{node}
	http.HandleFunc("/GetCurrentHeight", handler.GetCurrentHeight)
	http.HandleFunc("/CreateAccount",handler.CreateAccount)

	http.ListenAndServe(":"+strconv.Itoa(node.network.ServicePort), nil)
}

func (handler *HTTPHandler)GetCurrentHeight(res http.ResponseWriter, req *http.Request){
	defer req.Body.Close()

	msg := BasicMessage{Code: 200, Message: "SUCCESS",Result: handler.Node.mongo.Height}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}

	res.Write(data)
}

func (handler *HTTPHandler)CreateAccount(res http.ResponseWriter, req *http.Request){
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Read failed:", err)
	}
	defer req.Body.Close()
	var request CreatAccountMsg
	err = json.Unmarshal(b, &request)
	if err != nil {
		fmt.Println(err)
	}
	var basic BasicMessage
	if request.Address == "" || request.Pubkey == "" ||
		request.Timestamp == "" || request.Sig == "" {
		basic.Code = 401
		basic.Message = "PARAMETER WRONG"
		basic.Result = nil

	}

	_, existed := handler.Node.BalanceTable.Load(request.Address)
	if existed {
		basic.Code = 402
		basic.Message = "Address existed"
		basic.Result = nil
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
	handler.Node.SignResultCache.Store(txhash,1)

	if handler.Node.txPool.IsFull() == true{
		basic.Code = 500
		basic.Message = "Tx Pool is full"
		basic.Result = nil
	}

	//验证公钥和钱包地址是否对应
	temp := request
	temp.Sig = ""
	temp2, _ := json.Marshal(temp)
	hash := KeyManager.GetHash(temp2)
	ok, err := handler.Node.keymanager.Verify(hash,request.Sig,request.Pubkey)
	if err != nil || !ok {
		basic.Code = 403
		basic.Message = "key error"
		basic.Result = nil

	}

	ok, err = handler.Node.keymanager.VerifyAddressWithPubkey(request.Pubkey,request.Address)
	if err != nil || !ok {
		basic.Code = 403
		basic.Message = "address error"
		basic.Result = nil

	}

	timestamp,err :=strconv.Atoi(transaction.Timestamp)
	if err != nil{
		return
	}
	recordNode := uint32(timestamp % 4)
	//fmt.Println(p.Port)
	if handler.Node.accountManager.WorkerNumberSet[recordNode] == handler.Node.config.MyPubkey{

		handler.Node.txPool.PushbackTransactionFromTxByte(item)
	}else{
		//handler.Node.txPool.PushbackTransactionFromTxByte(item)
		if false{
			var blockmsg Message.BlockMsg //消息体
			blockmsg.Data = item

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = handler.Node.network.MyNodeInfo.ID
			msgheader.Receiver = handler.Node.accountManager.WorkerSet[handler.Node.accountManager.WorkerNumberSet[recordNode]]
			msgheader.Pubkey = handler.Node.config.MyPubkey
			msgheader.MsgType = Message.TransactionMsg
			handler.Node.SendMessage(msgheader, &blockmsg)
		}
	}

	basic.Code = 0
	basic.Message = "SUCCESS"
	result := txhash
	basic.Result = result

}