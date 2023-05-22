package main

func Prepend[T any](slice []T, head ...T) []T {
	return append(head, slice...)
}
