package product

import (
	ctrls "ordent/internal/controller"

	"ordent/internal/server/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controllers *ctrls.Controllers, jwt echo.MiddlewareFunc) {

  // public
	product := e.Group("/product")
  product.GET("/all", controllers.Product.GetProducts)
  product.GET("/tag", controllers.Product.GetProductsByType)
  product.GET("/one", controllers.Product.GetProduct)
  product.GET("/search", controllers.Product.SearchProduct)

  // need auth
  productAdmin := e.Group("/product", jwt, middleware.MidParseSession)
  productAdmin.POST("/insert", controllers.Product.InsertProduct)
  productAdmin.PUT("/update", controllers.Product.UpdateProduct)
  productAdmin.DELETE("/delete", controllers.Product.DeleteProduct)
}
