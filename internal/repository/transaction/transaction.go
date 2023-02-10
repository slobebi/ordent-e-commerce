package transaction

import (
	"context"
	"log"
	enTransaction "ordent/internal/entity/transaction"
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

func (r *Repository) CreateTransaction(ctx context.Context, form enTransaction.TransactionRequest) (int64, error) {
  result, err := r.database.ExecContext(ctx, `
    insert into transactions 
      (user_id, product_id, amount)
    values ($1, $2, $3)
  `, form.UserID, form.ProductID, form.ItemAmount)
  if err != nil {
    log.Printf("[CreateTransaction] failed to create transaction. err: %+v", err)
    return 0, err 
  }

  id, err := result.LastInsertId()
  if err != nil {
    log.Printf("[CreateTransaction] failed to create transaction. err: %+v", err)
    return 0, err 
  }

  return id, nil
}

func (r *Repository) GetTransactionsByUser(ctx context.Context, userID int64) ([]enTransaction.Transaction, error) {
  transaction := make([]enTransaction.Transaction, 0)

  err := r.database.SelectContext(ctx, &transaction, `
    select 
      t.id, t.user_id, t.product_id, t.item_amount,
      p.name as product_name, p."type" as product_type, p.price as product_price
    from transactions t
    inner join products p on t.product_id = p.id
    where t.user_id = $1
    order by t.created_time desc
    `, userID)
  if err != nil {
    log.Printf("[GetTransactionsByUser] failed to get transaction for user_id %d. err: %+v", userID, err)
    return transaction, err
  }

  return transaction, nil
}

func (r *Repository) GetAllTransactions(ctx context.Context) ([]enTransaction.Transaction, error) {
  transaction := make([]enTransaction.Transaction, 0)

  err := r.database.SelectContext(ctx, &transaction, `
    select 
      t.id, t.user_id, t.product_id, t.item_amount,
      p.name as product_name, p."type" as product_type, p.price as product_price
    from transactions t
    inner join products p on t.product_id = p.id
    order by t.created_time desc
    `) 
  if err != nil {
    log.Printf("[GetAllTransactions] failed to get transaction. err: %+v", err)
    return transaction, err
  }

  return transaction, nil
}
