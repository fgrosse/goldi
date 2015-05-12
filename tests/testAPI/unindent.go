package testAPI

import "bytes"

func Unindent(input string) string {
	indent := getFirstIndent(input)
	result := &bytes.Buffer{}
	indentIndex := 0
	for _, c := range input {
		if c == '\n' {
			indentIndex = 0
		} else if indentIndex < len(indent) {
			indentIndex++
			continue
		}

		result.WriteRune(c)
	}

	return result.String()
}

func getFirstIndent(input string) string {
	// search for the first real indent
	indent := &bytes.Buffer{}
	for _, c := range input {
		switch c {
		case ' ':
			fallthrough
		case '\t':
			indent.WriteRune(c)
		case '\n':
			indent.Reset()
		default:
			return indent.String()
		}
	}

	return indent.String()
}
