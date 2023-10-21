package middleware

import (
	"net/http"
	"strings"

	token "github.com/PainCodermax/FashionShop_Website_Backend/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("Authorization")
		if ClientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}
		tokenSplit := strings.Fields(ClientToken)
		tokenResult := strings.Join(tokenSplit[1:], " ")
		claims, err := token.ValidateToken(tokenResult)
		if err != "" {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}
