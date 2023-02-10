package transaction

import (
	"context"
	"net/http"
	enTransaction "ordent/internal/entity/transaction"
	enUser "ordent/internal/entity/user"

	"github.com/labstack/echo/v4"
)

type (
	userUsecase interface {
		GetUserSession(sess enUser.Session) *enUser.SessionData
	}

	transactionUsecase interface {
		CreateTransaction(ctx context.Context, form enTransaction.TransactionRequest) (int64, error)
		GetTransactionsByUser(ctx context.Context, userID int64) ([]enTransaction.Transaction, error)
		GetAllTransactions(ctx context.Context) ([]enTransaction.Transaction, error)
	}
)

type Controller struct {
	userUc        userUsecase
	transactionUc transactionUsecase
}

func NewController(
	userUc userUsecase,
	transactionUc transactionUsecase,
) *Controller {
	return &Controller{
		userUc:        userUc,
		transactionUc: transactionUc,
	}
}

func (c *Controller) CreateTransaction(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	form := enTransaction.TransactionRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	id, err := c.transactionUc.CreateTransaction(ctx.Request().Context(), form)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
			"ID":     id,
		},
	)
}

func (c *Controller) GetTransactionsByUser(ctx echo.Context) error {

  session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

  response, err := c.transactionUc.GetTransactionsByUser(ctx.Request().Context(), session.ID)
  if err != nil {
    return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
  }

  return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
			"Data":     response,
		},
	)
}

func (c *Controller) GetAllTransactions(ctx echo.Context) error {

  session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}
  
  if !session.IsAdmin {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Not Admin",
			},
		)
	}

  response, err := c.transactionUc.GetAllTransactions(ctx.Request().Context())
  if err != nil {
    return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
  }

  return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
			"Data":     response,
		},
	)
}
