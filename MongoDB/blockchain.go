package MongoDB

import (
	"hash/crc32"
	"log"
	"strconv"

	"ppov/MetaData"
	"ppov/lib/mgo/bson"
)

func (pl *Mongo) QueryHeight() int {
	session := ConnecToDB()

	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var blocks []MetaData.BlockGroup
	var height int = -1
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index + "-blockgroup")
	err := c.Find(nil).Sort("-height").Limit(1).All(&blocks)
	//err = c.Find(nil).All(&blocks)
	if err != nil {
		log.Println(err)
	}
	for _, x := range blocks {
		if x.Height > height {
			height = x.Height
		}
	}
	return height
}

func (pl *Mongo) GetAmount() int {
	return pl.QueryHeight() + 1
}
func (pl *Mongo) PushbackBlockToDatabase(bg MetaData.BlockGroup) {
	pl.Block = bg
	pl.Height = bg.Height
	pl.InsertToMogo(bg, pl.Pubkey)

	//if bg.Height == 0{
	//	jsonStr, _ := json.Marshal(bg)
	//	fmt.Println(string(jsonStr))
	//}
	//for _, v := range bg.Blocks {
	//	pl.InsertToMogoBlock(v,pl.Pubkey)
	//}
	//bg.Blocks = bg.Blocks[0:0]		//清空Blocks
	//
	//pl.InsertToMogo(bg, pl.Pubkey)
}

func (pl *Mongo) GetBlockFromDatabase(height int) MetaData.BlockGroup {
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var blockgroup MetaData.BlockGroup
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index + "-blockgroup")
	err := c.Find(bson.M{"height": height}).One(&blockgroup)
	if err != nil {
		log.Println(err)
	}

	var blocks []MetaData.Block
	c1 := session.DB("blockchain").C(index + "-block")
	err = c1.Find(bson.M{"height": height}).All(&blocks)
	if err != nil {
		log.Println(err)
	}
	//{
	//	jsonStr, _ := json.Marshal(blockgroup)
	//	fmt.Println(string(jsonStr))
	//}
	true_blocks := make([]MetaData.Block, len(blockgroup.CheckHeader))
	for _, v := range blocks {
		true_blocks[v.BlockNum] = v
	}

	blockgroup.Blocks = true_blocks
	return blockgroup
}

func (pl *Mongo) GetBlockByTxHashFromDatabase(hash []byte) MetaData.BlockGroup{
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	var block MetaData.Block
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	c := session.DB("blockchain").C(index + "-block")
	err := c.Find(bson.M{"transactionshash":bson.M{"$elemMatch":bson.M{"$eq":hash}}}).One(&block)
	if err != nil {
		log.Println(err)
	}

	var blockgroup MetaData.BlockGroup
	c1 := session.DB("blockchain").C(index + "-blockgroup")
	err = c1.Find(bson.M{"height": block.Height}).One(&blockgroup)
	if err != nil {
		log.Println(err)
	}

	blockgroup.Blocks = append(blockgroup.Blocks,block)
	return blockgroup
}
