package gotags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/MarvinJWendt/testza"
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

func testValidatorOk(value string) error {
	return nil
}

func testProcessorOk(field Field) error {
	return nil
}

func testValidatorErr(value string) error {
	return errors.New("some error")
}

func testProcessorErr(field Field) error {
	return errors.New("some error")
}

func testProcessorToUpper(field Field) error {
	if field.Kind != reflect.String {
		return errors.New("expected string kind")
	}

	field.SetValue(strings.ToUpper(field.Value.String()))
	return nil
}

func Test_NewTagSettings(t *testing.T) {
	testCases := []struct {
		Key
		Processor
		TagName          string
		Seperator        string
		Equal            string
		IncludeNotTagged bool
	}{
		{NewKey("isTest", true, false, nil), nil, "testTag", ",", "=", false},
		{NewKey("isTest", true, true, testValidatorOk), testProcessorOk, "sometag", ";", ":", true},
	}

	for _, v := range testCases {
		tagSettings := NewTagSettings(
			v.TagName,
			v.Seperator,
			v.Equal,
			v.Processor,
			v.IncludeNotTagged,
			v.Key,
		)

		testza.AssertEqual(t, tagSettings.Name, v.TagName, "unexpected tagName")
		testza.AssertEqual(t, tagSettings.Separator, v.Seperator, "unexpected seperator")
		testza.AssertEqual(t, tagSettings.Equals, v.Equal, "unexpected equals")
		testza.AssertEqual(t, tagSettings.IncludeNotTagged, v.IncludeNotTagged,
			"unexpected IncludeNotTagged")
		testza.AssertEqual(
			t,
			fmt.Sprintf("%p", tagSettings.Processor),
			fmt.Sprintf("%p", v.Processor),
			"unexpected processor",
		)

		testza.AssertLen(t, tagSettings.Keys, 1, "unexpected keys len")
		testza.AssertEqual(t, tagSettings.Keys[0].Name, v.Key.Name,
			"unexpected key name")
		testza.AssertEqual(t, tagSettings.Keys[0].IsBool, v.Key.IsBool,
			"unexpected key isBool value")
		testza.AssertEqual(t, tagSettings.Keys[0].IsRequired, v.Key.IsRequired,
			"unexpected key isRequired value")
		testza.AssertEqual(
			t,
			fmt.Sprintf("%p", tagSettings.Keys[0].Validator),
			fmt.Sprintf("%p", v.Key.Validator),
			"unexpected key validator value",
		)
	}
}

func Test_ParseStruct(t *testing.T) {
	t.Run("Basic test, tagged only", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"required"`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("required", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertLen(t, fields, 1, "unexpected fields len")

		testza.AssertEqual(t, fields[0].Name, "Name", "unexpected field name")
		testza.AssertEqual(t, fields[0].Kind, reflect.String, "unexpected field kind")
		testza.AssertEqual(t, fields[0].Value.String(), "Test", "unexpected field value")

		testza.AssertLen(t, fields[0].Tags, 1, "unexpected field tag len")
		testza.AssertEqual(t, fields[0].Tags[0].Key, "required", "unexpected field tag key")
		testza.AssertEqual(t, fields[0].Tags[0].Value, "", "unexpected field tag key value")
	})

	t.Run("Basic test, include not tagged only", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"required"`
			Age  int
		}{
			Name: "Test",
			Age:  55,
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			true, // Include every field, event untagged.
			NewKey("required", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertLen(t, fields, 2, "unexpected fields len")

		testza.AssertEqual(t, fields[0].Name, "Name", "unexpected field name")
		testza.AssertEqual(t, fields[0].Kind, reflect.String, "unexpected field kind")
		testza.AssertEqual(t, fields[0].Value.String(), "Test", "unexpected field value")

		testza.AssertLen(t, fields[0].Tags, 1, "unexpected field tag len")
		testza.AssertEqual(t, fields[0].Tags[0].Key, "required", "unexpected field tag key")
		testza.AssertEqual(t, fields[0].Tags[0].Value, "", "unexpected field tag key value")

		testza.AssertEqual(t, fields[1].Name, "Age", "unexpected field name")
		testza.AssertEqual(t, fields[1].Kind, reflect.Int, "unexpected field kind")
		testza.AssertEqual(t, fields[1].Value.Int(), int64(55), "unexpected field value")

		testza.AssertLen(t, fields[1].Tags, 0, "unexpected field tag len")
	})

	t.Run("Check string tag with value", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"value=Test1"`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("value", false, true, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertLen(t, fields, 1, "unexpected fields len")

		testza.AssertEqual(t, fields[0].Name, "Name", "unexpected field name")
		testza.AssertEqual(t, fields[0].Kind, reflect.String, "unexpected field kind")
		testza.AssertEqual(t, fields[0].Value.String(), "Test", "unexpected field value")

		testza.AssertLen(t, fields[0].Tags, 1, "unexpected field tag len")
		testza.AssertEqual(t, fields[0].Tags[0].Key, "value", "unexpected field tag key")
		testza.AssertEqual(t, fields[0].Tags[0].Value, "Test1", "unexpected field tag key value")
	})

	t.Run("Check is required tag defined", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"value=Test1,required"`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("value", false, true, nil),
			NewKey("required", true, true, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertLen(t, fields, 1, "unexpected fields len")

		testza.AssertEqual(t, fields[0].Name, "Name", "unexpected field name")
		testza.AssertEqual(t, fields[0].Kind, reflect.String, "unexpected field kind")
		testza.AssertEqual(t, fields[0].Value.String(), "Test", "unexpected field value")

		testza.AssertLen(t, fields[0].Tags, 2, "unexpected field tag len")

		testza.AssertEqual(t, fields[0].Tags[0].Key, "value", "unexpected field tag key")
		testza.AssertEqual(t, fields[0].Tags[0].Value, "Test1", "unexpected field tag key value")

		testza.AssertEqual(t, fields[0].Tags[1].Key, "required", "unexpected field tag key")
		testza.AssertEqual(t, fields[0].Tags[1].Value, "", "unexpected field tag key value")
	})

	t.Run("Check does validation work", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"value=Test1"`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("value", false, true, testValidatorOk),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertLen(t, fields, 1, "unexpected fields len")

		testza.AssertEqual(t, fields[0].Name, "Name", "unexpected field name")
		testza.AssertEqual(t, fields[0].Kind, reflect.String, "unexpected field kind")
		testza.AssertEqual(t, fields[0].Value.String(), "Test", "unexpected field value")

		testza.AssertLen(t, fields[0].Tags, 1, "unexpected field tag len")

		testza.AssertEqual(t, fields[0].Tags[0].Key, "value", "unexpected field tag key")
		testza.AssertEqual(t, fields[0].Tags[0].Value, "Test1", "unexpected field tag key value")
	})

	t.Run("Check does upper case processor work", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"required"`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			testProcessorToUpper,
			false,
			NewKey("required", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertLen(t, fields, 1, "unexpected fields len")

		testza.AssertEqual(t, fields[0].Name, "Name", "unexpected field name")
		testza.AssertEqual(t, fields[0].Kind, reflect.String, "unexpected field kind")
		testza.AssertEqual(t, fields[0].Value.String(), "TEST", "unexpected field value")

		testza.AssertLen(t, fields[0].Tags, 1, "unexpected field tag len")
	})

	t.Run("Check if empty tag does not return error", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:""`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("required", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error tag not defined", func(t *testing.T) {
		testStruct := struct {
			Name string `testtag:"value=Test1"`
			Age  int
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("required", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error tag not provided", func(t *testing.T) {
		testStruct := struct {
			Name string
			Age  int `testtag:"hmm"`
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("required", true, true, nil),
			NewKey("hmm", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error key does not take value", func(t *testing.T) {
		testStruct := struct {
			Name string
			Age  int `testtag:"hmm=1231"`
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("hmm", true, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error key takes value", func(t *testing.T) {
		testStruct := struct {
			Name string
			Age  int `testtag:"hmm"`
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("hmm", false, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error key takes value", func(t *testing.T) {
		testStruct := struct {
			Name string
			Age  int `testtag:"hmm"`
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("hmm", false, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error tag validator failed", func(t *testing.T) {
		testStruct := struct {
			Name string
			Age  int `testtag:"hmm=abc"`
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			nil,
			false,
			NewKey("hmm", false, false, testValidatorErr),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})

	t.Run("Check error processor failed", func(t *testing.T) {
		testStruct := struct {
			Name string
			Age  int `testtag:"hmm=abc"`
		}{
			Name: "Test",
		}

		tagSettings := NewTagSettings(
			"testtag",
			",",
			"=",
			testProcessorErr,
			false,
			NewKey("hmm", false, false, nil),
		)

		fields, err := tagSettings.ParseStruct(&testStruct)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, fields, "fields expected as nil")
	})
}
