package main

import (
	_ "embed"
	"log"
)

//go:embed LICENSE.md
var licenseString string

func main() {
	log.Default().Println(licenseString)
}
