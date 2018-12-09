package matchers

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

type codeParsingMatcher struct {
	source       []byte
	astFile      *ast.File
	compileError error
}

func (m *codeParsingMatcher) parse(src interface{}) (f *ast.File, isEmpty bool, err error) {
	m.source, err = readSource(src)
	if err != nil {
		return nil, false, err
	}

	if len(m.source) == 0 {
		return nil, true, nil
	}

	fileSet := token.NewFileSet()
	m.astFile, m.compileError = parser.ParseFile(fileSet, "test input", m.source, parser.ParseComments)
	return m.astFile, false, m.compileError
}

func readSource(src interface{}) ([]byte, error) {
	if src == nil {
		return nil, errors.New("input is nil")
	}

	switch s := src.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	case *bytes.Buffer:
		// is io.Reader, but src is already available in []byte form
		if s != nil {
			return s.Bytes(), nil
		}
	case io.Reader:
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, s); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return nil, errors.New("invalid source")
}

func (m *codeParsingMatcher) indentSource() string {
	return indent(m.source)
}

func indent(originalSource []byte) string {
	intendedSource := &bytes.Buffer{}
	lineCounter := 1
	fmt.Fprintf(intendedSource, "    %3d:  ", lineCounter)
	for _, c := range originalSource {
		if c == '\n' {
			lineCounter++
			fmt.Fprintf(intendedSource, "\n    %3d:  ", lineCounter)
		} else {
			intendedSource.WriteByte(c)
		}
	}
	return intendedSource.String()
}

func Dump(output *bytes.Buffer) {
	fmt.Printf("\n%s\n", indent(output.Bytes()))
}
