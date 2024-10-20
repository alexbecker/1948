package main

import (
	"github.com/alexbecker/1948/commands"
	"github.com/alexbecker/1948/serve"
	"os"
)

func init() {
	commands.Commands["serve"] = func(_ ...string) {
		serve.Serve()
	}
}

func main() {
	command := os.Args[1]
	f, ok := commands.Commands[command]
	if ok {
		f(os.Args[2:]...)
	} else {
		panic("Command not recognized.")
	}
}
