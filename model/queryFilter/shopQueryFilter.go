package queryfilter

import (
	"context"
	"fmt"
	"payhere/model"
	"payhere/util"

	"gorm.io/gorm"
)

// ShopQueryFilter ...
type ShopQueryFilter struct {
	BaseQueryFilter
	UID string `form:"uid"`
}

// GetQuery ...
func (q ShopQueryFilter) GetQuery(ctx context.Context, scope *gorm.DB) *gorm.DB {
	if q.UID != "" {
		scope = scope.Where("shops.uid = ?", q.UID)
	}
	return scope
}

// GetOrder ...
func (q ShopQueryFilter) GetOrder(ctx context.Context, scope *gorm.DB) *gorm.DB {
	return scope.Order("shops.id desc")
}

// IsAvailable ...
func (q ShopQueryFilter) IsAvailable(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return uid != "" && q.UID == uid
}

// Init ...
func (q *ShopQueryFilter) Init(ctx context.Context) {
	if q.UID != "" {
		return
	}
	uid, _ := ctx.Value(util.OwnerKey).(string)
	q.UID = uid
}

// GetCursor ...
func (q ShopQueryFilter) GetCursor(ctx context.Context, ss model.Shops) string {
	if q.IsEmptyCursor(ctx, len(ss)) {
		return ""
	}
	shop := ss[len(ss)-1]
	cursors := MultiPagingCursor{
		&PagingCursor{
			Table:  "shops",
			Column: "id",
			Value:  fmt.Sprintf("%d", shop.ID),
		},
	}
	return cursors.MakeCursor()
}
