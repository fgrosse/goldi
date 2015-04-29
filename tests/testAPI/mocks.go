package testAPI

type MockType struct{}

func NewMockType() *MockType {
	return &MockType{}
}

type MockTypeGenerator struct {
	HasBeenUsed bool
}

func (g *MockTypeGenerator) NewMockType() *MockType {
	g.HasBeenUsed = true
	return &MockType{}
}
