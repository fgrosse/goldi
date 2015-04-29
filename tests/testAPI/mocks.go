package testAPI

type MockType struct{}

func NewMockType() *MockType {
	return &MockType{}
}

func NewMockTypeWithArgs(_ string, _ bool) *MockType {
	return NewMockType()
}

type MockTypeGenerator struct {
	HasBeenUsed bool
}

func (g *MockTypeGenerator) NewMockType() *MockType {
	g.HasBeenUsed = true
	return &MockType{}
}
