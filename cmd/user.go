package cmd

import (
	"fmt"

	"github.com/dihanto/cobra/db"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  "A CLI tool to manage users",
	Run: func(cmd *cobra.Command, args []string) {
		// Your main user command logic goes here
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user",
	Long:  "Add a new user to the system",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")

		// Your 'add' user subcommand logic goes here
		fmt.Printf("Adding user '%s' with email '%s'...\n", name, email)
	},
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage database",
	Long:  "A CLI tool to manage database",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn, err := cmd.Flags().GetString("dsn")
		if err != nil {
			return err
		}

		// Connect to the database
		db, err := db.ConnectDB(dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		// Your 'db' command logic goes here
		fmt.Println("Database connection established successfully")
		return nil
	},
}

func init() {

	dbCmd.Flags().StringP("dsn", "", "", "The database connection string")
	dbCmd.MarkFlagRequired("dsn")

	RootCmd.AddCommand(dbCmd)

	addCmd.Flags().StringP("name", "n", "", "The name of the new user")
	addCmd.Flags().StringP("email", "e", "", "The email address of the new user")

	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("email")

	userCmd.AddCommand(addCmd)

	// Register the user command with Cobra
	RootCmd.AddCommand(userCmd)
}
