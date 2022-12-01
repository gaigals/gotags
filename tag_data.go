package gotags

import (
    "errors"
    "fmt"
    "strings"
)

type TagData struct {
    Key   string
    Value string
}

func newTagData(tag, equals string) (*TagData, error) {
    splitted := strings.Split(tag, equals)
    splittedLen := len(splitted)

    if splittedLen == 0 {
        return nil, errors.New("no keys defined")
    }

    if splittedLen > 2 {
        return nil, fmt.Errorf("unexpected tag format '%s'", tag)
    }

    tagData := TagData{
        Key: splitted[0],
    }

    if splittedLen == 2 {
        tagData.Value = splitted[1]
    }

    return &tagData, nil
}

func (td *TagData) validate(key *Key) error {
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
