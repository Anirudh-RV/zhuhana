package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Usage:
r.Group("some-api-group/", middleware.AuthMiddleware("http://localhost:8002/v1/microservice/authenticate/"))
*/

// AuthResponse is the expected structure of the authentication service response
type AuthResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	CallerService     string `json:"callerService,omitempty"`
	CalleeService     string `json:"calleeService,omitempty"`
}

// AuthMiddleware is the Gin middleware function
func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("AUTH_TOKEN")
		originService := c.GetHeader("ORIGIN_SERVICE")

		if authToken == "" || originService == "" {
			c.JSON(http.StatusUnauthorized, AuthResponse{Status: -1, StatusDescription: "Missing AUTH_TOKEN or ORIGIN_SERVICE"})
			c.Abort()
			return
		}
		req, err := http.NewRequest("POST", authServiceURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{Status: 0, StatusDescription: "Error creating auth request"})
			c.Abort()
			return
		}
		req.Header.Set("AUTH_TOKEN", authToken)
		req.Header.Set("ORIGIN_SERVICE", originService)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, AuthResponse{Status: -1, StatusDescription: "UASAM unreachable" + authServiceURL})
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

		var authResp AuthResponse
		if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{Status: 0, StatusDescription: "Invalid auth response format"})
			c.Abort()
			return
		}
		if authResp.Status != 1 {
			c.JSON(http.StatusUnauthorized, AuthResponse{Status: authResp.Status, StatusDescription: authResp.StatusDescription})
			c.Abort()
			return
		}

		c.Set("callerService", authResp.CallerService)
		c.Set("calleeService", authResp.CalleeService)
		c.Next()
	}
}
