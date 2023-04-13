package test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	rand0 "math/rand"
	"ppov/KeyManager"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand0.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestGenerateBLS(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100; i++ {
		km := KeyManager.KeyManager{}
		km.Init()
		km.GenRandomKeyPair()
	}
	end := time.Now()
	fmt.Println(end.Sub(start).Milliseconds())
	fmt.Println(end.Sub(start).Microseconds())
	fmt.Println(end.Sub(start).Nanoseconds())
}

func TestGenerateECDSA(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100; i++ {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalln(err)
		}

		privateKey.X.Bytes()
	}
	end := time.Now()
	fmt.Println(end.Sub(start).Milliseconds())
	fmt.Println(end.Sub(start).Microseconds())
	fmt.Println(end.Sub(start).Nanoseconds())
}

func TestBLSSignAndVerify(t *testing.T){
	base := 0
	randStr := ""
	km := KeyManager.KeyManager{}
	km.Init()
	km.GenRandomKeyPair()
	for i := 0; i < 21; i++ {
		randStr += RandStringBytes(base)
		start := time.Now()
		sign, err := km.Sign([]byte(randStr))
		if err != nil{
			panic(err)
		}
		end := time.Now()
		fmt.Printf("%.1f %f", 1.0*float64(base)/1024/1024, float64(1.0*float64(base))/1024/1024/end.Sub(start).Seconds())

		start = time.Now()
		_,err = km.Verify([]byte(randStr), sign, km.GetPubkey())
		if err != nil{
			panic(err)
		}
		end = time.Now()
		fmt.Printf(" %f\n", float64(1.0*float64(base))/1024/1024/end.Sub(start).Seconds())
		base += 1024*512
	}

}

func TestECDSASign(t *testing.T){
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	randStr := ""
	base := 0
	for i := 0; i < 21; i++ {
		randStr += RandStringBytes(base)
		start := time.Now()
		hashFunc := sha256.New()
		hashFunc.Write([]byte(randStr))
		r, s, err := ecdsa.Sign(rand.Reader, privateKey,hashFunc.Sum(nil))
		if err != nil{
			panic(err)
		}
		end := time.Now()
		fmt.Printf("%.1f %f", float64(base)/1024/1024,1.0*float64(base)/1024/1024/end.Sub(start).Seconds())

		start = time.Now()
		hashFunc2 := sha256.New()
		hashFunc2.Write([]byte(randStr))
		flag := ecdsa.Verify(&privateKey.PublicKey, hashFunc2.Sum(nil), r, s)
		if !flag{
			panic(err)
		}
		end = time.Now()
		fmt.Printf(" %f\n", float64(base)/1024/1024/end.Sub(start).Seconds())
		base += 1024*512
	}
}

func TestAggression(t *testing.T){
	for j := 1; j <= 100; j++ {
		kms := make([]KeyManager.KeyManager, j)
		pubStrs := make([]string,0, len(kms))
		for i := 0; i < len(kms); i++ {
			kms[i] = KeyManager.KeyManager{}
			kms[i].Init()
			kms[i].GenRandomKeyPair()
			pubStrs = append(pubStrs, kms[i].GetPubkey())
		}

		source :=[]byte(RandStringBytes(1024*1024))
		msgs := make([][]byte, 0, len(kms))
		signs := make([]string, 0, len(kms))

		for i := 0; i < len(kms); i++ {
			msgs = append(msgs, source)
			sign, _ := kms[i].Sign(source)
			signs = append(signs, sign)

			source = append(source, byte(1))
		}

		start := time.Now()
		aggSign, _ := kms[0].AggregateSign(signs)
		end := time.Now()
		fmt.Printf("%d %d", j, end.Sub(start).Milliseconds())


		start = time.Now()
		_, err :=kms[0].VerifyAggSign(aggSign, pubStrs, msgs)
		if err != nil{
			panic(err)
		}
		end = time.Now()
		fmt.Printf(" %d\n", end.Sub(start).Milliseconds())
	}
}
