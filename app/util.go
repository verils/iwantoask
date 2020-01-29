package app

import "strings"

func Capitalize(s string) string {
	split := strings.Split(s, " ")
	for i, str := range split {
		capitalized := strings.ToUpper(str[0:1]) + str[1:]
		split[i] = capitalized
	}
	return strings.Join(split, " ")
}
