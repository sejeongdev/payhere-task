package service

import (
	"context"
	"fmt"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
	"payhere/repository"

	"github.com/juju/errors"
)

type productUsecase struct {
	productRepo repository.ProductRepository
	shopRepo    repository.ShopRepository
}

// NewProductService ...
func NewProductService(
	productRepo repository.ProductRepository,
	shopRepo repository.ShopRepository,
) ProductService {
	u := &productUsecase{
		productRepo: productRepo,
		shopRepo:    shopRepo,
	}
	return u
}

// NewProduct ...
func (u productUsecase) NewProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.Init(ctx)
	if !product.CreateValidate(ctx) {
		return nil, errors.NotValid
	}
	var err error
	shop, err := u.shopRepo.GetShop(ctx, fmt.Sprintf("%d", product.ShopID))
	if err != nil {
		return nil, err
	}
	if !shop.IsOwner(product.UID) {
		return nil, errors.Unauthorized
	}
	if product, err = u.productRepo.NewProduct(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// GetProductCount ...
func (u productUsecase) GetProductCount(ctx context.Context, filter *queryfilter.ProductQueryFilter) (int64, error) {
	filter.Init(ctx)
	if !filter.IsAvailable(ctx) {
		return 0, errors.NotValid
	}
	return u.productRepo.GetProductCount(ctx, filter)
}

// GetProductList ...
func (u productUsecase) GetProductList(ctx context.Context, filter *queryfilter.ProductQueryFilter) (model.Products, error) {
	filter.Init(ctx)
	if !filter.IsAvailable(ctx) {
		return nil, errors.NotValid
	}
	return u.productRepo.GetProductList(ctx, filter)
}

// GetProduct ...
func (u productUsecase) GetProduct(ctx context.Context, productID string) (*model.Product, error) {
	product, err := u.productRepo.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !product.IsShowingAllowed(ctx) {
		return nil, errors.Unauthorized
	}
	return product, nil
}

// UpdateProduct ...
func (u productUsecase) UpdateProduct(ctx context.Context, productID string, product *model.Product) (*model.Product, error) {
	product.Init(ctx)
	if !product.UpdateValidate(ctx) {
		return nil, errors.NotValid
	}

	dbProduct, err := u.productRepo.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !product.CheckAuthorization(ctx) {
		return nil, errors.Unauthorized
	}
	updatedProduct := dbProduct.Update(product)
	if err = u.productRepo.UpdateProduct(ctx, dbProduct.ID, updatedProduct); err != nil {
		return nil, err
	}
	return updatedProduct, nil
}

// DeleteProduct ...
func (u productUsecase) DeleteProduct(ctx context.Context, productID string) error {
	dbProduct, err := u.productRepo.GetProduct(ctx, productID)
	if err != nil {
		return err
	}
	if !dbProduct.CheckAuthorization(ctx) {
		return errors.Unauthorized
	}
	if err = u.productRepo.DeleteProduct(ctx, dbProduct.ID); err != nil {
		return err
	}
	return nil
}
