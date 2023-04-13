package MetaData

import "fmt"

const (
	ADD = "add"
)

type KVRecord struct {
	Type   	   string     `msg:"type"`
	Key    	   string     `msg:"key"`
	Value  	   string     `msg:"value"`
}

//go:generate msgp
type Record struct {
	Type   	   string     `msg:"type"`
	Key    	   string     `msg:"key"`
	Value  	   string     `msg:"value"`
	Timestamp  string 	  `msg:"timestamp"`
	Sender     string     `msg:"sender"`
	Sig        string     `msg:"sig"`
	Command    string     `msg:"command"`
}

func (record *Record) ToByteArray() []byte {
	data, _ := record.MarshalMsg(nil)
	return data
}

func (record *Record) FromByteArray(data []byte) {
	_, err := record.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
