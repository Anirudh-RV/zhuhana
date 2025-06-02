package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func UserAuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.GetHeader("USER_TOKEN")

		if userToken == "" {
			c.JSON(http.StatusUnauthorized, UserAuthenticateResponse{Status: -1, StatusDescription: "Missing USER_TOKEN"})
			c.Abort()
			return
		}
		req, err := http.NewRequest("POST", authServiceURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserAuthenticateResponse{Status: 0, StatusDescription: "Error creating auth request"})
			c.Abort()
			return
		}
		req.Header.Set("USER_TOKEN", userToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, UserAuthenticateResponse{Status: -1, StatusDescription: "UASAM unreachable" + authServiceURL})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserAuthenticateResponse{Status: 0, StatusDescription: "Error reading auth response"})
			c.Abort()
			return
		}

		var authResp UserAuthenticateResponse
		if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
			c.JSON(http.StatusInternalServerError, UserAuthenticateResponse{Status: 0, StatusDescription: "Invalid auth response format"})
			c.Abort()
			return
		}
		if authResp.Status != 1 {
			c.JSON(http.StatusUnauthorized, UserAuthenticateResponse{Status: authResp.Status, StatusDescription: authResp.StatusDescription})
			c.Abort()
			return
		}

		c.Set("USER_ID", authResp.User.ID)
		c.Set("USER_FIRST_NAME", authResp.User.FirstName)
		c.Set("USER_MIDDLE_NAME", authResp.User.MiddleName)
		c.Set("USER_LAST_NAME", authResp.User.LastName)
		c.Set("USER_EMAIL_ID", authResp.User.EmailID)
		c.Set("USER_CREATED_AT", authResp.User.CreatedAt)
		c.Set("USER_UPDATED_AT", authResp.User.UpdatedAt)
		c.Next()
	}
}
