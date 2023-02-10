package transaction

import (
	"context"
	"errors"
	"fmt"
	enProduct "ordent/internal/entity/product"
	enTransaction "ordent/internal/entity/transaction"
)

type (
	transactionRepository interface {
		CreateTransaction(ctx context.Context, form enTransaction.TransactionRequest) (int64, error)
		GetTransactionsByUser(ctx context.Context, userID int64) ([]enTransaction.Transaction, error)
		GetAllTransactions(ctx context.Context) ([]enTransaction.Transaction, error)
	}

	productRepository interface {
		UpdateSoldAndStock(ctx context.Context, sold, stock int64, productID int64) error
    GetProduct(ctx context.Context, productID int64) (*enProduct.Product, error)
	}
)

type Usecase struct {
	transactionRepo transactionRepository
	productRepo     productRepository
}

func NewUsecase(
	transactionRepo transactionRepository,
	productRepo productRepository,
) *Usecase {
	return &Usecase{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

func (uc *Usecase) CreateTransaction(ctx context.Context, form enTransaction.TransactionRequest) (int64, error) {
  product, err := uc.productRepo.GetProduct(ctx, form.ProductID)
  if err != nil {
    return 0, errors.New(fmt.Sprintf("Failed to get product. err: %+v", err))
  }

  if product.ID == 0 {
    return 0, errors.New(fmt.Sprintf("Product not exist. err: %+v", err))
  }

  transactionID, err := uc.transactionRepo.CreateTransaction(ctx, form)
  if err != nil {
    return 0, errors.New(fmt.Sprintf("failed to create transaction. err: %+v", err))
  }

  err = uc.productRepo.UpdateSoldAndStock(ctx, form.ItemAmount, form.ItemAmount, form.ProductID)
  if err != nil {
    return 0, errors.New(fmt.Sprintf("failed to create transaction. err: %+v", err))
  }

  return transactionID, nil
}

func (uc *Usecase) GetTransactionsByUser(ctx context.Context, userID int64) ([]enTransaction.Transaction, error) {
  transactions, err := uc.transactionRepo.GetTransactionsByUser(ctx, userID)
  if err != nil {
    return make([]enTransaction.Transaction, 0), errors.New(fmt.Sprintf("failed to get transactions. err: %+v", err))
  }

  return transactions, nil
}

func (uc *Usecase) GetAllTransactions(ctx context.Context) ([]enTransaction.Transaction, error) {
  transactions, err := uc.transactionRepo.GetAllTransactions(ctx)
  if err != nil {
    return make([]enTransaction.Transaction, 0), errors.New(fmt.Sprintf("failed to get transactions. err: %+v", err))
  }

  return transactions, nil
}
