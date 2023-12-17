package auth

import (
	"net/http"
	"payhere/api/middleware"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/model"
	"payhere/repository"
	"payhere/service"
	"payhere/util"

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
	authRepo repository.AuthRepository,
	authSvc service.AuthService,
) {
	handler := &authHTTPHandler{
		conf:    conf,
		authSvc: authSvc,
	}

	auth := payhere.Group("/auth")

	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)

	authToken := auth.Use(middleware.JWTValidate(conf, authRepo))
	authToken.POST("/logout", handler.Logout)
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

	token, err := h.authSvc.Login(ctx, auth)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}

	cutil.Response(c, http.StatusOK, "ok", token)
}

// Logout ...
func (h *authHTTPHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	uid, _ := ctx.Value(util.OwnerKey).(string)

	err := h.authSvc.Logout(ctx, uid)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}

	cutil.Response(c, http.StatusOK, "ok")
}
