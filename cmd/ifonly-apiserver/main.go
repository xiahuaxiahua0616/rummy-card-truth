package main

import (
	"os"

	"github.com/xiahua/ifonly/cmd/ifonly-apiserver/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	command := app.NewIfOnlyCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
