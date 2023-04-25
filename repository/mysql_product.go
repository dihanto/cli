package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dihanto/cli/entity"
)

type mysqlProductRepository struct {
	Conn *sql.DB
}

func NewMysqlProductRepostitory(conn *sql.DB) ProductRepository {
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
