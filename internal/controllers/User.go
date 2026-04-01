package controllers

import (
	"CesarRodriguezPardo/template-go/internal/forms"
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/internal/services"

	logger "CesarRodriguezPardo/template-go/infra/logger"
	response "CesarRodriguezPardo/template-go/infra/response"

	"github.com/gin-gonic/gin"
)

const (
	createdUserHtml = "mailer/templates/confirmCreatedUser.html"
)

// postgres

// CreateUserControllerPostgres
// @Title CreateUserControllerPostgres
// @Description Permite crear un usuario en el sistema
// @Summary Crea un usuario
// @Tags Usuario
// @Accept json
// @Produce json
// @Success 200 {object} forms.UserFormPostgres "Usuario creado con exito."
// @Router /user/postgres [post]
func CreateUserControllerPostgres(c *gin.Context) {
	var userForm *forms.UserForm
	err := c.BindJSON(&userForm)

	if err != nil {
		response.JsonResponse(c, 500, "Error en los datos del usuario.", userForm)
		return
	}

	newUser := &models.User{
		Name:       userForm.Name,
		MiddleName: userForm.MiddleName,
		Email:      userForm.Email,
		Password:   userForm.Password,
		Phone:      userForm.Phone,
		Role:       userForm.Role,
	}

	returnedUser, err := services.CreateUserServicePostgres(newUser)
	if err != nil {

		logger.Info("Intento de creación de usuario erroneo desde: " + c.ClientIP())

		response.JsonResponse(c, 500, err.Error(), newUser)
		return
	}

	logger.Info("Creacion exitosa de usuario " + newUser.Email + " desde ip: " + c.ClientIP())
	response.JsonResponse(c, 201, "Usuario creado con exito.", returnedUser)
}

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
