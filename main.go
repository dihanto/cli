package main

import (
	"github.com/dihanto/cli/cmd"
	"github.com/dihanto/cli/db"
)

func main() {
	db.GetDB()
	cmd.Execute()
}
