package utils

import "strings"

func FilterCharacters(in string) (out string) {
	out = in
	for c, r := range replaceCharacters {
		out = strings.ReplaceAll(out, c, r)
	}
	return
}
