package main

import (
	"log"

	"github.com/wecoding/iam/cmd/apiserver/app"
)

func main() {
	cmd := app.NewAPIServerCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
