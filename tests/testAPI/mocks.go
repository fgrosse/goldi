package testAPI

type MockType struct {
	StringParameter string
	BoolParameter   bool
}

func NewMockType() *MockType {
	return &MockType{}
}

func NewMockTypeWithArgs(stringParameter string, boolParameter bool) *MockType {
	return &MockType{stringParameter, boolParameter}
}

type MockTypeFactory struct {
	HasBeenUsed bool
	NrOfCalls   int
}

func (g *MockTypeFactory) NewMockType() *MockType {
	g.HasBeenUsed = true
	g.NrOfCalls++
	return &MockType{}
}
