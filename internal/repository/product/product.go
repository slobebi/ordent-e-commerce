package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	enProduct "ordent/internal/entity/product"
	"ordent/internal/pkg/redigo"

	"github.com/jmoiron/sqlx"
)

type (
	redis interface {
		Del(keys ...string) error
		Get(key string) *redigo.Result
		Keys(key string) *redigo.Result
		Setex(key string, expireTime int, value interface{}) error
	}
)

type Repository struct {
	database *sqlx.DB
	redis    redis
}

func NewRepository(
	db *sqlx.DB,
	redis redis,
) *Repository {
	return &Repository{
		database: db,
		redis:    redis,
	}
}

const (
  expireTime = 3600
  productKey = "product-%d"
)

func (r *Repository) InsertProduct(ctx context.Context, form enProduct.ProductRequest) (int64, error) {
  result, err := r.database.ExecContext(ctx, `
    insert into products
      (name, "type", price, stock)
    values ($1, $2, $3, $4)
  `, form.Name, form.Type, form.Price, form.Stock)
  if err != nil {
    log.Printf("[InsertProduct] failed to insert product. err: %+v", err)
    return 0, err 
  }

  id, err := result.LastInsertId()
  if err != nil {
    log.Printf("[InsertProduct] failed to insert product. err: %+v", err)
    return 0, err 
  }

  return id, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, form enProduct.Product) error {
  _, err := r.database.ExecContext(ctx, `
    update products
     set name=$1, "type"=$2, price=$3, stock=$4, sold=$5
    where id = $6
  `, form.Name, form.Type, form.Price, form.Stock, form.Sold, form.ID)
  if err != nil {
    log.Printf("[InsertProduct] failed to insert product. err: %+v", err)
    return err 
  }

  key := fmt.Sprintf(productKey, form.ID)

  err = r.redis.Del(key)
  if err != nil {
    log.Printf("Failed to delete redis for key %s. err: %v", key, err)
  }

  return nil
}

func (r *Repository) UpdateSoldAndStock(ctx context.Context, sold, stock int64, productID int64) error {
  _, err := r.database.ExecContext(ctx, `
    update products
     set stock=stock-$1, sold=sold+$2
    where id = $3
  `, stock, sold, productID)
  if err != nil {
    log.Printf("[InsertProduct] failed to insert product. err: %+v", err)
    return err 
  }

  key := fmt.Sprintf(productKey, productID)

  err = r.redis.Del(key)
  if err != nil {
    log.Printf("Failed to delete redis for key %s. err: %v", key, err)
  }

  return nil
}

func (r *Repository) DeleteProduct(ctx context.Context, productID int64) error {
  _, err := r.database.ExecContext(ctx, `
    delete from products 
    where id = $1
  `, productID)
  if err != nil {
    log.Printf("[InsertProduct] failed to insert product. err: %+v", err)
    return err 
  }

  key := fmt.Sprintf(productKey, productID)

  err = r.redis.Del(key)
  if err != nil {
    log.Printf("Failed to delete redis for key %s. err: %v", key, err)
  }

  return nil

}

func (r *Repository) GetProducts(ctx context.Context, limit, offset int) ([]enProduct.Product, error) {
  products := make([]enProduct.Product, 0)

  err := r.database.SelectContext(ctx, &products, `
    select
      id, name, price, stock, "type", sold
    from products
    order by sold desc 
    limit $1
    offset $2
  `, limit, offset)
  if err != nil {
    if err == sql.ErrNoRows {
      return products, nil
    }

    log.Printf("[GetProducts] Failed to get products. err: %+v", err)
    return products, err
  }

  return products, nil
}

func (r *Repository) GetProductsByType(ctx context.Context, limit, offset int, productType enProduct.ProductType) ([]enProduct.Product, error) {
  products := make([]enProduct.Product, 0)

  err := r.database.SelectContext(ctx, &products, `
    select
      id, name, price, stock, "type", sold
    from products
    where "type"=$1
    order by sold desc 
    limit $2
    offset $3
  `, productType, limit, offset)
  if err != nil {
    if err == sql.ErrNoRows {
      return products, nil
    }

    log.Printf("[GetProducts] Failed to get products. err: %+v", err)
    return products, err
  }

  return products, nil
}

func (r *Repository) GetProduct(ctx context.Context, productID int64) (*enProduct.Product, error) {
  product := &enProduct.Product{}

  key := fmt.Sprintf(productKey, productID)
  cache := r.redis.Get(key)
  cacheByte, ok := cache.Value.([]byte)
  if ok {
    err := json.Unmarshal(cacheByte, product)
    if err != nil {
      log.Printf("[GetProduct] failed to unmarshal product")
      return product, nil
    }
  }

  err := r.database.GetContext(ctx, product, `
    select * from products where id=$1
  `, productID)

  if err != nil {
    if err == sql.ErrNoRows {
      return product, nil
    }

    log.Printf("[GetProduct] failed to get product. err: %+v", err)
    return product, err
  }

  productsByte, _ := json.Marshal(*product)

  err = r.redis.Setex(key, expireTime, productsByte)
  if err != nil {
    log.Printf("Failed to get redis for key %s. err: %v", key, err)
  }

  return product, nil
}

func (r *Repository) SearchProduct(ctx context.Context, query string, limit, offset int) ([]enProduct.Product, error) {
  products := make([]enProduct.Product, 0)

  err := r.database.SelectContext(ctx, &products, `
    select
      id, name, price, stock, "type", sold
    from products
    where name ilike '%$1%'
    order by sold desc 
    limit $2
    offset $3
  `, query, limit, offset)
  if err != nil {
    if err == sql.ErrNoRows {
      return products, nil
    }

    log.Printf("[GetProducts] Failed to get products. err: %+v", err)
    return products, err
  }

  return products, nil
}
