package main

import (
	"fmt"
	"os"
	// "regexp"
	"strconv"
)

func main() {
	str := os.Args[1]
	raw := strconv.Quote(str)
	fmt.Println(str)
	fmt.Println(raw)
}
