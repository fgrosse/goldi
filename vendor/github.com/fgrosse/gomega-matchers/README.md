[![GoDoc](https://img.shields.io/badge/gowalker-doc-blue.svg)][3]

A set of custom [gomega][1] matchers to test generated go code.

Example of included new matchers:

```go
myCode := `
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, 世界")
	}
`

Expect(myCode).To(BeValidGoCode())
Expect(myCode).To(DeclarePackage("main"))
Expect(myCode).To(ImportPackage("fmt"))
Expect(myCode).To(ContainCode(`fmt.Println("Hello, 世界")`))
Expect("2006-01-02T15:04").To(EqualTime(time.Now()))
```

These matchers are actively used in the [goldi][2] test suit.

[1]: http://onsi.github.io/gomega/
[2]: http://github.com/fgrosse/goldi/
[3]: https://gowalker.org/github.com/fgrosse/gomega-matchers
