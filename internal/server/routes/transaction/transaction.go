package transaction

import (
	ctrls "ordent/internal/controller"

	"ordent/internal/server/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controllers *ctrls.Controllers, jwt echo.MiddlewareFunc) {

	transaction := e.Group("/transaction", jwt, middleware.MidParseSession)

  transaction.POST("/create", controllers.Transcation.CreateTransaction)
  transaction.GET("/user", controllers.Transcation.GetTransactionsByUser)
  transaction.GET("/all", controllers.Transcation.GetAllTransactions)
}
