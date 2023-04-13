package bls

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

var blsMgr0 = NewBlsManager()

type CntBlock struct {
	BlockNum     uint32   `msg:"block_num"`
	MerkleRoot   string   `msg:"merkle_root"`
	Transactions [][]byte `msg:"transactions"`
	Timestamp    float64   `msg:"timestamp"`
}

type CntVoteTicket struct {
	VoteResult  int     `msg:"VoteResult"`
	Timestamp   float64  `msg:"Timestamp"`
	VoteNum     int
}

type CntBlockGroup struct {
	Height       int          `msg:"height"`
	PreviousHash string          `msg:"preHash"`
	VoteTickets  []CntVoteTicket `msg:"-"`
	Timestamp    float64         `msg:"timestamp"`
	Sig          string       `msg:"Sig"`
	BlockSig     string
	VoteSig      string
	Blocks       []CntBlock `msg:"-"`
}

func getCntBlock(n int, tr []byte) CntBlock {
	bg := CntBlock{}
	bg.BlockNum = 1
	bg.MerkleRoot = "468835bf033a5947d69b156e8effd2b0ec6fff0ff2fe221bc26b22f8cb5bf76e"
	bg.Timestamp = float64(time.Now().UTC().UnixNano()) / 1e6
	bg.Transactions = make([][]byte, n)
	for i := 0; i < n ; i++{
		bg.Transactions[i] = tr
	}
	return bg
}

func getCntVotes(n int) CntVoteTicket {
	vt := CntVoteTicket{}
	vt.VoteResult = 1

	vt.Timestamp = float64(time.Now().UTC().UnixNano()) / 1e6
	return vt
}

func GetCntBlockGroup( n int,num int, length int) (int, []SecretKey, CntBlockGroup, int64){
	pris := make([]SecretKey, n)
	for i := 0 ; i < n ; i++{
		pris[i], _ = blsMgr0.GenerateKey()
	}

	bg := CntBlockGroup{}
	bg.Height = 1
	bg.PreviousHash = "468835bf033a5947d69b156e8effd2b0ec6fff0ff2fe221bc26b22f8cb5bf76e"
	bg.Timestamp = float64(time.Now().UTC().UnixNano()) / 1e6

	bg.VoteTickets = make([]CntVoteTicket, n)
	for i := 0 ; i < n ; i++{
		bg.VoteTickets[i] = getCntVotes(n)
	}

	str := RandStringBytes(length)
	bg.Blocks = make([]CntBlock, n)
	for i := 0 ; i < n; i++{
		bg.Blocks[i] = getCntBlock(num,[]byte(str))
	}

	var totalTime int64
	blockDsigs := make([]Signature,0, n)
	for i := 0 ; i < n ; i++{

		buf, _ := json.Marshal(bg.Blocks[i])
		hash := sha256.New()
		hash.Write(buf)
		b := hash.Sum(nil)
		t1 := time.Now()
		sig := pris[i].Sign(b)
		t2 := time.Now()
		totalTime += t2.Sub(t1).Nanoseconds()
		//fmt.Println("block", t2.Sub(t1).Nanoseconds())
		blockDsigs = append(blockDsigs, sig)
	}

	t1 := time.Now()
	sig, err := blsMgr0.Aggregate(blockDsigs)
	if err != nil{
		fmt.Println(err)
	}
	t2 := time.Now()
	totalTime += t2.Sub(t1).Nanoseconds()

	bg.BlockSig = string(sig.Compress().Bytes())

	voteDsigs := make([]Signature, 0, n)

	for i := 0 ; i < n ; i++{

		buf, _ := json.Marshal(bg.VoteTickets[i])
		hash := sha256.New()
		hash.Write(buf)
		b := hash.Sum(nil)
		t11 := time.Now()
		sig := pris[i].Sign(b)
		t22 := time.Now()
		totalTime += t22.Sub(t11).Nanoseconds()

		//fmt.Println("vote", t22.Sub(t11).Nanoseconds())
		voteDsigs = append(voteDsigs, sig)
	}
	t111 := time.Now()
	sig, _ = blsMgr0.Aggregate(voteDsigs)
	t222 := time.Now()
	totalTime += t222.Sub(t111).Nanoseconds()
	bg.VoteSig = string(sig.Compress().Bytes())

	buf, _ := json.Marshal(bg)
	hash := sha256.New()
	hash.Write(buf)
	b := hash.Sum(nil)
	t11 := time.Now()
	bg.Sig = string(pris[0].Sign(b).Compress().Bytes())
	t22 := time.Now()
	totalTime += t22.Sub(t11).Nanoseconds()
	//fmt.Println("cnt blocks", SizeOf(bg.Blocks))
	//fmt.Println("cnt vote", SizeStruct(bg.VoteTickets))
	//fmt.Println("cnt ont vote", SizeOf(bg.VoteTickets[0]))
	//fmt.Println("cnt vote sig:",SizeStruct(bg.VoteSig))
	//fmt.Println("cnt block sig", SizeStruct(bg.BlockSig))
	//fmt.Println("cnt previous hash",SizeStruct(bg.PreviousHash))
	//fmt.Println("cnt sig", SizeStruct(bg.Sig))
	//fmt.Println("cnt len :",SizeStruct(bg))
	return SizeStruct(bg), pris,bg, totalTime
}