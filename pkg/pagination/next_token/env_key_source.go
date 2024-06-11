package next_token

import (
	"encoding/base64"
	"errors"
	"os"
)

var ErrMissingKeyVariable = errMissingKeyVariable()

type EnvironmentKeySource struct {
	VariableName string
}

func (s *EnvironmentKeySource) GetKey() ([]byte, error) {
	v := os.Getenv(s.VariableName)
	if v == "" {
		return nil, ErrMissingKeyVariable
	}
	key, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, ErrKeyDecode
	}

	return key, nil
}

func errMissingKeyVariable() error {
	return errors.New("missing key variable")
}
