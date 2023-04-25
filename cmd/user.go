package cmd

import (
	"log"

	db "github.com/dihanto/cli/config"
	"github.com/dihanto/cli/repository"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Inserts a new product into the database",
	RunE:  insertProduct,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows all products in the database",
	RunE:  showProducts,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existing product in the database",
	RunE:  updateProduct,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a product from the database",
	RunE:  deleteProduct,
}
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select a product from the database",
	RunE:  selectProduct,
}

var name string
var price float64
var quantity int
var mysqlRepo repository.ProductRepository
var id int

func init() {

	var err error
	mysqlRepo = repository.NewMysqlProductRepostitory(db.GetDB())
	if err != nil {
		log.Fatal(err)
	}

	insertCmd.Flags().StringVarP(&name, "name", "n", "", "Product name")
	insertCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	insertCmd.Flags().IntVar(&quantity, "quantity", 0, "Product quantity")

	updateCmd.Flags().StringVar(&name, "name", "", "Product name")
	updateCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	updateCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	updateCmd.Flags().IntVar(&quantity, "quantity", 0, "Product Quantity")
	updateCmd.MarkFlagRequired("id")

	deleteCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	deleteCmd.MarkFlagRequired("id")

	selectCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	selectCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(insertCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(selectCmd)
}
