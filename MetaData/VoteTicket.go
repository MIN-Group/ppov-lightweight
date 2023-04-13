package MetaData

//go:generate msgp

type VoteTicket struct {
	VoteResult  []uint8    `msg:"VoteResult"`
	BlockHashes [][]byte `msg:"hashes"`
	Voter       uint32   `msg:"Voter"`
	Timestamp   float64  `msg:"Timestamp"`
	Sig         string   `msg:"Sig"`
}