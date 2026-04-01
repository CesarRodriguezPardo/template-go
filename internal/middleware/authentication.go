package middleware

import (
	"CesarRodriguezPardo/template-go/config"
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/internal/services"
	"net/http"
	"slices"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	logger "CesarRodriguezPardo/template-go/infra/logger"
	response "CesarRodriguezPardo/template-go/infra/response"
)

type UserClaims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"roles"`
}

// LoadJWTAuth: funcion que implementa un JWT.
// Los atributos de la funcion vienen dados por:

// Realm string : String para mostrar al usuario. (Requerido)
// SigningAlgorithm string : Algoritmo para firmar. (Opcional)
// 					  Los valores posibles son HS256, HS384, HS512, RS256. Por default toma HS526
// Key []byte : Llave secreta usada para firmar. (Requerido)
// Timeout time.Duration : Duracion del token jwt. (Opcional)
// 					  Default una hora.
// MaxRefresh time.Duration : Permite al cliente hacer refresh al token hasta que el tiempo MaxRefresh haya pasado. (Opcional)
// 							  El tiempo maximo del token viene dado por MaxRefresh + Timeout.
// 					  Default cero, es decir, no se puede hacer refresh.
// Authenticator(c *gin.Context) (interface{}, error) : Llama a la funcion de autenticacion basada en la informacion del login. En este caso es: (Requerido)
// 														AuthenticatorFunc
// Authorizator(data interface{}, c *gin.context) bool : Funcion que realiza la autorizacion del usuario autenticado. (Opcional)
// 					 Default success.
// PayloadFunc(data interface{}) MapClaims : Funcion que sera llamada durante en login, es posible añadir mas informacion al payload data del token (Opcional)
//										     La data va a ser accesible a traves de c.Get("JWT_PAYLOAD") y no esta encriptada.
// 					  Default no añade informacion.
// Unauthorized (*gin.Context, int, string) : Funcion que maneja los casos no autorizados, se puede definir de manera personalizada. (Opcional)
// LoginResponse (*gin.Context, int, string, time.Time) : Funcion que maneja la respuesta del login, se puede definir de manera personalizada. (Opcional)
// RefreshResponse (*gin.Context, int, string, time.Time) : Funcion que maneja la respuesta del refresh, se puede definir de manera personalizada. (Opcional)
// IdentityHandler (jwt.MapClaims, interface{}) : Funcion que maneja las claims del jwt, se puede definir de manera personalizada. (Opcional)
// TokenLookup string : Es de la forma "<source>:<name>" y es utilizado para extraer el token de la request. (Opcional)
// 					  Default "header:Authorization".
// 					  Posibles valores:
// 						- "header:<name>"
// 						- "query:<name>"
// 						- "cookie:<name>"
// TokenHeadName string: Es el string en el header. (Opcional)
//					Default "Bearer"
// TimeFunc() time.Time : Entrega el tiempo actual, se puede sobreescribir para usar otro tiempo. (Opcional)
// HTTPStatusMessageFunc(e error, c *gin.Context) string : Es un handler de HTTP messages cuando algo en el middleware falla.
// PrivKeyFile string : Llave privada para algoritmos asimetricos (Opcional)
// PubKeyFile string : Llave privada para algoritmos asimetricos (Opcional)
// SendCookie bool : De manera opcional retorna el token como una cookie (Opcional)
// SecureCookie bool : Permite cookies inseguras para desarrollo sobre http. (Opcional)
// SendAuthorization bool : Permite retornar el header de la autorizacion en cada peticion. (Opcional)
// Si se envia por cookie existen los siguientes atributos extras:
// SendCookie:       true,
// SecureCookie:     false, // for non-HTTPS dev environments
// CookieHTTPOnly:   true,  // JS can't modify
// CookieDomain:     "localhost:8080",
// CookieName:       "token", // default jwt
// TokenLookup:      "cookie:token",
// CookieSameSite:   http.SameSiteDefaultMode, // SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode

func LoadJWTAuth() *jwt.GinJWTMiddleware {
	tokenLookup := "header: Authorization"
	if config.Cfg.Server.Mode == "debug" { // o release como gin trabaja
		tokenLookup = "cookie:token"
	}

	signingAlg := config.Cfg.JWT.SigningAlg

	key := config.Cfg.JWT.Key
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		Key:              []byte(key),
		SigningAlgorithm: signingAlg,
		Authenticator:    AuthenticatorFunc,
		Authorizator:     AuthorizatorFunc,
		PayloadFunc:      PayloadFunc,
		Unauthorized:     UnauthorizedHandlerFunc,
		LoginResponse:    LoginResponse,
		IdentityHandler:  IdentityHandlerFunc,
		TokenLookup:      tokenLookup,
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
		SendCookie:       true,
		CookieName:       "token",
		CookieSameSite:   http.SameSiteLaxMode,
	})
	if err != nil {
		logger.Fatal("Error en el middleware", err)
	}

	return authMiddleware
}

// SetRoles : funcion que define los roles que pueden realizar las peticiones.
// Se implementa sobre las rutas para definir que rol puede ocupar el servicio
func SetRoles(roles ...models.Role) gin.HandlerFunc {
	if slices.Contains(roles, models.ALL) {
		roles = models.ALL_ROLE
	}
	if slices.Contains(roles, models.ADMIN) {
		roles = append(roles, models.ADMIN_ROLE...)
	}

	return func(c *gin.Context) {
		c.Set("roles", roles)
		c.Next()
	}
}

// AuthorizatorFunc : funcion que define si el usuario esta autorizado a utilizar un servicio
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {
	userData := data.(map[string]interface{})

	activeRole := models.Role(userData["active_rol"].(string))

	roles, exists := c.Get("roles")
	if !exists {
		return true
	}
	for _, r := range roles.([]models.Role) {
		if activeRole == r {
			return true
		}
	}
	return false
}

// AuthenticatorFunc: funcion que valida credenciales y setea un usuario en el contexto. (es basicamente un login..)
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

	logger.Info("Login exitoso para usuario " + loginEmail + " desde ip: " + c.ClientIP())

	c.Set("user", userClaims)
	return userClaims, nil
}

// PayloadFunc: funcion que setea las claims del token jwt.
func PayloadFunc(data interface{}) jwt.MapClaims {
	userClaims := data.(UserClaims)
	// si se quisiera usar con postgres, se deberia usar:
	// userClaims := data.(UserClaimsPostgres)
	claims := jwt.MapClaims{
		"user": userClaims,
	}
	return claims
}

// IdentityHandlerFunc: funcion que retorna las claims del token jwt.
func IdentityHandlerFunc(c *gin.Context) interface{} {
	jwtClaims := jwt.ExtractClaims(c)
	claims := jwtClaims["user"]
	return claims
}

// UnauthorizedHandlerFunc: funcion que maneja una peticion no autorizada.

// En algun punto seria bueno poder generalizar los mensajes, debido a que los que retorna el middleware
// 1.- estan en ingles (no es un problema directamente. . )
// 2.- son poco especificos ante la situacion, si intento hacer una peticion directamente
// sin estar logueado, no quiero que el mensaje de error sea "cookie token is empty",
// en cambio, creo y solo creo que seria mejor dar algo mas contextualizado.
func UnauthorizedHandlerFunc(c *gin.Context, code int, message string) {
	if message == "cookie token is empty" {
		message = "No autorizado: token no encontrado o inválido"
	}

	userClaims := IdentityHandlerFunc(c)
	if userClaims == nil {
		logger.Info("Intento de petición no autenticada desde ip: " + c.ClientIP())
		response.JsonResponse(c, 403, message, nil)
		return
	}

	// si se quisiera usar postgres y printear el id, la conversion seria distinta
	// de uid -> string.
	userMap, _ := userClaims.(map[string]interface{})
	userIDStr, _ := userMap["_id"].(string)

	logger.Info("Intento de petición: " + c.Request.URL.Path + ". No autorizada por usuario con id: " + userIDStr + " e ip: " + c.ClientIP())
	response.JsonResponse(c, code, message, nil)
}

// LoginResponse: funcion que setea respuestas del jwt
func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	user, ok := c.Get("user")
	if !ok {
		response.JWTResponse(c, code, "Login fallido.", token, expire, nil)
		return
	}
	response.JWTResponse(c, code, "Login exitoso.", token, expire, user)
}

// LogoutResponse: funcion que setea respuestas del jwt
func LogoutResponse(c *gin.Context, code int) {
	if code != http.StatusOK {
		c.JSON(code, gin.H{
			"code":    code,
			"message": "No se pudo realizar el logout.",
			"status":  "failed",
		})
	}
	c.JSON(code, gin.H{
		"code":    code,
		"message": "Logout realizado con éxito.",
		"status":  "success",
	})
}
