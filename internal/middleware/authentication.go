package middleware

import (
	"CesarRodriguezPardo/template-go/config"
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/internal/services"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	logger "CesarRodriguezPardo/template-go/infra/logger"
	response "CesarRodriguezPardo/template-go/infra/response"
)

type UserClaims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

// authMiddleware es la instancia singleton del middleware JWT.
// Se inicializa una sola vez con InitJWTAuth() al arrancar la app.
var authMiddleware *jwt.GinJWTMiddleware

// InitJWTAuth inicializa la instancia singleton del middleware JWT.
// Debe llamarse una sola vez al arrancar la aplicación.
func InitJWTAuth() {
	tokenLookup := "header:Authorization,cookie:token"

	signingAlg := config.Cfg.JWT.SigningAlg
	key := config.Cfg.JWT.Key

	isRelease := config.Cfg.Server.Mode == "release"

	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "api",
		Key:              []byte(key),
		SigningAlgorithm: signingAlg,
		Timeout:          time.Hour * 1,
		MaxRefresh:       time.Hour * 24,
		Authenticator:    AuthenticatorFunc,
		Authorizator:     AuthorizatorFunc,
		PayloadFunc:      PayloadFunc,
		Unauthorized:     UnauthorizedHandlerFunc,
		LoginResponse:    LoginResponse,
		LogoutResponse:   LogoutResponse,
		IdentityHandler:  IdentityHandlerFunc,
		TokenLookup:      tokenLookup,
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
		SendCookie:       true,
		CookieName:       "token",
		CookieHTTPOnly:   true,
		SecureCookie:     isRelease,
		CookieSameSite:   http.SameSiteLaxMode,
	})
	if err != nil {
		logger.Fatal("could not setup JWT middleware", err)
	}

	authMiddleware = middleware
}

// GetJWTAuth retorna la instancia singleton del middleware JWT.
func GetJWTAuth() *jwt.GinJWTMiddleware {
	if authMiddleware == nil {
		logger.Fatal("JWT middleware not initialized. Call InitJWTAuth() first.", nil)
	}
	return authMiddleware
}

func AuthenticatorFunc(c *gin.Context) (interface{}, error) {
	var loginValues models.Login

	err := c.ShouldBind(&loginValues)
	if err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	loginEmail := loginValues.Email
	loginPassword := loginValues.Password

	user, err := services.AuthenticateUser(c, loginEmail, loginPassword)
	if err != nil {
		return nil, err
	}

	userClaims := UserClaims{
		ID:   user.ID,
		Role: user.Role,
	}

	logger.Info("Logged user with email: " + loginEmail + ". From ip: " + c.ClientIP())

	c.Set("user", userClaims)
	return userClaims, nil
}

func SetRoles(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("roles", roles)
		c.Next()
	}
}

func AuthorizatorFunc(data interface{}, c *gin.Context) bool {
	userData, ok := data.(map[string]interface{})
	if !ok {
		return false
	}

	roleValue, ok := userData["role"].(string)
	if !ok {
		return false
	}

	activeRole := models.Role(roleValue)

	roles, exists := c.Get("roles")
	if !exists {
		return true
	}

	allowedRoles, ok := roles.([]models.Role)
	if !ok {
		return false
	}

	for _, r := range allowedRoles {
		if activeRole == r {
			return true
		}
	}

	return false
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	userClaims := data.(UserClaims)
	claims := jwt.MapClaims{
		"user": userClaims,
	}
	return claims
}

func IdentityHandlerFunc(c *gin.Context) interface{} {
	jwtClaims := jwt.ExtractClaims(c)
	claims := jwtClaims["user"]
	return claims
}

func UnauthorizedHandlerFunc(c *gin.Context, code int, message string) {
	if message == "cookie token is empty" {
		message = "Not authorized."
	}

	userClaims := IdentityHandlerFunc(c)
	if userClaims == nil {
		logger.Info("Non authenticated request from ip: " + c.ClientIP())
		response.JsonResponse(c, 403, message, nil)
		return
	}

	userMap, ok := userClaims.(map[string]interface{})
	if !ok {
		logger.Info("Non authorized request from ip: " + c.ClientIP())
		response.JsonResponse(c, code, message, nil)
		return
	}

	userIDStr, _ := userMap["id"].(string)

	logger.Info("Request: " + c.Request.URL.Path + " Not authorized from user with id: " + userIDStr + " and ip: " + c.ClientIP())
	response.JsonResponse(c, code, message, nil)
}

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	user, ok := c.Get("user")
	if !ok {
		response.JWTResponse(c, code, "Login failed.", token, expire, nil)
		return
	}
	response.JWTResponse(c, code, "Login succesful.", token, expire, user)
}

func LogoutResponse(c *gin.Context, code int) {
	if code != http.StatusOK {
		c.JSON(code, gin.H{
			"code":    code,
			"message": "could not logout.",
			"status":  "failed",
		})
		return
	}
	c.JSON(code, gin.H{
		"code":    code,
		"message": "logout succesful.",
		"status":  "success",
	})
}
