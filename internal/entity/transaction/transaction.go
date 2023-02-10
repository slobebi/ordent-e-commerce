package transaction


type Transaction struct {
  Id int64 `json:"id" db:"id"`
  UserID int64 `json:"userID" db:"user_id"`
  ProductID int64 `json:"productID" db:"product_id"`
  ItemAmount int64 `json:"itemAmount" db:"item_amount"`
  ProductName string `json:"productName" db:"product_name"`
  ProductType string `json:"productType" db:"product_type"`
  ProductPrice string `json:"productPrice" db:"product_price"`
}

type TransactionRequest struct {
  UserID int64 `json:"userID"`
  ProductID int64 `json:"productID"`
  ItemAmount int64 `json:"itemAmount"`
}
