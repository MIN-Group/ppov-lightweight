package MetaData

const (
	Zero             = 0
	Genesis          = 1
	IdTransformation = 2
	ElectNewWorker   = 3
	Records		     = 4
	CreatACCOUNT	 = 5	//新建钱包账户
	TransferMONEY	 = 6	//转账
)

type TransactionInterface interface {
	ToByteArray() []byte
	FromByteArray(data []byte)
}

//go:generate msgp
type TransactionHeader struct {
	TXType int    `msg:"tx"`
	Data   []byte `msg:"data"`
}

func EncodeTransaction(header TransactionHeader, transactionInterface TransactionInterface) (data []byte) {
	data = transactionInterface.ToByteArray()
	header.Data = data
	data, _ = header.MarshalMsg(nil)
	return data
}

func DecodeTransaction(data []byte) (header TransactionHeader, transactionInterface TransactionInterface) {
	data, _ = header.UnmarshalMsg(data)
	data = header.Data
	switch header.TXType {
	case Zero:
		var zt ZeroTransaction
		zt.FromByteArray(data)
		transactionInterface = &zt
	case Genesis:
		var gt GenesisTransaction
		gt.FromByteArray(data)
		transactionInterface = &gt
	case IdTransformation:
		var idt IdentityTransformation
		idt.FromByteArray(data)
		transactionInterface = &idt
	case ElectNewWorker:
		var emwt ElectNewWorkerTeam
		emwt.FromByteArray(data)
		transactionInterface = &emwt
	case Records:
		var record Record
		record.FromByteArray(data)
		transactionInterface = &record
	case CreatACCOUNT:
		var ca CreatAccount
		ca.FromByteArray(data)
		transactionInterface = &ca
	case TransferMONEY:
		var tm TransferMoney
		tm.FromByteArray(data)
		transactionInterface = &tm
	}
	return
}
