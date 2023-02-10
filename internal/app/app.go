package app

import (
	"net/http"

	"ordent/internal/config"

	"ordent/internal/server"

	// Repositories
	productRepo "ordent/internal/repository/product"
	userRepo "ordent/internal/repository/user"
  transactionRepo "ordent/internal/repository/transaction"

	// Usecases
	userUsc "ordent/internal/usecase/user"
  productUsc "ordent/internal/usecase/product"
  transactionUsc "ordent/internal/usecase/transaction"

	// Controllers
	ctrls "ordent/internal/controller"
	userCtrl "ordent/internal/controller/user"
  productCtrl "ordent/internal/controller/product"
  transactionCtrl "ordent/internal/controller/transaction"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func InitHTTPServer(cfg config.Config) server.HTTPServerItf {
  // Drivers
  db := connectDatabase(cfg.Database)
  redis := connectRedis(cfg.Redis)

  // Initialize Repositories
  userRepository := userRepo.NewRepository(db, redis, cfg.JWT)
  productRepository := productRepo.NewRepository(db, redis)
  transactionRepository := transactionRepo.NewRepository(db, redis)


  // Initialize Usecases
  userUsecase := userUsc.NewUsecase(userRepository)
  productUsecase := productUsc.NewUsecase(productRepository)
  transactionUsecase := transactionUsc.NewUsecase(transactionRepository, productRepository)

  // Initialize Controllers
  userController := userCtrl.NewController(userUsecase)
  productController := productCtrl.NewController(userUsecase, productUsecase)
  transactionController := transactionCtrl.NewController(userUsecase, transactionUsecase)


  controllers := ctrls.NewControllers(
    userController,
    productController,
    transactionController,
  )

  preMiddlewares := []echo.MiddlewareFunc{
		echoMid.RemoveTrailingSlashWithConfig(echoMid.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
	}

	// Middlewares that runs after the router
	allMiddlewares := []echo.MiddlewareFunc{
		echoMid.Gzip(),
		echoMid.CORSWithConfig(echoMid.CORSConfig{
			Skipper:      echoMid.DefaultSkipper,
			AllowOrigins: []string{"*"}, 
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}),
		echoMid.Recover(),
		echoMid.RequestID(),
	}

  httpServerItf := server.NewHTTPServer(cfg, controllers, preMiddlewares, allMiddlewares)
  return httpServerItf
} 
