package gotags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Processor can be used to do some custom stuff for each field (if defined)
// and gets triggered after key validation (if passed).
type Processor func(fieldData Field) error

// TagSettings holds data about tag.
type TagSettings struct {
	Name      string
	Separator string
	Equals    string
	Keys      []Key
	Processor // Optional

	keysRequired []string
}

// AddKey adds new key to TagSettings
func (tg *TagSettings) AddKey(name string, isBool, isRequired bool, validator Validator, allowedKinds ...reflect.Kind) {
	tg.Keys = append(tg.Keys, NewKey(name, isBool, isRequired, validator, allowedKinds...))
	tg.keysRequired = tg.requiredKeys()
}

// ParseStruct parses passed struct and triggers validators if defined
// and field processors if defined.
func (tg *TagSettings) ParseStruct(data any) ([]Field, error) {
	structure, err := tg.unpackPtr(data)
	if err != nil {
		return nil, err
	}

	fieldData, err := tg.unpackStruct(structure)
	if err != nil {
		return nil, err
	}

	err = tg.runProcessor(fieldData)
	if err != nil {
		return nil, err
	}

	return fieldData, nil
}

func (tg *TagSettings) runProcessor(fieldData []Field) error {
	if tg.Processor == nil {
		return nil
	}

	for _, field := range fieldData {
		err := tg.Processor(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tg *TagSettings) unpackPtr(data any) (reflect.Value, error) {
	valueOf := reflect.ValueOf(data)

	if valueOf.Kind() != reflect.Ptr || valueOf.IsNil() {
		return reflect.Value{}, errors.New("passed value must be valid pointer")
	}

	return valueOf.Elem(), nil
}

func (tg *TagSettings) unpackStruct(valueOf reflect.Value) ([]Field, error) {
	typeOf := reflect.TypeOf(valueOf)

	if typeOf.Kind() != reflect.Struct {
		return nil, errors.New("passed value must be pointer of struct")
	}

	if typeOf.NumField() == 0 {
		return nil, nil
	}

	return tg.parseFields(valueOf)
}

func (tg *TagSettings) parseFields(valueOf reflect.Value) ([]Field, error) {
	typeOf := reflect.TypeOf(valueOf.Interface())

	fields := make([]Field, typeOf.NumField())

	for i := 0; i < typeOf.NumField(); i++ {
		structField := typeOf.Field(i)
		if !structField.IsExported() || structField.Tag == "" {
			continue
		}

		tags := tg.readTagContent(structField.Tag)
		if len(tags) == 0 {
			continue
		}

		tagData, err := tg.convertAsTagData(tags)
		if err != nil {
			return nil, err
		}

		err = tg.validateTagData(tagData)
		if err != nil {
			return nil, fmt.Errorf("field '%s': %w", structField.Name, err)
		}

		fields[i] = Field{
			Value:   valueOf.Field(i),
			Name:    structField.Name,
			Kind:    structField.Type.Kind(),
			TagData: tagData,
		}

		err = tg.hasRequiredKeys(fields[i])
		if err != nil {
			return nil, err
		}
	}

	return fields, nil
}

func (tg *TagSettings) readTagContent(tag reflect.StructTag) []string {
	tagString, ok := tag.Lookup(tg.Name)
	if !ok { // No pkg tag key, ignore this struct field.
		return nil
	}

	return strings.Split(tagString, tg.Separator)
}

func (tg *TagSettings) convertAsTagData(tags []string) ([]Tag, error) {
	tagsData := make([]Tag, len(tags))

	for k, v := range tags {
		tagData, err := newTagData(v, tg.Equals)
		if err != nil {
			return nil, err
		}

		tagsData[k] = tagData
	}

	return tagsData, nil
}

func (tg *TagSettings) validateTagData(tagData []Tag) error {
	for _, tag := range tagData {
		key := tg.findMatchingKey(tag.Key)
		if key == nil {
			return fmt.Errorf("tag '%s' does not exist", tag.Key)
		}

		err := tag.validate(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tg *TagSettings) findMatchingKey(key string) *Key {
	for _, v := range tg.Keys {
		if key == v.Name {
			return &v
		}
	}

	return nil
}

func (tg *TagSettings) hasRequiredKeys(fieldData Field) error {
	for _, v := range tg.keysRequired {
		if fieldData.HasKey(v) {
			continue
		}

		return fmt.Errorf("%s: key '%s' is required but not found",
			fieldData.Name, v)
	}

	return nil
}

func (tg *TagSettings) requiredKeys() []string {
	required := make([]string, 0)

	for _, v := range tg.Keys {
		if v.IsRequired {
			required = append(required, v.Name)
		}
	}

	return required
}
