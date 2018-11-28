package matchers

import (
	"fmt"

	"github.com/onsi/gomega/types"
)

// DeclarePackage succeeds if actual can be converted or read into to byte[], is valid go code and
// declares the given package name.
func DeclarePackage(packageName string) types.GomegaMatcher {
	return &declarePackageMatcher{expectedPackageName: packageName}
}

type declarePackageMatcher struct {
	beValidGoCodeMatcher
	expectedPackageName string
}

func (m *declarePackageMatcher) Match(actual interface{}) (success bool, err error) {
	isCompilable, err := m.beValidGoCodeMatcher.Match(actual)
	if isCompilable == false || err != nil {
		return isCompilable, err
	}

	return m.astFile.Name.Name == m.expectedPackageName, nil
}

func (m *declarePackageMatcher) FailureMessage(actual interface{}) (message string) {
	if m.outIsEmpty || m.compileError != nil {
		return m.beValidGoCodeMatcher.FailureMessage(actual)
	}

	return fmt.Sprintf("Expected output:\n%s\nto be in package %q", m.indentSource(), m.expectedPackageName)
}

func (m *declarePackageMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected output:\n%s\nnot to to be in package %q", m.indentSource(), m.expectedPackageName)
}
