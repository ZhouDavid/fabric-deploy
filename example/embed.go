package main

import (
	_ "embed"
	"fmt"
)

//go:embed test.yaml
var tests string

func main() {
	fmt.Println(tests)
}
