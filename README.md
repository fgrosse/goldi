![Goldi](http://fgrosse.github.io/goldi/images/goldi_wide.png)

[![Build Status](https://secure.travis-ci.org/fgrosse/goldi.png?branch=master)](http://travis-ci.org/fgrosse/goldi)
[![Coverage Status](https://coveralls.io/repos/fgrosse/goldi/badge.svg?branch=master)](https://coveralls.io/r/fgrosse/goldi?branch=master)
[![GoDoc](https://godoc.org/github.com/fgrosse/goldi?status.svg)](https://godoc.org/github.com/fgrosse/goldi)
[![license](https://img.shields.io/badge/license-MIT-4183c4.svg)](https://github.com/fgrosse/goldi/blob/master/LICENSE)

**Go**ldi: **l**azy **d**ependency **i**njection framework for go.

This library enables you to build your applications based on a dependency injection container.
It helps to make your code modular, flexible and ensures that you can reuse components easily.

If you are familiar with the [Symfony][5] dependency injection framework you should feel at home here.

## The goldi API
Use `go get` to get the goldi API:
```
$ go get github.com/fgrosse/goldi
```

No additional dependencies are required to use the library.
The full documentation is available at [godoc.org][3]. It is almost complete and includes a lot of examples on how to use goldi.

### Usage
First you need to define the types you are going to use later

```go
import (
    "github.com/fgrosse/goldi"
    "github.com/fgrosse/goldi/validation"
)

// create a new container when your application loads
registry := goldi.NewTypeRegistry()
config := map[string]interface{}{
    "some_parameter": "Hello World",
    "timeout":        42.7,
}
container := goldi.NewContainer(registry, config)

// now define the types you want to build using the di container
// you can use simple structs
container.RegisterType("logger", &SimpleLogger{})
container.RegisterType("api.geo.client", new(GeoClient), "http://example.com/geo:1234")

// you can also use factory functions and parameters
container.RegisterType("acme_corp.mailer", NewAwesomeMailer, "first argument", "%some_parameter%")

// dynamic or static parameters and references to other services can be used as arguments
container.RegisterType("renderer", NewRenderer, "@logger")

// closures and functions are also possible
container.Register("http_handler", goldi.NewFuncType(func(w http.ResponseWriter, r *http.Request) {
    // do amazing stuff
}))

// once you are done registering all your types you should probably validate the container
validator := validation.NewContainerValidator()
validator.MustValidate(container) // will panic, use validator.Validate to get the error

// whoever has access to the container can request these types now
logger := container.MustGet("logger").(LoggerInterface)
logger.DoStuff("...")

// in the tests you might want to exchange the registered types with mocks or other implementations
container.RegisterType("logger", NewNullLogger)

// if you already have an instance you want to be used you can inject it directly
myLogger := NewNullLogger()
container.InjectInstance("logger", myLogger)
```

The types are build lazily. This means that the `logger` will only be created when you ask the container for it the first time. Also all built types are singletons. This means that if you call `container.Get("typeID")`two times you will always get the same instance of whatever `typeID` stands for.

More detailed usage examples and a list of features will be available eventually.

## The goldigen binary

If you are used to frameworks like Symfony you might want to define your types in an easy to maintain yaml file.
You can do this using goldigen.

Use `go get` to install the goldigen binary:
```
$ go get github.com/fgrosse/goldi/goldigen
```
Goldigen depends on [gopkg.in/yaml.v2][4] (LGPLv3) for the parsing of the yaml files and [Kingpin][6] (MIT licensed) for the command line flag parsing.

You then need to define your types like this:

```yaml
types:
    logger:
        package: github.com/fgrosse/goldi-example/lib
        type: SimpleLogger

    my_fancy.client:
        package: github.com/fgrosse/goldi-example/lib
        type: Client
        factory: NewDefaultClient
        arguments:
            - "%client_base_url%"   # As in the API you can use parameters here
            - "@logger"             # You can also reference other types 

    time.clock:
        package: github.com/fgrosse/goldi-example/lib/mytime
        type: Clock
        factory: NewSystemClock
        
    http_handler:
        package: github.com/fgrosse/servo/example
        func:    HandleHTTP         # You can register functions as types using the "func" keyword
```

Now you have your type configuration file you can use goldigen like this:

```
$ goldigen --in config/types.yml --out lib/dependency_injection.go
```

This will generate the following output and write it to `lib/dependency_injection.go`:

```go
//go:goldigen --in "../config/types.yml" --out "dependency_injection.go" --package github.com/fgrosse/goldi-example/lib --function RegisterTypes --overwrite --nointeraction
package lib

import (
	"github.com/fgrosse/goldi"
	"github.com/fgrosse/goldi-example/lib/mytime"
	"github.com/fgrosse/servo/example"
)

// RegisterTypes registers all types that have been defined in the file "../config/types.yml"
//
// DO NOT EDIT THIS FILE: it has been generated by goldigen v0.9.7.
// It is however good practice to put this file under version control.
// See https://github.com/fgrosse/goldi for what is going on here.
func RegisterTypes(types goldi.TypeRegistry) {
	types.RegisterAll(map[string]goldi.TypeFactory{
        "http_handler":    goldi.NewFuncType(example.HandleHTTP),
        "logger":          goldi.NewStructType(new(SimpleLogger)),
        "my_fancy.client": goldi.NewType(NewDefaultClient, "%client_base_url%", "@logger"),
        "time.clock":      goldi.NewType(mytime.NewSystemClock),
	})
}
```

Ask you might have noticed goldigen has created a [go generate][7] comment for you.
Next time you want to update `dependency_injection.go` you can simply run `go generate`.

Goldigen tries its best to determine the output files package by looking into your `GOPATH`.
In certain situations this might not be enough so you can set a package explicitly using the `--package` parameter.

For a full list of goldigens flags and parameters try:

```
$ goldigen --help
```

Now all you need to to is to create the di container as you would just using the goldi API and then somewhere in the bootstrapping of your application call.

```go
RegisterTypes(registry)
```

If you have a serious error in your type registration (like returning more than one result from your type factory method)
goldi defers error handling by return an invalid type. You can check for invalid types with the `ContainerValidator`
or by using `goldi.IsValid(TypeFactory)` directly.
Using the [`ContainerValidator`][8] is always the preferred option since it will check for a wide variety of bad configurations
like undefined parameters or circular type dependencies.

Note that using goldigen is completely optional. If you do not like the idea of having an extra build step for your application just use goldis API directly.

### License

Goldi is licensed under the the MIT license. Please see the LICENSE file for details.

## Contributing

Any contributions are always welcome (use pull requests).
For each pull request make sure that you covered your changes and additions with ginkgo tests. If you are unsure how
to write those just drop me a message.

Please keep in mind that I might not always be able to respond immediately but I usually try to react within the week ☺.

[1]: http://onsi.github.io/ginkgo/
[2]: http://onsi.github.io/gomega/
[3]: http://godoc.org/github.com/fgrosse/goldi
[4]: https://github.com/go-yaml/yaml/tree/v2
[5]: http://symfony.com/doc/current/components/dependency_injection/introduction.html
[6]: https://github.com/alecthomas/kingpin/tree/v1.3.6
[7]: http://blog.golang.org/generate
[8]: https://github.com/fgrosse/goldi/blob/master/container_validator.go
