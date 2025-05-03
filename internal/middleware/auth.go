package middleware

import (
	"net/http"

	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := utility.GetTokenFromRequest(c.Request)
		if err != nil {
			utility.RespondWithError(c, http.StatusUnauthorized, "missing or invalid token")
			c.Abort()
			return
		}

		token, err := utility.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			utility.RespondWithError(c, http.StatusUnauthorized, "User not authorized")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utility.RespondWithError(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			utility.RespondWithError(c, http.StatusUnauthorized, "userID not found in token")
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
