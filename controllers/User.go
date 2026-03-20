package controllers

import (
	"citiaps/golang-backend-template/forms"
	"citiaps/golang-backend-template/mailer"
	"citiaps/golang-backend-template/models"
	"citiaps/golang-backend-template/services"
	"citiaps/golang-backend-template/utils"

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
	var newUserFormsPostgres *forms.UserFormPostgres
	err := c.BindJSON(&newUserFormsPostgres)

	if err != nil {
		utils.JsonResponse(c, 500, "Error en los datos del usuario.", newUserFormsPostgres)
		return
	}

	newUser := &models.UserPostgres{
		Name:       newUserFormsPostgres.Name,
		MiddleName: newUserFormsPostgres.MiddleName,
		Email:      newUserFormsPostgres.Email,
		Password:   newUserFormsPostgres.Password,
		Phone:      newUserFormsPostgres.Phone,
		Roles:      newUserFormsPostgres.Roles,
		ActiveRol:  newUserFormsPostgres.ActiveRol,
	}

	returnedUser, err := services.CreateUserServicePostgres(newUser)
	if err != nil {

		utils.Info("Intento de creación de usuario erroneo desde: " + c.ClientIP())

		utils.JsonResponse(c, 500, err.Error(), newUser)
		return
	}

	utils.Info("Creacion exitosa de usuario " + newUser.Email + " desde ip: " + c.ClientIP())
	utils.JsonResponse(c, 201, "Usuario creado con exito.", returnedUser)
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
		utils.JsonResponse(c, 500, "Error al obtener todos los usuarios.", returnedUsers)
		return
	}
	if returnedUsers == nil {
		utils.JsonResponse(c, 200, "No se han encontrado usuarios.", returnedUsers)
		return
	}
	utils.JsonResponse(c, 200, "Usuarios obtenidos con exito.", returnedUsers)
}
