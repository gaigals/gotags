package gotags

import (
	"fmt"
	"reflect"
)

// Field contains all information about struct field.
type Field struct {
	Value   reflect.Value // Pointer to field
	Name    string        // Field name
	Kind    reflect.Kind  // Field type/kind
	TagData []Tag         // Field tag data
}

// KeyValueBool acquires tag key value.
// Returns ok(true) if key exists.
func (fd Field) KeyValueBool(key string) (value string, ok bool) {
	for _, tag := range fd.TagData {
		if tag.Key == key {
			return tag.Value, true
		}
	}

	return "", false
}

// KeyValue returns tag key value.
func (fd Field) KeyValue(key string) string {
	value, _ := fd.KeyValueBool(key)
	return value
}

// HasKey checks if field contains tag key.
func (fd Field) HasKey(key string) bool {
	_, ok := fd.KeyValueBool(key)
	return ok
}

// HasType checks if field has passed type.
func (fd Field) HasType(targetType reflect.Type) bool {
	return reflect.TypeOf(fd.Value.Interface()) == targetType
}

// SetValue sets new value for field.
func (fd Field) SetValue(value any) error {
	// Safety checks ...
	if !fd.Value.CanSet() {
		return fmt.Errorf("%s: cannot be changed", fd.Name)
	}
	if !fd.HasType(reflect.TypeOf(value)) {
		return fmt.Errorf("%s(%s): cannot apply value of type %v\n",
			reflect.TypeOf(fd.Value.Interface()), reflect.TypeOf(value), value)
	}

	fd.Value.Set(reflect.ValueOf(value))
	return nil
}

// TagDataFormatted formats all field keys and values in provided format.
// For example, format, `%s=%s` will result in `key=value`.
func (fd Field) TagDataFormatted(format string) []string {
	slice := make([]string, 0)

	for _, v := range fd.TagData {
		slice = append(slice, v.StringFormatted(format))
	}

	return slice
}

func (fd Field) TagMap() map[string]string {
	tagMap := make(map[string]string)

	for _, v := range fd.TagData {
		tagMap[v.Key] = v.Value
	}

	return tagMap
}
