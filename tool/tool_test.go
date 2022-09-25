package tool

import "testing"

type replaceTest struct {
	s string
	replace string
	replaceLeft bool
	answer string
}

var replacetest []replaceTest = []replaceTest{
	{"f/c", "j+d", true, "(j+d)/c"},
	{"f/e", "j/d", true, "(j/d)/e"},
	{"f/e", "j/d", false, "f/(j/d)"},
	{"f-e", "j+d", false, "f-(j+d)"},
}

func TestValidSample(t *testing.T) {
	for _, ele := range replacetest {
		res := ReplaceExpression(ele.s, ele.replace, ele.replaceLeft)
		if res != ele.answer {
			t.Errorf("Incorrect, got: %s, want: %s.", res, ele.answer)
		}
	}
	
}