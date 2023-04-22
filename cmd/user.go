package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	db "github.com/dihanto/cli/config"
	"github.com/dihanto/cli/entity"
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
var mysqlRepo entity.ProductRepository
var updateName string
var updatePrice float64
var updateID int
var updateQuantity int
var id int

func init() {

	var err error
	mysqlRepo = repository.NewMysqlProductRepostitory(db.GetDB())
	if err != nil {
		log.Fatal(err)
	}

	insertCmd.Flags().StringVar(&name, "name", "", "Product name")
	insertCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	insertCmd.Flags().IntVar(&quantity, "quantity", 0, "Product quantity")

	updateCmd.Flags().StringVar(&updateName, "name", "", "Product name")
	updateCmd.Flags().Float64Var(&updatePrice, "price", 0.0, "Product price")
	updateCmd.Flags().IntVar(&updateID, "id", 0, "Product ID")
	updateCmd.Flags().IntVar(&updateQuantity, "quantity", 0, "Product Quantity")
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

func insertProduct(cmd *cobra.Command, args []string) error {
	if name == "" || price == 0.0 {
		return fmt.Errorf("missing required flags name")
	}

	product := &entity.Product{Name: name, Price: price, Quantity: quantity}
	err := mysqlRepo.Insert(context.Background(), product)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	fmt.Printf("Product with ID %d inserted successfully\n", product.ID)
	return nil
}

func showProducts(cmd *cobra.Command, args []string) error {
	products, err := mysqlRepo.Show(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve products: %w", err)
	}

	fmt.Println("Products:")
	for _, p := range products {
		fmt.Printf("- ID: %d, Name: %s, Price: %f, Quantity %d\n", p.ID, p.Name, p.Price, p.Quantity)
	}

	return nil
}
func updateProduct(cmd *cobra.Command, args []string) error {

	product := &entity.Product{Name: updateName, Price: updatePrice, Quantity: updateQuantity, ID: int64(updateID)}
	err := mysqlRepo.Update(context.Background(), product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	fmt.Printf("Product with ID %d updated successfully\n", updateID)
	return nil
}
func deleteProduct(cmd *cobra.Command, args []string) error {
	repository := repository.NewMysqlProductRepostitory(db.GetDB())
	productId, err := strconv.Atoi(cmd.Flag("id").Value.String())
	if err != nil {
		return fmt.Errorf("invalid id value")
	}

	err = repository.Delete(context.Background(), productId)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	fmt.Printf("Product with ID %d deleted successfully\n", productId)
	return nil
}
func selectProduct(cmd *cobra.Command, args []string) error {
	repository := repository.NewMysqlProductRepostitory(db.GetDB())
	productId, err := strconv.Atoi(cmd.Flag("id").Value.String())
	if err != nil {
		return fmt.Errorf("invalid id value")
	}

	err = repository.Select(context.Background(), productId)
	if err != nil {
		return fmt.Errorf("failed to retrieve product: %w", err)
	}

	return nil
}
