package middleware

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gboliknow/bildwerk/internal/cache"
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

func RateLimitMiddleware(rate time.Duration, limit uint, cacheClient *cache.RedisCache) gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: cacheClient.Client,
		Rate:        rate,
		Limit:       limit,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(http.StatusTooManyRequests, "Too many requests. Try again in %s", time.Until(info.ResetTime).String())
}
