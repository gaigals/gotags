package gotags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Processor can be used to do some custom stuff for each field (if defined)
// and gets triggered after key validation (if passed).
type Processor func(field Field) error

// TagSettings holds data about tag.
type TagSettings struct {
	Name                 string
	Separator            string
	Equals               string
	Keys                 []Key
	Processor                 // Optional
	IncludeNotTagged     bool // Include not tagged fields
	DisableKeyValueLogic bool // Disable key/value support, default false.
	keysRequired         []string
}

func NewSettings(name string) *TagSettings {
	return &TagSettings{
		Name:      name,
		Separator: defaultSeparator,
		Equals:    defaultEquals,
	}
}

// WithNoKeyValueLogic tells TagSettings that you won't use key/value logic
// and expect only value, like, `myTagName:"myValue".
// This method is completely opposite of TagSettings.WithCustomSeparators().
func (tg *TagSettings) WithNoKeyValueLogic() *TagSettings {
	tg.DisableKeyValueLogic = true
	return tg
}

// WithCustomSeparators can be used to set custom separator and equals key to
// your desired characters.
// By default, TagSettings use separator - ";" and equals - ":"
// ("key:value;otherKey:otherValue").
func (tg *TagSettings) WithCustomSeparators(separator, equals string) *TagSettings {
	tg.Separator = separator
	tg.Equals = equals
	tg.DisableKeyValueLogic = false
	return tg
}

// WithProcessor adds field processor.
// Processor can be used to some custom stuff for each parsed field.
// This is optional and is not required.
// Processor gets called after tag validation.
func (tg *TagSettings) WithProcessor(processor Processor) *TagSettings {
	tg.Processor = processor
	return tg
}

// IncludeUntaggedFields tells TagSettings to parse and include in results
// not tagged struct fields.
func (tg *TagSettings) IncludeUntaggedFields() *TagSettings {
	tg.IncludeNotTagged = true
	return tg
}

// AddKeys can be used to add new keys to TagSettings.
// Note: this method does not check for duplicates.
func (tg *TagSettings) AddKeys(keys ...Key) *TagSettings {
	for _, key := range keys {
		_ = tg.AddKey(key)
	}

	return tg
}

// AddKey adds new key to TagSettings.
// Note: this method does not check for duplicates.
func (tg *TagSettings) AddKey(key Key) *TagSettings {
	tg.Keys = append(tg.Keys, key)
	if key.IsRequired {
		tg.keysRequired = append(tg.keysRequired, key.Name)
	}
	return tg
}

// RemoveKey removes key from registered keys if exists.
func (tg *TagSettings) RemoveKey(name string) {
	for idx, key := range tg.Keys {
		if key.Name != name {
			continue
		}

		tg.Keys = append(tg.Keys[:idx], tg.Keys[idx+1:]...)
	}
}

// ParseStruct parses passed struct and triggers validators if defined
// and field processors if defined.
func (tg *TagSettings) ParseStruct(data any) ([]Field, error) {
	structure, err := tg.unpackPtr(data)
	if err != nil {
		return nil, err
	}

	fields, err := tg.unpackStruct(structure)
	if err != nil {
		return nil, err
	}

	err = tg.runProcessor(fields)
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func (tg *TagSettings) runProcessor(fields []Field) error {
	if tg.Processor == nil {
		return nil
	}

	for _, field := range fields {
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
	valueOf, err := tg.tryUnpackInterface(valueOf)
	if err != nil {
		return nil, err
	}

	typeOf := reflect.TypeOf(valueOf.Interface())

	fields := make([]Field, typeOf.NumField())
	addedFields := 0

	for i := 0; i < typeOf.NumField(); i++ {
		structField := typeOf.Field(i)
		if !structField.IsExported() ||
			(structField.Tag == "" && !tg.IncludeNotTagged) {
			continue
		}

		tagsSplitted := tg.readTagContent(structField.Tag)
		if len(tagsSplitted) == 0 && !tg.IncludeNotTagged {
			continue
		}

		tags, err := tg.convertAsTags(tagsSplitted)
		if err != nil {
			return nil, err
		}

		err = tg.validateTags(tags)
		if err != nil {
			return nil, fmt.Errorf("field '%s': %w", structField.Name, err)
		}

		fields[addedFields] = Field{
			Value: valueOf.Field(i),
			Name:  structField.Name,
			Kind:  structField.Type.Kind(),
			Tags:  tags,
		}

		err = tg.hasRequiredKeys(fields[addedFields])
		if err != nil {
			return nil, err
		}

		addedFields++
	}

	if addedFields != typeOf.NumField() {
		fields = fields[:addedFields]
	}

	return fields, nil
}

func (tg *TagSettings) tryUnpackInterface(valueOf reflect.Value) (reflect.Value, error) {
	if valueOf.Kind() == reflect.Struct {
		return valueOf, nil
	}

	if valueOf.Kind() != reflect.Interface && valueOf.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("underlying value must be struct")
	}

	return tg.tryUnpackInterface(valueOf.Elem())
}

func (tg *TagSettings) readTagContent(tag reflect.StructTag) []string {
	tagString, ok := tag.Lookup(tg.Name)
	if !ok { // No pkg tag key, ignore this struct field.
		return nil
	}

	if tg.DisableKeyValueLogic {
		return []string{tagString}
	}

	return strings.Split(tagString, tg.Separator)
}

func (tg *TagSettings) convertAsTags(tags []string) ([]Tag, error) {
	tagsSlice := make([]Tag, len(tags))

	for k, v := range tags {
		if tg.DisableKeyValueLogic {
			tagsSlice[k] = newTagFromValueString(v)
			continue
		}

		tag, err := NewTagFromString(v, tg.Equals)
		if err != nil {
			return nil, err
		}

		tagsSlice[k] = tag
	}

	return tagsSlice, nil
}

func (tg *TagSettings) validateTags(tags []Tag) error {
	for _, tag := range tags {
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

func (tg *TagSettings) hasRequiredKeys(field Field) error {
	for _, v := range tg.keysRequired {
		if field.HasKey(v) {
			continue
		}

		return fmt.Errorf("%s: key '%s' is required but not found",
			field.Name, v)
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
