# GOTags

Parse struct tags into `[]Field`.

## Install

```bash
go get github.com/gaigals/gotags@latest
```

## Example

Default syntax uses `;` between tags and `:` between key/value.

```go
package main

import (
	"fmt"
	"log"

	"github.com/gaigals/gotags"
)

const (
	tagKeyRequired = "required"
	tagKeyEquals   = "eq"
	tagKeyGT       = "gt"
	tagKeyLT       = "lt"
)

var (
	keyRequired = gotags.NewKey(tagKeyRequired, true, false, nil)
	keyEquals   = gotags.NewKey(tagKeyEquals, false, false, nil)
	keyGT       = gotags.NewKey(tagKeyGT, false, false, nil)
	keyLT       = gotags.NewKey(tagKeyLT, false, false, nil)
)

// Default separator (;), default equals (:).
var settings = gotags.NewTagSettingsDefault(
	"validator",
	tagProcessor, // Optional - can be nil.
	keyRequired,
	keyEquals,
	keyGT,
	keyLT,
)

type User struct {
	Name    string `validator:"required"`
	Age     uint   `validator:"gt:10;lt:130"`
	Country string `validator:"eq:Latvia"`
}

func tagProcessor(field gotags.Field) error {
	// Optional: runs after parsing and validation.
	// Use it for extra checks, transforms, or side effects.
	return nil
}

func main() {
	user := User{
		Name:    "John",
		Age:     22,
		Country: "Latvia",
	}

	fields, err := settings.ParseStruct(&user)
	if err != nil {
		log.Fatalln(err)
	}

	for _, field := range fields {
		fmt.Println(field.Name)
		for _, tag := range field.Tags {
			fmt.Printf("  %s = %q\n", tag.Key, tag.Value)
		}
	}

	// Name
	//   required = ""
	// Age
	//   gt = "10"
	//   lt = "130"
	// Country
	//   eq = "Latvia"
}
```


Parse with `ParseStruct(&value)`.

## Setup

```go
gotags.NewSettings("validator")
gotags.NewSettings("validator").WithEscapeCharacter('\\')
gotags.NewTagSettingsDefault("validator", nil, keys...)
gotags.NewTagSettings("validator", ",", "=", nil, false, keys...)
```

- `WithProcessor(fn)` runs after validation.
- `IncludeUntaggedFields()` keeps exported fields without the tag.
- `WithNoKeyExistValidation()` accepts dynamic tags.
- `WithEscapeCharacter('\\')` enables escape parsing.

## Custom Separators

```go
var settings = gotags.NewSettings("gotags").
	WithCustomSeparators(",", "=").
	AddKeys(
		gotags.NewKey("required", true, false, nil),
		gotags.NewKey("min", false, false, nil),
		gotags.NewKey("max", false, false, nil),
	)

type User struct {
	Name string `gotags:"required,min=2,max=255"`
}
```

## Dynamic Tags

```go
var settings = gotags.NewSettings("myTag").
	WithNoKeyExistValidation()

type Item struct {
	Value string `myTag:"rawValue"`
}

// Tag{Key: "rawValue", Value: ""}
```

## Direct Tag Parsing

```go
tag, err := gotags.NewTagFromString("min:2", ":")

escapedTag, err := gotags.NewTagFromStringWithEscape(
	`replace=old\,value|new\|value`,
	"=",
	'\\',
)
```

## Escaping

Escaping is off by default.\
Enable it per `TagSettings` only when values need parser syntax chars.
It also works with custom separators and custom equals chars.

```go
var validatorSettings = gotags.NewSettings("validator").
	WithEscapeCharacter('\\')

var gotagsSettings = gotags.NewSettings("gotags").
	WithCustomSeparators(",", "=").
	WithEscapeCharacter('\\')

type Rules struct {
	Regex      string `validator:"regex:^foo\,bar$"`
	RegexDots  string `validator:"regex:^\d+\.\d+$"`
	Replace    string `gotags:"replace=old\,value|new\|value"`
	RequiredIf string `gotags:"requiredIf=Type:admin\|user|Role"`
}
```

- `\\` => `\`
- `\,` => `,`
- `\|` => `|`
- `\:` => `:`
- `\=` => `=`
- unknown escapes stay as-is: `\d`, `\w`, `\.`
- trailing naked `\` returns an error

## Useful Field Helpers

```go
// Check whether a parsed field contains a specific tag key.
field.HasKey("required")

// Get a tag value by key. Missing keys return an empty string.
field.KeyValue("gt")

// Get both the tag value and whether that key was found.
value, ok := field.KeyValueBool("gt")

// Read the first parsed tag on the field.
tag := field.FirstTag()

// Update the struct field value through reflection.
err := field.SetValue("new value")
```
