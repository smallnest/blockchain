package main

import (
	"log"

	"github.com/smallnest/blockchain/wallet"
)

func main() {
	priKey, wif, pubKey, address := wallet.GenerateKeys()
	log.Println("===============生成公私钥===============")
	log.Printf("private key: %s\n", priKey)
	log.Printf("wif        : %s\n", wif)
	log.Printf("public  key: %s\n", pubKey)
	log.Printf("address    : %s\n\n", address)

	mnemonic, priKey, wif, pubKey, address := wallet.GenerateBIP39("this is a test")
	log.Println("===============根据BIP-39生成公私钥===============")
	log.Printf("mnemonic   : %s\n", mnemonic)
	log.Printf("private key: %s\n", priKey)
	log.Printf("wif        : %s\n", wif)
	log.Printf("public  key: %s\n", pubKey)
	log.Printf("address    : %s\n\n", address)

	priKey, wif, pubKey, address = wallet.RecoverBIP39(mnemonic, "this is a test")
	log.Println("===============根据BIP-39恢复公私钥===============")
	log.Printf("private key: %s\n", priKey)
	log.Printf("wif        : %s\n", wif)
	log.Printf("public  key: %s\n", pubKey)
	log.Printf("address    : %s\n", address)
}
