# bnm-go v2

A modern Go client for retrieving official exchange rates from the [National Bank of Moldova (BNM)](https://bnm.md).

[![Build Status](https://app.travis-ci.com/OsoianMarcel/bnm-go.svg?branch=master)](https://app.travis-ci.com/OsoianMarcel/bnm-go)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/OsoianMarcel/bnm-go/v2)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/OsoianMarcel/bnm-go/blob/master/LICENSE)

This version is a complete redesign of `v1`, built with improved Go practices, full test coverage, and extensibility for real-world usage.

## Features

- **100% unit test coverage** – fully tested and reliable.
- **Pluggable caching** – support for custom cache adapters.
- **Context support** – cancel or timeout ongoing requests with `context.Context`.
- **Concurrent safe** – designed for multi-goroutine usage.
- **Flexible configuration** – functional options to customize cache, HTTP client, unmarshaler, logging, etc.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/OsoianMarcel/bnm-go/v2"
)

func main() {
    client := bnm.NewClient(
        bnm.WithCache(bnm.NewMemoryCache()),
        bnm.WithWarnError(func(err error) {
            log.Printf("warn: %v", err)
        }),
    )

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    resp, err := client.Fetch(ctx, bnm.NewQuery(time.Now(), bnm.LANG_EN))
    if err != nil {
        log.Fatalf("failed to fetch rates: %v", err)
    }

    fmt.Printf("Rates: %+v\n", resp)
}
```

## Configuration Options

- **WithCache(cache Cache)** – provide a cache implementation.
- **WithWarnError(fn WarnFunc)** – handle non-critical errors gracefully.
- **WithGetRequest(fn GetRequestFunc)** – override HTTP request logic.
- **WithUnmarshaler(fn UnmarshalerFunc)** – customize response unmarshaling.

## Testing

This project follows Go testing best practices.

Run all tests with:

```bash
go test ./...
```

Generate a detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

✅ v2 is fully covered by unit tests.

## Contribute

Contributions to the package are always welcome!

* Report any bugs or issues you find on the [issue tracker].
* You can grab the source code at the package's [Git repository].

## License

All contents of this package are licensed under the [MIT license].

[issue tracker]: https://github.com/OsoianMarcel/bnm-go/issues
[Git repository]: https://github.com/OsoianMarcel/bnm-go
[MIT license]: LICENSE