package gotags

import (
	"fmt"
	"reflect"
)

type FieldData struct {
	Self    reflect.Value
	Name    string
	Kind    reflect.Kind
	TagData []TagData
}

func (fd FieldData) KeyValueBool(key string) (value string, ok bool) {
	for _, tag := range fd.TagData {
		if tag.Key == key {
			return tag.Value, true
		}
	}

	return "", false
}

func (fd FieldData) KeyValue(key string) string {
	value, _ := fd.KeyValueBool(key)
	return value
}

func (fd FieldData) HasKey(key string) bool {
	_, ok := fd.KeyValueBool(key)
	return ok
}

func (fd FieldData) HasSameType(targetType reflect.Type) bool {
	return reflect.TypeOf(fd.Self.Interface()) == targetType
}

func (fd FieldData) ApplySelfValue(value any) error {
	dataValueOf := reflect.ValueOf(value)

	// Safety checks ...
	if !fd.Self.CanSet() {
		return fmt.Errorf("%s: cannot be changed", fd.Name)
	}
	if !fd.HasSameType(reflect.TypeOf(value)) {
		return fmt.Errorf("%s(%s): cannot apply value of type %s\n",
			reflect.TypeOf(fd.Self.Interface()), reflect.TypeOf(value))
	}

	fd.Self.Set(dataValueOf)
	return nil
}
