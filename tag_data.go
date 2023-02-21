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

func newTagData(tag, equals string) (Tag, error) {
	splitted := strings.SplitN(tag, equals, 2)
	splittedLen := len(splitted)

	if splittedLen == 0 {
		return Tag{}, errors.New("no keys defined")
	}

	if splittedLen > 2 {
		return Tag{}, fmt.Errorf("unexpected tag format '%s'", tag)
	}

	tagData := Tag{
		Key: splitted[0],
	}

	// If tag has value ...
	if splittedLen == 2 {
		tagData.Value = splitted[1]
	}

	return tagData, nil
}

func (td *Tag) validate(key *Key) error {
	if key.IsBool && td.Value != "" {
		return fmt.Errorf("tag '%s' does not take any arguments", td.Key)
	}
	if !key.IsBool && td.Value == "" {
		return fmt.Errorf("tag '%s' requires argument", td.Key)
	}

	if key.Validator == nil {
		return nil
	}

	return key.Validator(td.Value)
}

// StringFormatted formats key and value in provided format.
// For example, format, `%s=%s` will result in `key=value`.
func (td *Tag) StringFormatted(format string) string {
	return fmt.Sprintf(format, td.Key, td.Value)
}
