package entity

import "context"

type Product struct {
	ID    int64
	Name  string
	Price float64
}

type ProductRepository interface {
	Insert(ctx context.Context, product *Product) (err error)
	Show(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) (err error)
	Delete(ctx context.Context, id int) (err error)
}

type ProductUsecase interface {
	Insert(ctx context.Context, product *Product) (err error)
	Show(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) (err error)
	Delete(ctx context.Context, id int) (err error)
}