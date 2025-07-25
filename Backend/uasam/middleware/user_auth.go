package middleware

import (
	"net/http"
	"time"
	"uasam/commonutils"
	"uasam/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

/*
Usage:
r.Group("some-api-group/", middleware.AuthMiddleware("http://localhost:8002/v1/microservice/authenticate/"))
*/

// AuthResponse is the expected structure of the authentication service response
type UserObject struct {
	ID         uuid.UUID `db:"id"`
	FirstName  string    `db:"first_name"`
	MiddleName *string   `db:"middle_name"`
	LastName   string    `db:"last_name"`
	EmailID    string    `db:"email_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type UserAuthenticateResponse struct {
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription"`
	User              *UserObject `json:"user,omitempty"`
}

// AuthMiddleware is the Gin middleware function
func UserAuthMiddleware(jwtService *commonutils.JWTService, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.GetHeader("USER_TOKEN")

		if userToken == "" {
			c.JSON(http.StatusUnauthorized, UserAuthenticateResponse{Status: -1, StatusDescription: "Missing USER_TOKEN"})
			c.Abort()
			return
		}

		userID, err := jwtService.ParseUserJWT(userToken)
		if err != nil {
			go logger.Warning("unable to authenticate user", zap.String("execution level", "UserAuthMiddleware"), zap.String("Error", err.Error()))
			c.JSON(http.StatusUnauthorized, UserAuthenticateResponse{Status: -1, StatusDescription: "Error Authenticating User"})
			c.Abort()
			return
		}

		c.Set("USER_ID", userID)
		c.Next()
	}
}
