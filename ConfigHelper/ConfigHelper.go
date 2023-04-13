package ConfigHelper

import (
	"fmt"
	"strconv"
	_ "strconv"
	"strings"

	"github.com/larspensjo/config"
	"ppov/KeyManager"
)

type AddressPair struct {
	IP   string
	Port int
}
type Config struct {
	WorkerList          []AddressPair
	WorkerCandidateList []AddressPair
	VoterList           []AddressPair
	MyAddress           AddressPair //单台服务器第一个节点地址
	ServicePort			int
	PubkeyList          []string
	PrikeyList          []string
	SingleServerNodeNum int    //单台服务器运行节点数量
	MyPubkey            string //公钥
	MyPrikey            string //私钥
	ManagementServerIP	string	//管理服务器IP
	ManagementServerPort	int	//管理服务器Port
	GenesisDutyWorker   int    //代理记账节点
	WorkerNum           int    //记账节点数量
	VotedNum            int    //记账节点数量投票节点数量

	TxPoolSize			int		//交易池大小
	BlockGroupPerCycle int     //每个公式周期产生的区块组数量
	Tcut               float64 //每轮时间
	GenerateBlockPeriod		float64		//区块生成周期
	DropDatabase		bool	//清空数据库
	ByzantineNode		bool	//拜占庭节点
	Hostname 			string
	RouterIP string
	RouterPort int

}

func CreateConfigs(){
	var conf Config
	var servers=[]AddressPair{
		AddressPair{"127.0.0.1", 5010},
		AddressPair{"127.0.0.1", 5011},
		AddressPair{"127.0.0.1", 5012},
		AddressPair{"127.0.0.1", 5013}}
	conf.SingleServerNodeNum = 1
	conf.ServicePort = 8010
	conf.ManagementServerIP = "127.0.0.1"
	conf.ManagementServerPort = 9999
	//var total_nodes=conf.SingleServerNodeNum* len(servers)
	for i:=0;i<len(servers);i++{
		for j:=0;j<conf.SingleServerNodeNum;j++ {
			var addr = AddressPair{servers[i].IP, servers[i].Port+j}
			conf.WorkerList=append(conf.WorkerList, addr)
			conf.WorkerCandidateList=append(conf.WorkerCandidateList, addr)
			conf.VoterList=append(conf.VoterList, addr)
		}
	}
	conf.GenesisDutyWorker = 0
	conf.WorkerNum = 4
	conf.VotedNum = 4
	conf.BlockGroupPerCycle = 100000
	conf.Tcut = 5
	conf.GenerateBlockPeriod = 0.5
	conf.TxPoolSize = 10000
	conf.ByzantineNode = false
	conf.DropDatabase = true

	var hostname = []string{"主机19", "主机20", "主机21","主机22"}


	//对每个服务器进行特定化配置
	for i:=0 ; i< len(servers);i++{
		var config=conf
		config.MyAddress = servers[i]
		config.ServicePort = conf.ServicePort	//change in need
		config.Hostname = hostname[i]
		var config_name="config_"+servers[i].IP+"_"+strconv.Itoa(servers[i].Port)+"_"+strconv.Itoa(config.SingleServerNodeNum)
		for i := 0; i < config.SingleServerNodeNum; i++ {
			var keyManager KeyManager.KeyManager
			keyManager.Init()
			keyManager.GenRandomKeyPair()
			config.PubkeyList = append(config.PubkeyList, keyManager.GetPubkey())
			key, _ := keyManager.GetPriKey()
			config.PrikeyList = append(config.PrikeyList, key)
		}
		config.WriteFile(config_name)
		fmt.Println(config)
	}
}

func CreateConfig() Config {
	var conf Config
	conf.WorkerList = []AddressPair{
		AddressPair{"127.0.0.1", 5010},
		AddressPair{"127.0.0.1", 5011},
		AddressPair{"127.0.0.1", 5012},
		AddressPair{"127.0.0.1", 5013}}
	conf.WorkerCandidateList = []AddressPair{
		AddressPair{"127.0.0.1", 5010},
		AddressPair{"127.0.0.1", 5011},
		AddressPair{"127.0.0.1", 5012},
		AddressPair{"127.0.0.1", 5013}}
	conf.VoterList = []AddressPair{
		AddressPair{"127.0.0.1", 5010},
		AddressPair{"127.0.0.1", 5011},
		AddressPair{"127.0.0.1", 5012},
		AddressPair{"127.0.0.1", 5013}}
	conf.MyAddress = AddressPair{"127.0.0.1", 5010}
	conf.ManagementServerIP = "121.15.171.86"
	conf.ManagementServerPort = 9999
	conf.SingleServerNodeNum = 4
	conf.Hostname = "本地"

	for i := 0; i < conf.SingleServerNodeNum; i++ {
		var keyManager KeyManager.KeyManager
		keyManager.Init()
		keyManager.GenRandomKeyPair()
		conf.PubkeyList = append(conf.PubkeyList, keyManager.GetPubkey())
		key, _ := keyManager.GetPriKey()
		conf.PrikeyList = append(conf.PrikeyList, key)
	}
	conf.ServicePort=8010
	conf.GenesisDutyWorker = 0
	conf.WorkerNum = 4
	conf.VotedNum = 4
	conf.BlockGroupPerCycle = 100000
	conf.Tcut = 5
	conf.GenerateBlockPeriod = 0.5
	conf.TxPoolSize = 100
	conf.ByzantineNode = false
	conf.DropDatabase = true
	return conf
}

func (conf Config) WriteFile(file string) {
	c := config.NewDefault()
	c.AddSection("network")
	//设置WorkList
	str := ""
	for i := 0; i < len(conf.WorkerList); i++ {
		str += conf.WorkerList[i].IP + ":" + strconv.Itoa(conf.WorkerList[i].Port)
		if i != len(conf.WorkerList)-1 {
			str += ","
		}
	}
	c.AddOption("network", "WorkerList", str)
	//设置WorkerCandidateList
	str = ""
	for i := 0; i < len(conf.WorkerCandidateList); i++ {
		str += conf.WorkerCandidateList[i].IP + ":" + strconv.Itoa(conf.WorkerCandidateList[i].Port)
		if i != len(conf.WorkerCandidateList)-1 {
			str += ","
		}
	}
	c.AddOption("network", "WorkerCandidateList", str)
	//设置VoterList
	str = ""
	for i := 0; i < len(conf.VoterList); i++ {
		str += conf.VoterList[i].IP + ":" + strconv.Itoa(conf.VoterList[i].Port)
		if i != len(conf.VoterList)-1 {
			str += ","
		}
	}
	c.AddOption("network", "VoterList", str)

	//设置SingleServerNodeNum
	c.AddOption("network", "SingleServerNodeNum", strconv.Itoa(conf.SingleServerNodeNum))
	//设置IP
	c.AddOption("network", "IP", conf.MyAddress.IP)
	c.AddOption("network", "Port", strconv.Itoa(conf.MyAddress.Port))
	c.AddOption("network", "ServicePort", strconv.Itoa(conf.ServicePort))
	c.AddOption("network", "ManagementServerIP", conf.ManagementServerIP)
	c.AddOption("network", "ManagementServerPort", strconv.Itoa(conf.ManagementServerPort))
	c.AddOption("network", "Hostname", conf.Hostname)

	pubkey := ""
	prikey := ""
	for i := 0; i < conf.SingleServerNodeNum; i++ {
		pubkey += conf.PubkeyList[i]
		if i != conf.SingleServerNodeNum-1 {
			pubkey += ","
		}
		prikey += conf.PrikeyList[i]
		if i != conf.SingleServerNodeNum-1 {
			prikey += ","
		}
	}

	c.AddSection("Consensus")
	c.AddOption("Consensus", "PubkeyList", pubkey)
	c.AddOption("Consensus", "PrikeyList", prikey)
	c.AddOption("Consensus", "MyPubkey", conf.MyPubkey)
	c.AddOption("Consensus", "MyPrikey", conf.MyPrikey)
	c.AddOption("Consensus", "GenesisDutyWorker", strconv.Itoa(conf.GenesisDutyWorker))
	c.AddOption("Consensus", "WorkerNum", strconv.Itoa(conf.WorkerNum))
	c.AddOption("Consensus", "VotedNum", strconv.Itoa(conf.VotedNum))
	c.AddOption("Consensus", "BlockGroupPerCycle", strconv.Itoa(conf.BlockGroupPerCycle))
	c.AddOption("Consensus", "Tcut", strconv.FormatFloat(conf.Tcut, 'f', 2, 64))
	c.AddOption("Consensus", "GenerateBlockPeriod", strconv.FormatFloat(conf.GenerateBlockPeriod, 'f', 2, 64))
	c.AddOption("Consensus", "TxPoolSize", strconv.Itoa(conf.TxPoolSize))
	c.AddOption("Consensus", "ByzantineNode", strconv.FormatBool(conf.ByzantineNode))
	c.AddOption("Consensus", "DropDatabase", strconv.FormatBool(conf.DropDatabase))

	c.WriteFile(file, 0644, "config file")
}

func (conf *Config) ReadFile(file string) {
	c, err := config.ReadDefault(file)
	if err != nil {
		fmt.Println(err)
	}
	//读取SingleServerNodeNum
	conf.SingleServerNodeNum, err = c.Int("network", "SingleServerNodeNum")
	//读取IP、端口
	conf.MyAddress.IP, err = c.String("network", "IP")
	conf.MyAddress.Port, err = c.Int("network", "Port")
	conf.ServicePort, err = c.Int("network", "ServicePort")
	conf.ManagementServerIP, err = c.String("network", "ManagementServerIP")
	conf.ManagementServerPort, err = c.Int("network", "ManagementServerPort")
	conf.Hostname, err = c.String("network", "Hostname")

	conf.MyPubkey, err = c.String("Consensus", "MyPubkey")
	conf.MyPrikey, err = c.String("Consensus", "MyPrikey")
	conf.GenesisDutyWorker, err = c.Int("Consensus", "GenesisDutyWorker")
	conf.WorkerNum, err = c.Int("Consensus", "WorkerNum")
	conf.VotedNum, err = c.Int("Consensus", "VotedNum")
	conf.BlockGroupPerCycle, err = c.Int("Consensus", "BlockGroupPerCycle")
	tcut,err:=c.String("Consensus", "Tcut")
	GenerateBlockPeriod,err :=c.String("Consensus", "GenerateBlockPeriod")
	conf.TxPoolSize,err = c.Int("Consensus", "TxPoolSize")
	conf.ByzantineNode,err = c.Bool("Consensus", "ByzantineNode")
	conf.DropDatabase, err = c.Bool("Consensus", "DropDatabase")

	conf.Tcut, err = strconv.ParseFloat(tcut,64)
	conf.GenerateBlockPeriod, err = strconv.ParseFloat(GenerateBlockPeriod,64)

	//读取WorkerList
	str, err := c.String("network", "WorkerList")
	str_nodes := strings.Split(str, ",")
	conf.WorkerList = []AddressPair{}
	for _, str_node := range str_nodes {
		var addr AddressPair
		split_str_node := strings.Split(str_node, ":")
		addr.IP = split_str_node[0]
		addr.Port, _ = strconv.Atoi(split_str_node[1])
		conf.WorkerList = append(conf.WorkerList, addr)
	}

	//读取WorkerCandidateList
	str, err = c.String("network", "WorkerCandidateList")
	str_nodes = strings.Split(str, ",")
	conf.WorkerCandidateList = []AddressPair{}
	for _, str_node := range str_nodes {
		var addr AddressPair
		split_str_node := strings.Split(str_node, ":")
		addr.IP = split_str_node[0]
		addr.Port, _ = strconv.Atoi(split_str_node[1])
		conf.WorkerCandidateList = append(conf.WorkerCandidateList, addr)
	}

	//读取VoterList
	str, err = c.String("network", "VoterList")
	str_nodes = strings.Split(str, ",")
	conf.VoterList = []AddressPair{}
	for _, str_node := range str_nodes {
		var addr AddressPair
		split_str_node := strings.Split(str_node, ":")
		addr.IP = split_str_node[0]
		addr.Port, _ = strconv.Atoi(split_str_node[1])
		conf.VoterList = append(conf.VoterList, addr)
	}
	//读取公私钥列表
	pubkeys, _ := c.String("Consensus", "PubkeyList")
	prikeys, _ := c.String("Consensus", "PrikeyList")
	for i := 0; i < conf.SingleServerNodeNum; i++ {
		conf.PubkeyList = strings.Split(pubkeys, ",")
		conf.PrikeyList = strings.Split(prikeys, ",")
	}
}

func ConfigHelperTest() {
	conf := CreateConfig()
	conf.WriteFile("conf")
	//var conf Config
	conf.ReadFile("conf")
	fmt.Println(conf)
}