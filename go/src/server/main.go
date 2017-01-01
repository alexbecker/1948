package main

import (
	"server/auth"
	"server/serve"
	"os"
)

func main() {
	command := os.Args[1]

	switch command {
	case "serve":
		serve.Serve()
	case "users":
		auth.Manage(os.Args[2], os.Args[3:]...)
	default:
		panic("Command must be one of \"serve\" or \"users\"")
	}
}
