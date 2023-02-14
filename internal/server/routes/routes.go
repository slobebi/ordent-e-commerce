package routes

import (
	"ordent/internal/config"
	ctrls "ordent/internal/controller"
	enUser "ordent/internal/entity/user"
	"ordent/internal/server/routes/user"
  "ordent/internal/server/routes/product"
  "ordent/internal/server/routes/transaction"

  "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controller *ctrls.Controllers, cfg config.JWT) {
  jwtMiddleware := echojwt.WithConfig(echojwt.Config{
    SigningKey: []byte(cfg.Secret),
    ContextKey: enUser.SessionContextKey,
    NewClaimsFunc: func(c echo.Context) jwt.Claims {
      return new(enUser.TokenClaim)
    },
  }) 
  user.Register(e, controller, jwtMiddleware)
  product.Register(e, controller, jwtMiddleware)
  transaction.Register(e, controller, jwtMiddleware)
}
