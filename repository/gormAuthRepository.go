package repository

import (
	"context"
	"payhere/model"
	"payhere/util"
	"time"

	"github.com/juju/errors"
	"gorm.io/gorm"
)

type gormAuthRepository struct {
	Conn       *gorm.DB
	ctxTimeout time.Duration
}

// NewGormAuthRepository ...
func NewGormAuthRepository(
	conn *gorm.DB,
	timeout time.Duration,
) AuthRepository {
	conn.AutoMigrate(&model.UserAuth{})
	return &gormAuthRepository{
		Conn:       conn,
		ctxTimeout: timeout,
	}
}

// Register ...
func (g gormAuthRepository) Register(ctx context.Context, auth *model.UserAuth) (err error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	if err = scope.Create(&auth).Error; err != nil {
		return err
	}
	return nil
}

// GetAuthByPhone ...
func (g gormAuthRepository) GetAuthByPhone(ctx context.Context, phone string) (auth *model.UserAuth, err error) {
	inCtx, cancel := util.WithTimeout(ctx, g.ctxTimeout)
	defer cancel()

	scope := g.Conn.WithContext(inCtx)
	scope = scope.Model(&model.UserAuth{}).
		Where("phone = ?", phone).
		Find(&auth)
	if err = scope.Error; err != nil {
		return nil, err
	}
	if scope.RowsAffected == 0 {
		return nil, errors.NotFoundf("auth phone[%s]", phone)
	}
	return auth, nil
}
