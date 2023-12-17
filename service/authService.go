package service

import (
	"context"
	"payhere/config"
	"payhere/model"
	"payhere/repository"

	"github.com/juju/errors"
)

type authUsecase struct {
	conf     *config.ViperConfig
	authRepo repository.AuthRepository
}

// NewAuthService ...
func NewAuthService(conf *config.ViperConfig, authRepo repository.AuthRepository) AuthService {
	u := &authUsecase{
		conf:     conf,
		authRepo: authRepo,
	}
	return u
}

// Register ...
func (u authUsecase) Register(ctx context.Context, auth *model.UserAuth) (err error) {
	auth.InitRegister()

	return u.authRepo.Register(ctx, auth)
}

// Login ...
func (u authUsecase) Login(ctx context.Context, auth *model.UserAuth) (rauth *model.UserAuth, err error) {
	dbauth, err := u.authRepo.GetAuthByPhone(ctx, auth.Phone)
	if err != nil {
		return nil, err
	}

	if !dbauth.LoginValidate(auth.Password) {
		return nil, errors.Unauthorizedf("로그인에 실패했습니다.")
	}

	dbauth.SetToken(u.conf.GetString("jwt_secret_key"))

	return dbauth, nil
}
