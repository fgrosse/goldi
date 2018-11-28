package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// BeValidGoCode succeeds if actual can be converted or read into to byte[] and contains valid go code.
func BeValidGoCode() types.GomegaMatcher {
	return &beValidGoCodeMatcher{}
}

type beValidGoCodeMatcher struct {
	codeParsingMatcher
	outIsEmpty bool
}

func (m *beValidGoCodeMatcher) Match(actual interface{}) (success bool, err error) {
	_, m.outIsEmpty, m.compileError = m.parse(actual)
	return m.compileError == nil, nil
}

func (m *beValidGoCodeMatcher) FailureMessage(_ interface{}) (message string) {
	if m.outIsEmpty {
		return "Expected output to contain compilable go code but it was empty"
	}

	return fmt.Sprintf("Expected output:\n%s\nto contain compilable go code but got parser error:\n    %s", indent(m.source), m.compileError)
}

func (m *beValidGoCodeMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return format.Message(fmt.Sprintf("%s", m.source), "not to contain compilable go code")
}
