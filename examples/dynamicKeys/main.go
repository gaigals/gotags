package main

import (
	"fmt"
	"log"

	"github.com/gaigals/gotags"
)

const tagName = "myTag"

// Use only value logic - `myTag:"myValue"`
var tagSettings = gotags.NewSettings(tagName).
	WithNoKeyExistValidation()

type MyData struct {
	Name    string `myTag:"myValue"`
	Age     uint
	Country string
	Postal  string
}

// Will return error on TagSettings.ParseStruct() if fails.
func tagProcessor(field gotags.Field) error {
	// Do some custom stuff for each field if required.

	// value := field.Value.Interface()
	// rules := field.TagDataFormatted("%s=%s")
	// errs := validator.Field(value, field.TagDataFormatted("%s=%s"))
	// ...
	return nil
}

func main() {
	myData := MyData{
		Name:    "John",
		Age:     22,
		Country: "Latvia",
	}

	// tagSettings.IncludeNotTagged = true

	// Parses all tags, triggers field processor if defined and validators
	// if defined.
	fields, err := tagSettings.ParseStruct(&myData)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(fields)

	// Do some additional stuff with fields if required.
	// ...
}
