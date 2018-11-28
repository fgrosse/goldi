package matchers

import (
	"fmt"

	"github.com/onsi/gomega/types"
)

// ImportPackage succeeds if actual can be converted or read into to byte[], is valid go code and
// contains an import statement for the given package
func ImportPackage(expected string) types.GomegaMatcher {
	return &importPackageMatcher{expectedPackage: fmt.Sprintf("%q", expected)}
}

type importPackageMatcher struct {
	beValidGoCodeMatcher
	expectedPackage   string
	foundMoreThanOnce bool
}

func (m *importPackageMatcher) Match(actual interface{}) (success bool, err error) {
	isCompilable, err := m.beValidGoCodeMatcher.Match(actual)
	if isCompilable == false || err != nil {
		return isCompilable, err
	}

	m.foundMoreThanOnce = false
	importFound := false
	for _, importSpec := range m.astFile.Imports {
		if importSpec.Path.Value == m.expectedPackage {
			if importFound {
				m.foundMoreThanOnce = true
				return false, nil
			}

			importFound = true
		}
	}

	return importFound, nil
}

func (m *importPackageMatcher) FailureMessage(actual interface{}) (message string) {
	if m.outIsEmpty || m.compileError != nil {
		return m.beValidGoCodeMatcher.FailureMessage(actual)
	}

	if m.foundMoreThanOnce {
		return fmt.Sprintf("Expected output:\n%s\nto import package %s exactly once but found multiple times", m.indentSource(), m.expectedPackage)
	}

	return fmt.Sprintf("Expected output:\n%s\nto import package %s", m.indentSource(), m.expectedPackage)
}

func (m *importPackageMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected output:\n%s\nnot to import package %s", m.indentSource(), m.expectedPackage)
}
