package controllers

import (
	"net/http"
	"uasam/logger"
	"uasam/news/models"
	"uasam/news/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NewsArticlesController struct {
	log                *logger.Logger
	newsArticleService *services.NewsArticleService
}

func NewNewsArticlesController(log *logger.Logger, newsArticleService *services.NewsArticleService) *NewsArticlesController {
	return &NewsArticlesController{
		log:                log,
		newsArticleService: newsArticleService,
	}
}

func (nac *NewsArticlesController) GetNewsArticle(c *gin.Context) {
	rawUserID, _ := c.Get("USER_ID")
	userIDStr, ok := rawUserID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, &models.NewsArticleResponse{
			Status:            -1,
			StatusDescription: "unable to find USER_ID",
		})
		return
	}

	// Parse to UUID
	_, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.NewsArticleResponse{
			Status:            -1,
			StatusDescription: "unable to find USER_ID",
		})
		return
	}

	articleQuery := c.Query("query")
	page := c.Query("page")

	newsArticleData, err := nac.newsArticleService.GetNewsArticle(articleQuery, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.NewsArticleResponse{
			Status:            -1,
			StatusDescription: "Error getting articles",
		})
		return
	}

	c.JSON(http.StatusOK, models.NewsArticleResponse{
		Status:            1,
		StatusDescription: "Articles fetch successful",
		Data:              *newsArticleData,
	})
}
