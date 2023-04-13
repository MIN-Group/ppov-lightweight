package KeyManager

import (
	"crypto/sha256"
	"errors"
	"ppov/lib/bls/bls"
	"ppov/utils"
)

type KeyManager struct {
	priv *bls.SecretKey
	pub  *bls.PublicKey
}

var blsMgr = bls.NewBlsManager()

func (km *KeyManager) Init() {
	km.priv = new(bls.SecretKey)
	km.pub = new(bls.PublicKey)
}

func (km *KeyManager) GenRandomKeyPair() {
	sk, pk := blsMgr.GenerateKey()
	if sk == nil || pk == nil{
		panic("Generate key failed")
	}

	km.pub = &pk
	km.priv = &sk
}

func (km *KeyManager) GetPubkey() string {
	if km.pub == nil{
		return ""
	}

	pubStr := (*km.pub).Compress()
	return pubStr.String()
}

func (km *KeyManager) GetPriKey() (string, error) {
	if km.priv == nil{
		return "", nil
	}

	priStr := (*km.priv).Compress()
	return priStr.String(), nil
}

func (km *KeyManager)GetAddress() string{
	pubkey := km.GetPubkey()
	b, err :=utils.HexToBytes(pubkey)
	if err != nil{
		return ""
	}

	return encodeAddress(b)
}

func (km *KeyManager)VerifyAddressWithPubkey(pubkey, address string)(bool, error){
	b, err := utils.HexToBytes(pubkey)
	if err != nil{
		return false,err
	}

	return encodeAddress(b) == address, nil
}

func (km *KeyManager) SetPriKey(data string) {
	priv, err := blsMgr.DecSecretKeyHex(data)
	if err != nil{
		panic("设置私钥错误")
	}

	km.priv = &priv
}

func (km *KeyManager) SetPubkey(data string) {
	pub, err := blsMgr.DecPublicKeyHex(data)
	if err != nil{
		panic("设置公钥错误")
	}

	km.pub = &pub
}

func (km *KeyManager) Sign(text []byte) (string, error) {
	if km.priv == nil{
		return "", errors.New("the private key hasn't been  initialized")
	}
	hashFunc := sha256.New()
	hashFunc.Write(text)
	sig := (*km.priv).Sign(hashFunc.Sum(nil))
	return string(sig.Compress().Bytes()), nil
}

func (km *KeyManager) VerifyWithSelfPubkey(text []byte, signature string) (bool, error) {
	sig, err := blsMgr.DecSignature([]byte(signature))
	if err != nil{
		return false, err
	}
	hashFunc := sha256.New()
	hashFunc.Write(text)
	err = (*km.pub).Verify(hashFunc.Sum(nil), sig)
	if err != nil{
		return false, nil
	}else{
		return true, nil
	}
}

func (km *KeyManager) SignWithPriKey(text []byte, data string) (string, error){
	return km.Sign(text)
}

func (km *KeyManager) Verify(text []byte, signature string, pubkey string) (bool, error) {
	pub, err := blsMgr.DecPublicKeyHex(pubkey)
	if err != nil{
		return false, errors.New("set pubkey failed: "+ err.Error())
	}

	sig, err := blsMgr.DecSignature([]byte(signature))
	if err != nil{
		return false, err
	}

	hashFunc := sha256.New()
	hashFunc.Write(text)
	err = pub.Verify(hashFunc.Sum(nil), sig)
	if err != nil{
		return false, nil
	}else{
		return true, nil
	}
}

func (km *KeyManager) AggregateSign(signStrs []string) (string, error){
	signs := make([]bls.Signature, 0, len(signStrs))
	for _, signStr := range signStrs{
		sig, err := blsMgr.DecSignature([]byte(signStr))
		if err != nil{
			return "", err
		}
		signs = append(signs, sig)
	}

	aggSign, err := blsMgr.Aggregate(signs)
	if err != nil{
		return "", err
	}
	return string(aggSign.Compress().Bytes()), nil
}

func (km *KeyManager) VerifyAggSign(aggSignStr string, pubkeyStr []string, msgs [][]byte) (bool, error){
	aggSign, err := blsMgr.DecSignature([]byte(aggSignStr))
	if err != nil{
		return false, err
	}

	pubkeys := make([]bls.PublicKey, 0, len(pubkeyStr))
	for _, pubStr := range pubkeyStr{
		pubkey, err := blsMgr.DecPublicKeyHex(pubStr)
		if err != nil{
			return false, err
		}
		pubkeys = append(pubkeys, pubkey)
	}

	messages := make([]bls.Message, 0, len(msgs))
	for _, m := range msgs{
		hashFunc := sha256.New()
		hashFunc.Write(m)
		messages = append(messages, bls.Message(hashFunc.Sum(nil)))
	}
	err = blsMgr.VerifyAggregatedN(pubkeys, messages, aggSign)

	if err != nil{
		return false, err
	}else {
		return true, nil
	}
}