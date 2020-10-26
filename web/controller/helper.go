package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const SecretKey = "123hn312r9fxh28739dn182gdahs987tgd56afs"

type jwtClaims struct {
	ID     uint
	Access bool
	Expire int64
}

func (c jwtClaims) Valid() error {
	if time.Now().Unix() > c.Expire || c.Expire == 0 {
		return fmt.Errorf("auth token is expred")
	} else if !c.Access {
		return fmt.Errorf("not an access token")
	}
	return nil
}

// Passing a handler that User is the first variable
// Then is auth is success, the handler will be called
func RequrieAuth(handler func(uint, *gin.Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.GetHeader("Authorization")
		if len(tokenStr) == 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "login required"})
			return
		}
		token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired Token"})
			return
		}
		if claims, ok := token.Claims.(*jwtClaims); token.Valid && ok {
			handler(claims.ID, context)

		} else {
			//authorized and pass the context
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad Token"})
		}
	}
}
