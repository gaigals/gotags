package gotags

import (
	"fmt"
	"reflect"
)

// FieldData contains all information about struct field.
type FieldData struct {
	Value   reflect.Value // Pointer to field
	Name    string        // Field name
	Kind    reflect.Kind  // Field type/kind
	TagData []TagData     // Field tag data
}

// KeyValueBool acquires tag key value.
// Returns ok(true) if key exists.
func (fd FieldData) KeyValueBool(key string) (value string, ok bool) {
	for _, tag := range fd.TagData {
		if tag.Key == key {
			return tag.Value, true
		}
	}

	return "", false
}

// KeyValue returns tag key value.
func (fd FieldData) KeyValue(key string) string {
	value, _ := fd.KeyValueBool(key)
	return value
}

// HasKey checks if field contains tag key.
func (fd FieldData) HasKey(key string) bool {
	_, ok := fd.KeyValueBool(key)
	return ok
}

// HasType checks if field has passed type.
func (fd FieldData) HasType(targetType reflect.Type) bool {
	return reflect.TypeOf(fd.Value.Interface()) == targetType
}

// SetValue sets new value for field.
func (fd FieldData) SetValue(value any) error {
	dataValueOf := reflect.ValueOf(value)

	// Safety checks ...
	if !fd.Value.CanSet() {
		return fmt.Errorf("%s: cannot be changed", fd.Name)
	}
	if !fd.HasType(reflect.TypeOf(value)) {
		return fmt.Errorf("%s(%s): cannot apply value of type %s\n",
			reflect.TypeOf(fd.Value.Interface()), reflect.TypeOf(value))
	}

	fd.Value.Set(dataValueOf)
	return nil
}

// TagDataFormatted formats all field keys and values in provided format.
// For example, format, `%s=%s` will result in `key=value`.
func (fd FieldData) TagDataFormatted(format string) []string {
	slice := make([]string, 0)

	for _, v := range fd.TagData {
		slice = append(slice, v.StringFormatted(format))
	}

	return slice
}
