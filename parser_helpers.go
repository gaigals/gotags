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
		if input[index] == escapeCharacter {
			if index+1 >= len(input) {
				return nil, fmt.Errorf("%w in %q", errTrailingBackslash, input)
			}

			if isEscapedReservedCharacter(input[index+1], escapeCharacter) {
				index += 2
				continue
			}

			index++
			continue
		}

		if strings.HasPrefix(input[index:], separator) {
			part, err := unescapeReservedCharacters(
				input[startIndex:index],
				escapeCharacter,
			)
			if err != nil {
				return nil, err
			}

			parts = append(parts, part)
			index += len(separator)
			startIndex = index
			continue
		}

		index++
	}

	part, err := unescapeReservedCharacters(input[startIndex:], escapeCharacter)
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
		if input[index] == escapeCharacter {
			if index+1 >= len(input) {
				return nil, fmt.Errorf("%w in %q", errTrailingBackslash, input)
			}

			if isEscapedReservedCharacter(input[index+1], escapeCharacter) {
				index += 2
				continue
			}

			index++
			continue
		}

		if strings.HasPrefix(input[index:], separator) {
			before, err := unescapeReservedCharacters(input[:index], escapeCharacter)
			if err != nil {
				return nil, err
			}

			after, err := unescapeReservedCharacters(
				input[index+len(separator):],
				escapeCharacter,
			)
			if err != nil {
				return nil, err
			}

			return []string{before, after}, nil
		}

		index++
	}

	unescaped, err := unescapeReservedCharacters(input, escapeCharacter)
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
	if separator == "" || containsEscapeCharacter(input, escapeCharacter) {
		parts, err := splitFirstWithOptionalEscapes(
			input,
			separator,
			escapeCharacter,
		)
		if err != nil {
			return "", "", false, err
		}

		if len(parts) == 0 {
			return "", "", false, nil
		}

		if len(parts) == 1 {
			return parts[0], "", false, nil
		}

		return parts[0], parts[1], true, nil
	}

	index := strings.Index(input, separator)
	if index < 0 {
		return input, "", false, nil
	}

	return input[:index], input[index+len(separator):], true, nil
}

func unescapeReservedCharacters(
	input string,
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

		if isEscapedReservedCharacter(input[index+1], escapeCharacter) {
			builder.WriteByte(input[index+1])
			index += 2
			continue
		}

		builder.WriteByte(input[index])
		index++
	}

	return builder.String(), nil
}

func isEscapedReservedCharacter(char, escapeCharacter byte) bool {
	switch char {
	case escapeCharacter, ',', '|', ':', '=':
		return true
	default:
		return false
	}
}

func containsEscapeCharacter(input string, escapeCharacter byte) bool {
	if escapeCharacter == 0 {
		return false
	}

	return strings.IndexByte(input, escapeCharacter) >= 0
}
