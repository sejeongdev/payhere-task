package shop

import (
	"net/http"
	"payhere/api/middleware"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
	"payhere/service"

	"github.com/gin-gonic/gin"
)

type shopHTTPHandler struct {
	conf    *config.ViperConfig
	shopSvc service.ShopService
}

func NewHTTPShopHandler(
	conf *config.ViperConfig,
	payhere *gin.RouterGroup,
	shopSvc service.ShopService,
) {
	handler := &shopHTTPHandler{
		conf:    conf,
		shopSvc: shopSvc,
	}

	shop := payhere.Group("/shop")

	shop.Use(middleware.JWTValidate(conf))

	shop.POST("", handler.NewShop)
	shop.GET("/count", handler.GetShopCount)
	shop.GET("", handler.GetShopList)
	shop.GET("/:shopID", handler.GetShop)
	shop.PUT("/:shopID", handler.UpdateShop)
	shop.DELETE("/:shopID", handler.DeleteShop)
}

// NewShop ...
func (h *shopHTTPHandler) NewShop(c *gin.Context) {
	ctx := c.Request.Context()

	shop := &model.Shop{}
	if err := c.BindJSON(shop); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	shop, err := h.shopSvc.NewShop(ctx, shop)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", shop)
}

// GetShopCount ...
func (h *shopHTTPHandler) GetShopCount(c *gin.Context) {
	ctx := c.Request.Context()

	filter := &queryfilter.ShopQueryFilter{}
	if err := c.BindQuery(filter); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	count, err := h.shopSvc.GetShopCount(ctx, filter)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", count)
}

// GetShopList ...
func (h *shopHTTPHandler) GetShopList(c *gin.Context) {
	ctx := c.Request.Context()

	filter := &queryfilter.ShopQueryFilter{}
	if err := c.BindQuery(filter); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	shops, err := h.shopSvc.GetShopList(ctx, filter)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", shops, filter.GetCursor(ctx, shops))
}

// GetShop ...
func (h *shopHTTPHandler) GetShop(c *gin.Context) {
	ctx := c.Request.Context()

	shopID := c.Param("shopID")
	if shopID == "" {
		cutil.Response(c, http.StatusBadRequest, "invalid shopID")
		return
	}

	shop, err := h.shopSvc.GetShop(ctx, shopID)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", shop)
}

// UpdateShop ...
func (h *shopHTTPHandler) UpdateShop(c *gin.Context) {
	ctx := c.Request.Context()

	shopID := c.Param("shopID")
	if shopID == "" {
		cutil.Response(c, http.StatusBadRequest, "invalid shopID")
		return
	}

	shop := &model.Shop{}
	if err := c.BindJSON(shop); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	shop, err := h.shopSvc.UpdateShop(ctx, shopID, shop)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", shop)
}

// DeleteShop ...
func (h *shopHTTPHandler) DeleteShop(c *gin.Context) {
	ctx := c.Request.Context()

	shopID := c.Param("shopID")
	if shopID == "" {
		cutil.Response(c, http.StatusBadRequest, "invalid shopID")
		return
	}

	err := h.shopSvc.DeleteShop(ctx, shopID)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok")
}
