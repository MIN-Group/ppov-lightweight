package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/Hyperledger-TWGC/ccs-gm/sm2"
	"github.com/Hyperledger-TWGC/ccs-gm/sm4"
	"github.com/Hyperledger-TWGC/ccs-gm/utils"
	"github.com/Hyperledger-TWGC/ccs-gm/x509"
	"math/big"
)

func mainn(){
	pempk := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAE876orfD5QwiatcFcGSpzYvCviW0N
eHrjtq3dwZbQbnb5QKD/pHijq8KqQEQlECKyC3jfVA8+RQsW1TR3XG/fvg==
-----END PUBLIC KEY-----`
	normalPk, err := utils.PEMtoPublicKey([]byte(pempk), nil)
	if err != nil{
		fmt.Println(err)
	}

	normalPk = &sm2.PublicKey{
		Curve: sm2.P256(),
	}

	x,_ := new(big.Int).SetString("A5893C5BA9C73888710E80F25F7680375EAAC156CA94C3095548AB9F47D705EE",16)
	normalPk.X = x
	fmt.Println(normalPk.X)
	y,_ := new(big.Int).SetString("7A6BC88C86489996DA0E55C2B7C67D0C7F309224D76F2B5ACDE30549D45237F7",16)
	normalPk.Y = y

	b, err :=hex.DecodeString("3046022100C9E8476A7720C4BD7D5436C0FA3D5A6D07F40D1B56E110BE08D4A544D3024AA8022100CAD5F6C0C5CCD19E36A405E2A8BC369BEAAD99645F25E19053B1AEBE7EDF45D0")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(b)

	result := normalPk.Verify([]byte("hellokdjasfkjdksajfkdsjkaf") ,b)
	fmt.Println(result)
}

func main()  {
	pempk := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAE876orfD5QwiatcFcGSpzYvCviW0N
eHrjtq3dwZbQbnb5QKD/pHijq8KqQEQlECKyC3jfVA8+RQsW1TR3XG/fvg==
-----END PUBLIC KEY-----`
	normalPk, err := utils.PEMtoPublicKey([]byte(pempk), nil)

	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(normalPk.X)
	pkb,err := x509.MarshalPKIXPublicKey(normalPk)
	fmt.Println("pk byte", hex.EncodeToString(pkb))
	fmt.Println("x", hex.EncodeToString(normalPk.X.Bytes()))
	fmt.Println("y", hex.EncodeToString(normalPk.Y.Bytes()))

	pemSk := `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgrx6HiMUfpBbWLb/T
QQq6BBNbvfFtmp6yFGLQjMtcQ/mgCgYIKoEcz1UBgi2hRANCAATzvqit8PlDCJq1
wVwZKnNi8K+JbQ14euO2rd3BltBudvlAoP+keKOrwqpARCUQIrILeN9UDz5FCxbV
NHdcb9++
-----END PRIVATE KEY-----`
	normalSk, err := utils.PEMtoPrivateKey([]byte(pemSk), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("D,",hex.EncodeToString(normalSk.D.Bytes()))

	pib, err := x509.MarshalECPrivateKey(normalSk)
	fmt.Println("pri byte", hex.EncodeToString(pib))

	b, err :=hex.DecodeString("3046022100cc45b32f80a64046a827dce1ddea8d593c755c99715a74f0e5a7b7189a64d4dd022100f955341562e77a68efc3de26512eeb335aced64dd6e3204af4abc68fd191ffdf")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(b)

	result := normalPk.Verify([]byte("hello") ,b)
	fmt.Println(result)

	msg := []byte("wzx")
	sig, err := normalSk.Sign(rand.Reader, msg, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sig,", hex.EncodeToString(sig))
	fmt.Println(sig)

	b1, err := hex.DecodeString("041133A83BAD727DC44F15B30B40E3A0EB84700BFEF4A59E5555ADC3DFFD09ABEFB595EF64DAD24AF06DF62746D952C0EDEA7CF333590382D29B8530F0EB4CB3399D330A75406C11AE828793518C2827B8E5F7A256D4DF0E1754EE25728520224038E0B7F788E78DE943B8AC842584E5F393")
	if err !=nil{
		fmt.Println(err)
	}

	//b1 ,err = base64.StdEncoding.DecodeString("BM38Gm0MEJj/KZx8qQ0n1sj3g7nIZ7H92lq0TOyfVfSU+jOkhbrJ5rzfvjTiTJeSFfdlxiDtwyrZJKFOxjf+5tEsMzSeyJhRfVqeabfLy12jx4LcSDGJ6qOLfdruqHJy4pHwTQXHy87WMNJbQbozVfj2mUGt9JZq")
	res, err :=sm2.Decrypt(b1, normalSk)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(res))
	msg = []byte("sm2 encryption standard")

	//test encryption
	cipher, err := sm2.Encrypt(rand.Reader, normalPk, msg)
	fmt.Println(cipher)
	fmt.Println( "cipher,", base64.StdEncoding.EncodeToString(cipher))

	key:= []byte("1234567890123456")
	if err !=nil{
		fmt.Println(err)
	}
	in, err := hex.DecodeString("16213395BC613851D96868649320FE8C")
	if err !=nil{
		fmt.Println(err)
	}
	//in = []byte("1234567890123456")

	out, err :=sm4.Sm4Ecb(key, in,sm4.DEC)
	if err !=nil{
		fmt.Println("sm4",err)
	}
	fmt.Println("sm4",string(out))
	fmt.Println(hex.EncodeToString(out))
}



