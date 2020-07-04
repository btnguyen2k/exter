package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/btnguyen2k/consu/checksum"
)

func main() {
	privKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	msg := []byte("message to sign")
	// msgHash := sha256.New()
	// _, err := msgHash.Write(msg)
	// if err != nil {
	// 	panic(err)
	// }
	msgHashSum := checksum.Sha256HashFunc(msg)
	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Signature: ", signature)

	privKey2, _ := rsa.GenerateKey(rand.Reader, 1024)

	err = rsa.VerifyPSS(&privKey2.PublicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
	} else {
		fmt.Println("signature verified")
	}

	err = rsa.VerifyPSS(&privKey.PublicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
	} else {
		fmt.Println("signature verified")
	}
}
