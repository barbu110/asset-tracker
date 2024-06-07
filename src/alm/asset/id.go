package asset

import (
	"crypto/rand"
	"fmt"
)

const IdSize = 8

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
