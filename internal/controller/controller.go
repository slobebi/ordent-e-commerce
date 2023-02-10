package controller

import (
  "ordent/internal/controller/user"
  "ordent/internal/controller/product"
  "ordent/internal/controller/transaction"
)

type Controllers struct {
  User *user.Controller
  Product *product.Controller
  Transcation *transaction.Controller
}

func NewControllers(
  user *user.Controller,
  product *product.Controller,
  transaction *transaction.Controller,
) *Controllers {
  return &Controllers{
    User: user,
    Product: product,
    Transcation: transaction,
  }
}
