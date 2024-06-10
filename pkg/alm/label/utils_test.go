package label

import (
	"reflect"
	"testing"
	"time"
)

func ensureTimeout(d time.Duration) {
	<-time.After(d)
	panic("Test timed out.")
}

func TestSplitIdInGroups(t *testing.T) {
	go ensureTimeout(5 * time.Second)

	type testCase struct {
		input          string
		maxGroupSize   int
		expectedOutput []string
	}

	cases := []testCase{
		{"12345678", 4, []string{"1234", "5678"}},
		{"1234567", 4, []string{"1234", "567"}},
		{"12345678", 3, []string{"123", "456", "78"}},
		{"", 3, []string{}},
		{"a", 3, []string{"a"}},
	}

	for _, c := range cases {
		if output := splitIdInGroups(c.input, c.maxGroupSize); !reflect.DeepEqual(output, c.expectedOutput) {
			t.Errorf("Expected equality of: %q and %q", output, c.expectedOutput)
		}
	}
}
