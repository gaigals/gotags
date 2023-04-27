package gotags

import (
	"fmt"
	"reflect"
)

// Field contains all information about struct field.
type Field struct {
	Value reflect.Value // Pointer to field
	Name  string        // Field name
	Kind  reflect.Kind  // Field type/kind
	Tags  []Tag         // Field tag data
}

// KeyValueBool acquires tag key value.
// Returns ok(true) if key exists.
func (field Field) KeyValueBool(key string) (value string, ok bool) {
	for _, tag := range field.Tags {
		if tag.Key == key {
			return tag.Value, true
		}
	}

	return "", false
}

// KeyValue returns tag key value.
func (field Field) KeyValue(key string) string {
	value, _ := field.KeyValueBool(key)
	return value
}

// HasKey checks if field contains tag key.
func (field Field) HasKey(key string) bool {
	_, ok := field.KeyValueBool(key)
	return ok
}

// HasType checks if field has passed type.
func (field Field) HasType(targetType reflect.Type) bool {
	return reflect.TypeOf(field.Value.Interface()) == targetType
}

// SetValue sets new value for field.
func (field Field) SetValue(value any) error {
	// Safety checks ...
	if !field.Value.CanSet() {
		return fmt.Errorf("%s: cannot be changed", field.Name)
	}
	if !field.HasType(reflect.TypeOf(value)) {
		return fmt.Errorf("%s(%s): cannot apply value of type %v\n",
			reflect.TypeOf(field.Value.Interface()), reflect.TypeOf(value), value)
	}

	field.Value.Set(reflect.ValueOf(value))
	return nil
}

func (field Field) FirstTag() Tag {
	if len(field.Tags) == 0 {
		return Tag{}
	}

	return field.Tags[0]
}
