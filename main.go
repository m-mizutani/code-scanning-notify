package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, Hello, Hello")
	for _, v := range os.Environ() {
		fmt.Println(v)
	}
}
