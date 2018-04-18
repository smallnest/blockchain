package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"math/rand"
	"time"

	secp256k1 "github.com/toxeus/go-secp256k1"
)

// Sign 对数据data进行签名.
func Sign(privateKey string, data []byte) (sign []byte, err error) {
	priKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	var privateKeyBytes32 [32]byte
	copy(privateKeyBytes32[:], priKey)

	shaHash := sha256.New()
	shaHash.Write(data)
	var hash = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	dhash := shaHash2.Sum(nil)

	var dhashed [32]byte
	copy(dhashed[:], dhash)

	secp256k1.Start()
	signed, success := secp256k1.Sign(dhashed, privateKeyBytes32, generateNonce())
	defer secp256k1.Stop()
	if !success {
		return nil, errors.New("failed to sign data")
	}

	return signed, nil
}

// Verify 校验签名的真实性.
func Verify(publicKey string, signed []byte, data []byte) bool {
	pubKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return false
	}
	shaHash := sha256.New()
	shaHash.Write(data)
	var hash = shaHash.Sum(nil)

	shaHash2 := sha256.New()
	shaHash2.Write(hash)
	dhash := shaHash2.Sum(nil)

	var dhashed [32]byte
	copy(dhashed[:], dhash)

	secp256k1.Start()
	verified := secp256k1.Verify(dhashed, signed, pubKey)
	defer secp256k1.Stop()
	return verified
}
func generateNonce() *[32]byte {
	var bytes [32]byte
	for i := 0; i < 32; i++ {
		//This is not "cryptographically random"
		bytes[i] = byte(randInt(0, math.MaxUint8))
	}
	return &bytes
}

func randInt(min int, max int) uint8 {
	rand.Seed(time.Now().UTC().UnixNano())
	return uint8(min + rand.Intn(max-min))
}
