package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Token  string      `json:"token"`
	Expire time.Time   `json:"expire"`
	Claims interface{} `json:"claims"`
}

func JsonResponse(c *gin.Context, code int, message string, data interface{}) {
	switch {
	case code >= 100 && code < 200:
		c.JSON(code, gin.H{
			"status":  "information",
			"message": message,
			"data":    data,
		})

	case code >= 200 && code < 300:
		c.JSON(code, gin.H{
			"status":  "success",
			"message": message,
			"data":    data,
		})

	case code >= 300 && code < 400:
		c.JSON(code, gin.H{
			"status":  "redirect",
			"message": message,
		})

	case code >= 400 && code < 500:
		c.JSON(code, gin.H{
			"status":  "client error",
			"message": message,
		})

	default:
		c.JSON(code, gin.H{
			"status":  "server error",
			"message": message,
		})
	}
}

func JWTResponse(c *gin.Context, code int, message string, token string, expire time.Time, claims interface{}) {
	data := &Data{
		Token:  token,
		Expire: expire,
		Claims: claims,
	}
	JsonResponse(c, code, message, data)
}

type PaginatedData struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func PaginatedJsonResponse(c *gin.Context, code int, message string, data interface{}, page, limit, total int) {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	paginated := PaginatedData{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	JsonResponse(c, code, message, paginated)
}
