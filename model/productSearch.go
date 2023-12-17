package model

import (
	"strings"

	"github.com/daangn/gorean"
)

// ProductSearch ...
type ProductSearch struct {
	ProductID uint64                         `json:"productID" gorm:"primaryKey"`
	Names     StrinArrayWithSpecialDelimiter `json:"names" gorm:"type:longtext"`
}

// Init ...
func (ps *ProductSearch) Init(name string) {
	ps.initNames(name)
}

func (ps *ProductSearch) initNames(name string) {
	ps.Names = append(ps.Names, name)
	c, _ := gorean.Chosung(name)
	str := strings.Join(c, "")
	if str == "" {
		return
	}
	ps.Names = append(ps.Names, str)
}
