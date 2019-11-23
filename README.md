# errors[![GoDoc](https://godoc.org/github.com/arikkfir/go-errors?status.svg)](http://godoc.org/github.com/arikkfir/go-errors) [![Report card](https://goreportcard.com/badge/github.com/arikkfir/go-errors)](https://goreportcard.com/report/github.com/arikkfir/go-errors) [![Sourcegraph](https://sourcegraph.com/github.com/arikkfir/go-errors/-/badge.svg)](https://sourcegraph.com/github.com/arikkfir/go-errors?badge)

This package is a drop-in replacement to the standard Golang `errors` package. The idea is to add missing constructs (in my opinion) that are useful in robust error handling.

This package is a fork of `github.com/pkg/errors` but incorporates useful ideas from `github.com/go-playground/errors` - many thanks to both projects!

## Usage

```go
// These errors will contain a message and a stacktrace
errors.New("failed")
errors.Errorf("bad ID: %d", id)
```

## Wrapping & causation

```go
// The wrapping error contains the given prefix and stacktrace, but also
// provides access to the wrapped error
err := errors.Wrap(fmt.Errorf("bad bad bad"), "Failed doing something")

// And to access it - this will print "bad bad bad"
fmt.Println(errors.Unwrap(err).Error())
fmt.Println(err.Cause().Error())
```

## Tags

You can attach _tags_ to created errors, which is essentially a way to attach metadata about the error, like the ID being updated, the name of the current user, etc.

```go
// These errors will contain a message and a stacktrace
err := errors.New("failed reading accounts").AddTag("source", "/file.csv")

// Prints "Source: /file.csv"
fmt.Printf("Source: %s\n", errors.LookupTag(err, "source"))

// Get all tags
var tags map[string]interface{} = errors.Tags(err) 
```

Keep in mind that tags apply to the error wrapping hierarchy - meaning that if one error is wrapping another error, and the wrapped error has a tag, looking up that tag on the wrapping error will provide that tag. Here's an example:

```go
// Notice how only the "inner" error has the "source" tag.
inner := errors.New("failed reading accounts").AddTag("source", "/file.csv")
outer := errors.Wrap(inner, "oops")

// Prints "Source: /file.csv"
fmt.Printf("Source: %s\n", errors.LookupTag(outer, "source"))

// Get all tags for both "inner" and "outer"
var tags map[string]interface{} = errors.Tags(err) 
```

## Types

You can mark certain errors by tainting them with _types_ - essentially enabling you to ask whether a certain error _is of a certain type or not_.

```go
// These errors will contain a message and a stacktrace
err := errors.New("failed reading accounts").AddType("persistent")

// Prints "Persistent: true"
fmt.Printf("Persistent: %b\n", errors.HasType(err, "persistent"))
// Prints "Transient: false"
fmt.Printf("Transient: %b\n", errors.HasType(err, "transient"))

// Get all tags
var types []string = errors.Types(err) 
```

Keep in mind that types apply to the error wrapping hierarchy - meaning that if one error is wrapping another error, and the wrapped error has a type, looking up that type on the wrapping error will provide that type. Here's an example:

```go
// Notice how only the "inner" error has the "persistent" type.
inner := errors.New("failed reading accounts").AddType("persistent")
outer := errors.Wrap(inner, "oops")

// Prints "Persistent: true"
fmt.Printf("Persistent: %b\n", errors.HasType(outer, "persistent"))
// Prints "Transient: false"
fmt.Printf("Transient: %b\n", errors.HasType(outer, "transient"))

// Get all types for both "inner" and "outer"
var tags map[string]interface{} = errors.Tags(err) 
```

## Contributing

Please read the [Code of Conduct](./docs/CODE_OF_CONDUCT.md) & [Contributing](./docs/CONTRIBUTING.md) documents.

## License

[GNUv3](./LICENSE)
