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
	query := `INSERT products SET name=?, price=?, quantity=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, product.Name, product.Price, product.Quantity)
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
	query := `SELECT  id, name, price, quantity FROM products`
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
			&product.Quantity,
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
	query := `UPDATE products SET name=?, price=?, quantity=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, product.Name, product.Price, product.Quantity, product.ID)
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

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d does not exist", id)
	}

	return
}
func (m *mysqlProductRepository) Select(ctx context.Context, id int) (err error) {
	query := `SELECT * FROM products WHERE id=?`
	rows, err := m.Conn.QueryContext(ctx, query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	type Product struct {
		ID       int
		Name     string
		Quantity int
		Price    float64
	}

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
		if err != nil {
			return
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return
	}

	if len(products) == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	for _, p := range products {
		fmt.Printf("ID: %d, Name: %s, Price: %f, Quantity: %d\n", p.ID, p.Name, p.Price, p.Quantity)
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
	mysqlRepo = NewMysqlProductRepostitory(db.GetDB())
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
func selectProduct(cmd *cobra.Command, args []string) error {
	repository := NewMysqlProductRepostitory(db.GetDB())
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
