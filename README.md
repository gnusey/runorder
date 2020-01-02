## Run Order

### Overview

Go package for determining the concurrent run order of jobs defined in a map of string slices; where the keys of the map represent the jobs and the values represent the jobs dependencies.

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
	o, err := runorder.Calculate(m, false)
	if err != nil {
		log.Fatal("F*ck!", err)
	}
	fmt.Println(o)
	// Output:
	// [[c d] [b] [a]]
}
```