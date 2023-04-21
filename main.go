/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/dihanto/cobra/cmd"
	"github.com/dihanto/cobra/db"
)

func main() {

	// Initialize the database connection pool
	db, err := db.ConnectDB("root:@tcp(localhost:3306)/cobra")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Execute the CLI tool
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
