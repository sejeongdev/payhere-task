package service

import (
	"context"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
)

// AuthService ...
type AuthService interface {
	Register(ctx context.Context, auth *model.UserAuth) (err error)
	Login(ctx context.Context, auth *model.UserAuth) (rauth *model.UserAuth, err error)
}

// UserService ...
type UserService interface {
	NewUser(ctx context.Context, user *model.User) (*model.User, error)
}

// ProductService ...
type ProductService interface {
	NewProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProductCount(ctx context.Context, filter *queryfilter.ProductQueryFilter) (int64, error)
	GetProductList(ctx context.Context, filter *queryfilter.ProductQueryFilter) (model.Products, error)
	GetProduct(ctx context.Context, productID string) (*model.Product, error)
	UpdateProduct(ctx context.Context, productID string, product *model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
}

// ShopService ...
type ShopService interface {
	NewShop(ctx context.Context, shop *model.Shop) (*model.Shop, error)
	GetShopCount(ctx context.Context, filter *queryfilter.ShopQueryFilter) (int64, error)
	GetShopList(ctx context.Context, filter *queryfilter.ShopQueryFilter) (model.Shops, error)
	GetShop(ctx context.Context, shopID string) (*model.Shop, error)
	UpdateShop(ctx context.Context, shopID string, shop *model.Shop) (*model.Shop, error)
	DeleteShop(ctx context.Context, shopID string) error
}
