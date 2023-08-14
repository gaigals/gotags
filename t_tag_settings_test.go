package gotags

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string `gotags:"required,min=2,max=255"`
	Age     int    `gotags:"required,min=10,max=130"`
	Country string `gotags:"required,country"`
	Phone   string `gotags:"required,min=5,max=50,phone"`
}

var keys = []Key{
	NewKey("required", true, false, nil),
	NewKey("min", false, false, nil),
	NewKey("max", false, false, nil),
	NewKey("country", true, false, nil),
	NewKey("phone", true, false, nil),
}

func Benchmark_ParseStruct(b *testing.B) {
	person := Person{
		Name:    "Jimmy",
		Age:     10,
		Country: "Estonia",
		Phone:   "2213123112",
	}

	tagSettings := NewTagSettings("gotags", ",", "=", nil, false, keys...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = tagSettings.ParseStruct(&person)
	}

	b.ReportAllocs()
}

type User struct {
	Name    string `other:"Name" customtag:"tag1;tag2:value;tag3"`
	Country string `other:"Country" customtag:"tag4"`
	Age     int    `other:"Age" customtag:"tag1"`
	Address string `other:"Address"`
	Data    struct {
		ehh bool
	}

	json string
}

func compareFields(t *testing.T, actual, expected []Field) {
	if len(actual) != len(expected) {
		t.Fatalf("Expected %d parsed fields, but got %d", len(expected), len(actual))
	}

	for i, expected := range expected {
		field := actual[i]

		if field.Value != expected.Value || field.Name != expected.Name ||
			field.Kind != expected.Kind || len(field.Tags) != len(expected.Tags) {
			t.Errorf("Expected field %+v, but got %+v", expected, field)
		}

		for j, expectedTag := range expected.Tags {
			tag := field.Tags[j]
			if tag.Key != expectedTag.Key || tag.Value != expectedTag.Value {
				t.Errorf("Expected tag %+v, but got %+v", expectedTag, tag)
			}
		}
	}
}

func TestTagSettings_ParseStruct(t *testing.T) {
	mockValidator := func(value string) error { return nil }
	key1 := NewKey("tag1", true, false, mockValidator, reflect.String)
	key2 := NewKey("tag2", false, false, mockValidator, reflect.String)
	key3 := NewKey("tag3", true, false, mockValidator, reflect.String)
	key4 := NewKey("tag4", true, false, mockValidator, reflect.String)

	t.Run("Valid: ParseStruct", func(t *testing.T) {
		tagSettings := NewTagSettingsDefault("customtag", nil, key1, key2, key3, key4)
		data := User{
			Name:    "John",
			Country: "USA",
			Age:     30,
			Address: "123 Main St",
		}

		fields, err := tagSettings.ParseStruct(&data)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		expectedFields := []Field{
			{Value: reflect.ValueOf(&data.Name).Elem(), Name: "Name", Kind: reflect.String, Tags: []Tag{{Key: "tag1"}, {Key: "tag2", Value: "value"}, {Key: "tag3"}}},
			{Value: reflect.ValueOf(&data.Country).Elem(), Name: "Country", Kind: reflect.String, Tags: []Tag{{Key: "tag4"}}},
			{Value: reflect.ValueOf(&data.Age).Elem(), Name: "Age", Kind: reflect.Int, Tags: []Tag{{Key: "tag1"}}},
		}

		compareFields(t, fields, expectedFields)
	})

	t.Run("Valid: ParseStruct; Include not tagged", func(t *testing.T) {
		tagSettings := NewTagSettingsDefault("customtag", nil, key1, key2, key3, key4)
		tagSettings.IncludeUntaggedFields()

		data := User{
			Name:    "John",
			Country: "USA",
			Age:     30,
			Address: "123 Main St",
		}

		fields, err := tagSettings.ParseStruct(&data)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		expectedFields := []Field{
			{Value: reflect.ValueOf(&data.Name).Elem(), Name: "Name", Kind: reflect.String, Tags: []Tag{{Key: "tag1"}, {Key: "tag2", Value: "value"}, {Key: "tag3"}}},
			{Value: reflect.ValueOf(&data.Country).Elem(), Name: "Country", Kind: reflect.String, Tags: []Tag{{Key: "tag4"}}},
			{Value: reflect.ValueOf(&data.Age).Elem(), Name: "Age", Kind: reflect.Int, Tags: []Tag{{Key: "tag1"}}},
			{Value: reflect.ValueOf(&data.Address).Elem(), Name: "Address", Kind: reflect.String, Tags: nil},
			{Value: reflect.ValueOf(&data.Data).Elem(), Name: "Data", Kind: reflect.Struct, Tags: nil},
		}

		compareFields(t, fields, expectedFields)
	})

	t.Run("Valid: ParseStruct Dynamic tags", func(t *testing.T) {
		tagSettings := NewTagSettingsDefault("other", nil)
		tagSettings.WithNoKeyExistValidation()
		data := User{
			Name:    "John",
			Country: "USA",
			Age:     30,
			Address: "123 Main St",
		}

		fields, err := tagSettings.ParseStruct(&data)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		expectedFields := []Field{
			{Value: reflect.ValueOf(&data.Name).Elem(), Name: "Name", Kind: reflect.String, Tags: []Tag{{Key: "Name"}}},
			{Value: reflect.ValueOf(&data.Country).Elem(), Name: "Country", Kind: reflect.String, Tags: []Tag{{Key: "Country"}}},
			{Value: reflect.ValueOf(&data.Age).Elem(), Name: "Age", Kind: reflect.Int, Tags: []Tag{{Key: "Age"}}},
			{Value: reflect.ValueOf(&data.Address).Elem(), Name: "Address", Kind: reflect.String, Tags: []Tag{{Key: "Address"}}},
		}

		compareFields(t, fields, expectedFields)
	})

	t.Run("Valid: ParseStruct Dynamic tags; Include All fields", func(t *testing.T) {
		tagSettings := NewTagSettingsDefault("other", nil)
		tagSettings.WithNoKeyExistValidation()
		data := User{
			Name:    "John",
			Country: "USA",
			Age:     30,
			Address: "123 Main St",
		}

		fields, err := tagSettings.ParseStruct(&data)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		expectedFields := []Field{
			{Value: reflect.ValueOf(&data.Name).Elem(), Name: "Name", Kind: reflect.String, Tags: []Tag{{Key: "Name"}}},
			{Value: reflect.ValueOf(&data.Country).Elem(), Name: "Country", Kind: reflect.String, Tags: []Tag{{Key: "Country"}}},
			{Value: reflect.ValueOf(&data.Age).Elem(), Name: "Age", Kind: reflect.Int, Tags: []Tag{{Key: "Age"}}},
			{Value: reflect.ValueOf(&data.Address).Elem(), Name: "Address", Kind: reflect.String, Tags: []Tag{{Key: "Address"}}},
			{Value: reflect.ValueOf(&data.Data).Elem(), Name: "Data", Kind: reflect.Struct, Tags: nil},
		}

		compareFields(t, fields, expectedFields)
	})

	t.Run("Invalid: Non-Pointer Value", func(t *testing.T) {
		tagSettings := NewTagSettingsDefault("customtag", nil, key1, key2, key3, key4)

		data := User{
			Name:    "John",
			Country: "USA",
			Age:     30,
			Address: "123 Main St",
		}

		_, err := tagSettings.ParseStruct(data)
		if err == nil {
			t.Fatalf("Expected an error for non-pointer value, but got no error")
		}
	})
}
