package gotags

import (
	"testing"
)

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
