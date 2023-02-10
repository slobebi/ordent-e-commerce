package product

import (
	"context"
	"net/http"
	enProduct "ordent/internal/entity/product"
	enUser "ordent/internal/entity/user"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	userUsecase interface {
		GetUserSession(sess enUser.Session) *enUser.SessionData
	}

	productUsecase interface {
		InsertProduct(ctx context.Context, form enProduct.ProductRequest) (int64, error)
		UpdateProduct(ctx context.Context, form enProduct.Product) error
		DeleteProduct(ctx context.Context, productID int64) error
		GetProduct(ctx context.Context, productID int64) (*enProduct.Product, error)
		GetProducts(ctx context.Context, page, limit int) ([]enProduct.Product, error)
		GetProductsByType(ctx context.Context, page, limit int, productType enProduct.ProductType) ([]enProduct.Product, error)
		SearchProduct(ctx context.Context, page, limit int, query string) ([]enProduct.Product, error)
	}
)

type Controller struct {
	userUsc    userUsecase
	productUsc productUsecase
}

func NewController(
	userUsc userUsecase,
	productUsc productUsecase,
) *Controller {
	return &Controller{
		userUsc:    userUsc,
		productUsc: productUsc,
	}
}

func (c *Controller) InsertProduct(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUsc.GetUserSession(session)
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

	form := enProduct.ProductRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	productId, err := c.productUsc.InsertProduct(ctx.Request().Context(), form)
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
			"ID":     productId,
		},
	)
}

func (c *Controller) UpdateProduct(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUsc.GetUserSession(session)
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

	form := enProduct.Product{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	err := c.productUsc.UpdateProduct(ctx.Request().Context(), form)
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
		},
	)
}

func (c *Controller) DeleteProduct(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUsc.GetUserSession(session)
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

	productID := ctx.QueryParam("productID")
	if productID == "" {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	productIDInt, err := strconv.ParseInt(productID, 0, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	err = c.productUsc.DeleteProduct(ctx.Request().Context(), productIDInt)
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
		},
	)
}

func (c *Controller) GetProduct(ctx echo.Context) error {
	productID := ctx.QueryParam("productID")
	if productID == "" {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	productIDInt, err := strconv.ParseInt(productID, 0, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	product, err := c.productUsc.GetProduct(ctx.Request().Context(), productIDInt)
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
			"Data":   product,
		},
	)
}

func (c *Controller) GetProducts(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	product, err := c.productUsc.GetProducts(ctx.Request().Context(), pageInt, limitInt)
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
			"Data":   product,
		},
	)
}

func (c *Controller) GetProductsByType(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	productType := ctx.QueryParam("type")
	if productType == "" {
		productType = string(enProduct.HatsProduct)
	}

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	product, err := c.productUsc.GetProductsByType(ctx.Request().Context(), pageInt, limitInt, enProduct.ProductType(productType))
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
			"Data":   product,
		},
	)
}

func (c *Controller) SearchProduct(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	query := ctx.QueryParam("query")

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	product, err := c.productUsc.SearchProduct(ctx.Request().Context(), pageInt, limitInt, query)
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
			"Data":   product,
		},
	)
}
