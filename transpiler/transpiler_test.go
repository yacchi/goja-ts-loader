package transpiler

import (
	"strings"
	"testing"
)

func TestTranspile(t *testing.T) {
	tp, err := NewTranspiler()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		ModuleName string
		Source     string
		Expect     string
	}{
		{
			ModuleName: "case1",
			Source:     `const a = 1`,
			Expect:     `var a = 1;`,
		},
		{
			ModuleName: "case2",
			Source:     `const f = () => { return 1 }`,
			Expect:     `var f = function () { return 1; };`,
		},
	}

	for idx, c := range cases {
		s, err := tp.Transpile(c.ModuleName, c.Source)
		s = strings.TrimSuffix(s, "\r\n")
		if err != nil {
			t.Fatal(err)
		} else {
			if c.Expect != s {
				t.Errorf("case: %d expected: '%s' actual: '%s' \n", idx+1, c.Expect, s)
			}
		}
	}
}
