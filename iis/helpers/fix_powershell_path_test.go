package helpers

import (
	"log"
	"testing"
)

func TestFixPowerShellPath(t *testing.T) {
	data := []struct {
		input    string
		expected string
	}{
		{
			input:    "C:\\\\Windows",
			expected: "C:\\Windows",
		},
		{
			input:    "C:/Windows",
			expected: "C:/Windows",
		},
		{
			input:    "C:\\Windows\\\\System32",
			expected: "C:\\Windows\\System32",
		},
		{
			input:    "//some/unc-path.on.the.network",
			expected: "//some/unc-path.on.the.network",
		},
	}

	for _, v := range data {
		log.Printf("Testing %q", v.input)
		actual := FixPowerShellPath(v.input)
		if actual != v.expected {
			t.Errorf("Exepcted %q but got %q", v.expected, actual)
		}
	}

}
