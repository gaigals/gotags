package gotags

import (
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
