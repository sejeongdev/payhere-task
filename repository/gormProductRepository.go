package repository

import (
	"context"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
	"payhere/util"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormProductRepository struct {
	Conn       *gorm.DB
	ctxTimeout time.Duration
}

// NewGormRepository ...
func NewGormProductRepository(
	conn *gorm.DB,
	timeout time.Duration,
) ProductRepository {
	conn.AutoMigrate(&model.Product{}, &model.ProductSearch{})
	return &gormProductRepository{
		Conn:       conn,
		ctxTimeout: timeout,
	}
}

func (g gormProductRepository) upsertProductSearch(ctx context.Context, scope *gorm.DB, productID uint64, search *model.ProductSearch) error {
	if err := scope.Model(&model.ProductSearch{}).
		Where("product_id = ?", productID).
		Delete(map[string]any{}).Error; err != nil {
		return err
	}

	if search == nil {
		return nil
	}

	search.ProductID = productID
	if err := scope.Create(search).Error; err != nil {
		return err
	}
	return nil
}

// NewProduct ...
func (g gormProductRepository) NewProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx).Begin()
	if err := scope.
		Omit(clause.Associations).
		Create(&product).Error; err != nil {
		scope.Rollback()
		return nil, err
	}

	if err := g.upsertProductSearch(ctx, scope, product.ID, product.Search); err != nil {
		scope.Rollback()
		return nil, err
	}

	scope.Commit()
	return product, nil
}

// GetProductCount ...
func (g gormProductRepository) GetProductCount(ctx context.Context, filter *queryfilter.ProductQueryFilter) (int64, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var count int64
	scope = filter.GetQuery(ctx, scope)
	if err := scope.Model(&model.Product{}).
		Distinct("products.id").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetProductList ...
func (g gormProductRepository) GetProductList(ctx context.Context, filter *queryfilter.ProductQueryFilter) (model.Products, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var products model.Products
	scope = filter.GetQuery(ctx, scope)
	scope = filter.GetOrder(ctx, scope)
	scope = filter.GetPagination(ctx, scope)
	if err := scope.Model(&model.Product{}).
		Select("products.*").
		Group("products.id").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProduct ...
func (g gormProductRepository) GetProduct(ctx context.Context, productID string) (*model.Product, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var product *model.Product
	if err := scope.Model(&model.Product{}).
		Where("id = ?", productID).
		Find(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct ...
func (g gormProductRepository) UpdateProduct(ctx context.Context, productID uint64, product *model.Product) error {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx).Begin()
	if err := scope.Model(&model.Product{}).
		Omit(clause.Associations).
		Select("*").
		Where("id = ?", productID).
		Updates(product).Error; err != nil {
		scope.Rollback()
		return err
	}

	if err := g.upsertProductSearch(ctx, scope, productID, product.Search); err != nil {
		scope.Rollback()
		return err
	}

	scope.Commit()
	return nil
}

// DeleteProduct ...
func (g gormProductRepository) DeleteProduct(ctx context.Context, productID uint64) error {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	if err := scope.Model(&model.Product{}).
		Where("id = ?", productID).
		Delete(map[string]any{}).Error; err != nil {
		return err
	}
	return nil
}
