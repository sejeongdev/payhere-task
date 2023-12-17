package controller

import (
	"os"
	"payhere/api/controller/auth"
	"payhere/api/controller/product"
	"payhere/api/controller/shop"
	"payhere/api/controller/user"
	"payhere/config"
	"payhere/repository"
	"payhere/service"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitHandler ...
func InitHandler(conf *config.ViperConfig, r *gin.Engine, db *gorm.DB, cancel <-chan os.Signal) (err error) {
	timeout := time.Duration(conf.GetInt("timeout")) * time.Second

	// repo
	authRepo := repository.NewGormAuthRepository(db, timeout)
	userRepo := repository.NewGormUserRepository(db, timeout)
	productRepo := repository.NewGormProductRepository(db, timeout)
	shopRepo := repository.NewGormShopRepository(db, timeout)

	// service
	authSvc := service.NewAuthService(conf, authRepo)
	userSvc := service.NewUserService(userRepo)
	productSvc := service.NewProductService(productRepo, shopRepo)
	shopSvc := service.NewShopService(shopRepo, userRepo)

	// controller
	payhere := r.Group("/payhere")
	// payhere.Use(middleware.JWTValidate(conf))

	auth.NewHTTPAuthHandler(conf, payhere, authRepo, authSvc)
	user.NewHTTPUserHandler(conf, payhere, authRepo, userSvc)
	product.NewHTTPProductHandler(conf, payhere, authRepo, productSvc)
	shop.NewHTTPShopHandler(conf, payhere, authRepo, shopSvc)

	return nil
}
