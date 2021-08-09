package main

import (
	"os"

	"github.com/m-mizutani/cs-alert-notify/pkg/controller"
)

func main() {
	controller.New().CLI(os.Args)
}
