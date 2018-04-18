package blockchain

import (
	"encoding/hex"
	"testing"

	"github.com/smallnest/blockchain/wallet"
)

func TestVerify(t *testing.T) {
	privateKey, _, publicKey, _ := wallet.GenerateKeys()

	data := []byte("飞鸽传输")
	signed, err := Sign(privateKey, data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("signed: %s", hex.EncodeToString(signed))
	ok := Verify(publicKey, signed, data)
	if !ok {
		t.Errorf("verify failed")
	}
}
