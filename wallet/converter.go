package wallet

import (
	"encoding/hex"

	"github.com/smallnest/blockchain/wallet/base58check"
)

// GetPublicKey 根据私钥得到公钥和p2pkh地址.
func GetPublicKey(privateKey string) (publicKey, p2pkh string) {
	priKey, _ := hex.DecodeString(privateKey)
	pubKey, ripeHashedBytes := generatePublicKey(priKey)
	publicKey = hex.EncodeToString(pubKey)
	publicKeyP2PKH := base58check.Encode(publicKeyPrefix, ripeHashedBytes)

	return publicKey, publicKeyP2PKH
}

// PrivateKey2Wif 根据私钥得到wif.
func PrivateKey2Wif(privateKey string) (wif string) {
	priKey, _ := hex.DecodeString(privateKey)
	privateKeyWif := base58check.Encode(privateKeyPrefix, priKey)
	return privateKeyWif
}

// PublicKey2P2PKH 根据公钥生成p2pkh地址.
func PublicKey2P2PKH(publicKey string) (p2pkh string) {
	pubKey, _ := hex.DecodeString(publicKey)
	return base58check.Encode(publicKeyPrefix, pubKey)
}
