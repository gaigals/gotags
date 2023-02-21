package gotags

import (
	"reflect"
)

const (
	defaultSeparator = ";"
	defaultEquals    = ":"
)

// NewKey creates new tag key.
// name - key name.
// isBool - is key boolean (does not require value).
// isRequired - is this key required if tag is defined.
// validator - key value validator interface. (optional)
// allowedKinds - allowed types/kinds. (optional)
func NewKey(name string, isBool, isRequired bool, validator Validator, allowedKinds ...reflect.Kind) Key {
	return Key{name, isBool, isRequired, allowedKinds, validator}
}

// NewTagSettings creates new Tag with custom separator and equals.
// name - tag name, for example, "validator" (`validator:"gt=10"`).
// separator - char which will separate keys, like, gt=10,lt=20.
// equals - char which defines key value, like, gt=10.
// processor - processor which will process each tag field. (optional)
// keys - enabled keys.
func NewTagSettings(name, separator, equals string, processor Processor, keys ...Key) TagSettings {
	tg := TagSettings{name, separator, equals, keys, processor, nil}
	tg.keysRequired = tg.requiredKeys()
	return tg
}

// NewTagSettingsDefault creates new Tag with default separator(;) and equals(:).
// name - tag name, for example, "validator" (`validator:"gt:10"`).
// processor - processor which will process each tag field. (optional)
// keys - enabled keys.
func NewTagSettingsDefault(name string, processor Processor, keys ...Key) TagSettings {
	return NewTagSettings(name, defaultSeparator, defaultEquals, processor, keys...)
}
