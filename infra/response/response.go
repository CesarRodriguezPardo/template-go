package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Token  string      `json:"token" bson:"token"`
	Expire time.Time   `json:"time" bson:"time"`
	Claims interface{} `json:"claims" bson:"claims"`
}

// JsonResponse: funcion que maneja distintos casos de respuesta, los entrega en formato json.
func JsonResponse(c *gin.Context, code int, message string, data interface{}) {
	// caso 200 OK
	// caso  >= 400 y <500 bad request y mas
	// caso >= 500 errores y mas
	// por añadir mas casos

	switch {
	case code < 200:
		c.JSON(code, gin.H{
			"status":  "information",
			"message": message,
			"data":    data,
		})
	case code < 300 && code >= 200:
		c.JSON(code, gin.H{
			"status":  "success",
			"message": message,
			"data":    data,
		})
	case code < 400 && code >= 300:
		c.JSON(code, gin.H{
			"status":  "redirect",
			"message": message,
			"data":    data,
		})
	case code < 500 && code > 400:
		c.JSON(code, gin.H{
			"status":  "client error response",
			"message": message,
			"data":    data,
		})
	default:
		c.JSON(code, gin.H{
			"status":  "server error response",
			"message": message,
			"data":    data,
		})
	}
}

func JWTResponse(c *gin.Context, code int, message string, token string, expire time.Time, data interface{}) {
	data = &Data{
		Token:  token,
		Expire: expire,
		Claims: data,
	}
	JsonResponse(c, code, message, data)
}
