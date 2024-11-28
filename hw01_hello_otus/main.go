package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	fmt.Println(reverseStr("Hello, OTUS!"))
}

func reverseStr(str string) string {
	return reverse.String(str)
}
