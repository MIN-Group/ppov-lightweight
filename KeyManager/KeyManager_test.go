package KeyManager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	 km := KeyManager{}
	 km.Init()
	 km.GenRandomKeyPair()

	 message := []byte("hello world")
	 sig,err := km.Sign(message)
	 fmt.Println(len(sig))
	 assert.NoError(t, err)

	 res, err := km.VerifyWithSelfPubkey(message, sig)
	 assert.NoError(t, err)
	 fmt.Println(res)
	 assert.True(t, true, res)
}

func TestDecompress(t *testing.T){
	km := KeyManager{}
	km.Init()
	km.GenRandomKeyPair()
	pubStr := km.GetPubkey()
	priStr,_ := km.GetPriKey()
	message := []byte("hello world")
	sign, _ := km.Sign(message)

	km2 := KeyManager{}
	km2.SetPubkey(pubStr)
	km2.SetPriKey(priStr)
	verifyResult, _ := km2.VerifyWithSelfPubkey(message, sign)
	fmt.Println(verifyResult)
	assert.Equal(t, km.GetPubkey(), km2.GetPubkey())
	assert.True(t, verifyResult)
}

func TestAggregate(t *testing.T) {
	kms := make([]KeyManager, 3)
	pubStrs := make([]string,0, len(kms))
	for i := 0; i < len(kms); i++ {
		kms[i] = KeyManager{}
		kms[i].Init()
		kms[i].GenRandomKeyPair()
		pubStrs = append(pubStrs, kms[i].GetPubkey())
	}

	source := []byte("hello world")
	msgs := make([][]byte, 0, len(kms))
	signs := make([]string, 0, len(kms))

	for i := 0; i < len(kms); i++ {
		msgs = append(msgs, source)
		sign, _ := kms[i].Sign(source)
		signs = append(signs, sign)

		source = append(source, byte(1))
	}

	aggSign, _ := kms[0].AggregateSign(signs)

	verRes, err :=kms[0].VerifyAggSign(aggSign, pubStrs, msgs)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("verify result: ",verRes)

	km := KeyManager{}
	km.Init()
	km.GenRandomKeyPair()

	signs = make([]string, 0, 2)
	sign, _ := km.Sign(source)
	signs = append(signs, sign)
	signs = append(signs, aggSign)
	aggSign2,err := km.AggregateSign(signs)

	pubStrs = append(pubStrs, km.GetPubkey())
	msgs = append(msgs, source)
	verRes, err =kms[0].VerifyAggSign(aggSign2, pubStrs, msgs)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(verRes)
}