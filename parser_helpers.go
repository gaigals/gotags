package gotags

import (
	"errors"
	"fmt"
	"strings"
)

var errTrailingBackslash = errors.New("trailing naked backslash")

// SplitWithEscape splits by the current layer separator while respecting the
// configured escape character. It only unescapes the current separator and
// escaped backslashes, leaving deeper escapes for later parsing layers.
func SplitWithEscape(
	input,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	return splitWithOptionalEscapes(input, separator, escapeCharacter)
}

// SplitFirstWithEscape splits by the first current layer separator while
// respecting the configured escape character. It only unescapes the current
// separator and escaped backslashes, leaving deeper escapes for later parsing
// layers.
func SplitFirstWithEscape(
	input,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	return splitFirstWithOptionalEscapes(input, separator, escapeCharacter)
}

func splitWithOptionalEscapes(
	input,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	if separator == "" || !containsEscapeCharacter(input, escapeCharacter) {
		return strings.Split(input, separator), nil
	}

	return splitEscaped(input, separator, escapeCharacter)
}

func splitFirstWithOptionalEscapes(
	input,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	if separator == "" || !containsEscapeCharacter(input, escapeCharacter) {
		return strings.SplitN(input, separator, 2), nil
	}

	return splitFirstEscaped(input, separator, escapeCharacter)
}

func splitEscaped(
	input,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	startIndex := 0
	parts := make([]string, 0, strings.Count(input, separator)+1)

	for index := 0; index < len(input); {
		tokenLength, err := nextCurrentLayerEscapedTokenLength(
			input,
			index,
			separator,
			escapeCharacter,
		)
		if err != nil {
			return nil, err
		}
		if tokenLength > 0 {
			index += 1 + tokenLength
			continue
		}
		if input[index] == escapeCharacter {
			index++
			continue
		}

		if !strings.HasPrefix(input[index:], separator) {
			index++
			continue
		}

		part, err := unescapeCurrentLayerCharacters(
			input[startIndex:index],
			separator,
			escapeCharacter,
		)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
		index += len(separator)
		startIndex = index
	}

	part, err := unescapeCurrentLayerCharacters(
		input[startIndex:],
		separator,
		escapeCharacter,
	)
	if err != nil {
		return nil, err
	}

	return append(parts, part), nil
}

func splitFirstEscaped(
	input,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	for index := 0; index < len(input); {
		tokenLength, err := nextCurrentLayerEscapedTokenLength(
			input,
			index,
			separator,
			escapeCharacter,
		)
		if err != nil {
			return nil, err
		}
		if tokenLength > 0 {
			index += 1 + tokenLength
			continue
		}
		if input[index] == escapeCharacter {
			index++
			continue
		}

		if !strings.HasPrefix(input[index:], separator) {
			index++
			continue
		}

		return unescapeSplitParts(input, index, separator, escapeCharacter)
	}

	unescaped, err := unescapeCurrentLayerCharacters(
		input,
		separator,
		escapeCharacter,
	)
	if err != nil {
		return nil, err
	}

	return []string{unescaped}, nil
}

func splitTagKeyValue(input, separator string) (
	key string,
	value string,
	hasValue bool,
	err error,
) {
	return splitTagKeyValueWithEscape(input, separator, 0)
}

// splitTagKeyValueWithEscape keeps the old fast path for plain input, then
// falls back to an escape-aware scan only when the configured escape character
// is present. It first finds the real key/value separator and only then
// unescapes only the current layer (`\\` and the active equals token).
func splitTagKeyValueWithEscape(
	input,
	separator string,
	escapeCharacter byte,
) (
	key string,
	value string,
	hasValue bool,
	err error,
) {
	// Keep the plain-input path identical to the old parser.
	if separator != "" && !containsEscapeCharacter(input, escapeCharacter) {
		index := strings.Index(input, separator)
		if index < 0 {
			return input, "", false, nil
		}

		return input[:index], input[index+len(separator):], true, nil
	}

	// Preserve existing empty-separator behavior through the shared splitter.
	if separator == "" {
		parts, err := splitFirstWithOptionalEscapes(
			input,
			separator,
			escapeCharacter,
		)
		if err != nil {
			return "", "", false, err
		}

		key, value, hasValue = splitPartsToKeyValue(parts)
		return key, value, hasValue, nil
	}

	// Find the first real separator while skipping escaped characters.
	for index := 0; index < len(input); {
		tokenLength, err := nextCurrentLayerEscapedTokenLength(
			input,
			index,
			separator,
			escapeCharacter,
		)
		if err != nil {
			return "", "", false, err
		}
		if tokenLength > 0 {
			index += 1 + tokenLength
			continue
		}
		if input[index] == escapeCharacter {
			index++
			continue
		}

		if !strings.HasPrefix(input[index:], separator) {
			index++
			continue
		}

		key, err = unescapeCurrentLayerCharacters(
			input[:index],
			separator,
			escapeCharacter,
		)
		if err != nil {
			return "", "", false, err
		}

		value, err = unescapeCurrentLayerCharacters(
			input[index+len(separator):],
			separator,
			escapeCharacter,
		)
		if err != nil {
			return "", "", false, err
		}

		return key, value, true, nil
	}

	key, err = unescapeCurrentLayerCharacters(input, separator, escapeCharacter)
	if err != nil {
		return "", "", false, err
	}

	return key, "", false, nil
}

// unescapeCurrentLayerCharacters returns the original string unchanged unless
// it finds an escape that belongs to the current layer. This keeps deeper
// escapes intact until a later parser chooses to use them, while still
// rejecting a trailing naked backslash.
func unescapeCurrentLayerCharacters(
	input string,
	separator string,
	escapeCharacter byte,
) (string, error) {
	if !containsEscapeCharacter(input, escapeCharacter) {
		return input, nil
	}

	for index := 0; index < len(input); {
		if input[index] != escapeCharacter {
			index++
			continue
		}

		if index+1 >= len(input) {
			return "", fmt.Errorf("%w in %q", errTrailingBackslash, input)
		}

		tokenLength := currentLayerEscapedTokenLength(
			input[index+1:],
			separator,
			escapeCharacter,
		)
		if tokenLength == 0 {
			index++
			continue
		}

		// Start rebuilding only after the first supported escape is found.
		return unescapeCurrentLayerCharactersFromIndex(
			input,
			index,
			tokenLength,
			separator,
			escapeCharacter,
		)
	}

	return input, nil
}

// Once the first current-layer escape is found, rebuild the rest of the string
// and unwrap only the active separator/equals token and escaped backslashes.
func unescapeCurrentLayerCharactersFromIndex(
	input string,
	startIndex int,
	tokenLength int,
	separator string,
	escapeCharacter byte,
) (string, error) {
	var builder strings.Builder
	builder.Grow(len(input))
	builder.WriteString(input[:startIndex])
	builder.WriteString(input[startIndex+1 : startIndex+1+tokenLength])

	for index := startIndex + 1 + tokenLength; index < len(input); {
		if input[index] != escapeCharacter {
			builder.WriteByte(input[index])
			index++
			continue
		}

		if index+1 >= len(input) {
			return "", fmt.Errorf("%w in %q", errTrailingBackslash, input)
		}

		tokenLength = currentLayerEscapedTokenLength(
			input[index+1:],
			separator,
			escapeCharacter,
		)
		if tokenLength == 0 {
			builder.WriteByte(input[index])
			index++
			continue
		}

		builder.WriteString(input[index+1 : index+1+tokenLength])
		index += 1 + tokenLength
	}

	return builder.String(), nil
}

func currentLayerEscapedTokenLength(
	input,
	separator string,
	escapeCharacter byte,
) int {
	if input == "" {
		return 0
	}

	if separator != "" && strings.HasPrefix(input, separator) {
		return len(separator)
	}

	if input[0] == escapeCharacter {
		return 1
	}

	return 0
}

func nextCurrentLayerEscapedTokenLength(
	input string,
	index int,
	separator string,
	escapeCharacter byte,
) (int, error) {
	if input[index] != escapeCharacter {
		return 0, nil
	}

	if index+1 >= len(input) {
		return 0, fmt.Errorf("%w in %q", errTrailingBackslash, input)
	}

	return currentLayerEscapedTokenLength(
		input[index+1:],
		separator,
		escapeCharacter,
	), nil
}

func unescapeSplitParts(
	input string,
	index int,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	before, err := unescapeCurrentLayerCharacters(
		input[:index],
		separator,
		escapeCharacter,
	)
	if err != nil {
		return nil, err
	}

	after, err := unescapeCurrentLayerCharacters(
		input[index+len(separator):],
		separator,
		escapeCharacter,
	)
	if err != nil {
		return nil, err
	}

	return []string{before, after}, nil
}

func splitPartsToKeyValue(parts []string) (
	key string,
	value string,
	hasValue bool,
) {
	if len(parts) == 0 {
		return "", "", false
	}

	if len(parts) == 1 {
		return parts[0], "", false
	}

	return parts[0], parts[1], true
}

func containsEscapeCharacter(input string, escapeCharacter byte) bool {
	if escapeCharacter == 0 {
		return false
	}

	return strings.IndexByte(input, escapeCharacter) >= 0
}
