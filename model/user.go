package model

import (
	"context"
	"encoding/json"
	"payhere/util"
	"time"

	"gorm.io/gorm"
)

// User ...
type User struct {
	UID       string         `json:"uid" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"<-:create;autoCreateTime;not null;index"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;not null"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Name      string         `json:"name" gorm:""`
	Type      UserType       `json:"type" gorm:"index"`
}

// UserType ...
type UserType int32

// UserTypeConst ...
const (
	UserTypeNone UserType = iota
	UserTypeNormal
	UserTypeShopkeeper
)

var userTypeStr []string = []string{"None", "Normal", "Shopkeeper"}
var userTypeMap = map[string]int{
	"none":       int(UserTypeNone),
	"normal":     int(UserTypeNormal),
	"shopkeeper": int(UserTypeShopkeeper),
}

// String ...
func (ut UserType) String() string {
	return customTypeToStr(userTypeStr, int(ut))
}

// MarshalJSON ...
func (ut *UserType) MarshalJSON() (data []byte, err error) {
	return json.Marshal(ut.String())
}

// UnmarshalJSON ...
func (ut *UserType) UnmarshalJSON(data []byte) (err error) {
	*ut = UserType(unmarshalCustomType(data, userTypeMap, int(UserTypeNone)))
	return nil
}

func (ut UserType) isNone() bool {
	return ut == UserTypeNone
}

func (ut UserType) isShopkeeper() bool {
	return ut == UserTypeShopkeeper
}

// CreateValidate ...
func (u User) CreateValidate(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return u.isOwner(uid) && u.hasType() && u.Name != ""
}

// MakeShopAvailable ...
func (u User) MakeShopAvailable() bool {
	return u.Type.isShopkeeper()
}

// Init ...
func (u *User) Init(ctx context.Context) {
	if u.UID != "" {
		return
	}
	uid, _ := ctx.Value(util.OwnerKey).(string)
	u.UID = uid
}

func (u User) isOwner(uid string) bool {
	return uid != "" && u.UID == uid
}

func (u User) hasType() bool {
	return !u.Type.isNone()
}

// GetHTTPResponse ...
func (u *User) GetHTTPResponse() map[string]any {
	return map[string]any{
		"user": u,
	}
}
