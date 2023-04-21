package cmd

import (
	"fmt"

	"github.com/dihanto/cobra/db"
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:   "createuser",
	Short: "Create a new user",
	Long:  `Create a new user with a name and email address`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		createUser(name, email)
	},
}

func init() {
	createUserCmd.Flags().StringP("name", "n", "", "The name of the user")
	createUserCmd.Flags().StringP("email", "e", "", "The email address of the user")
	createUserCmd.MarkFlagRequired("name")
	createUserCmd.MarkFlagRequired("email")

	rootCmd.AddCommand(createUserCmd)
}

func createUser(name string, email string) {
	db := db.GetDB()

	insert, err := db.Prepare("INSERT INTO users(name, email) VALUES(?, ?)")
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	_, err = insert.Exec(name, email)
	if err != nil {
		panic(err)
	}

	fmt.Println("User created successfully with name:", name)
}
