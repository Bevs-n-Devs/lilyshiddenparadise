package handlers

import (
	"strings"
)

// Contains checks if a string contains another string
// This is exported so it can be used in tests
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
