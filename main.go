package main

import (
	"sidewarslobby/pkg/rest"
)

func main() {
	rest := rest.Create()

	rest.Listen(":3000")
}
