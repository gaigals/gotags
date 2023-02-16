package gotags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Processor can be used to do some custom stuff for each field (if defined)
// and gets triggered after key validation (if passed).
type Processor func(fieldData FieldData) error

// TagSettings holds data about tag.
type TagSettings struct {
	Name      string
	Separator string
	Equals    string
	Keys      []Key
	Processor // Optional
}

// ParseStruct parses passed struct and triggers validators if defined
// and field processors if defined.
func (tg *TagSettings) ParseStruct(data any) ([]FieldData, error) {
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

func (tg *TagSettings) runProcessor(fieldData []FieldData) error {
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

func (tg *TagSettings) unpackStruct(valueOf reflect.Value) ([]FieldData, error) {
	typeOf := reflect.TypeOf(valueOf)

	if typeOf.Kind() != reflect.Struct {
		return nil, errors.New("passed value must be pointer of struct")
	}

	nrOfFields := typeOf.NumField()
	if nrOfFields == 0 {
		return nil, nil
	}

	return tg.parseFields(valueOf)
}

func (tg *TagSettings) parseFields(valueOf reflect.Value) ([]FieldData, error) {
	typeOf := reflect.TypeOf(valueOf.Interface())

	fields := make([]FieldData, 0)

	for i := 0; i < typeOf.NumField(); i++ {
		structField := typeOf.Field(i)
		if !structField.IsExported() || structField.Tag == "" {
			continue
		}

		tags, err := tg.readTagValue(structField.Tag)
		if err != nil {
			return nil, err
		}
		if len(tags) == 0 {
			continue
		}

		tagData, err := tg.convertAsTagData(tags)
		if err != nil {
			return nil, err
		}

		err = tg.validateTagData(tagData)
		if err != nil {
			return nil, err
		}

		fieldTypeOf := structField.Type
		fieldName := structField.Name
		kind := fieldTypeOf.Kind()

		field := FieldData{valueOf.Field(i), fieldName, kind, tagData}

		err = tg.hasRequiredKeys(field)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func (tg *TagSettings) readTagValue(tag reflect.StructTag) ([]string, error) {
	tagString, ok := tag.Lookup(tg.Name)
	if !ok { // No pkg tag key, ignore this struct field.
		return nil, nil
	}

	return strings.Split(tagString, tg.Separator), nil
}

func (tg *TagSettings) convertAsTagData(tags []string) ([]TagData, error) {
	tagsData := make([]TagData, 0)

	for _, v := range tags {
		tagData, err := newTagData(v, tg.Equals)
		if err != nil {
			return nil, err
		}

		tagsData = append(tagsData, *tagData)
	}

	return tagsData, nil
}

func (tg *TagSettings) validateTagData(tagData []TagData) error {
	for _, tag := range tagData {
		key, err := tg.findMatchingKey(tag.Key)
		if err != nil {
			return err
		}

		err = tag.validate(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tg *TagSettings) findMatchingKey(key string) (*Key, error) {
	for _, v := range tg.Keys {
		if key == v.Name {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("tag '%s' does not exist", key)
}

func (tg *TagSettings) hasRequiredKeys(fieldData FieldData) error {
	required := tg.requiredKeys()

	for _, v := range required {
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
