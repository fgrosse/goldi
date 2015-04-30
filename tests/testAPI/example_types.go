package testAPI

type Foo struct{}
type Bar struct{}

type Baz struct {
	parameter1, parameter2 string
}

func NewFoo() *Foo {
	return &Foo{}
}

func NewBar() *Bar {
	return &Bar{}
}

func NewBaz(parameter1, parameter2 string) *Baz {
	return &Baz{parameter1, parameter2}
}
