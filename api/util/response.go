package util

import (
	"payhere/model"

	"github.com/gin-gonic/gin"
)

// PayhereResponse ...
type PayhereResponse struct {
	Meta *ResponseMeta  `json:"meta"`
	Data map[string]any `json:"data"`
}

func (pr *PayhereResponse) setData(result any) {
	pr.Data = make(map[string]any)
	if r, ok := result.(int64); ok {
		pr.Data["count"] = r
	} else if c, ok := result.(model.Responser); ok {
		pr.Data = c.GetHTTPResponse()
	}
}

func (pr *PayhereResponse) setCursor(cursor any) {
	if pr.Data == nil {
		pr.Data = make(map[string]any)
	}
	pr.Data["cursor"] = cursor.(string)
}

// ResponseMeta ...
type ResponseMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// // ResponseData ...
// type ResponseData map[string]any

// Response ...
func Response(c *gin.Context, code int, msg string, result ...any) {
	res := PayhereResponse{
		Meta: &ResponseMeta{
			Code:    code,
			Message: msg,
		},
	}
	if result != nil {
		res.setData(result[0])
		if len(result) > 1 {
			res.setCursor(result[1])
		}
	}

	c.JSON(code, res)
}
