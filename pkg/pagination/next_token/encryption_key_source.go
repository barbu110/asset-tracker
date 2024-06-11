package next_token

import "errors"

var ErrKeyDecode = errors.New("key decoding failed")

type EncryptionKeySource interface {
	GetKey() ([]byte, error)
}
