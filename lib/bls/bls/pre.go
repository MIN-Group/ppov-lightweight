package bls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"time"
)

type PreBlock struct {
	Height       int      `msg:"height"`
	BlockNum     uint32   `msg:"block_num"`
	MerkleRoot   string   `msg:"merkle_root"`
	Transactions [][]byte `msg:"transactions"`
	Sig 		 string   `msg:"sig"`
	Timestamp    float64   `msg:"timestamp"`
}

type PreVoteTicket struct {
	BlockNum    int
	VoteResult  int    `msg:"VoteResult"`
	Timestamp   float64  `msg:"Timestamp"`
	Sig         string   `msg:"Sig"`
}

type PreBlockGroup struct {
	Height       int          `msg:"height"`
	PreviousHash string          `msg:"preHash"`
	VoteTickets  []PreVoteTicket `msg:"-"`
	Timestamp    float64         `msg:"timestamp"`
	Sig          string       `msg:"Sig"`

	Blocks       []PreBlock `msg:"-"`
}


func getPreBlock(pri *ecdsa.PrivateKey,n int, tr []byte) (PreBlock,int64) {
	bl := PreBlock{}
	bl.Height = 1
	bl.BlockNum = 1
	bl.Transactions = nil
	bl.Timestamp = float64(time.Now().UTC().UnixNano()) / 1e6
	bl.MerkleRoot = "a55b531527c35f26def4e8577ca9c8930436e197555499fa2120edba38de7260"


	buf,_ := json.Marshal(bl)
	hash := sha256.New()
	hash.Write(buf)
	b := hash.Sum(nil)
	t1 := time.Now()
	r, s, _ :=  ecdsa.Sign(rand.Reader,pri, b)
	t2 := time.Now()

	tmp := append(r.Bytes(), s.Bytes()[:]...)
	bl.Sig = string(tmp)

	bl.Transactions = make([][]byte, n)
	for i := 0; i < n ; i++{
		bl.Transactions[i] = tr
	}
	return bl, t2.Sub(t1).Nanoseconds()
}

func getPreVote (pri *ecdsa.PrivateKey, n int) (PreVoteTicket, int64) {
	vt := PreVoteTicket{}
	vt.BlockNum = 1
	vt.VoteResult = 1


	buf,_ := json.Marshal(vt)
	hash := sha256.New()
	hash.Write(buf)
	b := hash.Sum(nil)
	t1 := time.Now()
	r, s, _ :=  ecdsa.Sign(rand.Reader,pri, b)
	t2 := time.Now()


	tmp := append(r.Bytes(), s.Bytes()[:]...)
	vt.Sig = string(tmp)

	return vt, t2.Sub(t1).Nanoseconds()
}

func GetPreBlockGroup(n int, num int, length int) (int,PreBlockGroup, int64) {
	pri,_:= ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
	bg := PreBlockGroup{}
	bg.Height = 1
	bg.PreviousHash = "468835bf033a5947d69b156e8effd2b0ec6fff0ff2fe221bc26b22f8cb5bf76e"

	bg.VoteTickets = make([]PreVoteTicket, n)
	total := 0
	var totalTime int64 = 0
	for i := 0 ; i < n ; i++{
		tmp, tmp2 := getPreVote(pri,n)
		total += SizeStruct(bg.VoteTickets[i].Sig)
		bg.VoteTickets[i] = tmp
		//.Println("vote", tmp2)
		totalTime += tmp2
	}
	bg.Timestamp =  float64(time.Now().UTC().UnixNano()) / 1e6

	bg.Blocks = make([]PreBlock, n)

	str := RandStringBytes(length)
	for i :=0 ; i < n ; i++{
		var tmp int64
		bg.Blocks[i], tmp = getPreBlock(pri, num, []byte(str))
		totalTime += tmp
		//fmt.Println("block", tmp)
	}


	buf,_ := json.Marshal(bg)
	hash := sha256.New()
	hash.Write(buf)
	b := hash.Sum(nil)
	t1 := time.Now()
	r, s, _ :=  ecdsa.Sign(rand.Reader,pri, b)
	t2 := time.Now()

	totalTime += t2.Sub(t1).Nanoseconds()
	tmp := append(r.Bytes(), s.Bytes()[:]...)
	bg.Sig = string(tmp)

	return SizeStruct(bg), bg, totalTime
}