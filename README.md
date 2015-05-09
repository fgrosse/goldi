Goldi
========

[![Build Status](https://secure.travis-ci.org/fgrosse/goldi.png?branch=master)](http://travis-ci.org/fgrosse/goldi)
[![GoDoc](https://godoc.org/github.com/fgrosse/goldi?status.svg)](https://godoc.org/github.com/fgrosse/goldi)

A go dependency injection framework.

**Note: This library is at the very early stages of its development.**

## Installation

Use `go get` to install goldi:
```
go get github.com/fgrosse/goldi
```

No additional dependencies are required to use the library.
If you want to run the tests you need [ginkgo][1] and [gomega][2]

## Documentation

A generated documentation is available at [godoc.org][3]

## Usage

For usage examples have a look at the [functional tests](tests).
Its probably best to start by looking at the [container_test.go](tests/container_test.go).

## Running the tests

Goldi uses the awesome [ginkgo][1] framework for its tests.
You can execute the tests running:
```
ginkgo tests
```

If you prefer to use `go test` directly you can either switch into the `./tests` directory and run it there or
run the following from the repository root directory:
```
go test ./tests
```

## Contributing

Any contributions are always welcome (use pull requests).
Please keep in mind that I might not always be able to respond immediately but I usually try to react within the week ☺.

[1]: http://onsi.github.io/ginkgo/
[2]: http://onsi.github.io/gomega/
[3]: http://godoc.org/github.com/fgrosse/goldi
