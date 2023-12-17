package service

import (
	"context"
	"payhere/model"
	"payhere/repository"

	"github.com/juju/errors"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserService ...
func NewUserService(userRepo repository.UserRepository) UserService {
	u := &userUsecase{
		userRepo: userRepo,
	}
	return u
}

// NewUser ...
func (u userUsecase) NewUser(ctx context.Context, user *model.User) (ruser *model.User, err error) {
	user.Init(ctx)
	if !user.CreateValidate(ctx) {
		return nil, errors.NotValidf("user")
	}
	return u.userRepo.NewUser(ctx, user)
}
