package Node

import (
	"encoding/json"
	"fmt"
	"ppov/MetaData"
	"ppov/utils"
)



type TransactionMsgForManagerment struct {
	Type      int
	Pubkey    string
	Height    int
	Agreement int
	Txs_num   int
}
func (node *Node) SendTransactionMsgToManagementServer(bg MetaData.BlockGroup) {
	var msg []TransactionMsgForManagerment
	var num_of_trans = 0
	for _, eachBlock := range bg.Blocks {
		num_of_trans += len(eachBlock.Transactions)
	}
	deVoteTickets := utils.DecompressToIntSlice(bg.VoteResult)
	for i, eachticket := range bg.VoteTickets {
		var one_msg TransactionMsgForManagerment
		one_msg.Type = 0
		one_msg.Pubkey = node.accountManager.VoterNumberSet[eachticket.Voter]
		one_msg.Height = bg.Height
		one_msg.Agreement = deVoteTickets[i]
		one_msg.Txs_num = num_of_trans
		msg = append(msg, one_msg)
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	go node.network.SendPacket(data, node.config.ManagementServerIP, node.config.ManagementServerPort)
}

type NormalMsgForManagerment struct {
	Type                int
	Pubkey              string
	Name                string
	IP                  string
	Is_butler_candidate bool
	Is_butler           bool
	Is_commissioner     bool
}
func (node *Node) SendNormalMsgToManagementServer() {
	if node.mongo.Height%20 == 0 || node.mongo.Height < 100 && node.mongo.Height%5 == 0 {
		var msg NormalMsgForManagerment
		msg.Type = 1
		msg.Pubkey = node.config.MyPubkey
		msg.Name = node.config.Hostname
		msg.IP = node.network.MyNodeInfo.IP
		_, msg.Is_butler = node.accountManager.WorkerSet[msg.Pubkey]
		_, msg.Is_butler_candidate = node.accountManager.WorkerCandidateSet[msg.Pubkey]
		_, msg.Is_commissioner = node.accountManager.VoterSet[msg.Pubkey]
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
		}
		go node.network.SendPacket(data, node.config.ManagementServerIP, node.config.ManagementServerPort)
	}
}