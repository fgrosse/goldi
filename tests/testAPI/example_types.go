package testAPI

type Foo struct{}
type Bar struct{}
type Baz struct{}

func NewFoo() *Foo {
	return &Foo{}
}

func NewBar() *Bar {
	return &Bar{}
}

func NewBaz() *Baz {
	return &Baz{}
}
