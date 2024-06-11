package next_token

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
)

type EncryptionEngine struct {
	KeySource EncryptionKeySource
}

func (e *EncryptionEngine) NewToken(raw string) NextToken {
	return NextToken{Raw: raw}
}

func (e *EncryptionEngine) Encrypt(token NextToken) ([]byte, error) {
	key, err := e.KeySource.GetKey()
	if err != nil {
		return nil, fmt.Errorf("key retrieval failed: %w", err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %w", err)
	}

	out := bytes.Buffer{}
	c.Encrypt(out.Bytes(), []byte(token.Raw))

	return out.Bytes(), nil
}

func (e *EncryptionEngine) Decrypt(ciphertext []byte) (NextToken, error) {
	key, err := e.KeySource.GetKey()
	if err != nil {
		return NextToken{}, fmt.Errorf("key retrieval failed: %w", err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return NextToken{}, fmt.Errorf("cipher creation failed: %w", err)
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return e.NewToken(string(pt)), nil
}

func (e *EncryptionEngine) EncryptToString(token NextToken) (string, error) {
	b, err := e.Encrypt(token)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func (e *EncryptionEngine) DecryptFromString(s string) (NextToken, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return NextToken{}, fmt.Errorf("base64 decode failed: %w", err)
	}

	t, err := e.Decrypt(b)
	if err != nil {
		return NextToken{}, fmt.Errorf("decryption failed: %w", err)
	}

	return t, nil
}
