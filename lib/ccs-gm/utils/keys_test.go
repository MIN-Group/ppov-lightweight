package utils

//func TestPEM2Key(t *testing.T) {
//	iniSk, _ := sm2.GenerateKey(rand.Reader)
//	iniPk := iniSk.PublicKey
//
//	pemSk, err := PrivateKeyToPEM(iniSk, nil)
//	fmt.Println(string(pemSk))
//	if err != nil {
//		t.Errorf("private key to pem error %t", err)
//	}
//
//	pemPk, err := PublicKeyToPEM(&iniPk, nil)
//	if err != nil {
//		t.Errorf("public key to pem error %t", err)
//	}
//	fmt.Println(string(pemPk))
//
//	normalSk, err := PEMtoPrivateKey(pemSk, nil)
//	if err != nil {
//		t.Errorf("pem to private key error %t", err)
//	}
//
//	normalPk, err := PEMtoPublicKey(pemPk, nil)
//	if err != nil {
//		t.Errorf("pem to public key error %t", err)
//	}
//	testMsg := []byte("123456")
//	signedData, _ := normalSk.Sign(rand.Reader, testMsg, nil)
//	ok := normalPk.Verify(testMsg, signedData)
//	if !ok {
//		t.Error("key verify error")
//	}
//}
