package matchers

import (
	"bytes"
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func BeValidGoCode() types.GomegaMatcher {
	return &BeValidGoCodeMatcher{}
}

type BeValidGoCodeMatcher struct {
	codeParsingMatcher
	outIsEmpty bool
}

func (m *BeValidGoCodeMatcher) Match(actual interface{}) (success bool, err error) {
	var output *bytes.Buffer
	switch a := actual.(type) {
	case *bytes.Buffer:
		output = a
	default:
		return false, fmt.Errorf("BeValidGoCode expects actual to be an instance of *bytes.Buffer")
	}

	if output.Len() == 0 {
		m.outIsEmpty = true
		return false, nil
	}

	_, m.compileError = m.parse(output)
	return m.compileError == nil, nil
}

func (m *BeValidGoCodeMatcher) FailureMessage(_ interface{}) (message string) {
	if m.outIsEmpty {
		return "Expected output to contain compilable go code but it was empty"
	}

	return fmt.Sprintf("Expected output:\n%s\nto contain compilable go code but got parser error:\n    %s", indent(m.source), m.compileError)
}

func (m *BeValidGoCodeMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return format.Message(fmt.Sprintf("%s", m.source), "not to contain compilable go code")
}
