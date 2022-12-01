package gotags

import (
    "reflect"
)

type FieldData struct {
    Self    any
    Name    string
    Kind    reflect.Kind
    TagData []TagData
}

func (fd FieldData) HasKey(key string) bool {
    for _, tag := range fd.TagData {
        if key == tag.Key {
            return true
        }
    }

    return false
}
