package MetaData

//用于辅助生成区块头的bytes
type BlockHeader struct {
	Height       int    `msg:"height"`
	BlockNum     uint32 `msg:"block_num"`
	Generator    string `msg:"generator"`
	MerkleRoot   []byte `msg:"merkle_root"`
	PreviousHash []byte `msg:"previous"`
}

//go:generate msgp
type Block struct {
	Height       int      `msg:"height"`
	BlockNum     uint32   `msg:"block_num"`
	Generator    string   `msg:"generator"`
	MerkleRoot   []byte   `msg:"merkle_root"`
	Timestamp	float64   `msg:"timestamp"`
	Sig          string   `msg:"sig"`
	PreviousHash []byte   `msg:"previous"`
	Transactions [][]byte `msg:"transactions"`
	TransactionsHash [][]byte `msg:"transactionshash"`
	Transactions_s []string `msg:"-"`
	IsSet        bool     `msg:"-"`
}

func (b *Block) ExtensionType() int8 {
	return 10
}

func (b *Block) Len() int {
	return b.Msgsize()
}

func (b *Block) MarshalBinaryTo(data []byte) error {
	data, err := b.MarshalMsg(nil)
	return err
}

func (b *Block) UnmarshalBinary(data []byte) error {
	data, err := b.UnmarshalMsg(data)
	return err
}

func (b *Block) GetTransactionsBytes() []byte {
	var helper DoubleByteArrayHelper
	helper.Data = b.Transactions
	data, _ := helper.MarshalMsg(nil)
	return data
}

func (b *Block) GetBlockHeaderBytes() []byte {
	var header BlockHeader
	header.Generator = b.Generator
	header.BlockNum = b.BlockNum
	header.Height = b.Height
	header.MerkleRoot = b.MerkleRoot
	header.PreviousHash = b.PreviousHash
	data, _ := header.MarshalMsg(nil)
	return data
}
