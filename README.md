# GOTags

API for managing struct tags - defining and parsing.

### Example:

Create new tag `validator` with tags `required` (bool), `eq` (any), `gt`
(int) and `lt` (int) with default separator `;` and default equals `;`

```go
package main

import (
    "log"
    "reflect"

    "github.com/gaigals/gotags"
)

const (
    tagKeyRequired    = "required"
    tagKeyEquals      = "eq"
    tagKeyGreaterThan = "gt"
    tagKeyLessThan    = "lt"
)

var (
    keyRequired    = gotags.NewKey(tagKeyRequired, true, false, nil)
    keyEquals      = gotags.NewKey(tagKeyEquals, false, false, nil)
    keyGreaterThan = gotags.NewKey(tagKeyGreaterThan, false, false, nil)
    keyLessThan    = gotags.NewKey(tagKeyLessThan, false, false, nil)
)

// Default separator (;), default equals (:)
var tagSettings = gotags.NewTagSettingsDefault(
    "validator",
    tagProcessor, // Optional - can be nil
    keyRequired,
    keyEquals,
    keyGreaterThan,
    keyLessThan,
)

type MyData struct {
    Name    string `validator:"required"`
    Age     uint   `validator:"gt:10;lt:130"`
    Country string `validator:"eq:Latvia"`
}

// Will return error on TagSettings.ParseStruct() if fails.
func tagProcessor(field gotags.Field) error {
    // Do some custom stuff for each field if required.
    // ...
    return nil
}

func main() {
    myData := MyData{
        Name:    "John",
        Age:     22,
        Country: "Latvia",
    }

    // Parses all tags, triggers field processor if defined and validators
    // if defined.
    fields, err := tagSettings.ParseStruct(&myData)
    if err != nil {
        log.Fatalln(err)
    }

    // Do some additional stuff with fields if required.
    // ...
}

```
