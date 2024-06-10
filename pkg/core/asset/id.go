package asset

import (
	"crypto/rand"
	"encoding"
	"encoding/hex"
	"fmt"
)

const IdSize = 16

type Id []byte

func (id Id) MarshalBinary() (data []byte, err error) {
	return id, nil
}

func RandomId() Id {
	buf := make([]byte, IdSize)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to generate asset ID: %w", err))
	}

	return buf
}

func ParseId(id string) (Id, error) {
	return hex.DecodeString(id)
}

func EncodeIdToString(id encoding.BinaryMarshaler) string {
	b, _ := id.MarshalBinary()
	return hex.EncodeToString(b)
}
