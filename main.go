package main

import (
	"github.com/dihanto/cobra/cmd"
	"github.com/dihanto/cobra/db"
)

func main() {
	db.GetDB()
	cmd.Execute()
}
