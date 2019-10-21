## Run Order

### Overview

Go package for determining the concurrent run order of jobs defined in a map of string slices.

### Usage

``` go
package main

import (
	"fmt"
	"log"

	"github.com/gnusey/runorder"
)

func main() {
	m := map[string][]string{
		"a": []string{"b"},
		"b": []string{"c"},
		"c": []string{},
		"d": []string{},
	}
	o, err := runorder.New(m, false)
	if err != nil {
		log.Fatal("F*ck!", err)
	}
	fmt.Println(o)
	// Output:
	// [[c d] [b] [a]]
}
```