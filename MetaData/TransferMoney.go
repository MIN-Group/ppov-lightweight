package MetaData

import "fmt"

//go:generate msgp
type TransferMoney struct {
	From       string     `msg:"from" json:"from"`
	To     	   string 	  `msg:"to" json:"to"`
	Pubkey	   string 	  `msg:"pubkey" json:"pubkey"`
	Amount     int 	      `msg:"amount" json:"amount"`
	Timestamp  string     `msg:"timestamp" json:"timestamp"`
	Sig        string     `msg:"sig" json:"sig"`
}

func (tm *TransferMoney) ToByteArray() []byte {
	data, _ := tm.MarshalMsg(nil)
	return data
}

func (tm *TransferMoney) FromByteArray(data []byte) {
	_, err := tm.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
