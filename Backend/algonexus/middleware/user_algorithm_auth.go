package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

/*
Usage:
r.Group("some-api-group/", middleware.AuthMiddleware("http://localhost:8002/v1/microservice/authenticate/"))
*/

type UserAlgorithmAuthenticateResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	UserAlgorithmID   string `json:"userAlgorithmID,omitempty"`
}

// AuthMiddleware is the Gin middleware function
func UserAlgorithmAuthMiddleware(authServiceURL string, microserviceAuthenticator *MicroSeviceAuthenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.GetHeader("USER_ALGORITHM_TOKEN")
		if userToken == "" {
			c.JSON(http.StatusUnauthorized, UserAuthenticateResponse{Status: -1, StatusDescription: "Missing USER_ALGORITHM_TOKEN"})
			c.Abort()
			return
		}
		req, err := http.NewRequest("POST", authServiceURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserAlgorithmAuthenticateResponse{Status: 0, StatusDescription: "Error creating auth request"})
			c.Abort()
			return
		}
		req.Header.Set("USER_ALGORITHM_TOKEN", userToken)
		req.Header.Set("ORIGIN_SERVICE", os.Getenv("ORIGIN_SERVICE"))
		req.Header.Set("AUTH_TOKEN", microserviceAuthenticator.ALL_SERVICE_JWT_TOKENS["uasam"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, UserAlgorithmAuthenticateResponse{Status: -1, StatusDescription: "UASAM unreachable" + authServiceURL})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserAlgorithmAuthenticateResponse{Status: 0, StatusDescription: "Error reading auth response"})
			c.Abort()
			return
		}

		var authResp UserAlgorithmAuthenticateResponse
		if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
			c.JSON(http.StatusInternalServerError, UserAlgorithmAuthenticateResponse{Status: 0, StatusDescription: "Invalid auth response format"})
			c.Abort()
			return
		}
		if authResp.Status != 1 {
			c.JSON(http.StatusUnauthorized, UserAlgorithmAuthenticateResponse{Status: authResp.Status, StatusDescription: authResp.StatusDescription})
			c.Abort()
			return
		}

		c.Set("USER_ALGORITHM_ID", authResp.UserAlgorithmID)
		c.Next()
	}
}
