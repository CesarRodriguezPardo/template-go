package controllers

import (
	"citiaps/golang-backend-template/middleware"

	"github.com/gin-gonic/gin"
)

// LoginFunc : Función que permite hacer login en la aplicación y conseguir un token jwt
//
//	@Tags			Auth
//	@Title			Login
//	@Description	Permite validar el usuario y contraseña en el sistema y generar un token jwt. La solicitud entregara el token jwt, su fecha de expiración y el usuario autenticado.
//	@Accept			json
//	@Produce		json
//	@Param			loginData	body		models.Login	true	"Credenciales de autenticación usach"
//	@Success		200			{string}	string			"ok"
//	@Router			/auth/login [post]
func LoginFunc(c *gin.Context) {
	middleware.LoadJWTAuth().LoginHandler(c)
}

// RefreshToken : Función que permite refrescar el token jwt
//
//	@Tags			Auth
//	@Title			RefreshToken
//	@Description	Permite refrescar el token jwt en caso de que este se encuentre caducado
//	@Accept			json
//	@Produce		json
//	@Router			/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	middleware.LoadJWTAuth().RefreshHandler(c)
}

// Logout : Función que permite cerrar la sesión del usuario
//
//	@Tags			Auth
//	@Title			Logout
//	@Description	Permite cerrar la sesión del usuario, invalidando el token jwt y eliminandolo de los token
//	@Accept			json
//	@Produce		json
//	@Router			/auth/logout [post]
func Logout(c *gin.Context) {
	middleware.LoadJWTAuth().LogoutHandler(c)
}
