package utils

import "strings"

var replaceCharacters = map[string]string{
	`?`: "",
	`*`: "",
	`:`: "",
	`"`: "",
	`>`: "",
	`<`: "",
	`\`: "",
	`/`: "-",
	`|`: "-",
	`｜`: "-",
	`丨`: "-",
}

func FilterCharacters(in string) (out string) {
	out = strings.Trim(in, " ")
	for c, r := range replaceCharacters {
		out = strings.ReplaceAll(out, c, r)
	}
	return
}
