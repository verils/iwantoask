package app

import "testing"

func TestCapitalize(t *testing.T) {
	s := "abc def-ghi_jkl"
	capitalized := capitalize(s)
	if capitalized != "Abc Def-ghi_jkl" {
		t.Fail()
	}
}
