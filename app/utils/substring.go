package utils

// Get a defined substring with x characters from the start.
func Substring(s string, i int) string {
	return s[:min(len(s), i)]
}
