package model

import (
	"context"
	"encoding/json"
	"payhere/util"
	"time"
)

// Product ...
type Product struct {
	BaseModelField
	UID           string       `json:"uid" gorm:"type:varchar(36);index"`
	ShopID        uint64       `json:"shopID" gorm:"index"`
	State         ProductState `json:"state" gorm:"default:0"`
	Name          string       `json:"name" gorm:"type:varchar(191);index"`
	Category      string       `json:"category" gorm:""`
	Price         uint64       `json:"price" gorm:""`
	OriginalPrice *uint64      `json:"originalPrice" gorm:""`
	Description   *string      `json:"description" gorm:"type:longtext"`
	Barcode       *string      `json:"barcode" gorm:"type:text"`
	ExpireDate    *time.Time   `json:"expireDate" gorm:""`
	Size          ProductSize  `json:"size" gorm:"default:0"`

	Search *ProductSearch `json:"-" gorm:"->;foreignKey:ProductID;references:ID"`
}

// ProductState ...
type ProductState int32

// ProductStateConst ...
const (
	ProductStateNone ProductState = iota
	ProductStateNew
	ProductStateSale
	ProductStateSoldout
	ProductStateEnd = 999999
)

var productStateStr []string = []string{"None", "New", "Sale", "Soldout", "End"}
var productStateEnd = map[string]int{
	"none":    int(ProductStateNone),
	"new":     int(ProductStateNew),
	"sale":    int(ProductStateSale),
	"soldout": int(ProductStateSoldout),
	"end":     int(ProductStateEnd),
}

// String ...
func (ps ProductState) String() string {
	return customTypeToStr(productStateStr, int(ps))
}

// MarshalJSON ...
func (ps *ProductState) MarshalJSON() (data []byte, err error) {
	return json.Marshal(ps.String())
}

// UnmarsalJSON ...
func (ps *ProductState) UnmarshalJSON(data []byte) (err error) {
	*ps = ProductState(unmarshalCustomType(data, productStateEnd, int(ProductStateNone)))
	return nil
}

func (ps ProductState) isNone() bool {
	return ps == ProductStateNone
}

// ProductSize ...
type ProductSize int32

// ProductSizeConst ...
const (
	ProductSizeNone ProductSize = iota
	ProductSizeSmall
	ProductSizeLarge
)

var productSizeStr []string = []string{"None", "Small", "Large"}
var productSizeMap = map[string]int{
	"none":  int(ProductSizeNone),
	"small": int(ProductSizeSmall),
	"large": int(ProductSizeLarge),
}

// String ...
func (ps ProductSize) String() string {
	return customTypeToStr(productSizeStr, int(ps))
}

// MarshalJSON ...
func (ps *ProductSize) MarshalJSON() (data []byte, err error) {
	return json.Marshal(ps.String())
}

// UnmarshalJSON ...
func (ps *ProductSize) UnmarshalJSON(data []byte) (err error) {
	*ps = ProductSize(unmarshalCustomType(data, productSizeMap, int(ProductSizeNone)))
	return nil
}

func (ps ProductSize) isNone() bool {
	return ps == ProductSizeNone
}

// CreateValidate ...
func (p Product) CreateValidate(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	// price 0 is free product
	return p.isOwner(uid) && p.hasState() && p.Name != "" && p.hasOriginalPrice() && p.hasDescription() && p.hasBarcode() && p.hasExpireDate() && p.hasSize() && p.Category != "" && p.ShopID != 0
}

// UpdateValidate ...
func (p Product) UpdateValidate(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return p.isOwner(uid) && p.Name != ""
}

// CheckAuthorization ...
func (p Product) CheckAuthorization(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return p.isOwner(uid)
}

// IsShowingAllowed ...
func (p Product) IsShowingAllowed(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return p.isOwner(uid)
}

// Update ...
func (p Product) Update(np *Product) *Product {
	if p.Name != np.Name {
		p.Name = np.Name
		p.Search = np.Search
	}
	if p.Price != np.Price {
		p.Price = np.Price
	}
	if np.Description != nil && (p.Description == nil || *p.Description != *np.Description) {
		if *np.Description == "" {
			p.Description = nil
		} else {
			p.Description = np.Description
		}
	}
	return &p
}

// GetHTTPResponse ...
func (p Product) GetHTTPResponse() map[string]any {
	return map[string]any{
		"product": p,
	}
}

// Init ...
func (p *Product) Init(ctx context.Context) {
	if p.UID != "" {
		return
	}
	uid, _ := ctx.Value(util.OwnerKey).(string)
	p.UID = uid

	if p.Name != "" {
		p.Search = &ProductSearch{}
		p.Search.Init(p.Name)
	}
}

func (p Product) isOwner(uid string) bool {
	return uid != "" && p.UID == uid
}

func (p Product) hasOriginalPrice() bool {
	// original price can be free
	return p.OriginalPrice != nil
}

func (p Product) hasDescription() bool {
	return p.Description != nil && *p.Description != ""
}

func (p Product) hasBarcode() bool {
	return p.Barcode != nil && *p.Barcode != ""
}

func (p Product) hasExpireDate() bool {
	return p.ExpireDate != nil && !(*p.ExpireDate).IsZero()
}

func (p Product) hasSize() bool {
	return !p.Size.isNone()
}

func (p Product) hasState() bool {
	return !p.State.isNone()
}

// Products ...
type Products []*Product

// GetHTTPResponse ...
func (ps Products) GetHTTPResponse() map[string]any {
	return map[string]any{
		"products": ps,
	}
}
