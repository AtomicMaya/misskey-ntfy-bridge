package utils

import "strings"

func SanitizeString(str string) string {
	targets := []string{"\u202A", "\u202B", "\u202C", "\u202D", "\u202E", "\u200E", "\u200F"}

	edited_str := str
	for _, c := range targets {
		edited_str = strings.Replace(edited_str, c, "", -1)
	}
	return edited_str
}
