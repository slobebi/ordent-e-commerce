package product

type Product struct {
	ID    int64       `json:"id" db:"id"`
	Name  string      `json:"name" db:"name"`
	Type  ProductType `json:"type" db:"type"`
	Price int64       `json:"price" db:"price"`
	Stock int64       `json:"stock" db:"stock"`
	Sold  int64       `json:"sold" db:"sold"`
}

type ProductRequest struct {
	Name  string      `json:"name"`
	Type  ProductType `json:"type"`
	Price int64       `json:"price"`
	Stock int64       `json:"stock"`
}

type ProductType string

const (
	HatsProduct  ProductType = "hats"
	TopsProduct  ProductType = "tops"
	ShortsProcut ProductType = "shorts"
)
