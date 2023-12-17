package model

import (
	"context"
	"payhere/util"
)

// Shop ...
type Shop struct {
	BaseModelField
	UID  string `json:"uid" gorm:"index"`
	Name string `json:"name" gorm:""`
}

// CreateValidate ...
func (s Shop) CreateValidate(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return s.IsOwner(uid) && s.Name != ""
}

// UpdateValiate ...
func (s Shop) UpdateValidate(ctx context.Context) bool {
	return s.CreateValidate(ctx)
}

// CheckAuthorization ...
func (s Shop) CheckAuthorization(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return s.IsOwner(uid)
}

// IsShowingAllowed ...
func (s Shop) IsShowingAllowed(ctx context.Context) bool {
	uid, _ := ctx.Value(util.OwnerKey).(string)
	return s.IsOwner(uid)
}

// Update ...
func (s Shop) Update(ns *Shop) (*Shop, map[string]any) {
	update := map[string]any{}
	if s.Name != ns.Name {
		s.Name = ns.Name
		update["name"] = ns.Name
	}
	return &s, update
}

// GetHTTPResponse ...
func (s Shop) GetHTTPResponse() map[string]any {
	return map[string]any{
		"shop": s,
	}
}

// Init ...
func (s *Shop) Init(ctx context.Context) {
	if s.UID != "" {
		return
	}
	uid, _ := ctx.Value(util.OwnerKey).(string)
	s.UID = uid
}

// IsOwner ...
func (s Shop) IsOwner(uid string) bool {
	return uid != "" && s.UID == uid
}

// Shops ...
type Shops []*Shop

// GetHTTPResponse ...
func (ss Shops) GetHTTPResponse() map[string]any {
	return map[string]any{
		"shops": ss,
	}
}
