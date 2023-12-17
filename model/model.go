package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BaseModelField ...
type BaseModelField struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"<-:create;autoCreateTime;not null;index"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;not null"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func customTypeToStr(strs []string, idx int) string {
	if idx == 999999 {
		return strs[len(strs)-1]
	}
	if len(strs) <= idx {
		return strs[0]
	}
	return strs[idx]
}

func unmarshalCustomType(data []byte, mapdata map[string]int, defaultVal int) (val int) {
	str := strings.Trim(string(data), "\"")
	if str == "" {
		return defaultVal
	}
	return mapdata[strings.ToLower(str)]
}

// StringArrayithSpecialDelimiter ...
type StrinArrayWithSpecialDelimiter []string

// Scan ...
func (a *StrinArrayWithSpecialDelimiter) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}
	str := string(b)

	if len(str) == 0 {
		return nil
	}

	*a = strings.Split(str, "␞")
	return nil
}

// Value ...
func (a StrinArrayWithSpecialDelimiter) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "", nil
	}
	return strings.Join(a, "␞"), nil
}

// Responser ...
type Responser interface {
	GetHTTPResponse() map[string]any
}
