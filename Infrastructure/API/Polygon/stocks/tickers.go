package stocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTickersV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List of users"})
}
