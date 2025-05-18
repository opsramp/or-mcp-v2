package common

import "fmt"

// Errorf wraps errors with a formatted message
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
