package cmd

import (
	"context"
	"fmt"
	"strconv"

	db "github.com/dihanto/cli/config"
	"github.com/dihanto/cli/entity"
	"github.com/dihanto/cli/repository"
	"github.com/spf13/cobra"
)

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

	product := &entity.Product{Name: name, Price: price, Quantity: quantity, ID: int64(id)}
	err := mysqlRepo.Update(context.Background(), product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	fmt.Printf("Product with ID %d updated successfully\n", id)
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
