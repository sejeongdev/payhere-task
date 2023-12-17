package queryfilter

import (
	"context"
	"fmt"
	"payhere/model"
	"payhere/util"

	"gorm.io/gorm"
)

// ProductQueryFilter ...
type ProductQueryFilter struct {
	BaseQueryFilter
	UID  string `form:"uid"`
	Name string `form:"name"`
}

// GetQuery ...
func (q ProductQueryFilter) GetQuery(ctx context.Context, scope *gorm.DB) *gorm.DB {
	if q.UID != "" {
		scope = scope.Where("products.uid = ?", q.UID)
	}
	if q.Name != "" {
		scope = scope.Joins("JOIN product_searches ON product_searches.product_id = products.id").
			Where("product_searches.names like ?", fmt.Sprintf("%%%s%%", q.Name))
	}
	return scope
}

// GetOrder ...
func (q ProductQueryFilter) GetOrder(ctx context.Context, scope *gorm.DB) *gorm.DB {
	return scope.Order("products.id desc")
}

// IsAvailable ...
func (q ProductQueryFilter) IsAvailable(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return uid != "" && q.UID == uid
}

// Init ...
func (q *ProductQueryFilter) Init(ctx context.Context) {
	if q.UID != "" {
		return
	}
	uid, _ := ctx.Value(util.OwnerKey).(string)
	q.UID = uid
}

// GetCursor ...
func (q ProductQueryFilter) GetCursor(ctx context.Context, ps model.Products) string {
	if q.IsEmptyCursor(ctx, len(ps)) {
		return ""
	}
	product := ps[len(ps)-1]
	cursors := MultiPagingCursor{
		&PagingCursor{
			Table:  "products",
			Column: "id",
			Value:  fmt.Sprintf("%d", product.ID),
		},
	}
	return cursors.MakeCursor()
}
