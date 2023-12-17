package user

import (
	"net/http"
	"payhere/api/middleware"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/model"
	"payhere/repository"
	"payhere/service"

	"github.com/gin-gonic/gin"
)

type userHTTPHandler struct {
	conf    *config.ViperConfig
	userSvc service.UserService
}

// NewHTTPUserHandler ...
func NewHTTPUserHandler(
	conf *config.ViperConfig,
	payhere *gin.RouterGroup,
	authRepo repository.AuthRepository,
	userSvc service.UserService,
) {
	handler := &userHTTPHandler{
		conf:    conf,
		userSvc: userSvc,
	}

	user := payhere.Group("/user")

	user.Use(middleware.JWTValidate(conf, authRepo))

	user.POST("", handler.NewUser)
}

// NewUser ...
func (h *userHTTPHandler) NewUser(c *gin.Context) {
	ctx := c.Request.Context()

	user := &model.User{}
	if err := c.BindJSON(user); err != nil {
		cutil.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userSvc.NewUser(ctx, user)
	if err != nil {
		rescode, msg := cutil.CauseError(err)
		cutil.Response(c, rescode, msg)
		return
	}
	cutil.Response(c, http.StatusOK, "ok", user)
}
