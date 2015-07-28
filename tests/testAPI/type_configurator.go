package testAPI
import "fmt"

type MockTypeConfigurator struct {
	StringParameter string
}

func NewMockTypeConfigurator(configuredStringParam string) *MockTypeConfigurator {
	return &MockTypeConfigurator{StringParameter: configuredStringParam}
}

func (c *MockTypeConfigurator) Configure(m *MockType) {
	m.StringParameter = c.StringParameter
}

type FailingMockTypeConfigurator struct {
	StringParameter string
}

func NewFailingMockTypeConfigurator() *FailingMockTypeConfigurator {
	return &FailingMockTypeConfigurator{}
}

func (c *FailingMockTypeConfigurator) Configure(m *MockType) error {
	return fmt.Errorf("this is the error message from the testAPI.MockTypeConfigurator")
}
