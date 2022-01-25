# bnm-go
Go library used to get official exchange rates of National bank of Moldova

[![Build Status](https://app.travis-ci.com/OsoianMarcel/bnm-go.svg?branch=master)](https://app.travis-ci.com/OsoianMarcel/bnm-go)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/OsoianMarcel/bnm-go)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/OsoianMarcel/bnm-go/blob/master/LICENSE)

# Basic example
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/OsoianMarcel/bnm-go"
)

func main() {
	// Create new bnm instance
	inst := bnm.NewBnm()

	// Request today's exchange rates in romanian language
	res, err := inst.Request(bnm.NewQuery("ro", time.Now()))
	if err != nil {
		log.Fatal(err)
	}

	// Find USD exchange rate
	if rate, ok := res.FindByCode("USD"); ok {
		fmt.Printf("%s (%s): %.2f\n", rate.Name, rate.Code, rate.Value)
        // Dolar S.U.A. (USD): 18.04
	} else {
		fmt.Printf("USD not found!\n")
	}
	
	// Print all exchange rates
	fmt.Printf("\n%+v\n", res.Rates)
}
```

## Contribute

Contributions to the package are always welcome!

* Report any bugs or issues you find on the [issue tracker].
* You can grab the source code at the package's [Git repository].

## License

All contents of this package are licensed under the [MIT license].

[issue tracker]: https://github.com/OsoianMarcel/bnm-go/issues
[Git repository]: https://github.com/OsoianMarcel/bnm-go
[MIT license]: LICENSE