package gotags

// Validator can be used to validate key value pair.
type Validator func(value string) error

// Key holds data about specific key.
type Key struct {
	Validator
	Name       string
	IsBool     bool
	IsRequired bool
}
