package repository

import (
	"context"
	"payhere/model"
	"payhere/util"
	"time"

	"github.com/juju/errors"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	Conn       *gorm.DB
	ctxTimeout time.Duration
}

// NewGormUserRepository ...
func NewGormUserRepository(
	conn *gorm.DB,
	timeout time.Duration,
) UserRepository {
	conn.AutoMigrate(&model.User{})
	return &gormUserRepository{
		Conn:       conn,
		ctxTimeout: timeout,
	}
}

// NewUser ...
func (g gormUserRepository) NewUser(ctx context.Context, user *model.User) (ruser *model.User, err error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	if err = scope.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser ...
func (g gormUserRepository) GetUser(ctx context.Context, uid string) (*model.User, error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	var user *model.User
	scope = scope.Model(&model.User{}).
		Where("uid = ?", uid).
		Find(&user)
	if err := scope.Error; err != nil {
		return nil, err
	}
	if scope.RowsAffected == 0 {
		return nil, errors.NotFoundf("user uid[%s]", uid)
	}
	return user, nil
}
