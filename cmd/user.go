package cmd

import (
	"fmt"

	"github.com/dihanto/cli/db"
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
var showByIdCmd = &cobra.Command{
	Use:   "showbyid",
	Short: "Show user by id",
	Long:  "Show user data by id from database",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")

		db := db.GetDB()

		stmt, err := db.Prepare("SELECT name,email FROM USERS WHERE id=?")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()
		rows, err := stmt.Query(id)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var name, email string
		if rows.Next() {
			err = rows.Scan(&name, &email)
			if err != nil {
				panic(err)
			}
			fmt.Println("Name:", name, "Email:", email)
		} else {
			fmt.Println("No rows returned")
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	},
}

var showUserCmd = &cobra.Command{
	Use:   "showuser",
	Short: "Show user",
	Long:  `Show user list from database`,
	Run: func(cmd *cobra.Command, args []string) {
		db := db.GetDB()

		stmt, err := db.Prepare("SELECT id, name, email FROM users")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var name, email string
			var id int
			err := rows.Scan(&id, &name, &email)
			if err != nil {
				panic(err)
			}
			fmt.Println("id:", id, "Name:", name, "Email:", email)
		}

		err = rows.Err()
		if err != nil {
			panic(err)
		}
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "updateuser",
	Short: "Update a user",
	Long:  `Update the name and/or email of a user in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")

		db := db.GetDB()

		stmt, err := db.Prepare("UPDATE users SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, id)
		if err != nil {
			panic(err)
		}

		fmt.Printf("User with ID %d updated successfully\n", id)
	},
}
var deleteUserCmd = &cobra.Command{
	Use:   "deleteuser",
	Short: "Delete a user",
	Long:  `Delete a user from the database`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")

		db := db.GetDB()

		stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			panic(err)
		}

		fmt.Printf("User with ID %d deleted successfully\n", id)
	},
}

func init() {
	createUserCmd.Flags().StringP("name", "n", "", "The name of the user")
	createUserCmd.Flags().StringP("email", "e", "", "The email address of the user")
	createUserCmd.MarkFlagRequired("name")
	createUserCmd.MarkFlagRequired("email")

	updateUserCmd.Flags().IntP("id", "i", 0, "The ID of the user to update")
	updateUserCmd.Flags().StringP("name", "n", "", "The new name of the user")
	updateUserCmd.Flags().StringP("email", "e", "", "The new email address of the user")
	updateUserCmd.MarkFlagRequired("id")

	deleteUserCmd.Flags().IntP("id", "i", 0, "The ID of the user to delete")
	deleteUserCmd.MarkFlagRequired("id")

	showByIdCmd.Flags().IntP("id", "i", 0, "The ID of the user to show")
	showByIdCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(showUserCmd)
	rootCmd.AddCommand(updateUserCmd)
	rootCmd.AddCommand(deleteUserCmd)
	rootCmd.AddCommand(showByIdCmd)
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
