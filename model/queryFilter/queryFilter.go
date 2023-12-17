package queryfilter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// BaseQueryFilter ...
type BaseQueryFilter struct {
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int64  `json:"limit" form:"limit"`
}

// QueryFilterConst ...
const (
	MinQueryLimit int64 = 10
	MaxQueryLimit int64 = 100
)

// GetPagination ...
func (f *BaseQueryFilter) GetPagination(ctx context.Context, scope *gorm.DB) *gorm.DB {
	f.Limit = f.getLimit()
	scope = scope.Limit(int(f.Limit))
	scope = scope.Where(f.getCursor())
	return scope
}

func (f BaseQueryFilter) getLimit() int64 {
	if f.Limit < MinQueryLimit {
		return MinQueryLimit
	} else if f.Limit > MaxQueryLimit {
		return MaxQueryLimit
	}
	return f.Limit
}

func (f BaseQueryFilter) getCursor() string {
	cursors := f.parseCursor()
	return cursors.GetQuery()
}

func (f BaseQueryFilter) parseCursor() (cursors MultiPagingCursor) {
	dec, err := base64.URLEncoding.DecodeString(f.Cursor)
	if err != nil {
		return nil
	}
	cursors = MultiPagingCursor{}
	if err := json.Unmarshal(dec, &cursors); err != nil {
		return nil
	}
	return cursors
}

// IsEmptyCursor ...
func (q BaseQueryFilter) IsEmptyCursor(ctx context.Context, size int) bool {
	return int(q.Limit) > size
}

// PagingCursor ...
type PagingCursor struct {
	Table  string `json:"table"`
	Column string `json:"column"`
	Value  string `json:"value"`
	Sort   string `json:"sort"`
}

func (pc PagingCursor) getColumnQuery() string {
	if pc.Table == "" {
		return pc.Column
	}
	return fmt.Sprintf("%s.%s", pc.Table, pc.Column)
}

func (pc PagingCursor) getEqualitySign() string {
	if pc.Sort == "asc" {
		return ">"
	}
	return "<"
}

// MultiPagingCursor ...
type MultiPagingCursor []*PagingCursor

// GetQuery ...
func (mc MultiPagingCursor) GetQuery() (query string) {
	for idx, cursor := range mc {
		// can be complicated ...
		// assume only one cursor ...
		if idx != 0 {
			query += " AND "
		}
		query = fmt.Sprintf("%s %s %s", cursor.getColumnQuery(), cursor.getEqualitySign(), cursor.Value)
	}
	return query
}

// MakeCursor ...
func (ms MultiPagingCursor) MakeCursor() string {
	str, err := json.Marshal(ms)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(str)
}
