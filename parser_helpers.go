package gotags

import (
	"errors"
	"fmt"
	"strings"
)

var errTrailingBackslash = errors.New("trailing naked backslash")

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
		tokenLength, err := nextEscapedTokenLength(
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

		part, err := unescapeReservedCharacters(
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

	part, err := unescapeReservedCharacters(
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
		tokenLength, err := nextEscapedTokenLength(
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

	unescaped, err := unescapeReservedCharacters(
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
	if separator != "" && !containsEscapeCharacter(input, escapeCharacter) {
		index := strings.Index(input, separator)
		if index < 0 {
			return input, "", false, nil
		}

		return input[:index], input[index+len(separator):], true, nil
	}

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

func unescapeReservedCharacters(
	input string,
	separator string,
	escapeCharacter byte,
) (string, error) {
	if !containsEscapeCharacter(input, escapeCharacter) {
		return input, nil
	}

	var builder strings.Builder
	builder.Grow(len(input))

	for index := 0; index < len(input); {
		if input[index] != escapeCharacter {
			builder.WriteByte(input[index])
			index++
			continue
		}

		if index+1 >= len(input) {
			return "", fmt.Errorf("%w in %q", errTrailingBackslash, input)
		}

		tokenLength := escapedTokenLength(
			input[index+1:],
			separator,
			escapeCharacter,
		)
		if tokenLength > 0 {
			builder.WriteString(input[index+1 : index+1+tokenLength])
			index += 1 + tokenLength
			continue
		}

		builder.WriteByte(input[index])
		index++
	}

	return builder.String(), nil
}

func escapedTokenLength(
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

	switch input[0] {
	case escapeCharacter, ',', '|', ':', '=':
		return 1
	default:
		return 0
	}
}

func nextEscapedTokenLength(
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

	return escapedTokenLength(input[index+1:], separator, escapeCharacter), nil
}

func unescapeSplitParts(
	input string,
	index int,
	separator string,
	escapeCharacter byte,
) ([]string, error) {
	before, err := unescapeReservedCharacters(
		input[:index],
		separator,
		escapeCharacter,
	)
	if err != nil {
		return nil, err
	}

	after, err := unescapeReservedCharacters(
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
