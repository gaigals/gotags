package gotags

import (
	"errors"
	"fmt"
)

type Tag struct {
	Key   string
	Value string
}

func NewTagFromString(tagStr, equals string) (Tag, error) {
	return NewTagFromStringWithEscape(tagStr, equals, 0)
}

func NewTagFromStringWithEscape(
	tagStr,
	equals string,
	escapeCharacter byte,
) (Tag, error) {
	return newTagFromString(tagStr, equals, escapeCharacter)
}

func newTagFromString(tagStr, equals string, escapeCharacter byte) (Tag, error) {
	key, value, hasValue, err := splitTagKeyValueWithEscape(
		tagStr,
		equals,
		escapeCharacter,
	)
	if err != nil {
		return Tag{}, err
	}

	if key == "" && tagStr == "" {
		return Tag{}, errors.New("no keys defined")
	}

	tag := Tag{
		Key: key,
	}

	// If tagStr has value ...
	if hasValue {
		tag.Value = value
	}

	return tag, nil
}

func (tag *Tag) validate(key *Key) error {
	if key.IsBool && tag.Value != "" {
		return fmt.Errorf("tag '%s' does not take any arguments", tag.Key)
	}
	if !key.IsBool && tag.Value == "" {
		return fmt.Errorf("tag '%s' requires argument", tag.Key)
	}

	if key.Validator == nil {
		return nil
	}

	return key.Validator(tag.Value)
}

// StringFormatted formats key and value in provided format.
// For example, format, `%s=%s` will result in `key=value`.
func (tag *Tag) StringFormatted(format string) string {
	return fmt.Sprintf(format, tag.Key, tag.Value)
}
