package Node

import (
	"log"
	"ppov/MetaData"
	"ppov/utils"
)

type Ticket struct {
	counter map[string]int
}

func (node *Node) VotingStatistics(item MetaData.BlockGroup) (MetaData.BlockGroup, bool) {
	var THRESHOLD = node.config.VotedNum / 3 * 2
	var THRESHOLD2 = node.config.VotedNum / 2

	counter_yes := make([]Ticket, node.config.WorkerNum)
	counter_no := make([]Ticket, node.config.WorkerNum)
	result := make([]Ticket, node.config.WorkerNum)

	CounterInit(counter_yes)
	CounterInit(counter_no)
	CounterInit(result)

	for _, x := range item.VoteTickets {
		rs := utils.DecompressToIntSlice(x.VoteResult)
		if x.BlockHashes == nil || rs == nil {
			continue
		}
		if len(x.BlockHashes) != len(rs) {
			log.Println("VotingStatistics : len(x.BlockHashes) != len(x.VoteResult)")
			return item, false
		}
		for i := 0; i < len(x.BlockHashes); i++ {
			block_hash_str := utils.BytesToHex(x.BlockHashes[i])
			_, is_exist := counter_yes[i].counter[block_hash_str]
			if !is_exist {
				counter_yes[i].counter[block_hash_str] = 0
			}

			_, is_exist = counter_no[i].counter[block_hash_str]
			if !is_exist {
				counter_no[i].counter[block_hash_str] = 0
			}
			if rs[i] == 1 {
				counter_yes[i].counter[block_hash_str] += 1
			} else if rs[i] == -1 {
				counter_no[i].counter[block_hash_str] += 1
			}
		}
	}

	check := 0
	if node.round < 2 {
		for i, x := range counter_yes {
			for k, v := range x.counter {
				if v > THRESHOLD {
					result[i].counter[k] = 1
					check += 1
					break
				}
				if counter_no[i].counter[k] > THRESHOLD {
					result[i].counter[k] = -1
					check += 1
					break
				}
			}
		}
		if check != node.config.WorkerNum {
			return item, false
		}
		//fmt.Println("vote statics success")
	} else {
		for i, x := range counter_yes {
			for k, v := range x.counter {
				if v > THRESHOLD2 {
					result[i].counter[k] = 1
					check += 1
					break
				}
				if counter_no[i].counter[k] > THRESHOLD2 {
					result[i].counter[k] = -1
					check += 1
					break
				}
			}
		}
		if check == 0 {
			return item, false
		}
		//fmt.Println("check=", check)
	}

	item.BlockHashes=make([][]byte,node.config.WorkerNum)
	voteResult := make([]int,node.config.WorkerNum)
	for i, x := range result {
		for k, v := range x.counter {
			item.BlockHashes[i], _ = utils.HexToBytes(k)
			voteResult[i]= v
/*			item.BlockHashes = append(item.BlockHashes, k)
			item.VoteResult = append(item.VoteResult, v)*/
		}
	}
	item.VoteResult = utils.CompressIntSlice(voteResult)
	return item, true
}

func CounterInit(counter []Ticket) {
	for i := 0; i < len(counter); i++ {
		counter[i].counter = make(map[string]int)
	}
}