package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Usage:
r.Group("some-api-group/", middleware.UserAuthMiddleware("http://localhost:8002/v1/microservice/authenticate/"))
*/

// AuthResponse is the expected structure of the authentication service response
type UserAuthResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	UserID            string `json:"UserID,omitempty"`
}

// UserAuthMiddleware is the Gin middleware function
func UserAuthMiddleware(userAuthServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("AUTH_TOKEN")
		originService := c.GetHeader("ORIGIN_SERVICE")
		userToken := c.GetHeader("USER_TOKEN")

		if userToken == "" {
			c.JSON(http.StatusUnauthorized, UserAuthResponse{Status: -1, StatusDescription: "Missing AUTH_TOKEN or ORIGIN_SERVICE"})
			c.Abort()
			return
		}
		req, err := http.NewRequest("POST", userAuthServiceURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{Status: 0, StatusDescription: "Error creating auth request"})
			c.Abort()
			return
		}
		req.Header.Set("USER_TOKEN", userToken)
		req.Header.Set("AUTH_TOKEN", authToken)
		req.Header.Set("ORIGIN_SERVICE", originService)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, AuthResponse{Status: -1, StatusDescription: "UASAM unreachable" + userAuthServiceURL})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{Status: 0, StatusDescription: "Error reading auth response"})
			c.Abort()
			return
		}

		var userAuthResp UserAuthResponse
		if err := json.Unmarshal(bodyBytes, &userAuthResp); err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{Status: 0, StatusDescription: "Invalid auth response format"})
			c.Abort()
			return
		}
		if userAuthResp.Status != 1 {
			c.JSON(http.StatusUnauthorized, AuthResponse{Status: userAuthResp.Status, StatusDescription: userAuthResp.StatusDescription})
			c.Abort()
			return
		}

		fmt.Println("USER ID FROM AUTH IN AUTH: ", userAuthResp.UserID)
		c.Set("USER_ID", userAuthResp.UserID)
		c.Next()
	}
}
