package next_token

import (
	"crypto/rand"
	"reflect"
	"testing"
)

type inMemoryKeySource struct {
	Key []byte
}

func (i *inMemoryKeySource) GetKey() ([]byte, error) {
	return i.Key, nil
}

func TestEncryptionEngine_EncryptDecrypt(t *testing.T) {
	k := make([]byte, 32)
	_, _ = rand.Read(k)

	ks := inMemoryKeySource{k}
	engine := EncryptionEngine{KeySource: &ks}

	token := engine.NewToken([]byte("this is a next token"))
	encrypted, err := engine.Encrypt(token)
	if err != nil {
		t.Errorf("could not encrypt")
	}

	t.Logf("Raw token: %v", token.Raw)
	t.Logf("Encrypted: %v", encrypted)

	decrypted, err := engine.Decrypt(encrypted)
	if err != nil {
		t.Errorf("could not decrypt")
	}

	t.Logf("Decrypted: %v", []byte(decrypted.Raw))

	if !reflect.DeepEqual(decrypted.Raw, token.Raw) {
		t.Errorf("symmetry test failed")
	}
}

func TestEncryptionEngine_StringEncryptDecrypt(t *testing.T) {
	k := make([]byte, 32)
	_, _ = rand.Read(k)

	ks := inMemoryKeySource{k}
	engine := EncryptionEngine{KeySource: &ks}

	token := engine.NewToken([]byte("this is a next token"))
	encrypted, err := engine.EncryptToString(token)
	if err != nil {
		t.Errorf("could not encrypt")
	}

	t.Logf("Raw token: %v", token.Raw)
	t.Logf("Encrypted: %v", encrypted)

	decrypted, err := engine.DecryptFromString(encrypted)
	if err != nil {
		t.Errorf("could not decrypt")
	}

	t.Logf("Decrypted: %v", []byte(decrypted.Raw))

	if !reflect.DeepEqual(decrypted.Raw, token.Raw) {
		t.Errorf("symmetry test failed")
	}
}
