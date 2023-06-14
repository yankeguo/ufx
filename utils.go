package ufx

// jointPoint is a virtual type for ensuring invocation order
type jointPoint struct{}

// named arbitrary type with a name
type named[T any] struct {
	Name string
	Val  T
}
