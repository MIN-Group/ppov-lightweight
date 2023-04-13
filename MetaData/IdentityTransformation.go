package MetaData

import "fmt"

//go:generate msgp
type IdentityTransformation struct {
	Type   string    `msg:"type"`
	Pubkey string `msg:"pubkey"`
	NodeID uint64 `msg:"nodeid"`
	IPAddr string `msg:"ip"`
	Port   int    `msg:"port"`
}

func (itmsg IdentityTransformation) ToByteArray() ([]byte) {
	data, _ := itmsg.MarshalMsg(nil)
	return data
}

func (itmsg *IdentityTransformation) FromByteArray(data []byte)  {
	_, err := itmsg.UnmarshalMsg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
}
