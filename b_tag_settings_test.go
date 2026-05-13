package gotags

import (
	"testing"
)

type benchmarkParseStructNoEscape struct {
	Name    string `gotags:"required,min=2,max=255"`
	Age     int    `gotags:"required,min=10,max=130"`
	Country string `gotags:"required,country"`
	Phone   string `gotags:"required,min=5,max=50,phone"`
}

type benchmarkParseStructWithEscape struct {
	Regex            string `gotags:"regex=^foo\\,bar$"`
	RegexCount       string `gotags:"regexCount=a\\|b|2"`
	Replace          string `gotags:"replace=old\\,value|new\\|value"`
	NotContainsRegex string `gotags:"notContainsRegex=^test\\=value$"`
}

var benchmarkParseStructValueNoEscape = benchmarkParseStructNoEscape{
	Name:    "Janis",
	Age:     10,
	Country: "Latvia",
	Phone:   "+37122131231",
}

var benchmarkParseStructValueWithEscape = benchmarkParseStructWithEscape{
	Regex:            "^foo,bar$",
	RegexCount:       "a|b|2",
	Replace:          "old,value|new|value",
	NotContainsRegex: "^test=value$",
}

var benchmarkParseStructSettingsNoEscape = NewTagSettings(
	"gotags",
	",",
	"=",
	nil,
	false,
	NewKey("required", true, false, nil),
	NewKey("min", false, false, nil),
	NewKey("max", false, false, nil),
	NewKey("country", true, false, nil),
	NewKey("phone", true, false, nil),
)

var benchmarkParseStructSettingsWithEscape = func() TagSettings {
	tagSettings := NewTagSettings(
		"gotags",
		",",
		"=",
		nil,
		false,
		NewKey("regex", false, false, nil),
		NewKey("regexCount", false, false, nil),
		NewKey("replace", false, false, nil),
		NewKey("notContainsRegex", false, false, nil),
	)
	tagSettings.WithEscapeCharacter('\\')
	return tagSettings
}()

const (
	benchmarkTagStringNoEscape   = "replace=oldValue|newValue"
	benchmarkTagStringWithEscape = `replace=old\,value|new\|value`
)

// PAST benchmarks:
//
// Benchmark_ParseStruct/NoEscape-4     821473   1456 ns/op   672 B/op
// 5 allocs/op
// Benchmark_ParseStruct/WithEscape-4   318123   3534 ns/op   848 B/op
// 19 allocs/op
func Benchmark_ParseStruct(b *testing.B) {
	b.Run("NoEscape", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = benchmarkParseStructSettingsNoEscape.ParseStruct(
				&benchmarkParseStructValueNoEscape,
			)
		}
	})

	b.Run("WithEscape", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = benchmarkParseStructSettingsWithEscape.ParseStruct(
				&benchmarkParseStructValueWithEscape,
			)
		}
	})
}

// PAST benchmarks:
//
// Benchmark_NewTagFromString/NoEscape-4    100000000    11.63 ns/op
// 0 B/op   0 allocs/op
// Benchmark_NewTagFromString/WithEscape-4    6461421    185.7 ns/op
// 24 B/op   1 allocs/op
func Benchmark_NewTagFromString(b *testing.B) {
	b.Run("NoEscape", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = NewTagFromString(benchmarkTagStringNoEscape, "=")
		}
	})

	b.Run("WithEscape", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = NewTagFromStringWithEscape(
				benchmarkTagStringWithEscape,
				"=",
				'\\',
			)
		}
	})
}
