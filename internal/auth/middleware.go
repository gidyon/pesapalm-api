package auth

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(auth TokenInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj, act string, enforcer *casbin.Enforcer, tkMng TokenInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tkMng.TokenValid(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user hasn't logged in yet"})
			return
		}

		// Extract token
		metadata, err := tkMng.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		// casbin enforces polic
		ok, err := enforcer.Enforce(fmt.Sprint(metadata.UserId), obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error occurred while authorizing"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"message": "forbidden"})
			return
		}

		c.Next()
	}
}
