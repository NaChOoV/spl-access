package helpers

import "testing"

func TestMaskRun(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Standard RUN", "12345678-9", "12******-9"},
		{"Short RUN", "1234", "1234"},
		{"Empty RUN", "", ""},
		{"No Run", "admin", "admin"},
		{"Multiple hyphens", "12-345678-9", "12-345678-9"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := maskRun(tc.input)
			if actual != tc.expected {
				t.Errorf("maskRun(%q) = %q; want %q", tc.input, actual, tc.expected)
			}
		})
	}
}
