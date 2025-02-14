package utils

func Substring(s string, i int) string {
	return s[:min(len(s), i)]
}
