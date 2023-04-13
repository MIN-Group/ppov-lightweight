package MongoDB

import (
	"hash/crc32"
	"log"
	"strconv"

	"ppov/MetaData"
	"ppov/lib/mgo/bson"
)

func (pl *Mongo) SaveRecordToDatabase(typ string, item MetaData.Record) {
	temp := MetaData.KVRecord{Type: item.Type, Key: item.Key, Value: item.Value}
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	pl.InsertToMogoRecord(temp, subname)
}

func (pl *Mongo) UpdateRecordToDatabase(typ string, item map[string]interface{}) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	pl.UpdateRecordToMongo(item, subname)
}


func (pl *Mongo) GetResultFromDatabase(typ, mt,mv, k,kv string) map[string]interface{} {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	item := make(map[string]interface{})
	c := session.DB("blockchain").C(subname)
	err := c.Find(bson.M{mt:mv, k:kv}).One(&item)
	if err != nil {
		//log.Println(err)
	}
	return item
}

func (pl *Mongo) DeleteData(typ, key, value string) {
	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
	subname := index + "-" + typ
	session := ConnecToDB()
	//session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	c := session.DB("blockchain").C(subname)
	_, err := c.RemoveAll(bson.M{key: value})
	if err != nil {
		log.Println(err)
	}
}

//func (pl *Mongo) DeleteAllDataFromDatabase(typ string) {
//	index := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(pl.Pubkey))))
//	subname := index + "-" + typ
//	session := ConnecToDB()
//	session.SetMode(mgo.Monotonic, true)
//	defer session.Close()
//
//	c := session.DB("blockchain").C(subname)
//	_, err := c.RemoveAll(bson.M{})
//	if err != nil {
//		log.Println(err)
//	}
//}
