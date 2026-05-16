// Package utils provides common utility helpers for pointers, strings, and time formatting.
package utils

// Ptr returns a pointer to the given value.
func Ptr[T any](value T) *T {
	return &value
}
