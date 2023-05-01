package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dihanto/cli/entity"
	"github.com/spf13/cobra"
)

type ProductController struct {
	ProductUsecase entity.ProductUsecase
}

func NewProductController(productUsecase entity.ProductUsecase) *ProductController {
	return &ProductController{
		ProductUsecase: productUsecase,
	}
}

var name string
var price float64
var quantity int
var id int

func (pc *ProductController) Route(rootCmd *cobra.Command) {

	insertCmd := &cobra.Command{
		Use:   "insert",
		Short: "Insert a new product to the database",
		RunE:  pc.Insert(),
	}
	insertCmd.Flags().StringVarP(&name, "name", "n", "", "Product name")
	insertCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	insertCmd.Flags().IntVar(&quantity, "quantity", 0, "Product quantity")

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Shows all products in the database",
		RunE:  pc.Show(),
	}
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates an existing product in the database",
		RunE:  pc.Update(),
	}
	updateCmd.Flags().StringVar(&name, "name", "", "Product name")
	updateCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	updateCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	updateCmd.Flags().IntVar(&quantity, "quantity", 0, "Product Quantity")
	updateCmd.MarkFlagRequired("id")

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a product from the database",
		RunE:  pc.Delete(),
	}
	deleteCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	deleteCmd.MarkFlagRequired("id")

	selectCmd := &cobra.Command{
		Use:   "select",
		Short: "Select a product from the database",
		RunE:  pc.Select(),
	}
	selectCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	selectCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(insertCmd, showCmd, updateCmd, deleteCmd, selectCmd)

}

func (pc *ProductController) Insert() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if name == "" || price == 0.0 {
			return fmt.Errorf("missing required flags name")
		}

		product := &entity.Product{Name: name, Price: price, Quantity: quantity}

		if err := pc.ProductUsecase.Insert(context.Background(), product); err != nil {
			return err
		}

		fmt.Printf("Product with ID %d inserted successfully\n", product.ID)
		return nil
	}
}

func (pc ProductController) Show() func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {
		products, err := pc.ProductUsecase.Show(context.Background())
		if err != nil {
			return fmt.Errorf("failed to retrieve products: %w", err)
		}

		fmt.Println("Products:")
		for _, p := range products {
			fmt.Printf("- ID: %d, Name: %s, Price: %0.1f, Quantity %d\n", p.ID, p.Name, p.Price, p.Quantity)
		}
		return nil
	}
}

func (pc ProductController) Update() func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {

		product := &entity.Product{Name: name, Price: price, Quantity: quantity, ID: int64(id)}
		err := pc.ProductUsecase.Update(context.Background(), product)
		if err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}

		fmt.Printf("Product with ID %d updated successfully\n", id)
		return nil
	}
}

func (pc ProductController) Delete() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// repository := repository.NewMysqlProductRepostitory(db.GetDB())
		productId, err := strconv.Atoi(cmd.Flag("id").Value.String())
		if err != nil {
			return fmt.Errorf("invalid id value")
		}

		err = pc.ProductUsecase.Delete(context.Background(), productId)
		if err != nil {
			return fmt.Errorf("failed to delete product: %w", err)
		}

		fmt.Printf("Product with ID %d deleted successfully\n", productId)
		return nil
	}
}

func (pc ProductController) Select() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// repository := repository.NewMysqlProductRepostitory(db.GetDB())
		productId, err := strconv.Atoi(cmd.Flag("id").Value.String())
		if err != nil {
			return fmt.Errorf("invalid id value")
		}

		_, err = pc.ProductUsecase.Select(context.Background(), productId)
		if err != nil {
			return fmt.Errorf("failed to retrieve product: %w", err)
		}

		return nil
	}
}
