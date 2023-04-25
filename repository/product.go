package repository

import (
	"context"

	"github.com/dihanto/cli/entity"
)

type ProductRepository interface {
	Insert(ctx context.Context, product *entity.Product) (err error)
	Show(ctx context.Context) ([]entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (err error)
	Delete(ctx context.Context, id int) (err error)
	Select(ctx context.Context, id int) (err error)
}
