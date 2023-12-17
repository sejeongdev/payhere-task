package service

import (
	"context"
	"payhere/model"
	queryfilter "payhere/model/queryFilter"
	"payhere/repository"

	"github.com/juju/errors"
)

type shopUsecase struct {
	shopRepo repository.ShopRepository
	userRepo repository.UserRepository
}

// NewShopService ...
func NewShopService(
	shopRepo repository.ShopRepository,
	userRepo repository.UserRepository,
) ShopService {
	u := &shopUsecase{
		shopRepo: shopRepo,
		userRepo: userRepo,
	}
	return u
}

// NewShop ...
func (u shopUsecase) NewShop(ctx context.Context, shop *model.Shop) (*model.Shop, error) {
	shop.Init(ctx)
	if !shop.CreateValidate(ctx) {
		return nil, errors.NotValid
	}
	user, err := u.userRepo.GetUser(ctx, shop.UID)
	if err != nil {
		return nil, err
	}
	if !user.MakeShopAvailable() {
		return nil, errors.NotValid
	}
	shop, err = u.shopRepo.NewShop(ctx, shop)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

// GetShopCount ...
func (u shopUsecase) GetShopCount(ctx context.Context, filter *queryfilter.ShopQueryFilter) (int64, error) {
	filter.Init(ctx)
	if !filter.IsAvailable(ctx) {
		return 0, errors.NotValid
	}
	return u.shopRepo.GetShopCount(ctx, filter)
}

// GetShopList ...
func (u shopUsecase) GetShopList(ctx context.Context, filter *queryfilter.ShopQueryFilter) (model.Shops, error) {
	filter.Init(ctx)
	if !filter.IsAvailable(ctx) {
		return nil, errors.NotValid
	}
	return u.shopRepo.GetShopList(ctx, filter)
}

// GetShop ...
func (u shopUsecase) GetShop(ctx context.Context, shopID string) (*model.Shop, error) {
	dbShop, err := u.shopRepo.GetShop(ctx, shopID)
	if err != nil {
		return nil, err
	}
	if !dbShop.IsShowingAllowed(ctx) {
		return nil, errors.Unauthorized
	}
	return dbShop, nil
}

// UpdateShop ...
func (u shopUsecase) UpdateShop(ctx context.Context, shopID string, shop *model.Shop) (*model.Shop, error) {
	shop.Init(ctx)
	if !shop.UpdateValidate(ctx) {
		return nil, errors.NotValid
	}
	dbShop, err := u.shopRepo.GetShop(ctx, shopID)
	if err != nil {
		return nil, err
	}
	if !dbShop.CheckAuthorization(ctx) {
		return nil, errors.Unauthorized
	}
	updatedShop, updateMap := dbShop.Update(shop)
	if len(updateMap) == 0 {
		return dbShop, nil
	}
	if err = u.shopRepo.UpdateShop(ctx, dbShop.ID, updateMap); err != nil {
		return nil, err
	}
	return updatedShop, nil
}

// DeleteShop ...
func (u shopUsecase) DeleteShop(ctx context.Context, shopID string) error {
	dbShop, err := u.shopRepo.GetShop(ctx, shopID)
	if err != nil {
		return err
	}
	if !dbShop.CheckAuthorization(ctx) {
		return errors.Unauthorized
	}
	return u.shopRepo.DeleteShop(ctx, dbShop.ID)
}
