package middleware

import (
	"context"
	"fmt"
	"net/http"
	cutil "payhere/api/util"
	"payhere/config"
	"payhere/repository"
	"payhere/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTValidate ...
func JWTValidate(conf *config.ViperConfig, authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens := c.Request.Header["Token"]
		if len(tokens) == 0 {
			cutil.Response(c, http.StatusInternalServerError, "token not found")
			return
		}
		token := tokens[0]
		ctx := c.Request.Context()
		uid, valid := validateToken(ctx, conf, authRepo, token)
		if !valid {
			cutil.Response(c, http.StatusInternalServerError, "invalid token")
			return
		}
		if uid == "" {
			cutil.Response(c, http.StatusInternalServerError, "uhauthrized user")
			return
		}

		ctx = context.WithValue(ctx, util.OwnerKey, uid)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func validateToken(ctx context.Context, conf *config.ViperConfig, authRepo repository.AuthRepository, token string) (string, bool) {
	ptoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetString("jwt_secret_key")), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	claims, ok := ptoken.Claims.(jwt.MapClaims)
	if !ok || !ptoken.Valid {
		return "", false
	}
	session, _ := claims["session"].(string)
	uid := claims["uid"].(string)
	if session == "" || uid == "" {
		return "", false
	}
	if _, err = authRepo.GetUserAuthBySession(ctx, uid, session); err != nil {
		return "", false
	}
	return uid, true
}
