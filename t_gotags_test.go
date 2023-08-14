package gotags

import (
	"reflect"
	"slices"
	"testing"
)

func TestNewKey(t *testing.T) {
	// Valid Test Cases
	t.Run("Valid: Bool with Validator", func(t *testing.T) {
		key := NewKey("boolKey", true, true, nil, reflect.Bool)
		if key.Name != "boolKey" || !key.IsBool || !key.IsRequired || len(key.AllowedKinds) != 1 || key.AllowedKinds[0] != reflect.Bool || key.Validator != nil {
			t.Fatalf("Expected valid key, but got %+v", key)
		}
	})

	t.Run("Valid: Non-Bool with Validator", func(t *testing.T) {
		mockValidator := func(value string) error { return nil }
		key := NewKey("stringKey", false, false, mockValidator)
		if key.Name != "stringKey" || key.IsBool || key.IsRequired || len(key.AllowedKinds) != 0 || key.Validator == nil {
			t.Fatalf("Expected valid key, but got %+v", key)
		}
	})

	// Invalid Test Cases
	//t.Run("Invalid: Empty Name", func(t *testing.T) {
	//	key := NewKey("", true, true, nil, reflect.Bool)
	//	if key.Name == "" {
	//		t.Fatalf("Expected invalid key with empty name, but got %+v", key)
	//	}
	//})
}

func TestNewTagSettings(t *testing.T) {
	// Valid Test Cases
	t.Run("Valid: Basic Settings", func(t *testing.T) {
		tagSettings := NewTagSettings("settings", ",", "=", nil, true)
		if tagSettings.Name != "settings" || tagSettings.Separator != "," || tagSettings.Equals != "=" ||
			tagSettings.Processor != nil || !tagSettings.IncludeNotTagged || tagSettings.disableKeyValidation ||
			len(tagSettings.Keys) != 0 || len(tagSettings.keysRequired) != 0 {
			t.Fatalf("Expected valid tag settings, but got %+v", tagSettings)
		}
	})

	t.Run("Valid: With Keys", func(t *testing.T) {
		mockValidator := func(value string) error { return nil }
		key1 := NewKey("key1", false, true, mockValidator, reflect.String)
		key2 := NewKey("key2", true, true, nil, reflect.Bool)
		key3 := NewKey("key3", true, false, nil, reflect.Int)

		tagSettings := NewTagSettings("settings", ",", "=", nil, true, key1, key2, key3)
		if tagSettings.Name != "settings" || tagSettings.Separator != "," || tagSettings.Equals != "=" ||
			tagSettings.Processor != nil || !tagSettings.IncludeNotTagged || tagSettings.disableKeyValidation ||
			len(tagSettings.Keys) != 3 || len(tagSettings.keysRequired) != 2 ||
			!(slices.Contains(tagSettings.keysRequired, "key1") && slices.Contains(tagSettings.keysRequired, "key2")) {
			t.Fatalf("Expected valid tag settings with keys, but got %+v", tagSettings)
		}
	})

	// Invalid Test Cases
	//t.Run("Invalid: Empty Name", func(t *testing.T) {
	//	tagSettings := NewTagSettings("", ",", "=", nil, true)
	//	if tagSettings.Name == "" {
	//		t.Fatalf("Expected invalid tag settings with empty name, but got %+v", tagSettings)
	//	}
	//})

	//t.Run("Invalid: Empty Separator", func(t *testing.T) {
	//	tagSettings := NewTagSettings("settings", "", "=", nil, true)
	//	if tagSettings.Separator == "" {
	//		t.Fatalf("Expected invalid tag settings with empty separator, but got %+v", tagSettings)
	//	}
	//})

	//t.Run("Invalid: Empty Equals", func(t *testing.T) {
	//	tagSettings := NewTagSettings("settings", ",", "", nil, true)
	//	if tagSettings.Equals == "" {
	//		t.Fatalf("Expected invalid tag settings with empty equals, but got %+v", tagSettings)
	//	}
	//})
}

func TestNewTagSettingsDefault(t *testing.T) {
	// Valid Test Cases
	t.Run("Valid: With Processor and Keys", func(t *testing.T) {
		mockProcessor := func(field Field) error { return nil }
		tagSettings := NewTagSettingsDefault("settings", mockProcessor)
		if tagSettings.Name != "settings" || tagSettings.Separator != defaultSeparator ||
			tagSettings.Equals != defaultEquals || tagSettings.Processor == nil ||
			tagSettings.IncludeNotTagged || tagSettings.disableKeyValidation {
			t.Fatalf("Expected valid tag settings with default values, but got %+v", tagSettings)
		}
	})
}
