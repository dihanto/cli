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
	Command        *cobra.Command
}

func NewProductController(productUsecase entity.ProductUsecase, cmd *cobra.Command) *ProductController {
	return &ProductController{
		ProductUsecase: productUsecase,
		Command:        cmd,
	}
}

func (pc *ProductController) InsertCmd() *cobra.Command {
	var name string
	var quantity int
	var price float64

	insertCmd := &cobra.Command{
		Use:   "insert",
		Short: "Inserts a new product into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" || price == 0.0 {
				return fmt.Errorf("missing required flags name")
			}

			product := &entity.Product{Name: name, Price: price, Quantity: quantity}

			if err := pc.ProductUsecase.Insert(context.Background(), product); err != nil {
				return err
			}

			fmt.Printf("Product with ID %d inserted successfully\n", product.ID)
			return nil
		},
	}
	insertCmd.Flags().StringVarP(&name, "name", "n", "", "Product name")
	insertCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	insertCmd.Flags().IntVar(&quantity, "quantity", 0, "Product quantity")
	pc.Command.AddCommand(insertCmd)
	return insertCmd
}
func (pc ProductController) ShowCmd() *cobra.Command {

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows all products in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			products, err := pc.ProductUsecase.Show(context.Background())
			if err != nil {
				return fmt.Errorf("failed to retrieve products: %w", err)
			}

			fmt.Println("Products:")
			for _, p := range products {
				fmt.Printf("- ID: %d, Name: %s, Price: %f, Quantity %d\n", p.ID, p.Name, p.Price, p.Quantity)
			}

			return nil
		},
	}
	pc.Command.AddCommand(showCmd)
	return showCmd
}

func (pc ProductController) UpdateCmd() *cobra.Command {

	var name string
	var quantity int
	var price float64
	var id int

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates an existing product in the database",
		RunE: func(cmd *cobra.Command, args []string) error {

			product := &entity.Product{Name: name, Price: price, Quantity: quantity, ID: int64(id)}
			err := pc.ProductUsecase.Update(context.Background(), product)
			if err != nil {
				return fmt.Errorf("failed to update product: %w", err)
			}

			fmt.Printf("Product with ID %d updated successfully\n", id)
			return nil
		},
	}
	updateCmd.Flags().StringVar(&name, "name", "", "Product name")
	updateCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")
	updateCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	updateCmd.Flags().IntVar(&quantity, "quantity", 0, "Product Quantity")
	updateCmd.MarkFlagRequired("id")
	pc.Command.AddCommand(updateCmd)
	return updateCmd
}

func (pc ProductController) DeleteCmd() *cobra.Command {

	var id int
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Deletes a product from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}
	deleteCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	deleteCmd.MarkFlagRequired("id")
	pc.Command.AddCommand(deleteCmd)
	return deleteCmd
}
func (pc ProductController) SelectCmd() *cobra.Command {
	var id int
	var selectCmd = &cobra.Command{
		Use:   "select",
		Short: "Select a product from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}

	selectCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	selectCmd.MarkFlagRequired("id")
	pc.Command.AddCommand(selectCmd)
	return selectCmd
}
