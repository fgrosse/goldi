package matchers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type codeParsingMatcher struct {
	source       []byte
	astFile      *ast.File
	compileError error
}

func (m *codeParsingMatcher) parse(output *bytes.Buffer) (*ast.File, error) {
	m.source = output.Bytes()

	fileSet := token.NewFileSet()
	m.astFile, m.compileError = parser.ParseFile(fileSet, "test input", m.source, parser.ParseComments)
	return m.astFile, m.compileError
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
