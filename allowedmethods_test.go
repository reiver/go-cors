package cors

import (
	"testing"
)

func TestAllowedMethods(t *testing.T) {
	tests := []struct{
		Methods  []string
		Expected string
	}{
		{
			Expected: "GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE",
		},
		{
			Methods: []string(nil),
			Expected: "GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE",
		},
		{
			Methods: []string{},
			Expected: "GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE",
		},


		{
			Methods: []string{"apple"},
			Expected: "apple",
		},
		{
			Methods: []string{"apple","BANANA"},
			Expected: "apple, BANANA",
		},
		{
			Methods: []string{"apple","BANANA","Cherry"},
			Expected: "apple, BANANA, Cherry",
		},
	}

	for testNumber, test := range tests {
		actual := allowedMethods(test.Methods...)

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual allowed-method is not what was expected.", testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("METHODS: %#v", test.Methods)
			continue
		}
	}
}
