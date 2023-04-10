package main

import (
	"log"

	"github.com/coding-hui/iam/cmd/apiserver/app"
)

func main() {
	cmd := app.NewAPIServerCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
