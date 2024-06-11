package next_token

import (
	"errors"
	"reflect"
	"testing"
)

func TestEnvironmentKeySource_GetKey(t *testing.T) {
	s := EnvironmentKeySource{VariableName: "NEXT_TOKEN_KEY"}

	t.Run("ExistingKey", func(t *testing.T) {
		t.Setenv("NEXT_TOKEN_KEY", "qrvM3Q==")

		expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}

		key, err := s.GetKey()
		if err != nil {
			t.Errorf("no error expected")
		}

		if !reflect.DeepEqual(key, expected) {
			t.Errorf("key and expected don't match")
		}
	})

	t.Run("MissingKey", func(t *testing.T) {
		_, err := s.GetKey()
		if !errors.Is(err, ErrMissingKeyVariable) {
			t.Errorf("expected error")
		}
	})
}
