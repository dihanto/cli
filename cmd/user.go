package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	db "github.com/dihanto/cli/config"
	"github.com/dihanto/cli/entity"
	"github.com/spf13/cobra"
)

type mysqlProductRepository struct {
	Conn *sql.DB
}

func NewMysqlProductRepostitory(conn *sql.DB) entity.ProductRepository {
	return &mysqlProductRepository{conn}
}

func (m *mysqlProductRepository) Insert(ctx context.Context, product *entity.Product) (err error) {
	query := `INSERT products SET name=?, price=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, product.Name, product.Price)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	product.ID = lastID
	return
}

func (m *mysqlProductRepository) Show(ctx context.Context) (products []entity.Product, err error) {
	query := `SELECT  id, name, price FROM products`
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println(errRow)
		}
	}()
	for rows.Next() {
		product := entity.Product{}
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
		)
		if err != nil {
			log.Println(err)
			return
		}
		products = append(products, product)
	}
	return
}

func (m *mysqlProductRepository) Update(ctx context.Context, product *entity.Product) (err error) {
	query := `UPDATE products SET name=?, price=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, product.Name, product.Price, product.ID)
	if err != nil {
		return
	}
	return
}

func (m *mysqlProductRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM products WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}
	return
}

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

var name string
var price float64
var mysqlRepo entity.ProductRepository
var updateName string
var updatePrice float64
var updateID int
var id int

func init() {

	var err error
	mysqlRepo = NewMysqlProductRepostitory(db.GetDB())
	if err != nil {
		log.Fatal(err)
	}

	insertCmd.Flags().StringVar(&name, "name", "", "Product name")
	insertCmd.Flags().Float64Var(&price, "price", 0.0, "Product price")

	updateCmd.Flags().StringVar(&updateName, "name", "", "Product name")
	updateCmd.Flags().Float64Var(&updatePrice, "price", 0.0, "Product price")
	updateCmd.Flags().IntVar(&updateID, "id", 0, "Product ID")
	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("price")
	updateCmd.MarkFlagRequired("id")

	deleteCmd.Flags().IntVar(&id, "id", 0, "Product ID")
	deleteCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(insertCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(updateCmd)
}

func insertProduct(cmd *cobra.Command, args []string) error {
	if name == "" || price == 0.0 {
		return fmt.Errorf("missing required flags")
	}

	product := &entity.Product{Name: name, Price: price}
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
		fmt.Printf("- ID: %d, Name: %s, Price: %f\n", p.ID, p.Name, p.Price)
	}

	return nil
}
func updateProduct(cmd *cobra.Command, args []string) error {
	if updateID == 0 || updateName == "" || updatePrice == 0.0 {
		return fmt.Errorf("missing required flags")
	}

	product := &entity.Product{Name: updateName, Price: updatePrice, ID: int64(updateID)}
	err := mysqlRepo.Update(context.Background(), product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	fmt.Printf("Product with ID %d updated successfully\n", updateID)
	return nil
}
func deleteProduct(cmd *cobra.Command, args []string) error {
	repository := NewMysqlProductRepostitory(db.GetDB())
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
