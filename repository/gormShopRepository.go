package repository

import (
	"context"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
	"payhere/util"
	"time"

	"gorm.io/gorm"
)

type gormShopRepository struct {
	Conn       *gorm.DB
	ctxTimeout time.Duration
}

// NewGormShopRepository ...
func NewGormShopRepository(
	conn *gorm.DB,
	timeout time.Duration,
) ShopRepository {
	conn.AutoMigrate(&model.Shop{})
	return &gormShopRepository{
		Conn:       conn,
		ctxTimeout: timeout,
	}
}

// NewShop ...
func (g gormShopRepository) NewShop(ctx context.Context, shop *model.Shop) (*model.Shop, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	if err := scope.Create(&shop).Error; err != nil {
		return nil, err
	}
	return shop, nil
}

// GetShopCount ...
func (g gormShopRepository) GetShopCount(ctx context.Context, filter *queryfilter.ShopQueryFilter) (int64, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var count int64
	scope = filter.GetQuery(ctx, scope)
	if err := scope.Model(&model.Shop{}).
		Distinct("shops.id").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetShopList ...
func (g gormShopRepository) GetShopList(ctx context.Context, filter *queryfilter.ShopQueryFilter) (model.Shops, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var shops model.Shops
	scope = filter.GetQuery(ctx, scope)
	scope = filter.GetOrder(ctx, scope)
	scope = filter.GetPagination(ctx, scope)
	if err := scope.Model(&model.Shop{}).
		Group("shops.id").
		Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

// GetShop ...
func (g gormShopRepository) GetShop(ctx context.Context, shopID string) (*model.Shop, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var shop *model.Shop
	if err := scope.Model(&model.Shop{}).
		Where("id = ?", shopID).
		Find(&shop).Error; err != nil {
		return nil, err
	}
	return shop, nil
}

// UpdateShop ...
func (g gormShopRepository) UpdateShop(ctx context.Context, shopID uint64, update map[string]any) error {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	if err := scope.Model(&model.Shop{}).
		Where("id = ?", shopID).
		Updates(update).Error; err != nil {
		return err
	}
	return nil
}

// DeleteShop ...
func (g gormShopRepository) DeleteShop(ctx context.Context, shopID uint64) (err error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	if err := scope.Model(&model.Shop{}).
		Where("id = ?", shopID).
		Delete(map[string]any{}).Error; err != nil {
		return err
	}
	return nil
}
