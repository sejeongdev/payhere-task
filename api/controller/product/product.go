package product

import (
	"net/http"
	"payhere/api/middleware"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
	"payhere/repository"
	"payhere/service"

	"github.com/gin-gonic/gin"
)

type productHTTPHandler struct {
	conf *config.ViperConfig

	productSvc service.ProductService
}

// NewHTTPProductHandler ...
func NewHTTPProductHandler(
	conf *config.ViperConfig,
	payhere *gin.RouterGroup,
	authRepo repository.AuthRepository,
	productSvc service.ProductService,
) {
	handler := &productHTTPHandler{
		conf:       conf,
		productSvc: productSvc,
	}

	product := payhere.Group("/product")

	product.Use(middleware.JWTValidate(conf, authRepo))

	product.POST("", handler.NewProduct)
	product.GET("/count", handler.GetProductCount)
	product.GET("", handler.GetProductList)
	product.GET("/:productID", handler.GetProduct)
	product.PUT("/:productID", handler.UpdateProduct)
	product.DELETE("/:productID", handler.DeleteProduct)
}

// NewProduct ...
func (h *productHTTPHandler) NewProduct(c *gin.Context) {
	ctx := c.Request.Context()

	product := &model.Product{}
	if err := c.BindJSON(product); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productSvc.NewProduct(ctx, product)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", product)
}

// GetProductCount ...
func (h *productHTTPHandler) GetProductCount(c *gin.Context) {
	ctx := c.Request.Context()

	filter := &queryfilter.ProductQueryFilter{}
	if err := c.BindQuery(filter); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	count, err := h.productSvc.GetProductCount(ctx, filter)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", count)
}

// GetPrductList ...
func (h *productHTTPHandler) GetProductList(c *gin.Context) {
	ctx := c.Request.Context()

	filter := &queryfilter.ProductQueryFilter{}
	if err := c.BindQuery(filter); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.productSvc.GetProductList(ctx, filter)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", products, filter.GetCursor(ctx, products))
}

// GetProduct ...
func (h *productHTTPHandler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()

	productID := c.Param("productID")
	if productID == "" {
		cutil.Response(c, http.StatusBadRequest, "invalid productID")
		return
	}

	product, err := h.productSvc.GetProduct(ctx, productID)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", product)
}

// UpdateProduct ...
func (h *productHTTPHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	productID := c.Param("productID")
	if productID == "" {
		cutil.Response(c, http.StatusBadRequest, "invalid productID")
		return
	}

	product := &model.Product{}
	if err := c.BindJSON(product); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productSvc.UpdateProduct(ctx, productID, product)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", product)
}

// DeleteProduct ...
func (h *productHTTPHandler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	productID := c.Param("productID")
	if productID == "" {
		cutil.Response(c, http.StatusBadRequest, "invalid productID")
		return
	}

	err := h.productSvc.DeleteProduct(ctx, productID)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok")
}
