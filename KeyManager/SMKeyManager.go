package KeyManager

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"ppov/lib/ccs-gm/sm2"
	"ppov/lib/ccs-gm/sm3"
	"ppov/utils"
	"strings"
)

var AddressLength = 32

var (
	INDEXES  []int
	bigRadix = big.NewInt(58)
	bigZero  = big.NewInt(0)
)

const (
	ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type SMKeyManager struct {
	priv	 *sm2.PrivateKey
	pub      *sm2.PublicKey
}

func (km *SMKeyManager) Init() {
	km.priv=new(sm2.PrivateKey)
	km.pub=new(sm2.PublicKey)
}

func (km *SMKeyManager) GenRandomKeyPair() {
	prk, err := sm2.GenerateKey(rand.Reader) //根据随机数得到私钥
	if err != nil {
		fmt.Println(err)
	}
	km.priv = prk
	km.pub = &prk.PublicKey

	for len(km.GetPubkey()) != 130 {
		prk, err := sm2.GenerateKey(rand.Reader) //根据随机数得到私钥
		if err != nil {
			fmt.Println(err)
		}
		km.priv = prk
		km.pub = &prk.PublicKey
	}
}

func (km *SMKeyManager) GetPubkey() string {
	xStr := hex.EncodeToString(km.pub.X.Bytes())
	yStr := hex.EncodeToString(km.pub.Y.Bytes())

	return "04" + xStr + yStr
}

func (km *SMKeyManager) SignWithPriKey(text []byte, data string) (string, error){
	sig, err := km.priv.Sign(rand.Reader, text, nil)
	return base64.StdEncoding.EncodeToString(sig), err
}

func (km *SMKeyManager) GetPriKey() (string, error) {
	return hex.EncodeToString(km.priv.D.Bytes()), nil
}

func (km *SMKeyManager) SetPriKey(data string) {
	if len(data) != 64 {
		panic("设置私钥错误")
		return
	}

	sk, _ := new(big.Int).SetString(data, 16)
	pkx, pky := sm2.P256().ScalarBaseMult(sk.Bytes())

	priv := &sm2.PrivateKey{sm2.PublicKey{sm2.P256(), pkx, pky, nil}, sk, nil}

	km.priv = priv
}

func (km *SMKeyManager) SetPubkey(data string) {
	if len(data) != 130 {
		panic("设置公钥错误")
		return
	}

	x, flag := new(big.Int).SetString(data[2:66], 16)

	if flag == false {
		panic("设置公钥错误")
		return
	}

	y, flag := new(big.Int).SetString(data[66:130], 16)
	if flag != true {
		panic("设置公钥错误")
		return
	}

	pub := &sm2.PublicKey{
		Curve: sm2.P256(),
	}
	pub.X = x
	pub.Y = y
	km.pub = pub
}

func (km *SMKeyManager) Sign(text []byte) (string, error) {
	sig, err := km.priv.Sign(rand.Reader, text, nil)
	return base64.StdEncoding.EncodeToString(sig), err
}


func (km *SMKeyManager) Verify(text []byte, signature string, pubkey string) (bool, error) {
	t, err := base64.StdEncoding.DecodeString(string(signature))
	if err != nil {
		fmt.Errorf("base64编码错误")
		return false, err
	}

	if len(pubkey) != 130 {
		fmt.Errorf("设置公钥错误")
		return false,err
	}

	x, flag := new(big.Int).SetString(pubkey[2:66], 16)

	if flag == false {
		fmt.Errorf("设置公钥错误")
		return false,err
	}

	y, flag := new(big.Int).SetString(pubkey[66:130], 16)
	if flag != true {
		fmt.Errorf("设置公钥错误")
		return false,err
	}

	pub := &sm2.PublicKey{
		Curve: sm2.P256(),
	}
	pub.X = x
	pub.Y = y

	return pub.Verify(text, t), nil
}

func (km *SMKeyManager) VerifyWithSelfPubkey(text []byte, signature string) (bool, error) {
	return km.Verify(text, signature, km.GetPubkey())
}

func (km *SMKeyManager)GetAddress() string{
	pubkey := km.GetPubkey()
	b, err :=utils.HexToBytes(pubkey)
	if err != nil{
		return ""
	}

	return encodeAddress(b)
}

func (km *SMKeyManager)VerifyAddressWithPubkey(pubkey, address string)(bool, error){
	b, err :=utils.HexToBytes(pubkey)
	if err != nil{
		return false,err
	}

	return encodeAddress(b) == address, nil
}

func encodeAddress(hash []byte) string {
	tosum := make([]byte, 32)
	copy(tosum[0:15], hash[0:15])
	copy(tosum[16:],hash[len(hash)-16:])
	cksum := doubleHash(tosum)

	b := make([]byte, 25)
	copy(b[0:], hash)
	copy(b[12:], cksum[:13])

	return base58Encode(b)
}



/**
  证书分解
  通过hex解码，分割成数字证书r，s
*/
func getSign(signature string) (rint, sint big.Int, err error) {
	byterun, err := hex.DecodeString(signature)
	if err != nil {
		err = errors.New("decrypt error, " + err.Error())
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode = ", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return

}

/**
  校验文本内容是否与签名一致
  使用公钥校验签名和文本内容
*/
func verify(text []byte, signature string, key *sm2.PublicKey) (bool, error) {

	rint, sint, err := getSign(signature)
	if err != nil {
		return false, err
	}
	result := sm2.Verify(key,text,&rint,&sint)

	return result, nil

}

func GetHash(data []byte) []byte {
	h := sm3.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	return sum
}

func doubleHash(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// Base58Encode encodes a byte slice to a modified base58 string.
func base58Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		answer = append(answer, ALPHABET[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, ALPHABET[0])
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}