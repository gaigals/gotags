package gotags

import (
	"reflect"
)

// Validator can be used to validate key value pair.
type Validator func(value string) error

// Key holds data about specific key.
type Key struct {
	Name         string
	IsBool       bool
	IsRequired   bool
	AllowedKinds []reflect.Kind // Optional
	Validator                   // Optional
}
