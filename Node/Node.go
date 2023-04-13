package Node

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"ppov/AccountManager"
	"ppov/ConfigHelper"
	"ppov/KeyManager"
	"ppov/Message"
	"ppov/MetaData"
	"ppov/MongoDB"
	"ppov/Network"
	"ppov/utils"
	"sync"
	"time"
)

const (
	Sync    = 0
	Genesis = 1
	Normal  = 2
)

type CallBackInstance struct {
	run       func(msg Message.MessageInterface, header Message.MessageHeader)
	MsgType   int
	ChildType int
	StartTime float64
	WaitTime  float64
}

type Node struct {
	network        *Network.Network
	mongo          *MongoDB.Mongo
	keymanager     *KeyManager.KeyManager
	accountManager *AccountManager.AccountManager
	msgManager     *Message.MessagerManager
	config         *ConfigHelper.Config
	txPool         *utils.TransactionPool
	syncTool       *SynchronizeModule
	//各种控制状态
	state chan int //系统状态

	//共识状态变量
	dutyWorkerNumber   uint32 //轮值记账节点编号
	true_dutyWorkerNum uint32
	round              uint32  //当前经过轮数
	StartTime          float64 //当前共识开始时间
	Tcut               float64 //一轮的时间
	isTimeOut, isTimeout2          bool    //超时flag

	BlockGroups		 *sync.Map
	RWLock			*sync.RWMutex

	CallBackList []CallBackInstance

	BalanceTable	*sync.Map	//用户余额表
	SignResultCache	*sync.Map	//交易验签结果缓存

	//byj temp
	WorkerPubList          map[string]uint64 //记账节点公钥列表
	WorkerCandidatePubList map[string]uint64 //记账候选节点公钥列表
	VoterPubList           map[string]uint64 //投票节点公钥列表
	ElectNewWorkerList     []Message.ElectNewWorkerMsg

	GenesisBlockDone                    bool //是否生成了创世区块
	NormalGenerateBlockDone             bool
	NormalGenerateVoteDone              bool
	NormalGenerateBlockGroupHeaderDone  bool
	//统计用的变量
	TxsAmount           uint64
	TxsPeriodAmount     uint64
	StatStartTime       float64
	StatPeriodStartTime float64
	StatPeriod          float64
	PeriodTPS           float64
	TotalTPS            float64
	TotalLatencySum     float64
	PeriodLatencySum    float64
	PeriodLatency       float64
	TotalLatency        float64

	lastTime            int64

	//AddressCache       *ccache.Cache
	//SignCache          *ccache.Cache

	//wzx test
	BlockGenerateTime            time.Time
	GenerateVoteTime             time.Time
	GenerateBlockGroupHeaderTime time.Time
	CommitTime                   time.Time
}

func (node *Node) CreateBlockGroup() MetaData.BlockGroup {
	var group MetaData.BlockGroup

	group.VoteTickets = make([]MetaData.VoteTicket, node.config.VotedNum)
	group.Blocks = make([]MetaData.Block, node.config.WorkerNum)
	group.CheckTransactions = make([]int, node.config.WorkerNum)
	group.CheckHeader = make([]int, node.config.WorkerNum)
	group.ReceivedBlockGroupHeader = false
	for i, _ := range group.CheckTransactions {
		group.CheckTransactions[i] = 0
	}
	for i, _ := range group.CheckHeader {
		group.CheckHeader[i] = 0
	}
	return group
}

func (node *Node) SetConfig(config ConfigHelper.Config) {
	node.config = &config
	node.network.SetConfig(config)
	node.mongo.SetConfig(config)
	node.LoadBlockChain()
	node.keymanager.SetPriKey(config.MyPrikey)
	node.keymanager.SetPubkey(config.MyPubkey)
	node.txPool.Init(node.config.TxPoolSize, node.config.TxPoolSize)
	node.Tcut=config.Tcut
	node.msgManager.Pubkey=config.MyPubkey
	node.msgManager.ID=node.network.MyNodeInfo.ID
}

func (node *Node) LoadBlockChain(){
	fmt.Println("加载区块链")
	height:=node.mongo.QueryHeight()
	for i:=0;i<=height;i++ {
		fmt.Println("加载高度为",i,"的区块组")
		group:=node.mongo.GetBlockFromDatabase(i)
		group.ReceivedBlockGroupHeader=true
		//group.CheckTransactions = make([]int, node.config.WorkerNum)
		//group.CheckHeader = make([]int, node.config.WorkerNum)
		//for k,v:=range group.VoteResult{
		//	if v==1 {
		//		group.CheckTransactions[k]=1
		//		group.CheckHeader[k]=1
		//	}
		//}

		//同步区块时所有区块组都通过Commit函数执行提交操作，需要对创世区块组进行特殊处理
		node.mongo.Height=i
		if i == 0 {
			node.UpdateGenesisBlockVaribles(&group)
			if (len(node.state) == 1) {
				<-node.state
			}
		} else {
			//fmt.Println("高度为", i, "的区块组成功共识并保存")
			node.UpdateVariblesFromDisk(&group)
		}
	}
}

func (node *Node) Init() {
	//初始化网络模块
	node.syncTool=&SynchronizeModule{}
	node.syncTool.Init()
	node.msgManager = &Message.MessagerManager{}
	node.txPool = &utils.TransactionPool{}
	node.network = &Network.Network{
		MyNodeInfo: Network.NodeInfo{
			IP:   "",
			PORT: 0,
			ID:   0,
		},
		NodeList: make(map[Network.NodeID]Network.NodeInfo),
		NodeConnList: make(map[Network.NodeID]net.Conn),
		CB:       nil,
	}
	node.network.SetCB(node.HandleMessage)
	//初始化数据库模块
	node.mongo = &MongoDB.Mongo{}
	node.keymanager = &KeyManager.KeyManager{}
	node.keymanager.Init()
	node.keymanager.GenRandomKeyPair()
	node.state = make(chan int, 3)
	node.state <- Sync
	node.accountManager = &AccountManager.AccountManager{}
	node.accountManager.WorkerSet = make(map[string]uint64)
	node.accountManager.VoterSet = make(map[string]uint64)
	node.accountManager.WorkerCandidateSet = make(map[string]uint64)
	node.accountManager.WorkerNumberSet = make(map[uint32]string)
	node.accountManager.VoterNumberSet = make(map[uint32]string)
	node.accountManager.VoterSetNumber = make(map[string]uint32)
	node.accountManager.WorkerSetNumber = make(map[string]uint32)
	node.WorkerPubList = make(map[string]uint64)
	node.WorkerCandidatePubList = make(map[string]uint64)
	node.VoterPubList = make(map[string]uint64)
	node.GenesisBlockDone = false
	node.NormalGenerateBlockDone = false
	node.NormalGenerateVoteDone = false
	node.NormalGenerateBlockGroupHeaderDone = false
	node.StartTime = utils.GetCurrentTime()
	node.dutyWorkerNumber = 0
	//设置统计变量
	node.TxsAmount = 0
	node.TxsPeriodAmount = 0
	node.TotalTPS = 0
	node.PeriodTPS = 0
	node.TotalLatencySum = 0
	node.PeriodLatencySum = 0
	node.TotalLatency = 0
	node.PeriodLatency = 0
	//node.TotalTPS=0
	node.StatPeriod = 10
	node.StatStartTime = -1
	node.StatPeriodStartTime = -1
	node.RWLock = new(sync.RWMutex)
	node.BlockGroups = new(sync.Map)
	node.BalanceTable = new(sync.Map)
	node.SignResultCache = new(sync.Map)
	node.lastTime = time.Now().UnixNano()

	//node.AddressCache = ccache.New(ccache.Configure())
	//node.SignCache = ccache.New(ccache.Configure())
	//test code
	{
		node.BlockGenerateTime = time.Now()
		node.GenerateVoteTime = time.Now()
		node.GenerateBlockGroupHeaderTime = time.Now()
		node.CommitTime = time.Now()
	}
}

func (node *Node) Start() {
	go node.network.Start()
	//go node.StartRPC()
	go node.StartRPC()

	node.run()
}

//记账节点生成区块并发布
func (node *Node) GenerateBlock() {
	//fmt.Println("GenerateBlock start")
	//fmt.Println("node.NormalGenerateBlockDone", node.NormalGenerateBlockDone)
	if !node.NormalGenerateBlockDone {
		if utils.GetCurrentTimeMilli() - node.StartTime * 1e3 > node.config.GenerateBlockPeriod * 1e3 {
			node.BlockGenerateTime = time.Now()
			fmt.Println("commit ---> generate block ---> ", time.Since(node.CommitTime))
			time1 := time.Now()
			var block MetaData.Block
			block.Height = node.mongo.GetHeight() + 1
			pubkey := node.keymanager.GetPubkey()
			var block_num = node.GetMyWorkerNumber()
			block.BlockNum = block_num
			block.Generator = pubkey
			//设置前一区块组hash值
			//设置
			block.Transactions = node.txPool.GetCurrentTxsListDelete()
			for _, item := range block.Transactions {
				hash := KeyManager.GetHash(item)
				block.TransactionsHash = append(block.TransactionsHash, hash)
			}
			headerBytes,_:=node.mongo.Block.ToHeaderBytes(nil)
			block.PreviousHash = KeyManager.GetHash(headerBytes)
			block.MerkleRoot = KeyManager.GetHash(block.GetTransactionsBytes())
			block.Timestamp = utils.GetCurrentTime()
			block.Sig,_ = node.keymanager.Sign(append([]byte(block.MerkleRoot),block.PreviousHash... ))
			header, msg := node.msgManager.CreatePublishBlockMsg(block, 0)
			node.SendMessage(header, &msg)
			node.NormalGenerateBlockDone = true
			fmt.Println("generate", time.Since(time1))
		}
	}
}

//投票节点对区块投票并发送投票结果给轮值记账节点
func (node *Node) GenerateVote() {
	//fmt.Println("GenerateVote start")
	//fmt.Println("node.NormalGenerateVoteDone", node.NormalGenerateVoteDone)
	if !node.NormalGenerateVoteDone {
		height := node.mongo.GetHeight() + 1
		//检查是否存在当前共识所需高度的区块组
		value, existed := node.BlockGroups.Load(height)
		if !existed {
			return
		}
		item := value.(MetaData.BlockGroup)
		//设置投票结果
		var checkResult = make([]int, node.config.WorkerNum)
		if node.config.ByzantineNode && node.mongo.Height > 9 {		//拜占庭节点
			for i := 0; i < len(node.accountManager.WorkerNumberSet); i++ {
				if item.CheckTransactions[i] == 0 || item.CheckHeader[i] == 0 {
					checkResult[i] = -1
				} else {
					if item.CheckTransactions[i] == 1 && item.CheckHeader[i] == 1 {
						checkResult[i] = -1
					} else {
						checkResult[i] = 1
					}
				}
			}
		} else {														//正常节点
			for i := 0; i < len(node.accountManager.WorkerNumberSet); i++ {
				if item.CheckTransactions[i] == 0 || item.CheckHeader[i] == 0 {
					checkResult[i] = 0
				} else {
					if item.CheckTransactions[i] == 1 && item.CheckHeader[i] == 1 {
						checkResult[i] = 1
					} else {
						checkResult[i] = -1
					}
				}
			}
		}

		//在未超时时需要检查是否所有投票已经产生，超时后不需要
		if node.round==0 {
			if len(checkResult) < len(node.accountManager.WorkerNumberSet) {
				return
			}
			for _, value := range checkResult {
				if value == 0 {
					return
				}
			}
		} else {
			fmt.Println("超时投票！")
		}
		//设置各个区块hash值
		var hashes = make([][]byte, len(node.accountManager.WorkerNumberSet))
		for i, value := range checkResult {
			if value != 0 {
				block := item.Blocks[i]
				data := block.GetBlockHeaderBytes()
				hashes[i] = KeyManager.GetHash(data)
			}
		}
		/*		var hashes []string
				for _, block := range node.BlockGroups[height].Blocks {
					data := block.GetBlockHeaderBytes()
					hashes = append(hashes, KeyManager.GetHash(data))
				}*/
		node.GenerateVoteTime = time.Now()
		fmt.Println("generate block ---> generate vote --> ", time.Since(node.BlockGenerateTime))
		var ticket MetaData.VoteTicket
		ticket.VoteResult = utils.CompressIntSlice(checkResult)
		ticket.BlockHashes = hashes
		ticket.Timestamp = utils.GetCurrentTime()
		ticket.Voter = node.accountManager.VoterSetNumber[node.config.MyPubkey]
		ticket.Sig = ""
		data, _ := ticket.MarshalMsg(nil)
		data_hash := KeyManager.GetHash(data)
		ticket.Sig, _ = node.keymanager.Sign(data_hash)
		pubkey := node.accountManager.WorkerNumberSet[node.true_dutyWorkerNum]
		var receiver uint64 = node.accountManager.WorkerSet[pubkey]
		var block_num = node.GetMyVoterNumber()
		header, msg := node.msgManager.CreateNormalBlocksVoteMsg(ticket, receiver, height, block_num)
		node.SendMessage(header, &msg)
		node.NormalGenerateVoteDone = true
		//fmt.Println("GenerateVote end")
	}
}

//轮值记账节点生成区块组头部并发布
func (node *Node) GenerateBlockGroupHeader() {
	if !node.NormalGenerateBlockGroupHeaderDone {
		value, ok := node.BlockGroups.Load(node.mongo.GetHeight()+1)
		if ok {

			item := value.(MetaData.BlockGroup)
			var count int = 0
			for _, y := range item.VoteTickets {
				if y.BlockHashes != nil && y.VoteResult != nil {
					count += 1
				}
			}
			if !node.isTimeout2 && count < node.config.VotedNum {
				return
			}
			new_item, is_complete := node.VotingStatistics(item)
			if !is_complete {
				return
			}

			node.GenerateBlockGroupHeaderTime = time.Now()
			fmt.Println("generate vote --> generate header -->", time.Since(node.GenerateVoteTime))
			//var  sumTime float64 = 0
			//for _, eachVoteTic := range item.VoteTickets {
			//	sumTime += eachVoteTic.Timestamp
			//}
			new_item.NextDutyWorker = (node.dutyWorkerNumber + 1) % uint32(node.config.WorkerNum)

			new_item.Height = node.mongo.GetHeight()+1
			new_item.Generator = node.accountManager.WorkerSetNumber[node.config.MyPubkey]
			headerBytes,_:=node.mongo.Block.ToHeaderBytes(nil)
			new_item.PreviousHash = KeyManager.GetHash(headerBytes)
			new_item.Timestamp = utils.GetCurrentTime()

			signs := make([]string, 0, len(new_item.VoteTickets))
			for i := 0; i < len(new_item.VoteTickets); i++ {
				if new_item.VoteTickets[i].Sig != ""{
					signs = append(signs, new_item.VoteTickets[i].Sig)
					new_item.VoteTickets[i].Sig = ""
				}
			}

			aggSign,_ := node.keymanager.AggregateSign(signs)
			new_item.VoteAggSign = []byte(aggSign)

			var temp_header MetaData.BlockGroup
			temp_header=new_item
			temp_header.VoteTickets=nil
			tempHeaderBytes,_:=temp_header.ToHeaderBytes(nil)
			tempHeaderBytes_hash := KeyManager.GetHash(tempHeaderBytes)
			new_item.Sig,_ = node.keymanager.Sign(tempHeaderBytes_hash)

			//node.SendTransactionMsgToManagementServer(new_item)		//Front end
			//{
			//	jsonStr, _ := json.Marshal(new_item)
			//	fmt.Println(string(jsonStr))
			//}

			var blockgroupheadermsg Message.BlockGroupHeader
			blockgroupheadermsg.Data, _ = new_item.ToHeaderBytes(nil)

			//var group MetaData.BlockGroup
			//group.FromHeaderBytes(blockgroupheadermsg.Data)

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = node.network.MyNodeInfo.ID
			msgheader.Receiver = 0
			msgheader.Pubkey = node.config.MyPubkey
			msgheader.MsgType = Message.BlockGroupHeaderMsg
			node.SendMessage(msgheader, &blockgroupheadermsg)
			node.NormalGenerateBlockGroupHeaderDone = true
			//fmt.Println("GenerateBlockGroupHeader",new_item.Height,"done!")

		}
	}
}

//区块组提交
func (node *Node) Commit() {
	node.RWLock.RLock()
	height:=node.mongo.GetHeight() + 1
	value, ok := node.BlockGroups.Load(height)
	node.RWLock.RUnlock()
	flag :=true
	if ok {
		item := value.(MetaData.BlockGroup)
		deVoteResult := utils.DecompressToIntSlice(item.VoteResult)
		if item.ReceivedBlockGroupHeader {
			for k,v:=range deVoteResult{
				if v ==1 {
					if !item.Blocks[k].IsSet{
						flag=false
						break
					} else {
						block := item.Blocks[k]
						data := block.GetBlockHeaderBytes()
						if !bytes.Equal(item.BlockHashes[k],KeyManager.GetHash(data))  {
							item.Blocks[k].IsSet = false
							flag = false
							break
						}
					}
				}
			}
			if flag {
				signs := make([]string,0, len(item.Blocks))
				for i := 0 ; i< len(item.Blocks); i++{
					if item.Blocks[i].Sig != ""{
						signs = append(signs, item.Blocks[i].Sig)
					}
					item.Blocks[i].Sig = ""
				}
				aggSign, err := node.keymanager.AggregateSign(signs)
				if err != nil{
					panic("agg block sign failed:"+ err.Error())
				}
				item.BlockAggSign = []byte(aggSign)

				//同步区块时所有区块组都通过Commit函数执行提交操作，需要对创世区块组进行特殊处理
				node.CommitTime = time.Now()
				if node.accountManager.WorkerNumberSet[node.true_dutyWorkerNum] == node.keymanager.GetPubkey(){
					fmt.Println("generate header --> commit", time.Since(node.GenerateBlockGroupHeaderTime))
				}else{
					fmt.Println("generate vote --> commit", time.Since(node.GenerateVoteTime))
				}

				if height == 0{
					node.UpdateGenesisBlockVaribles(&item)
					if len(node.state) == 1 {
						<-node.state
					}
				} else {
					//node.SendNormalMsgToManagementServer()
					time1 := time.Now()
					node.UpdateVaribles(&item)
					fmt.Println("commit ", time.Since(time1))
				} //更新变量
				node.mongo.PushbackBlockToDatabase(item) //数据落盘

				node.NormalGenerateVoteDone = false
				node.NormalGenerateBlockDone = false
				node.NormalGenerateBlockGroupHeaderDone = false
				node.isTimeout2=false
				node.round = 0
				tmp := node.lastTime
				node.lastTime = time.Now().UnixNano()
				if len(item.ExecutionResult) > 0 || height == 1{
				tps :=  float32(len(item.ExecutionResult) * 1000000000)/float32(node.lastTime - tmp)

				sum :=0
				for _,value := range item.ExecutionResult{
					if value == true{
						sum++
					}
				}

				var u8 uint8
				fmt.Println("height:",utils.SizeStruct(item.Height))
				fmt.Println("generator",utils.SizeStruct(u8))
				fmt.Println("previoushash",utils.SizeStruct(item.PreviousHash))
				fmt.Println("merkleroot",utils.SizeStruct(item.Blocks[0].MerkleRoot))
				fmt.Println("timestamp",utils.SizeStruct(item.Timestamp))
				fmt.Println("voteaggsign",utils.SizeStruct(item.VoteAggSign))
				fmt.Println("blockaggsign",utils.SizeStruct(item.BlockAggSign))
				fmt.Println("高度为", height, "的区块组成功保存", len(item.ExecutionResult) , "个交易", "true :", sum)
				fmt.Println("高度为" , height , "的区块组生成中的TPS为" , tps)
				}
				fmt.Println("高度为", height, "的区块组成功保存")
			}
		} else {
			//fmt.Println("ReceivedBlockGroupHeader false")
			flag=false
		}
	}
	if !flag {
		if node.round>0 && node.isTimeOut {
			//fmt.Println("执行超时请求操作")
			node.NormalTimeOutProcess()
			node.isTimeOut=false
		}
	}
}

func (node *Node) SynchronizeBlockGroup() {
	//node.state <- Genesis
	node.syncBlockGroupSequence()
	node.StartTime = utils.GetCurrentTime()
}

func (node *Node) QueryAllPubkey() { //only for 生成创世节点的节点
	var mh Message.MessageHeader //消息头
	mh.MsgType = Message.QueryPubkey
	mh.Sender = node.network.MyNodeInfo.ID
	mh.Receiver = 0
	var gm Message.QueryPubkeyMsg //消息体
	gm.Type = 100
	node.SendMessage(mh, &gm)
}

func (node *Node) GenerateGenesisBlockGroup() {
	if node.config.MyAddress == node.config.WorkerList[0] {
		if !node.GenesisBlockDone {
			node.QueryAllPubkey() //请求公钥
			time.Sleep(time.Second)
			node.RWLock.Lock()

			if !(len(node.WorkerPubList) == len(node.config.WorkerList) &&
				len(node.WorkerCandidatePubList) == len(node.config.WorkerCandidateList) &&
				len(node.VoterPubList) == len(node.config.VoterList)) {
				node.state <- Genesis
				node.RWLock.Unlock()
				return
			}
			var genesisTransaction MetaData.GenesisTransaction //创世交易
			genesisTransaction.WorkerNum = node.config.WorkerNum
			genesisTransaction.VotedNum = node.config.VotedNum
			genesisTransaction.BlockGroupPerCycle = node.config.BlockGroupPerCycle
			genesisTransaction.Tcut = node.config.Tcut
			genesisTransaction.WorkerPubList = node.WorkerPubList
			genesisTransaction.WorkerCandidatePubList = node.WorkerCandidatePubList
			genesisTransaction.VoterPubList = node.VoterPubList

			for key, _ := range genesisTransaction.WorkerPubList {
				genesisTransaction.WorkerSet = append(genesisTransaction.WorkerSet, key)
			}

			for key, _ := range genesisTransaction.VoterPubList {
				genesisTransaction.VoterSet = append(genesisTransaction.VoterSet, key)
			}

			var transactionHeader MetaData.TransactionHeader //交易头
			transactionHeader.TXType = MetaData.Genesis

			var block MetaData.Block //区块
			block.Height = 0
			block.Generator = node.config.MyPubkey
			block.Transactions = append(block.Transactions, MetaData.EncodeTransaction(transactionHeader, &genesisTransaction))
			block.MerkleRoot = KeyManager.GetHash(block.GetTransactionsBytes())
			var blockgroup MetaData.BlockGroup //区块组
			blockgroup.Height = 0
			blockgroup.Generator = node.accountManager.WorkerSetNumber[node.config.MyPubkey]
			blockgroup.Timestamp = utils.GetCurrentTime()
			blockgroup.Blocks = append(blockgroup.Blocks, block)
			temp, _ := blockgroup.ToHeaderBytes(nil)
			temp_hash := KeyManager.GetHash(temp)
			blockgroup.Sig, _ = node.keymanager.Sign(temp_hash)

			var blockmsg Message.BlockMsg //消息体
			blockmsg.Data, _ = blockgroup.ToBytes(nil)

			var msgheader Message.MessageHeader //消息头
			msgheader.Sender = node.network.MyNodeInfo.ID
			msgheader.Receiver = 0
			msgheader.Pubkey = node.config.MyPubkey
			msgheader.MsgType = Message.GenesisBlock

			node.SendMessage(msgheader, &blockmsg)
			node.GenesisBlockDone = true //修改状态
			node.RWLock.Unlock()
			return
		}
	}
}

func (node *Node) Normal() {
	var height = node.mongo.GetHeight() + 1
	//删除过时的区块组
	if height % 199 == 0 {
		node.BlockGroups.Range(func(k, _ interface{}) bool {
			if k.(int) < height-10 {
				node.BlockGroups.Delete(k)
			}
			return true
		})
		node.SignResultCache = new(sync.Map)
	}
	////若还没创建当前所需共识的区块组，则创建一个
	for h:= height; h < height + 2; h++ {
		node.BlockGroups.LoadOrStore(h,node.CreateBlockGroup())
	}
	//更新当前轮数和轮值记账节点
	round := uint32((utils.GetCurrentTime() - node.StartTime) / node.Tcut)
	if round != node.round {
		node.isTimeOut = true
		node.isTimeout2 = true
		node.NormalGenerateVoteDone=false
		node.round=round
		fmt.Println("轮数切换为",round)
	}
	//经过3轮超时后才进行轮值记账节点的更换
	if round >= 3 {
		node.true_dutyWorkerNum = (node.dutyWorkerNumber + round - 1) % uint32(node.config.WorkerNum)
	} else {
		node.true_dutyWorkerNum = node.dutyWorkerNumber
	}
	//判断自己的身份
	pubkey := node.keymanager.GetPubkey()
	_, isWorkerCandidate := node.accountManager.WorkerCandidateSet[pubkey]
	_, isWorker := node.accountManager.WorkerSet[pubkey]
	_, isVoter := node.accountManager.VoterSet[pubkey]
	//所有节点都需要执行Commit提交区块组数据

	node.Commit()
	//在主线程检查区块头是否正确
	node.CheckBlocksHeader()
	//投票节点执行生成投票过程
	if isVoter {
		node.GenerateVote()
	}
	//记账节点生成区块
	if isWorker {
		node.GenerateBlock()
		//轮值记账节点产生区块组头
		//fmt.Println("node.true_dutyWorkerNum=",node.true_dutyWorkerNum)
		duty_pubkey, _ := node.accountManager.WorkerNumberSet[node.true_dutyWorkerNum]
		if duty_pubkey == pubkey {
			node.GenerateBlockGroupHeader()
		}
	}
	//候选记账节点不需要做任何事情
	if isWorkerCandidate {

	}
}

func (node *Node) run() {
	for {
		time.Sleep(50 * time.Millisecond)
		state := <-node.state
		switch state {
		case Sync:
			node.SynchronizeBlockGroup()
		case Genesis:
			node.GenerateGenesisBlockGroup()
		case Normal:
			node.Normal()
			node.state <- Normal
		}
	}
}

func (node *Node) SendMessage(header Message.MessageHeader, messageInterface Message.MessageInterface) {
	data, err := messageInterface.ToByteArray()
	if err != nil {
		fmt.Println("SendMessage：messageInterface.ToByteArray()错误！")
	}
	header.Data = data
	b, err := header.MarshalMsg(nil)
	if err != nil {
		fmt.Println("SendMessage：header.MarshalMsg()错误！")
	}
	go node.network.SendMessage(b, header.Receiver)
}

func (node *Node) HandleMessage(data []byte, conn net.Conn) {
	var header Message.MessageHeader
	data, err := header.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("HandleMessage：header.UnmarshalMsg(data)错误！")
		return
	}
	data = header.Data
	switch header.MsgType {
	case Message.Zero:
		var msg Message.ZeroMsg
		data, err = msg.UnmarshalMsg(data)
	case Message.QueryPubkey:
		var msg Message.QueryPubkeyMsg
		data, err = msg.UnmarshalMsg(data)
		node.HandleQueryPubMessage(msg, header.Sender)
	case Message.GenesisBlock:
		var msg Message.BlockMsg
		data, err = msg.UnmarshalMsg(data)
		node.HandleGenesisBlockPublishMessage(msg, header)
	case Message.TransactionMsg:
		var msg Message.TransactionMessage
		data, err = msg.UnmarshalMsg(data)
		node.HandleTransactionMessage(msg, header)
	case Message.NormalPublishBlock:
		var msg Message.PublishBlockMsg
		err = msg.FromByteArray(data)
		//fmt.Println(conn.RemoteAddr().String())
		node.HandlePublishBlockMessage(msg, header)
	case Message.NormalBlockVoteMsg:
		var msg Message.NormalBlocksVoteMsg
		err = msg.FromByteArray(data)
		node.HandleVoteTicketMessage(msg, header)
	case Message.RequestHeight:
		var msg Message.RequestHeightMsg
		err = msg.FromByteArray(data)
		node.HandleRequestHeightMessage(msg, header)
	case Message.RespondHeight:
		var msg Message.RequestHeightMsg
		err = msg.FromByteArray(data)
		node.HandleRespondHeightMessage(msg, header)
	case Message.RequestBlockGroup:
		var msg Message.RequestBlockGroupMsg
		err = msg.FromByteArray(data)
		node.HandleRequestBlockGroupMessage(msg, header)
	case Message.RespondBlockGroup:
		var msg Message.RespondBlockGroupMsg
		err = msg.FromByteArray(data)
		node.HandleRespondBlockGroupMessage(msg, header)
	case Message.RequestBlockGroupHeader:
		var msg Message.RequestBlockGroupHeaderMsg
		err = msg.FromByteArray(data)
		node.HandleRequestBlockGroupHeaderMessage(msg, header)
	case Message.ResponseBlockGroupHeader:
		var msg Message.RespondBlockGroupHeaderMsg
		err = msg.FromByteArray(data)
		node.HandleRespondBlockGroupHeaderMessage(msg, header)
	case Message.RequestBlock:
		var msg Message.RequestBlockMsg
		err = msg.FromByteArray(data)
		node.HandleRequestBlockMessage(msg, header)
	case Message.ResponseBlock:
		var msg Message.RespondBlockMsg
		err = msg.FromByteArray(data)
		node.HandleRespondBlockMessage(msg, header)
	case Message.ElectNewWorker:
		var msg Message.ElectNewWorkerMsg
		err = msg.FromByteArray(data)
		node.HandleElectNewWorkerMessage(msg, header)
	case Message.BlockGroupHeaderMsg:
		var msg Message.BlockGroupHeader
		err = msg.FromByteArray(data)
		node.HandleBlockGroupHeaderMessage(msg, header)

	}
}

func (node *Node) HandleQueryPubMessage(msg Message.QueryPubkeyMsg, pre_sender Network.NodeID) {
	if msg.Type == 100 {
		var mh Message.MessageHeader
		mh.MsgType = Message.QueryPubkey
		mh.Sender = node.network.MyNodeInfo.ID
		mh.Receiver = pre_sender
		var gm Message.QueryPubkeyMsg
		gm.Type = 101
		gm.Information = node.config.MyPubkey
		gm.NodeID = node.network.MyNodeInfo.ID
		node.SendMessage(mh, &gm)

	} else if msg.Type == 101 {
		node.RWLock.Lock()

		temp, ok := node.network.NodeList[pre_sender]
		if !ok {
			fmt.Println("HandleQueryPubMessage Wrong")
			node.RWLock.Unlock()
			return
		}
		for _, x := range node.config.WorkerList {
			if temp.IP == x.IP && temp.PORT == x.Port {
				_, ok := node.WorkerPubList[msg.Information]
				if !ok {
					node.WorkerPubList[msg.Information] = msg.NodeID
				}
				break
			}
		}
		for _, x := range node.config.WorkerCandidateList {
			if temp.IP == x.IP && temp.PORT == x.Port {
				_, ok := node.WorkerCandidatePubList[msg.Information]
				if !ok {
					node.WorkerCandidatePubList[msg.Information] = msg.NodeID
				}
				break
			}
		}
		for _, x := range node.config.VoterList {
			if temp.IP == x.IP && temp.PORT == x.Port {
				_, ok := node.VoterPubList[msg.Information]
				if !ok {
					node.VoterPubList[msg.Information] = msg.NodeID
				}
				break
			}
		}
		node.RWLock.Unlock()
	}
}

func (node *Node) HandleGenesisBlockPublishMessage(msg Message.BlockMsg, header Message.MessageHeader) {
	var bg MetaData.BlockGroup
	bg.FromBytes(msg.Data)

	//{
	//	jsonStr,_ := json.Marshal(bg)
	//	fmt.Println(string(jsonStr))
	//}

	sig := bg.Sig
	bg.Sig = ""
	temp, _ := bg.ToHeaderBytes(nil)
	temp_hash := KeyManager.GetHash(temp)

	if genePubkey, ok := node.accountManager.WorkerNumberSet[bg.Generator]; bg.Height != 0 && ok == false{
		return
	}else{
		if bg.Height == 0{
			bg.CheckHeader = []int{1}
			node.mongo.PushbackBlockToDatabase(bg)
			node.UpdateGenesisBlockVaribles(&bg)
			fmt.Println("高度为0的区块组成功保存")
			return
		}
		ok, err := node.keymanager.Verify(temp_hash,sig,genePubkey)
		if err != nil {
			fmt.Println(err)
		}
		if ok {
			node.mongo.PushbackBlockToDatabase(bg)
			node.UpdateGenesisBlockVaribles(&bg)
			fmt.Println("高度为0的区块组成功保存")
		}
	}

}

func (node *Node) HandleTransactionMessage(msg Message.TransactionMessage, header Message.MessageHeader) {
	head, transactionInterface := MetaData.DecodeTransaction(msg.Data)
	switch head.TXType {
	case MetaData.CreatACCOUNT:
		if transaction, ok := transactionInterface.(*MetaData.CreatAccount); ok {
			txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(msg.Data))
			if node.ValidateCreatAccountTransaction(transaction) {
				node.SignResultCache.Store(txhash,1)
			} else {
				node.SignResultCache.Store(txhash,0)
			}
			if _,ok := node.SignResultCache.Load(txhash); !ok{
				node.txPool.PushbackTransaction(head, transactionInterface)
			}
		}
	case MetaData.TransferMONEY:
		if transaction, ok := transactionInterface.(*MetaData.TransferMoney); ok {
			txhash := base64.StdEncoding.EncodeToString(KeyManager.GetHash(msg.Data))
			if node.ValidateTransferMoneyTransaction(transaction) {
				node.SignResultCache.Store(txhash,1)
			} else {
				node.SignResultCache.Store(txhash,0)
			}
			if _,ok := node.SignResultCache.Load(txhash); !ok{
				node.txPool.PushbackTransaction(head, transactionInterface)
			}
		}
	default:
		node.txPool.PushbackTransaction(head, transactionInterface)
	}
}

func (node *Node) HandlePublishBlockMessage(msg Message.PublishBlockMsg, header Message.MessageHeader) {
	height := msg.GetHeight()
	blockNum := msg.GetBlockNum()
	if height <= node.mongo.GetHeight() {
		return
	}

	//增加稳定性，防止程序崩溃
	if int(blockNum) >= node.config.WorkerNum {
		return
	}

	value,_ := node.BlockGroups.LoadOrStore(height,node.CreateBlockGroup())
	block := msg.GetBlock()
	block.IsSet = true

	item := value.(MetaData.BlockGroup)
	item.Blocks[blockNum] = block
	if !node.ValidateTransactions(&block.Transactions) {
		item.CheckTransactions[blockNum] = -1
	} else {
		item.CheckTransactions[blockNum] = 1
	}

	ok,_ := node.keymanager.Verify([]byte(append(block.MerkleRoot, block.PreviousHash...)), block.Sig, block.Generator)
	if !ok{
		return
	}
	node.BlockGroups.Store(height,item)
}

func (node *Node) HandleVoteTicketMessage(msg Message.NormalBlocksVoteMsg, header Message.MessageHeader) {
	height := msg.Height
	blockNum:=msg.BlockNum

	if height <= node.mongo.GetHeight() {
		return
	}

	var ticket MetaData.VoteTicket
	ticket.VoteResult = msg.Ticket.VoteResult
	ticket.BlockHashes = msg.Ticket.BlockHashes
	ticket.Timestamp = msg.Ticket.Timestamp
	ticket.Voter = msg.Ticket.Voter
	ticket.Sig = ""
	data, _ := ticket.MarshalMsg(nil)
	data_hash := KeyManager.GetHash(data)

	if voterPub, ok := node.accountManager.VoterNumberSet[msg.Ticket.Voter]; ok == false{
		return
	}else{
		verifyRes, err := node.keymanager.Verify(data_hash, msg.Ticket.Sig, voterPub)
		if err != nil || verifyRes == false{
			return
		}
	}


	//{
	//	fmt.Println("verifying true")
	//}

	//增加稳定性，防止程序崩溃
	if int(blockNum) >= node.config.VotedNum {
		return
	}
	//如果不存在BlockGroup，则创建一个
	value, _ := node.BlockGroups.LoadOrStore(height,node.CreateBlockGroup())

	item := value.(MetaData.BlockGroup)
	item.VoteTickets[blockNum] = msg.Ticket
	node.BlockGroups.Store(height,item)
}

func (node *Node) HandleRequestHeightMessage(msg Message.RequestHeightMsg, header Message.MessageHeader) {
	header, msg = node.msgManager.CreateRespondHeightMsg(header.Sender, node.mongo.GetHeight())
	node.SendMessage(header, &msg)
}

func (node *Node) HandleRespondHeightMessage(msg Message.RequestHeightMsg, header Message.MessageHeader) {
	height := msg.Height
	if height > node.mongo.GetHeight() {
		var s Syncer
		s.NodeID = header.Sender
		s.Height = height
		node.syncTool.Syncers <- s
		fmt.Println("收到height=",s.Height,"NodeID=",s.NodeID,"的高度回复")
	}

}

func (node *Node) HandleRequestBlockGroupMessage(msg Message.RequestBlockGroupMsg, header Message.MessageHeader) {
	height := msg.Height
	if node.mongo.GetHeight() < height {
		return
	}
	group := node.mongo.GetBlockFromDatabase(height)
	header, response := node.msgManager.CreateRespondBlockGroupMsg(header.Sender, height, group)
	node.SendMessage(header, &response)
}

func (node *Node) HandleRespondBlockGroupMessage(msg Message.RespondBlockGroupMsg, header Message.MessageHeader) {
	fmt.Println("接收到区块组回复")
	msg.Group.ReceivedBlockGroupHeader=true
	msg.Group.CheckTransactions = make([]int, node.config.WorkerNum)
	msg.Group.CheckHeader = make([]int, node.config.WorkerNum)
	deVoteResult := utils.DecompressToIntSlice(msg.Group.VoteResult)
	for k,v:=range deVoteResult{
		if int(v)==1 {
			msg.Group.CheckTransactions[k]=1
			msg.Group.CheckHeader[k]=1
			msg.Group.Blocks[k].IsSet = true
		}
	}
	node.RWLock.Lock()
	node.BlockGroups.Store(msg.Group.Height,msg.Group)
	node.RWLock.Unlock()
}

func (node *Node) HandleElectNewWorkerMessage(msg Message.ElectNewWorkerMsg, header Message.MessageHeader) {
	node.ElectNewWorkerList = append(node.ElectNewWorkerList, msg)
}

func (node *Node) HandleBlockGroupHeaderMessage(msg Message.BlockGroupHeader, header Message.MessageHeader) {
	var blockgroup_header, newBGHeader MetaData.BlockGroup
	blockgroup_header.FromHeaderBytes(msg.Data)
	newBGHeader = blockgroup_header
	newBGHeader.VoteTickets = nil
	newBGHeader.Sig = ""
	temp, _ := newBGHeader.ToHeaderBytes(nil)
	temp_hash := KeyManager.GetHash(temp)

	if genePubkey, ok := node.accountManager.WorkerNumberSet[blockgroup_header.Generator]; ok == false{
		return
	}else{
		ok, err := node.keymanager.Verify(temp_hash,blockgroup_header.Sig,genePubkey)
		if err != nil {
			fmt.Println(err)
		}
		if ok == false{
			return
		}
	}

		if blockgroup_header.Height >= node.mongo.Height + 1 {
			value, _ := node.BlockGroups.LoadOrStore(blockgroup_header.Height,node.CreateBlockGroup())

			item := value.(MetaData.BlockGroup)

			if node.mongo.Height != 0{
				pubs := make([]string, 0, len(blockgroup_header.VoteTickets))
				msgs := make([][]byte, 0, len(blockgroup_header.VoteTickets))
				for i := 0; i < len(blockgroup_header.VoteTickets); i++ {
					if voterPub, ok := node.accountManager.VoterNumberSet[blockgroup_header.VoteTickets[i].Voter]; ok == false{
						fmt.Println("The voter is illegal")
						return
					}else{
						if blockgroup_header.VoteTickets[i].BlockHashes == nil{
							continue
						}
						var ticket MetaData.VoteTicket
						ticket.VoteResult = blockgroup_header.VoteTickets[i].VoteResult
						ticket.BlockHashes = blockgroup_header.VoteTickets[i].BlockHashes
						ticket.Timestamp = blockgroup_header.VoteTickets[i].Timestamp
						ticket.Voter = blockgroup_header.VoteTickets[i].Voter
						ticket.Sig = ""
						data, _ := ticket.MarshalMsg(nil)
						data_hash := KeyManager.GetHash(data)

						pubs = append(pubs, voterPub)
						msgs = append(msgs, data_hash)

					}
				}

				verRes, err := node.keymanager.VerifyAggSign(string(blockgroup_header.VoteAggSign), pubs, msgs)
				if err != nil || verRes == false{
					fmt.Println("Check Block group header failed", err)
					return
				}
			}

			item.Height=blockgroup_header.Height
			item.Generator = blockgroup_header.Generator
			item.PreviousHash = blockgroup_header.PreviousHash
			item.MerkleRoot = blockgroup_header.MerkleRoot
			item.Timestamp = blockgroup_header.Timestamp
			item.Sig = blockgroup_header.Sig
			item.NextDutyWorker = blockgroup_header.NextDutyWorker
			item.BlockHashes = blockgroup_header.BlockHashes
			item.VoteResult = blockgroup_header.VoteResult
			item.VoteTickets=blockgroup_header.VoteTickets
			item.ReceivedBlockGroupHeader=true
			item.VoteAggSign = blockgroup_header.VoteAggSign
			//fmt.Println("HandleBlockGroupHeaderMessage change variable", blockgroup_header.Height)

			node.BlockGroups.Store(blockgroup_header.Height,item)
		}
}

func (node *Node) HandleRequestBlockGroupHeaderMessage(msg Message.RequestBlockGroupHeaderMsg, header Message.MessageHeader) {
	//fmt.Println("接收到高度为",msg.Height,"的区块组头请求")
	if msg.Height >= node.mongo.GetHeight() {
		value, ok := node.BlockGroups.Load(msg.Height)
		if ok {
			group := value.(MetaData.BlockGroup)
			if group.ReceivedBlockGroupHeader{
				//fmt.Println(node.network.MyNodeInfo.ID,"在内存中找到并发送高度为",msg.Height,"的区块组头")
				RespHeader, RespMsg := node.msgManager.CreateRespondBlockGroupHeaderMsg(header.Sender, msg.Height, &group)
				node.SendMessage(RespHeader, &RespMsg)
			}
		} else {
			//fmt.Println(node.network.MyNodeInfo.ID,"在内存中找不到高度为",msg.Height,"的区块组")
		}

	} else {
		group := node.mongo.GetBlockFromDatabase(msg.Height)
		RespHeader, RespMsg := node.msgManager.CreateRespondBlockGroupHeaderMsg(header.Sender, msg.Height, &group)
		//fmt.Println(node.network.MyNodeInfo.ID,"在数据库中找到并发送高度为",msg.Height,"的区块组头")
		node.SendMessage(RespHeader, &RespMsg)
	}
}

func (node *Node) HandleRespondBlockGroupHeaderMessage(msg Message.RespondBlockGroupHeaderMsg, header Message.MessageHeader) {
	var blockgroup_header MetaData.BlockGroup
	blockgroup_header.FromHeaderBytes(msg.BlockGroupHeaderBytes)
	newBGHeader := blockgroup_header
	newBGHeader.VoteTickets = nil
	newBGHeader.Sig = ""
	temp, _ := newBGHeader.ToHeaderBytes(nil)
	temp_hash := KeyManager.GetHash(temp)
	if genePubkey, ok := node.accountManager.WorkerNumberSet[blockgroup_header.Generator]; ok == false{
		return
	}else{
		ok, err := node.keymanager.Verify(temp_hash,blockgroup_header.Sig,genePubkey)
		if err != nil {
			fmt.Println(err)
		}
		if ok == false{
			return
		}
	}
		if blockgroup_header.Height >= node.mongo.Height + 1 {
			value, _ := node.BlockGroups.LoadOrStore(blockgroup_header.Height,node.CreateBlockGroup())

			item := value.(MetaData.BlockGroup)
			item.Height=blockgroup_header.Height
			item.Generator = blockgroup_header.Generator
			item.PreviousHash = blockgroup_header.PreviousHash
			item.MerkleRoot = blockgroup_header.MerkleRoot
			item.Timestamp = blockgroup_header.Timestamp
			item.Sig = blockgroup_header.Sig
			item.NextDutyWorker = blockgroup_header.NextDutyWorker
			item.BlockHashes = blockgroup_header.BlockHashes
			item.VoteResult = blockgroup_header.VoteResult
			item.VoteTickets=blockgroup_header.VoteTickets
			item.ReceivedBlockGroupHeader=true

			node.BlockGroups.Store(blockgroup_header.Height,item)
		}
}

func (node *Node) HandleRequestBlockMessage(msg Message.RequestBlockMsg, header Message.MessageHeader) {
	if msg.Height >= node.mongo.GetHeight()+1 {
		value, ok := node.BlockGroups.Load(msg.Height)
		group := value.(MetaData.BlockGroup)
		if ok {
			if group.Blocks[msg.BlockNum].IsSet {
				RespHeader, RespMsg := node.msgManager.CreateRequestBlockMsg(header.Sender, msg.Height, msg.BlockNum)
				node.SendMessage(RespHeader, &RespMsg)
			}
		}
	} else {
		group := node.mongo.GetBlockFromDatabase(msg.Height)

		RespHeader, RespMsg := node.msgManager.CreateRespondBlockMsg(header.Sender, msg.Height, msg.BlockNum, group.Blocks[msg.BlockNum])
		node.SendMessage(RespHeader, &RespMsg)
	}
}

func (node *Node) HandleRespondBlockMessage(msg Message.RespondBlockMsg, header Message.MessageHeader) {
	height := msg.Height
	blockNum := msg.BlockNum
	if height <= node.mongo.GetHeight()  {
		return
	}
	//增加稳定性，防止程序崩溃
	if int(blockNum) >= node.config.WorkerNum {
		return
	}
	value, _ := node.BlockGroups.LoadOrStore(height,node.CreateBlockGroup())

	block := msg.Block
	block.IsSet = true
	item := value.(MetaData.BlockGroup)
	item.Blocks[blockNum] = block
	if !node.ValidateTransactions(&block.Transactions) {
		item.CheckTransactions[blockNum] = -1
	} else {
		item.CheckTransactions[blockNum] = 1
	}
	node.BlockGroups.Store(height,item)
}

func (node *Node) GetMyWorkerNumber() uint32 {
	pubkey := node.keymanager.GetPubkey()
	var block_num uint32
	var find = false
	for block_num = 0; block_num < uint32(len(node.accountManager.WorkerNumberSet)); block_num++ {
		key, ok := node.accountManager.WorkerNumberSet[block_num]
		if ok && key == pubkey {
			find = true
			break
		}
	}
	if !find {
		fmt.Println("GenerateBlock--找不到记账节点编号")
	}
	return block_num
}

func (node *Node) GetMyVoterNumber() uint32 {
	pubkey := node.keymanager.GetPubkey()
	var block_num uint32
	var find = false
	for block_num = 0; block_num < uint32(len(node.accountManager.VoterNumberSet)); block_num++ {
		key, ok := node.accountManager.VoterNumberSet[block_num]
		if ok && key == pubkey {
			find = true
			break
		}
	}
	if !find {
		fmt.Println("GenerateBlock--找不到记账节点编号")
	}
	return block_num
}

func (node *Node) CheckBlocksHeader() {
	height := node.mongo.GetHeight() + 1
	value, _ := node.BlockGroups.LoadOrStore(height,node.CreateBlockGroup())
	item := value.(MetaData.BlockGroup)
	for i, value := range item.CheckTransactions {
		if value != 0 && item.CheckHeader[i] == 0 {
			if node.ValidateBlockHeader(&item.Blocks[i]) {
				item.CheckHeader[i] = 1
			} else {
				item.CheckHeader[i] = -1
			}
		}
	}
	node.BlockGroups.Store(height,item)
}

