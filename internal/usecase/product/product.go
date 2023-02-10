package product

import (
	"context"
	"errors"
	"fmt"

	enProduct "ordent/internal/entity/product"
)


type (
  productRepository interface {
    InsertProduct(ctx context.Context, form enProduct.ProductRequest) (int64, error)
    UpdateProduct(ctx context.Context, form enProduct.Product) error
    UpdateSoldAndStock(ctx context.Context, sold, stock int64, productID int64) error 
    DeleteProduct(ctx context.Context, productID int64) error
    GetProducts(ctx context.Context, limit, offset int) ([]enProduct.Product, error)
    GetProductsByType(ctx context.Context, limit, offset int, productType enProduct.ProductType) ([]enProduct.Product, error)
    GetProduct(ctx context.Context, productID int64) (*enProduct.Product, error)
    SearchProduct(ctx context.Context, query string, limit, offset int) ([]enProduct.Product, error)
  }
)

type Usecase struct {
  productRepo productRepository
}

func NewUsecase(
  productRepo productRepository,
) *Usecase {
  return &Usecase{
    productRepo: productRepo,
  }
}

func (uc *Usecase) InsertProduct(ctx context.Context, form enProduct.ProductRequest) (int64, error) {
  productID, err := uc.productRepo.InsertProduct(ctx, form)
  if err != nil {
    return 0, errors.New(fmt.Sprintf("Failed to insert product. err: %+v", err))
  }

  return productID, nil
}

func (uc *Usecase) UpdateProduct(ctx context.Context, form enProduct.Product) error {
  // check existing product
  product, err := uc.productRepo.GetProduct(ctx, form.ID)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to check product. err: %+v", err))
  }

  if product.ID == 0 {
    return errors.New(fmt.Sprintf("Product not existed"))
  }

  err = uc.productRepo.UpdateProduct(ctx, form)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to update product. err: %+v", err))
  }

  return nil
}

func (uc *Usecase) DeleteProduct(ctx context.Context, productID int64) error {
  // check existing product
  product, err := uc.productRepo.GetProduct(ctx, productID)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to check product. err: %+v", err))
  }

  if product.ID == 0 {
    return errors.New(fmt.Sprintf("Product not existed"))
  }

  err = uc.productRepo.DeleteProduct(ctx, productID)
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to delete product. err: %+v", err))
  }

  return nil
}

func (uc *Usecase) GetProduct(ctx context.Context, productID int64) (*enProduct.Product, error) {
  product, err := uc.productRepo.GetProduct(ctx, productID)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Failed to get product. err: %+v", err))
  }

  return product, nil
}

func (uc *Usecase) GetProducts(ctx context.Context, page, limit int) ([]enProduct.Product, error) {
  offset := (page - 1) * limit

  products, err := uc.productRepo.GetProducts(ctx, limit, offset)
  if err != nil {
    return make([]enProduct.Product, 0), errors.New(fmt.Sprintf("Failed to get all products. err: %+v", err))
  }

  return products, nil
}

func (uc *Usecase) GetProductsByType(ctx context.Context, page, limit int, productType enProduct.ProductType) ([]enProduct.Product, error) {
  offset := (page - 1) * limit

  products, err := uc.productRepo.GetProductsByType(ctx, limit, offset, productType)
  if err != nil {
    return make([]enProduct.Product, 0), errors.New(fmt.Sprintf("Failed to get products. err: %+v", err))
  }

  return products, nil
}

func (uc *Usecase) SearchProduct(ctx context.Context, page, limit int, query string) ([]enProduct.Product, error) {
  offset := (page - 1) * limit

  products, err := uc.productRepo.SearchProduct(ctx, query, limit, offset)
  if err != nil {
    return make([]enProduct.Product, 0), errors.New(fmt.Sprintf("Failed to get products. err: %+v", err))
  }

  return products, nil

}
