package gotags

import (
	"errors"
	"fmt"
	"strings"
)

type Tag struct {
	Key   string
	Value string
}

func NewTagFromString(tagStr, equals string) (Tag, error) {
	splitted := strings.SplitN(tagStr, equals, 2)
	splittedLen := len(splitted)

	if splittedLen == 0 {
		return Tag{}, errors.New("no keys defined")
	}

	if splittedLen > 2 {
		return Tag{}, fmt.Errorf("unexpected tagStr format '%s'", tagStr)
	}

	tag := Tag{
		Key: splitted[0],
	}

	// If tagStr has value ...
	if splittedLen == 2 {
		tag.Value = splitted[1]
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
