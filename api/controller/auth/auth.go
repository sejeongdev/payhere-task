package auth

import (
	"net/http"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/model"
	"payhere/service"

	"github.com/gin-gonic/gin"
)

type authHTTPHandler struct {
	conf    *config.ViperConfig
	authSvc service.AuthService
}

// NewHTTPAuthHandler ...
func NewHTTPAuthHandler(
	conf *config.ViperConfig,
	payhere *gin.RouterGroup,
	authSvc service.AuthService,
) {
	handler := &authHTTPHandler{
		conf:    conf,
		authSvc: authSvc,
	}

	auth := payhere.Group("/auth")

	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
}

// Register ...
func (h *authHTTPHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	auth := &model.UserAuth{}
	if err := c.BindJSON(auth); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if !auth.Validate() {
		cutil.Response(c, http.StatusBadRequest, "인증정보 오류입니다.")
		return
	}

	if err := h.authSvc.Register(ctx, auth); err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}

	cutil.Response(c, http.StatusOK, "ok")
}

// Login ...
func (h *authHTTPHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	auth := &model.UserAuth{}
	if err := c.BindJSON(auth); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	if !auth.Validate() {
		cutil.Response(c, http.StatusBadRequest, "인증정보 오류입니다.")
		return
	}

	auth, err := h.authSvc.Login(ctx, auth)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}

	cutil.Response(c, http.StatusOK, "ok", auth)
}
