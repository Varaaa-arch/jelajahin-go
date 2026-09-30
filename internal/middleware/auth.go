package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PersonalAccessToken struct {
	ID         uint64    `gorm:"primaryKey"`
	TokenableType string  `gorm:"index"`
	TokenableID   string  `gorm:"index"`
	Name       string
	Token      string    `gorm:"uniqueIndex"`
	Abilities  string
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization header"})
			return
		}

		token := parts[1]
		tokenID := strings.Split(token, "|")[0]

		var accessToken PersonalAccessToken
		result := db.Table("personal_access_tokens").
			Where("id = ?", tokenID).
			First(&accessToken)

		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		if accessToken.ExpiresAt != nil && accessToken.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token expired"})
			return
		}

		var user struct {
			ID    string
			Email string
			Role  string
		}
		result = db.Table("users").
			Select("id, email, role").
			Where("id = ?", accessToken.TokenableID).
			First(&user)

		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
		c.Next()
	}
}
