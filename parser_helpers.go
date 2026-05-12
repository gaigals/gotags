package gotags

import (
	"errors"
	"fmt"
	"strings"
)

const escapeCharacter = '\\'

var errTrailingBackslash = errors.New("trailing naked backslash")

func splitWithOptionalEscapes(input, separator string) ([]string, error) {
	if separator == "" || !strings.ContainsRune(input, escapeCharacter) {
		return strings.Split(input, separator), nil
	}

	return splitEscaped(input, separator)
}

func splitFirstWithOptionalEscapes(input, separator string) ([]string, error) {
	if separator == "" || !strings.ContainsRune(input, escapeCharacter) {
		return strings.SplitN(input, separator, 2), nil
	}

	return splitFirstEscaped(input, separator)
}

func splitEscaped(input, separator string) ([]string, error) {
	startIndex := 0
	parts := make([]string, 0, strings.Count(input, separator)+1)

	for index := 0; index < len(input); {
		if input[index] == escapeCharacter {
			if index+1 >= len(input) {
				return nil, fmt.Errorf("%w in %q", errTrailingBackslash, input)
			}

			if isEscapedReservedCharacter(input[index+1]) {
				index += 2
				continue
			}

			index++
			continue
		}

		if strings.HasPrefix(input[index:], separator) {
			part, err := unescapeReservedCharacters(input[startIndex:index])
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

	part, err := unescapeReservedCharacters(input[startIndex:])
	if err != nil {
		return nil, err
	}

	return append(parts, part), nil
}

func splitFirstEscaped(input, separator string) ([]string, error) {
	for index := 0; index < len(input); {
		if input[index] == escapeCharacter {
			if index+1 >= len(input) {
				return nil, fmt.Errorf("%w in %q", errTrailingBackslash, input)
			}

			if isEscapedReservedCharacter(input[index+1]) {
				index += 2
				continue
			}

			index++
			continue
		}

		if strings.HasPrefix(input[index:], separator) {
			before, err := unescapeReservedCharacters(input[:index])
			if err != nil {
				return nil, err
			}

			after, err := unescapeReservedCharacters(input[index+len(separator):])
			if err != nil {
				return nil, err
			}

			return []string{before, after}, nil
		}

		index++
	}

	unescaped, err := unescapeReservedCharacters(input)
	if err != nil {
		return nil, err
	}

	return []string{unescaped}, nil
}

func unescapeReservedCharacters(input string) (string, error) {
	if !strings.ContainsRune(input, escapeCharacter) {
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

		if isEscapedReservedCharacter(input[index+1]) {
			builder.WriteByte(input[index+1])
			index += 2
			continue
		}

		builder.WriteByte(input[index])
		index++
	}

	return builder.String(), nil
}

func isEscapedReservedCharacter(char byte) bool {
	switch char {
	case escapeCharacter, ',', '|', ':', '=':
		return true
	default:
		return false
	}
}
