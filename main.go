package main

import (
	"os"

	"github.com/m-mizutani/code-scanning-notify/pkg/controller"
)

func main() {
	controller.New().CLI(os.Args)
}
