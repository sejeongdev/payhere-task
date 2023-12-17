package repository

import (
	"context"
	"fmt"
	"os"
	"payhere/config"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB ...
func InitDB(conf *config.ViperConfig, prefix string) *gorm.DB {
	mainDial := getDialector(conf, prefix)
	dbConfig := &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	dbConn, err := gorm.Open(mainDial, dbConfig)
	if err != nil {
		os.Exit(1)
	}
	return dbConn
}

func getDialector(conf *config.ViperConfig, dbprefix string) gorm.Dialector {
	prefix := func(key string) string {
		if dbprefix == "" {
			return key
		}
		return fmt.Sprintf("%s_%s", dbprefix, key)
	}
	port := conf.GetInt(prefix("db_port"))
	if port == 0 {
		port = 3306
	}

	setting := "charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=UTC"
	dburi := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&%s",
		conf.GetString(prefix("db_user")),
		conf.GetString(prefix("db_pass")),
		conf.GetString(prefix("db_host")),
		port,
		conf.GetString(prefix("db_name")),
		setting,
	)

	main := mysql.Open(dburi)
	return main
}

// AuthRepository ...
type AuthRepository interface {
	Register(ctx context.Context, auth *model.UserAuth) error
	GetAuthByPhone(ctx context.Context, phone string) (*model.UserAuth, error)
	GetUserAuthBySession(ctx context.Context, uid, session string) (auth *model.UserAuth, err error)
	UpdateUserAuthSession(ctx context.Context, uid string, session string) (err error)
}

// UserRepository ...
type UserRepository interface {
	NewUser(ctx context.Context, user *model.User) (ruser *model.User, err error)
	GetUser(ctx context.Context, uid string) (*model.User, error)
}

// ProductRepository ...
type ProductRepository interface {
	NewProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProductCount(ctx context.Context, filter *queryfilter.ProductQueryFilter) (int64, error)
	GetProductList(ctx context.Context, filter *queryfilter.ProductQueryFilter) (model.Products, error)
	GetProduct(ctx context.Context, productID string) (*model.Product, error)
	UpdateProduct(ctx context.Context, productID uint64, product *model.Product) error
	DeleteProduct(ctx context.Context, productID uint64) error
}

// ShopRepository ...
type ShopRepository interface {
	NewShop(ctx context.Context, shop *model.Shop) (*model.Shop, error)
	GetShopCount(ctx context.Context, filter *queryfilter.ShopQueryFilter) (int64, error)
	GetShopList(ctx context.Context, filter *queryfilter.ShopQueryFilter) (model.Shops, error)
	GetShop(ctx context.Context, shopID string) (*model.Shop, error)
	UpdateShop(ctx context.Context, shopID uint64, update map[string]any) error
	DeleteShop(ctx context.Context, shopID uint64) (err error)
}
