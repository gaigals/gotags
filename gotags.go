package gotags

const (
	defaultSeparator = ";"
	defaultEquals    = ":"
)

// NewKey creates new tag key.
// name - key name.
// isBool - is key boolean (does not require value).
// isRequired - is this key required if tag is defined.
// validator - key value validator interface. (optional)
func NewKey(name string, isBool, isRequired bool, validator Validator) Key {
	return Key{
		Name:       name,
		IsBool:     isBool,
		IsRequired: isRequired,
		Validator:  validator,
	}
}

// NewTagSettings creates new Tag with custom separator and equals.
// name - tag name, for example, "validator" (`validator:"gt=10"`).
// separator - char which will separate keys, like, gt=10,lt=20.
// equals - char which defines key value, like, gt=10.
// processor - processor which will process each tag field. (optional)
// includeNotTagged - include fields not tagged with provided name.
// keys - enabled keys.
func NewTagSettings(
	name, separator, equals string,
	processor Processor,
	includeNotTagged bool,
	keys ...Key,
) TagSettings {
	tg := TagSettings{
		Name:             name,
		Separator:        separator,
		Equals:           equals,
		Keys:             keys,
		Processor:        processor,
		IncludeNotTagged: includeNotTagged,
		keysRequired:     nil,
	}
	tg.keysRequired = tg.requiredKeys()
	return tg
}

// NewTagSettingsDefault creates new Tag with default separator(;) and equals(:).
// name - tag name, for example, "validator" (`validator:"gt:10"`).
// processor - processor which will process each tag field. (optional)
// keys - enabled keys.
func NewTagSettingsDefault(name string, processor Processor, keys ...Key) TagSettings {
	return NewTagSettings(
		name,
		defaultSeparator,
		defaultEquals,
		processor,
		false,
		keys...,
	)
}
