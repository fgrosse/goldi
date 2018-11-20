package main

import "bytes"

type sanitizer struct {
	buf writer

	inQuotes  bool
	quoteChar byte
}

type writer interface {
	WriteString(string) (int, error)
	WriteByte(c byte) error
	Bytes() []byte
}

func newSanitizer() *sanitizer {
	return &sanitizer{
		buf:      &bytes.Buffer{},
		inQuotes: false,
	}
}

// Write escapes all @ signs that are not inside of a quoted string
func (s *sanitizer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		switch {
		case b == '@' && !s.inQuotes:
			var m int
			if m, err = s.buf.WriteString(`\@`); err != nil {
				return n, err
			}
			n = n + m
		case b == '\'':
			fallthrough
		case b == '"':
			if s.inQuotes == false {
				s.inQuotes = true
				s.quoteChar = b
			} else if s.quoteChar == b {
				s.inQuotes = false
			}
			fallthrough
		default:
			if err = s.buf.WriteByte(b); err != nil {
				return n, err
			}
			n = n + 1
		}
	}

	return n, nil
}

func (s *sanitizer) Bytes() []byte {
	return s.buf.Bytes()
}
