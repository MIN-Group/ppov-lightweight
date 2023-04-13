package Client

import (
	"ppov/KeyManager"
	"encoding/json"
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"sync"
	"sync/atomic"
	"time"
)

type KeyMessage struct {
	Prikey string
	Pubkey string
}

type BasicMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type Request struct {
}

/*
 *1.获取当前区块高度信息
 */
type GetHeightResult struct {
	Height int `json:"height"`
}


/*
 *2.获取指定高度的区块组
 */
type GetBlockByHeightRequest struct {
	Height int `json:"height"`
}

/*
*3.获取某个区块组的所有区块
 */
type GetBlocksInGroupRequest struct {
	Height int  `json:"height"`
}

/*
 *4.获取某个具体的区块信息
 */
type GetBlockRequest struct {
	Height int  `json:"height"`
	Num    int  `json:"num"`
}

/*
*5.获取某个范围内所有的区块组
 */
type GetBlockRangeRequest struct {
	From int  `json:"from"`
	To    int  `json:"to"`
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

/*
*7.根据公钥获取节点信息
 */
type GetNodeByPubkeyRequest struct {
	Pubkey string `json:"pubkey"`
}



/*
*12.上传一条记录
 */
type PostRecordRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

/*
*13.查询一条记录
 */
type GetRecordRequest struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

/*
*14.根据哈希查询一条交易
 */
type GetTransactionByHashRequest struct {
	Hash string `json:"hash"`
}

/*
*15.获取指定区组块内的交易信息
 */
type GetTransactionsInBlockGroupRequest struct {
	Height int `json:"height"`
}

/*
*16.获取指定区组块内指定高度的交易信息
 */
type GetTransactionsInBlockRequest struct {
	Height    int `json:"height"`
	BlockNums int `json:"blockNums"`
}

/*
*17.创建账户
 */
type CreatAccountMsg struct {
	Address   string `json:"address"`
	Pubkey    string `json:"pubkey"`
	Timestamp string `json:"timestamp"`
	Sig       string `json:"sig"`
}

/*
*18.转账
 */
type TransferMoneyMsg struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Pubkey    string `json:"pubkey"`
	Amount    int    `json:"amount"`
	Timestamp string `json:"timestamp"`
	Sig       string `json:"sig"`
}

/*
*19.查询余额
 */
type GetBalanceMsg struct {
	Address string `json:"address"`
}

func Client() {
	//连接远程rpc服务
	//这里使用jsonrpc.Dial
	rpc1, err := jsonrpc.Dial("tcp", "127.0.0.1:8010")
	if err != nil {
		log.Fatal(err)
	}
	var choice int
	fmt.Print("请选择需要的函数: \n 1.获取当前区块高度信息\n 2.获取指定高度的区块组\n 3.获取某个区块组的所有区块\n 4.获取某个具体的区块信息\n 5.获取某个范围内所有的区块组\n 6.获取所有节点信息\n 7.根据公钥获取节点信息\n 8.获取投票节点列表\n 9.获取记账节点列表\n 10.获取候选记账节点列表\n 11.获取当前轮值记账节点\n 12.上传一条记录\n 13.查询一条记录\n 14.根据哈希查询一条交易\n 15.获取指定区组块内的交易信息\n 16.获取指定区组块内指定高度的交易信息\n 17.创建账户\n 18.转账\n 19.查询余额\n 20.重复测试\n 21.并发测试\n")
	fmt.Scanln(&choice)
	var basic BasicMessage
	var request Request

	switch choice {
	case 1: //1.获取当前区块高度信息
		err2 := rpc1.Call("PPoV.GetCurrentHeight", request, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 2: //2.获取指定高度的区块组
		var r GetBlockByHeightRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&r.Height)

		err2 := rpc1.Call("PPoV.GetBlockGroupByHeight", r, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 3: //3.获取某个区块组的所有区块
		var r GetBlocksInGroupRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&r.Height)

		err2 := rpc1.Call("PPoV.GetBlocksInGroup", r, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 4: //4.获取某个具体的区块信息
		var r GetBlockRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&r.Height)
		fmt.Println("请输入Num：")
		fmt.Scanln(&r.Num)

		err2 := rpc1.Call("PPoV.GetBlock", r, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 5: //5.获取某个范围内所有的区块组
		var r GetBlockRangeRequest
		fmt.Println("请输入From：")
		fmt.Scanln(&r.From)
		fmt.Println("请输入To：")
		fmt.Scanln(&r.To)

		err2 := rpc1.Call("PPoV.GetBlockRange", r, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 6: //6.获取所有节点信息
		err2 := rpc1.Call("PPoV.GetNodeList", request, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 7: //7.根据公钥获取节点信息
		var r GetNodeByPubkeyRequest
		fmt.Println("请输入Pubkey：")
		fmt.Scanln(&r.Pubkey)
		err2 := rpc1.Call("PPoV.GetNodeByPubkey", r, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 8: //8.获取投票节点列表
		err2 := rpc1.Call("PPoV.GetVoterList", request, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 9: //9.获取记账节点列表
		err2 := rpc1.Call("PPoV.GetWorkerList", request, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 10: //10.获取候选记账节点列表
		err2 := rpc1.Call("PPoV.GetWorkerCandidateList", request, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 11: //11.获取当前轮值记账节点
		err2 := rpc1.Call("PPoV.GetCurrentDutyWorker", request, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 12: //12.上传一条记录
		var s PostRecordRequest
		fmt.Println("请输入Key：")
		fmt.Scanln(&s.Key)
		fmt.Println("请输入Value：")
		fmt.Scanln(&s.Value)
		fmt.Println("请输入Type：")
		fmt.Scanln(&s.Type)
		err2 := rpc1.Call("PPoV.PostRecord", &s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 13: //13.查询一条记录
		var s GetRecordRequest
		fmt.Println("请输入Key：")
		fmt.Scanln(&s.Key)
		fmt.Println("请输入Type：")
		fmt.Scanln(&s.Type)
		err2 := rpc1.Call("PPoV.GetRecord", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 14: //14.根据哈希查询一条交易
		var s GetTransactionByHashRequest
		fmt.Println("请输入Hash：")
		fmt.Scanln(&s.Hash)
		err2 := rpc1.Call("PPoV.GetTransactionByHash", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 15: //15.获取指定区组块内的交易信息
		var s GetTransactionsInBlockGroupRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&s.Height)
		err2 := rpc1.Call("PPoV.GetTransactionsInBlockGroup", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 16: //16.获取指定区组块内指定高度的交易信息
		var s GetTransactionsInBlockRequest
		fmt.Println("请输入Height：")
		fmt.Scanln(&s.Height)
		fmt.Println("请输入BlockNums：")
		fmt.Scanln(&s.BlockNums)
		err2 := rpc1.Call("PPoV.GetTransactionsInBlock", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 17: //17.创建账户
		var s CreatAccountMsg
		now := time.Now()
		s.Timestamp = fmt.Sprintf("%d", now.Unix())

		var km KeyManager.KeyManager
		km.Init()
		km.GenRandomKeyPair()
		s.Address = km.GetAddress()
		s.Pubkey = km.GetPubkey()
		s.Sig = ""
		temp2, _ := json.Marshal(s)
		s.Sig, err = km.Sign(KeyManager.GetHash(temp2))
		pri, _ := km.GetPriKey()
		fmt.Println("私钥:", pri)
		fmt.Println("公钥", s.Pubkey)
		fmt.Println("地址", s.Address)
		KeyTxtWrite("PubKey", s.Address, s.Pubkey) //写入公钥
		KeyTxtWrite("PriKey", s.Address, pri)      //写入私钥
		err2 := rpc1.Call("PPoV.CreatAccount", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 18: //18.转账
		var s TransferMoneyMsg
		fmt.Println("请输入From：")
		fmt.Scanln(&s.From)
		fmt.Println("请输入To：")
		fmt.Scanln(&s.To)
		s.Pubkey = KeyTxtRead("PubKey", s.From)
		pri := KeyTxtRead("PriKey", s.From)

		now := time.Now()
		s.Timestamp = fmt.Sprintf("%d", now.Unix())
		fmt.Println("请输入Amount：")
		fmt.Scanln(&s.Amount)

		var km KeyManager.KeyManager
		km.Init()
		//km.GenRandomKeyPair()
		km.SetPriKey(pri)
		km.SetPubkey(s.Pubkey)
		temp2, _ := json.Marshal(s)
		s.Sig = ""
		s.Sig, err = km.SignWithPriKey(KeyManager.GetHash(temp2), pri)
		//fmt.Println(s.Pubkey)
		//fmt.Println(pri)
		err2 := rpc1.Call("PPoV.TransferMoney", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 19: //19.查询余额
		var s GetBalanceMsg
		fmt.Println("请输入Address：")
		fmt.Scanln(&s.Address)
		err2 := rpc1.Call("PPoV.GetBalance", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

	case 20:
		var s TransferMoneyMsg
		s.Pubkey = "3059301306072a8648ce3d020106082a811ccf5501822d03420004c3a32777ca149e92fd4d3a5bd575c4b153082f43a704f120c6398e461409014f47a53a3943f8597af786027cb35085b24c33c1bd24e4e2af6786dc11925b27f6"
		pri := "308193020100301306072a8648ce3d020106082a811ccf5501822d0479307702010104203a2b913838c60da66d9370d1b6f90f8f91d524faf695c07832b388b609bfc688a00a06082a811ccf5501822da14403420004c3a32777ca149e92fd4d3a5bd575c4b153082f43a704f120c6398e461409014f47a53a3943f8597af786027cb35085b24c33c1bd24e4e2af6786dc11925b27f6"
		s.From = "LTMXzYiLnsEF3TNnJYcjMMXsfGmLJuWzb3"
		s.To = "LTMXzYiLnsEF3TNnTAkZtu7xDkFvsQL2st"
		now := time.Now()
		s.Timestamp = fmt.Sprintf("%d", now.Unix())
		s.Amount = 10

		var km KeyManager.KeyManager
		km.Init()
		km.GenRandomKeyPair()
		km.SetPriKey(pri)
		km.SetPubkey(s.Pubkey)
		temp2, _ := json.Marshal(s)
		s.Sig = ""
		s.Sig, err = km.SignWithPriKey(KeyManager.GetHash(temp2), pri)

		err2 := rpc1.Call("PPoV.TransferMoney", s, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ := json.Marshal(basic)
		fmt.Println(string(data))

		var s1 TransferMoneyMsg
		s1.Pubkey = "3059301306072a8648ce3d020106082a811ccf5501822d034200043c0293368962d7816e3d9094e049add334f183a3e612db98ecaf0980fae58389e34292854781a272ac9292d45c82aed78e585e6bbb20640617f50f8605f769fc"
		pri = "308193020100301306072a8648ce3d020106082a811ccf5501822d047930770201010420cc1e3d90162ea3191e5a24920b180d2b5656ab145bda80324439cdd8bb6f830da00a06082a811ccf5501822da144034200043c0293368962d7816e3d9094e049add334f183a3e612db98ecaf0980fae58389e34292854781a272ac9292d45c82aed78e585e6bbb20640617f50f8605f769fc"
		s1.From = "LTMXzYiLnsEF3TNnTAkZtu7xDkFvsQL2st"
		s1.To = "LTMXzYiLnsEF3TNnJYcjMMXsfGmLJuWzb3"
		now2 := time.Now()
		s1.Timestamp = fmt.Sprintf("%d", now2.Unix())
		s1.Amount = 10

		var km2 KeyManager.KeyManager
		km2.Init()
		km2.GenRandomKeyPair()
		km2.SetPriKey(pri)
		km2.SetPubkey(s1.Pubkey)
		temp2, _ = json.Marshal(s1)
		s1.Sig = ""
		s1.Sig, err = km2.SignWithPriKey(KeyManager.GetHash(temp2), pri)

		err2 = rpc1.Call("PPoV.TransferMoney", s1, &basic)
		if err2 != nil {
			log.Fatal(err2)
		}
		data, _ = json.Marshal(basic)
		fmt.Println(string(data))
	case 21:
		var cs int
		var gs int
		fmt.Println("请输入client数量:")
		fmt.Scan(&cs)
		fmt.Println("请输入每个client发送数量:")
		fmt.Scan(&gs)

		//addressList := []string{"10.0.0.19:8010","10.0.0.20:8010","10.0.0.21:8010","10.0.0.22:8010"}
		addressList := []string{"127.0.0.1:8010","127.0.0.1:8011","127.0.0.1:8012","127.0.0.1:8010"}
		rpc0, err := jsonrpc.Dial("tcp", addressList[0])
		rpc1, err := jsonrpc.Dial("tcp", addressList[1])
		rpc2, err := jsonrpc.Dial("tcp", addressList[2])
		//rpc3, err := jsonrpc.Dial("tcp", addressList[3])

		if err != nil {
			log.Fatal(err)
		}
		var num int32 = 0
		var wg sync.WaitGroup
		all := cs * 2
		wg.Add(all)

		var messages []CreatAccountMsg
		var priList  []string
		mutex := sync.Mutex{}
		for i := 0; i < cs*2; i++ {

			go func() {
				defer wg.Done()
				for j := 0; j < gs; j++ {
					var s CreatAccountMsg
					now := time.Now()
					s.Timestamp = fmt.Sprintf("%d", now.UnixNano())

					var km KeyManager.KeyManager
					km.Init()
					km.GenRandomKeyPair()
					s.Address = km.GetAddress()
					s.Pubkey = km.GetPubkey()
					s.Sig = ""
					temp2, _ := json.Marshal(s)
					s.Sig, err = km.Sign(KeyManager.GetHash(temp2))
					mutex.Lock()
					messages = append(messages, s)
					pri, _ := km.GetPriKey()
					priList = append(priList, pri)
					mutex.Unlock()
				}
			}()
		}
		wg.Wait()

		var wg0 sync.WaitGroup
		var transMsg []TransferMoneyMsg
		var num2 int32 = 0
		wg0.Add(cs)
		for i := 0; i < cs; i++ {

			go func() {
				defer wg0.Done()
				for j := 0; j < gs; j++ {
					tmp := atomic.AddInt32(&num2, 2)
					var s TransferMoneyMsg
					now := time.Now()
					s.Timestamp = fmt.Sprintf("%d", now.UnixNano())

					var km KeyManager.KeyManager
					km.Init()
					km.SetPriKey(priList[tmp-2])
					km.SetPubkey(messages[tmp-2].Pubkey)
					s.From = messages[tmp-2].Address
					s.To = messages[tmp-1].Address
					s.Pubkey = messages[tmp-2].Pubkey
					s.Amount = 10
					s.Sig = ""
					temp2, _ := json.Marshal(s)
					s.Sig, err = km.Sign(KeyManager.GetHash(temp2))
					mutex.Lock()
					transMsg = append(transMsg, s)
					mutex.Unlock()
				}
			}()
		}
		wg0.Wait()

		wg1 := sync.WaitGroup{}
		wg1.Add(cs)

		for i := 0; i < cs; i++ {

			go func() {
				defer wg1.Done()
				for j := 0; j < gs; j++ {
					now := atomic.AddInt32(&num, 1)
					err2 := rpc0.Call("PPoV.CreatAccount", messages[now-1], &basic)
					err2 = rpc1.Call("PPoV.CreatAccount", messages[now-1], &basic)
					err2 = rpc2.Call("PPoV.CreatAccount", messages[now-1], &basic)
					//err2 = rpc3.Call("PPoV.CreatAccount", messages[now-1], &basic)
					if err2 != nil {
						log.Fatal(err2)
					}
					//data, _ := json.Marshal(basic)
					//fmt.Println(string(data))
				}
			}()
		}
		wg1.Wait()
		time.Sleep(3000 * time.Millisecond)
		wg2 := sync.WaitGroup{}
		wg2.Add(cs)
		for i := 0; i < cs; i++ {

			go func() {
				defer wg2.Done()
				for j := 0; j < gs; j++ {
					now := atomic.AddInt32(&num, 1)
					err2 := rpc0.Call("PPoV.CreatAccount", messages[now-1], &basic)
					fmt.Println(basic)
					err2 = rpc1.Call("PPoV.CreatAccount", messages[now-1], &basic)
					err2 = rpc2.Call("PPoV.CreatAccount", messages[now-1], &basic)
				//	err2 = rpc3.Call("PPoV.CreatAccount", messages[now-1], &basic)
					if err2 != nil {
						log.Fatal(err2)
					}
					//data, _ := json.Marshal(basic)
					//fmt.Println(string(data))
				}
			}()
		}
		wg2.Wait()
		now := atomic.AddInt32(&num, 1)
		fmt.Println(now-1, "个创建账户交易发送完成")

		time.Sleep(3 * time.Second)

		num = 0
		wg3 := sync.WaitGroup{}
		wg3.Add(cs)
		for i := 0; i < cs; i++ {

			go func() {
				defer wg3.Done()
				for j := 0; j < gs; j++ {
					now := atomic.AddInt32(&num, 1)
					err2 := rpc0.Call("PPoV.TransferMoney", transMsg[now-1], &basic)
					fmt.Println(basic)
					err2 = rpc1.Call("PPoV.TransferMoney", transMsg[now-1], &basic)
					err2 = rpc2.Call("PPoV.TransferMoney", transMsg[now-1], &basic)
					//err2 = rpc3.Call("PPoV.TransferMoney",transMsg[now-1], &basic)
					if err2 != nil {
						log.Fatal(err2)
					}
					//data, _ := json.Marshal(basic)
					//fmt.Println(string(data))
				}
			}()
		}
		wg3.Wait()
		now = atomic.AddInt32(&num2, 1)
		fmt.Println(now-1, "个转账发送完成")
	}

}
