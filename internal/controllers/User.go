package controllers

import (
	"CesarRodriguezPardo/template-go/internal/forms"
	"CesarRodriguezPardo/template-go/internal/services"

	logger "CesarRodriguezPardo/template-go/infra/logger"
	response "CesarRodriguezPardo/template-go/infra/response"

	"github.com/gin-gonic/gin"
)

const (
	createdUserHtml = "mailer/templates/confirmCreatedUser.html"
)

// postgres

// CreateUser
// @Title CreateUser
// @Description Permite crear un usuario en el sistema
// @Summary Crea un usuario
// @Tags Usuario
// @Accept json9
// @Produce json
// @Success 200 {object} forms.UserFormPostgres "Usuario creado con exito."
// @Router /user/postgres [post]
func CreateUser(c *gin.Context) {
	var userForm *forms.UserForm
	if err := c.BindJSON(&userForm); err != nil {
		response.JsonResponse(c, 400, "invalid user data", nil)
		return
	}

	toCreateUser := forms.UserFormToUser(userForm)

	returnedUser, err := services.CreateUser(c, toCreateUser)
	if err != nil {
		logger.Info("failed user creation attempt from: " + c.ClientIP())
		response.JsonResponse(c, 500, err.Error(), toCreateUser)
		return
	}

	logger.Info("Created user with email: " + toCreateUser.Email + " from " + c.ClientIP())
	response.JsonResponse(c, 201, "user created", returnedUser)
}

/*

// GetAllUsersControllerPostgres
// @Title GetAllUsersControllerPostgres
// @Description Permite obtener todos los usuarios del sistema
// @Summary Obtiene todos los usuario
// @Tags Usuario
// @Accept json
// @Produce json
// @Success 200  "Usuario creado con exito."
// @Router /user/postgres [get]
func GetAllUsersControllerPostgres(c *gin.Context) {
	returnedUsers, err := services.GetAllUsersServicePostgres()
	if err != nil {
		response.JsonResponse(c, 500, "Error al obtener todos los usuarios.", returnedUsers)
		return
	}
	if returnedUsers == nil {
		response.JsonResponse(c, 200, "No se han encontrado usuarios.", returnedUsers)
		return
	}
	response.JsonResponse(c, 200, "Usuarios obtenidos con exito.", returnedUsers)
}

*/
