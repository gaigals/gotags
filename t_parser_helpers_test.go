package gotags

import (
	"strings"
	"testing"

	"github.com/MarvinJWendt/testza"
)

const testEscapeCharacter = '\\'

func Test_SplitWithOptionalEscapes(t *testing.T) {
	t.Run("No backslashes matches strings split", func(t *testing.T) {
		input := "min=1,max=10"

		actual, err := splitWithOptionalEscapes(input, ",", 0)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, strings.Split(input, ","), "unexpected split result")
	})

	t.Run("Escaped separator does not split", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(
			`one,two\,three,four`,
			",",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"one", "two,three", "four"},
			"unexpected split result")
	})

	t.Run("Inner escapes stay untouched at top level", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(
			`replace=old\,value|new\|value`,
			",",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{`replace=old,value|new\|value`},
			"unexpected split result")
	})

	t.Run("Escaped backslash is preserved", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(`one\\,two`, ",", testEscapeCharacter)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{`one\`, "two"},
			"unexpected split result")
	})

	t.Run("Unknown escapes remain untouched", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(`^\d+\.\d+$`, ",", testEscapeCharacter)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{`^\d+\.\d+$`},
			"unexpected split result")
	})

	t.Run("Trailing naked backslash returns error", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(`one,two\`, ",", testEscapeCharacter)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertNil(t, actual, "expected no split result")
	})

	t.Run("Value list splitting respects escaped pipes", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(
			`part1|part2\|part3|part4`,
			"|",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"part1", "part2|part3", "part4"},
			"unexpected split result")
	})

	t.Run("Custom separator does not split when escaped", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(
			`one#two\#three#four`,
			"#",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"one", "two#three", "four"},
			"unexpected split result")
	})

	t.Run("Custom separator keeps deeper custom equals escapes", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(
			`requiredIf@Type\@admin#replace@old\#value`,
			"#",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{
			`requiredIf@Type\@admin`,
			`replace@old#value`,
		}, "unexpected split result")
	})

	t.Run("Escape disabled keeps backslashes literal", func(t *testing.T) {
		actual, err := splitWithOptionalEscapes(`one,two\,three,four`, ",", 0)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"one", `two\`, "three", "four"},
			"unexpected split result")
	})
}

func Test_SplitFirstWithOptionalEscapes(t *testing.T) {
	t.Run("No backslashes matches strings splitn", func(t *testing.T) {
		input := "value=Test1"

		actual, err := splitFirstWithOptionalEscapes(input, "=", 0)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, strings.SplitN(input, "=", 2),
			"unexpected split result")
	})

	t.Run("Current separator stays in value", func(t *testing.T) {
		actual, err := splitFirstWithOptionalEscapes(
			`Key:value\:x`,
			":",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"Key", "value:x"},
			"unexpected split result")
	})

	t.Run("Parameterized value keeps escaped pipes", func(t *testing.T) {
		actual, err := splitFirstWithOptionalEscapes(
			`Type:admin\|user|Role`,
			":",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"Type", `admin\|user|Role`},
			"unexpected split result")
	})

	t.Run("Inner escapes stay untouched at current layer", func(t *testing.T) {
		actual, err := splitFirstWithOptionalEscapes(
			`Key:value\=x`,
			":",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"Key", `value\=x`},
			"unexpected split result")
	})

	t.Run("URL style values keep escaped colons", func(t *testing.T) {
		actual, err := splitFirstWithOptionalEscapes(
			`URL:https\://example.com`,
			":",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"URL", "https://example.com"},
			"unexpected split result")
	})

	t.Run("Custom equals stays in value when escaped", func(t *testing.T) {
		actual, err := splitFirstWithOptionalEscapes(
			`Key@value\@x`,
			"@",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"Key", "value@x"},
			"unexpected split result")
	})

	t.Run("Custom equals keeps deeper pipe escapes", func(t *testing.T) {
		actual, err := splitFirstWithOptionalEscapes(
			`requiredIf@Type\@admin\|user|Role`,
			"@",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, actual, []string{"requiredIf", `Type@admin\|user|Role`},
			"unexpected split result")
	})
}

func Test_NewTagFromStringWithEscape(t *testing.T) {
	t.Run("Escape disabled keeps separators literal", func(t *testing.T) {
		tag, err := NewTagFromString(`replace=old\,value|new\|value`, "=")
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, tag, Tag{
			Key:   "replace",
			Value: `old\,value|new\|value`,
		}, "unexpected tag")
	})

	t.Run("Escape enabled preserves inner escapes", func(t *testing.T) {
		tag, err := NewTagFromStringWithEscape(
			`replace=old\,value|new\|value`,
			"=",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, tag, Tag{
			Key:   "replace",
			Value: `old\,value|new\|value`,
		}, "unexpected tag")
	})

	t.Run("Escape enabled unescapes current separator", func(t *testing.T) {
		tag, err := NewTagFromStringWithEscape(
			`regex=^test\=value$`,
			"=",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, tag, Tag{
			Key:   "regex",
			Value: "^test=value$",
		}, "unexpected tag")
	})

	t.Run("Escape enabled preserves unknown escapes", func(t *testing.T) {
		tag, err := NewTagFromStringWithEscape(
			`regex=^\d+\.\d+$`,
			"=",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, tag, Tag{
			Key:   "regex",
			Value: `^\d+\.\d+$`,
		}, "unexpected tag")
	})

	t.Run("Trailing naked backslash returns error", func(t *testing.T) {
		tag, err := NewTagFromStringWithEscape(
			`regex=abc\`,
			"=",
			testEscapeCharacter,
		)
		testza.AssertNotNil(t, err, "expected error")
		testza.AssertEqual(t, tag, Tag{}, "unexpected tag")
	})

	t.Run("Custom equals works with escape", func(t *testing.T) {
		tag, err := NewTagFromStringWithEscape(
			`requiredIf@Type\@admin`,
			"@",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, tag, Tag{
			Key:   "requiredIf",
			Value: "Type@admin",
		}, "unexpected tag")
	})

	t.Run("Custom equals preserves deeper escapes", func(t *testing.T) {
		tag, err := NewTagFromStringWithEscape(
			`requiredIf@Type\@admin\|user|Role`,
			"@",
			testEscapeCharacter,
		)
		testza.AssertNoError(t, err, "unexpected error")
		testza.AssertEqual(t, tag, Tag{
			Key:   "requiredIf",
			Value: `Type@admin\|user|Role`,
		}, "unexpected tag")
	})
}

func Test_ParseStructWithEscapedValues(t *testing.T) {
	type testStruct struct {
		Plain            string `testtag:"min=1,max=10"`
		Regex            string `testtag:"regex=^foo\\,bar$"`
		RegexUnknown     string `testtag:"regex=^\\d+\\.\\d+$"`
		RegexCount       string `testtag:"regexCount=a\\|b|2"`
		Replace          string `testtag:"replace=old\\,value|new\\|value"`
		RequiredIf       string `testtag:"requiredIf=Type:admin\\|user|Role"`
		NotContainsRegex string `testtag:"notContainsRegex=^test\\=value$"`
	}

	tagSettings := NewTagSettings(
		"testtag",
		",",
		"=",
		nil,
		false,
		NewKey("min", false, false, nil),
		NewKey("max", false, false, nil),
		NewKey("regex", false, false, nil),
		NewKey("regexCount", false, false, nil),
		NewKey("replace", false, false, nil),
		NewKey("requiredIf", false, false, nil),
		NewKey("notContainsRegex", false, false, nil),
	)
	tagSettings.WithEscapeCharacter(testEscapeCharacter)

	fields, err := tagSettings.ParseStruct(&testStruct{})
	testza.AssertNoError(t, err, "unexpected error")
	testza.AssertLen(t, fields, 7, "unexpected fields len")

	testza.AssertEqual(t, fields[0].Tags, []Tag{
		{Key: "min", Value: "1"},
		{Key: "max", Value: "10"},
	}, "unexpected plain tags")

	testza.AssertEqual(t, fields[1].FirstTag(), Tag{Key: "regex", Value: "^foo,bar$"},
		"unexpected regex tag")
	testza.AssertEqual(t, fields[2].FirstTag(), Tag{Key: "regex", Value: `^\d+\.\d+$`},
		"unexpected regex tag")
	testza.AssertEqual(t, fields[3].FirstTag(), Tag{Key: "regexCount", Value: `a\|b|2`},
		"unexpected regexCount tag")
	testza.AssertEqual(t, fields[4].FirstTag(), Tag{Key: "replace", Value: `old,value|new\|value`},
		"unexpected replace tag")
	testza.AssertEqual(t, fields[5].FirstTag(), Tag{Key: "requiredIf", Value: `Type:admin\|user|Role`},
		"unexpected requiredIf tag")
	testza.AssertEqual(t, fields[6].FirstTag(), Tag{Key: "notContainsRegex", Value: "^test=value$"},
		"unexpected notContainsRegex style tag")
}

func Test_ParseStructWithEscapeCharacterAndRegexEscapesOnly(t *testing.T) {
	type testStruct struct {
		Regex string `testtag:"regex=^\\d+\\.\\d+$"`
	}

	tagSettings := NewTagSettings(
		"testtag",
		",",
		"=",
		nil,
		false,
		NewKey("regex", false, false, nil),
	)
	tagSettings.WithEscapeCharacter(testEscapeCharacter)

	fields, err := tagSettings.ParseStruct(&testStruct{})
	testza.AssertNoError(t, err, "unexpected error")
	testza.AssertLen(t, fields, 1, "unexpected fields len")
	testza.AssertEqual(t, fields[0].FirstTag(), Tag{
		Key:   "regex",
		Value: `^\d+\.\d+$`,
	}, "unexpected regex tag")
}

func Test_ParseStructWithCustomSeparatorsAndEscapes(t *testing.T) {
	type testStruct struct {
		Value string `testtag:"replace@old\\#value#requiredIf@Type\\@admin"`
	}

	tagSettings := NewTagSettings(
		"testtag",
		"#",
		"@",
		nil,
		false,
		NewKey("replace", false, false, nil),
		NewKey("requiredIf", false, false, nil),
	)
	tagSettings.WithEscapeCharacter(testEscapeCharacter)

	fields, err := tagSettings.ParseStruct(&testStruct{})
	testza.AssertNoError(t, err, "unexpected error")
	testza.AssertLen(t, fields, 1, "unexpected fields len")
	testza.AssertEqual(t, fields[0].Tags, []Tag{
		{Key: "replace", Value: "old#value"},
		{Key: "requiredIf", Value: "Type@admin"},
	}, "unexpected parsed tags")
}
