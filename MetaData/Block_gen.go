package MetaData

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Block) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "height":
			z.Height, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Height")
				return
			}
		case "block_num":
			z.BlockNum, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "BlockNum")
				return
			}
		case "generator":
			z.Generator, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Generator")
				return
			}
		case "merkle_root":
			z.MerkleRoot, err = dc.ReadBytes(z.MerkleRoot)
			if err != nil {
				err = msgp.WrapError(err, "MerkleRoot")
				return
			}
		case "timestamp":
			z.Timestamp, err = dc.ReadFloat64()
			if err != nil {
				err = msgp.WrapError(err, "Timestamp")
				return
			}
		case "sig":
			z.Sig, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Sig")
				return
			}
		case "previous":
			z.PreviousHash, err = dc.ReadBytes(z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		case "transactions":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Transactions")
				return
			}
			if cap(z.Transactions) >= int(zb0002) {
				z.Transactions = (z.Transactions)[:zb0002]
			} else {
				z.Transactions = make([][]byte, zb0002)
			}
			for za0001 := range z.Transactions {
				z.Transactions[za0001], err = dc.ReadBytes(z.Transactions[za0001])
				if err != nil {
					err = msgp.WrapError(err, "Transactions", za0001)
					return
				}
			}
		case "transactionshash":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "TransactionsHash")
				return
			}
			if cap(z.TransactionsHash) >= int(zb0003) {
				z.TransactionsHash = (z.TransactionsHash)[:zb0003]
			} else {
				z.TransactionsHash = make([][]byte, zb0003)
			}
			for za0002 := range z.TransactionsHash {
				z.TransactionsHash[za0002], err = dc.ReadBytes(z.TransactionsHash[za0002])
				if err != nil {
					err = msgp.WrapError(err, "TransactionsHash", za0002)
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Block) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "height"
	err = en.Append(0x89, 0xa6, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Height)
	if err != nil {
		err = msgp.WrapError(err, "Height")
		return
	}
	// write "block_num"
	err = en.Append(0xa9, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.BlockNum)
	if err != nil {
		err = msgp.WrapError(err, "BlockNum")
		return
	}
	// write "generator"
	err = en.Append(0xa9, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Generator)
	if err != nil {
		err = msgp.WrapError(err, "Generator")
		return
	}
	// write "merkle_root"
	err = en.Append(0xab, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.MerkleRoot)
	if err != nil {
		err = msgp.WrapError(err, "MerkleRoot")
		return
	}
	// write "timestamp"
	err = en.Append(0xa9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.Timestamp)
	if err != nil {
		err = msgp.WrapError(err, "Timestamp")
		return
	}
	// write "sig"
	err = en.Append(0xa3, 0x73, 0x69, 0x67)
	if err != nil {
		return
	}
	err = en.WriteString(z.Sig)
	if err != nil {
		err = msgp.WrapError(err, "Sig")
		return
	}
	// write "previous"
	err = en.Append(0xa8, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.PreviousHash)
	if err != nil {
		err = msgp.WrapError(err, "PreviousHash")
		return
	}
	// write "transactions"
	err = en.Append(0xac, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Transactions)))
	if err != nil {
		err = msgp.WrapError(err, "Transactions")
		return
	}
	for za0001 := range z.Transactions {
		err = en.WriteBytes(z.Transactions[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Transactions", za0001)
			return
		}
	}
	// write "transactionshash"
	err = en.Append(0xb0, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.TransactionsHash)))
	if err != nil {
		err = msgp.WrapError(err, "TransactionsHash")
		return
	}
	for za0002 := range z.TransactionsHash {
		err = en.WriteBytes(z.TransactionsHash[za0002])
		if err != nil {
			err = msgp.WrapError(err, "TransactionsHash", za0002)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Block) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "height"
	o = append(o, 0x89, 0xa6, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74)
	o = msgp.AppendInt(o, z.Height)
	// string "block_num"
	o = append(o, 0xa9, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d)
	o = msgp.AppendUint32(o, z.BlockNum)
	// string "generator"
	o = append(o, 0xa9, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72)
	o = msgp.AppendString(o, z.Generator)
	// string "merkle_root"
	o = append(o, 0xab, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74)
	o = msgp.AppendBytes(o, z.MerkleRoot)
	// string "timestamp"
	o = append(o, 0xa9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70)
	o = msgp.AppendFloat64(o, z.Timestamp)
	// string "sig"
	o = append(o, 0xa3, 0x73, 0x69, 0x67)
	o = msgp.AppendString(o, z.Sig)
	// string "previous"
	o = append(o, 0xa8, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73)
	o = msgp.AppendBytes(o, z.PreviousHash)
	// string "transactions"
	o = append(o, 0xac, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Transactions)))
	for za0001 := range z.Transactions {
		o = msgp.AppendBytes(o, z.Transactions[za0001])
	}
	// string "transactionshash"
	o = append(o, 0xb0, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x61, 0x73, 0x68)
	o = msgp.AppendArrayHeader(o, uint32(len(z.TransactionsHash)))
	for za0002 := range z.TransactionsHash {
		o = msgp.AppendBytes(o, z.TransactionsHash[za0002])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Block) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "height":
			z.Height, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Height")
				return
			}
		case "block_num":
			z.BlockNum, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlockNum")
				return
			}
		case "generator":
			z.Generator, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Generator")
				return
			}
		case "merkle_root":
			z.MerkleRoot, bts, err = msgp.ReadBytesBytes(bts, z.MerkleRoot)
			if err != nil {
				err = msgp.WrapError(err, "MerkleRoot")
				return
			}
		case "timestamp":
			z.Timestamp, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timestamp")
				return
			}
		case "sig":
			z.Sig, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Sig")
				return
			}
		case "previous":
			z.PreviousHash, bts, err = msgp.ReadBytesBytes(bts, z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		case "transactions":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Transactions")
				return
			}
			if cap(z.Transactions) >= int(zb0002) {
				z.Transactions = (z.Transactions)[:zb0002]
			} else {
				z.Transactions = make([][]byte, zb0002)
			}
			for za0001 := range z.Transactions {
				z.Transactions[za0001], bts, err = msgp.ReadBytesBytes(bts, z.Transactions[za0001])
				if err != nil {
					err = msgp.WrapError(err, "Transactions", za0001)
					return
				}
			}
		case "transactionshash":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "TransactionsHash")
				return
			}
			if cap(z.TransactionsHash) >= int(zb0003) {
				z.TransactionsHash = (z.TransactionsHash)[:zb0003]
			} else {
				z.TransactionsHash = make([][]byte, zb0003)
			}
			for za0002 := range z.TransactionsHash {
				z.TransactionsHash[za0002], bts, err = msgp.ReadBytesBytes(bts, z.TransactionsHash[za0002])
				if err != nil {
					err = msgp.WrapError(err, "TransactionsHash", za0002)
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Block) Msgsize() (s int) {
	s = 1 + 7 + msgp.IntSize + 10 + msgp.Uint32Size + 10 + msgp.StringPrefixSize + len(z.Generator) + 12 + msgp.BytesPrefixSize + len(z.MerkleRoot) + 10 + msgp.Float64Size + 4 + msgp.StringPrefixSize + len(z.Sig) + 9 + msgp.BytesPrefixSize + len(z.PreviousHash) + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.Transactions {
		s += msgp.BytesPrefixSize + len(z.Transactions[za0001])
	}
	s += 17 + msgp.ArrayHeaderSize
	for za0002 := range z.TransactionsHash {
		s += msgp.BytesPrefixSize + len(z.TransactionsHash[za0002])
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BlockHeader) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "height":
			z.Height, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Height")
				return
			}
		case "block_num":
			z.BlockNum, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "BlockNum")
				return
			}
		case "generator":
			z.Generator, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Generator")
				return
			}
		case "merkle_root":
			z.MerkleRoot, err = dc.ReadBytes(z.MerkleRoot)
			if err != nil {
				err = msgp.WrapError(err, "MerkleRoot")
				return
			}
		case "previous":
			z.PreviousHash, err = dc.ReadBytes(z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *BlockHeader) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "height"
	err = en.Append(0x85, 0xa6, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Height)
	if err != nil {
		err = msgp.WrapError(err, "Height")
		return
	}
	// write "block_num"
	err = en.Append(0xa9, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.BlockNum)
	if err != nil {
		err = msgp.WrapError(err, "BlockNum")
		return
	}
	// write "generator"
	err = en.Append(0xa9, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Generator)
	if err != nil {
		err = msgp.WrapError(err, "Generator")
		return
	}
	// write "merkle_root"
	err = en.Append(0xab, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.MerkleRoot)
	if err != nil {
		err = msgp.WrapError(err, "MerkleRoot")
		return
	}
	// write "previous"
	err = en.Append(0xa8, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.PreviousHash)
	if err != nil {
		err = msgp.WrapError(err, "PreviousHash")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BlockHeader) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "height"
	o = append(o, 0x85, 0xa6, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74)
	o = msgp.AppendInt(o, z.Height)
	// string "block_num"
	o = append(o, 0xa9, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d)
	o = msgp.AppendUint32(o, z.BlockNum)
	// string "generator"
	o = append(o, 0xa9, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72)
	o = msgp.AppendString(o, z.Generator)
	// string "merkle_root"
	o = append(o, 0xab, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x6f, 0x74)
	o = msgp.AppendBytes(o, z.MerkleRoot)
	// string "previous"
	o = append(o, 0xa8, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73)
	o = msgp.AppendBytes(o, z.PreviousHash)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockHeader) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "height":
			z.Height, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Height")
				return
			}
		case "block_num":
			z.BlockNum, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BlockNum")
				return
			}
		case "generator":
			z.Generator, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Generator")
				return
			}
		case "merkle_root":
			z.MerkleRoot, bts, err = msgp.ReadBytesBytes(bts, z.MerkleRoot)
			if err != nil {
				err = msgp.WrapError(err, "MerkleRoot")
				return
			}
		case "previous":
			z.PreviousHash, bts, err = msgp.ReadBytesBytes(bts, z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *BlockHeader) Msgsize() (s int) {
	s = 1 + 7 + msgp.IntSize + 10 + msgp.Uint32Size + 10 + msgp.StringPrefixSize + len(z.Generator) + 12 + msgp.BytesPrefixSize + len(z.MerkleRoot) + 9 + msgp.BytesPrefixSize + len(z.PreviousHash)
	return
}
