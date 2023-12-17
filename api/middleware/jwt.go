package middleware

import (
	"context"
	"net/http"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTValidate ...
func JWTValidate(conf *config.ViperConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens := c.Request.Header["Token"]
		if len(tokens) == 0 {
			cutil.Response(c, http.StatusInternalServerError, "token not found")
			return
		}
		token := tokens[0]
		uid, valid := validateToken(conf, token)
		if !valid {
			cutil.Response(c, http.StatusInternalServerError, "invalid token")
			return
		}
		if uid == "" {
			cutil.Response(c, http.StatusInternalServerError, "uhauthrized user")
			return
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, util.OwnerKey, uid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func validateToken(conf *config.ViperConfig, token string) (string, bool) {
	ptoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetString("jwt_secret_key")), nil
	})
	if err != nil {
		return "", false
	}
	uid := ""
	if claims, ok := ptoken.Claims.(jwt.MapClaims); ok && ptoken.Valid {
		uid = claims["uid"].(string)
	}
	return uid, true
}
