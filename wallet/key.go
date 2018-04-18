package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/smallnest/blockchain/wallet/base58check"
	secp256k1 "github.com/toxeus/go-secp256k1"
	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
)

var (
	privateKeyPrefix = "80"
	publicKeyPrefix  = "00"
)

// GenerateKeys 产生私钥、WIF地址、公钥、P2PKH地址
func GenerateKeys() (privateKey, wif string, publicKey, p2pkh string) {
	priKey := generatePrivateKey()
	privateKey = hex.EncodeToString(priKey)
	privateKeyWif := base58check.Encode(privateKeyPrefix, priKey)

	pubKey := generatePublicKey(priKey)
	publicKey = hex.EncodeToString(pubKey)
	publicKeyP2PKH := base58check.Encode(publicKeyPrefix, pubKey)

	return privateKey, privateKeyWif, publicKey, publicKeyP2PKH
}

// GenerateBIP39 根据BIP-39规范生成助记词、私钥、WIF地址、公钥、P2PKH地址
func GenerateBIP39(secretPassphrase string) (mnemonic string, privateKey, wif string, publicKey, p2pkh string) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ = bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, secretPassphrase)

	masterKey, _ := bip32.NewMasterKey(seed)

	priKey := masterKey.Key
	privateKey = hex.EncodeToString(priKey)
	privateKeyWif := base58check.Encode(privateKeyPrefix, priKey)

	pubKey := generatePublicKey(priKey)
	publicKey = hex.EncodeToString(pubKey)
	publicKeyP2PKH := base58check.Encode(publicKeyPrefix, pubKey)

	return mnemonic, privateKey, privateKeyWif, publicKey, publicKeyP2PKH
}

// RecoverBIP39 根据助记词和密码恢复私钥、WIF地址、公钥、P2PKH地址
func RecoverBIP39(mnemonic, secretPassphrase string) (privateKey, wif string, publicKey, p2pkh string) {
	seed := bip39.NewSeed(mnemonic, secretPassphrase)
	masterKey, _ := bip32.NewMasterKey(seed)

	priKey := masterKey.Key
	privateKey = hex.EncodeToString(priKey)
	privateKeyWif := base58check.Encode(privateKeyPrefix, priKey)

	pubKey := generatePublicKey(priKey)
	publicKey = hex.EncodeToString(pubKey)
	publicKeyP2PKH := base58check.Encode(publicKeyPrefix, pubKey)

	return privateKey, privateKeyWif, publicKey, publicKeyP2PKH
}

func generatePublicKey(privateKeyBytes []byte) []byte {
	var privateKeyBytes32 [32]byte
	copy(privateKeyBytes32[:], privateKeyBytes)

	secp256k1.Start()
	publicKeyBytes, success := secp256k1.Pubkey_create(privateKeyBytes32, false)
	if !success {
		log.Fatal("Failed to create public key.")
	}

	secp256k1.Stop()

	//Next we get a sha256 hash of the public key generated
	//via ECDSA, and then get a ripemd160 hash of the sha256 hash.
	shaHash := sha256.New()
	shaHash.Write(publicKeyBytes)
	shadPublicKeyBytes := shaHash.Sum(nil)

	ripeHash := ripemd160.New()
	ripeHash.Write(shadPublicKeyBytes)
	ripeHashedBytes := ripeHash.Sum(nil)

	return ripeHashedBytes
}

func generatePrivateKey() []byte {
	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		//This is not "cryptographically random"
		bytes[i] = byte(randInt(0, math.MaxUint8))
	}
	return bytes
}

func randInt(min int, max int) uint8 {
	rand.Seed(time.Now().UTC().UnixNano())
	return uint8(min + rand.Intn(max-min))
}
