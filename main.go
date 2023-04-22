package main

import (
	"github.com/dihanto/cli/cmd"
	db "github.com/dihanto/cli/config"
)

func main() {
	db.GetDB()
	cmd.Execute()
}
