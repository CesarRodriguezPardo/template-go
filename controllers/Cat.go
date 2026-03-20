package controllers

import (
	"citiaps/golang-backend-template/forms"
	"citiaps/golang-backend-template/middleware"
	"citiaps/golang-backend-template/models"
	"citiaps/golang-backend-template/services"
	"citiaps/golang-backend-template/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCatController
// @Title CreateCatController
// @Description Permite crear un gato en el sistema
// @Summary Crea un gato
// @Tags Gato
// @Accept json
// @Produce json
// @Success 200 {object} forms.CatForm "Gato creado con exito."
// @Router /cat/ [post]
func CreateCatController(c *gin.Context) {
	var newFormCat *forms.CatForm
	err := c.BindJSON(&newFormCat)

	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), newFormCat)
		return
	}

	userClaims := middleware.IdentityHandlerFunc(c)
	if userClaims == nil {
		utils.JsonResponse(c, 403, "Forbidden", nil)
		return
	}

	userMap, _ := userClaims.(map[string]interface{})
	userIDStr, _ := userMap["_id"].(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)

	newCat := &models.Cat{
		Name:  newFormCat.Name,
		Age:   newFormCat.Age,
		Owner: userID,
	}

	returnedCat, err := services.CreateCatService(newCat)
	if err != nil {
		utils.Info("Intento de creación de gato fallido del usuario con id: " + userIDStr + " desde ip: " + c.ClientIP())
		utils.JsonResponse(c, 500, err.Error(), newCat)
		return
	}

	utils.Info("Creacion exitosa de gato del usuario con id " + userIDStr + " desde ip: " + c.ClientIP())
	utils.JsonResponse(c, 201, "Gato creado con exito.", returnedCat)
}

// GetAllCatsController
// @Title GetAllCatsController
// @Description Permite obtener todos los gatos del sistema
// @Summary Obtiene todos los gatos
// @Tags Gato
// @Accept json
// @Produce json
// @Success 200 {object} []models.Cat "Gatos obtenidos con exito."
// @Router /cat/postgres [get]
func GetAllCatsController(c *gin.Context) {
	returnedCats, err := services.FindAllCatService()
	if err != nil {
		utils.JsonResponse(c, 500, "Error al obtener todos los gatos:"+err.Error(), returnedCats)
		return
	}
	if returnedCats == nil {
		utils.JsonResponse(c, 200, "No se han encontrado gatitos :(.", returnedCats)
		return
	}
	utils.JsonResponse(c, 200, "Gatos obtenidos con exito.", returnedCats)
}

// postgres

// CreateCatControllerPostgres
// @Title CreateCatControllerPostgres
// @Description Permite crear un gato en el sistema
// @Summary Crea un gato
// @Tags Gato
// @Accept json
// @Produce json
// @Success 200 {object} forms.CatForm "Gato creado con exito."
// @Router /cat/postgres [post]
func CreateCatControllerPostgres(c *gin.Context) {
	var newFormCat *forms.CatForm
	err := c.BindJSON(&newFormCat)

	if err != nil {
		utils.JsonResponse(c, 500, "Error en los datos del gato.", newFormCat)
		return
	}

	userClaims := middleware.IdentityHandlerFunc(c)
	if userClaims == nil {
		utils.JsonResponse(c, 401, "No autorizado: usuario no encontrado en el token.", nil)
		return
	}

	userMap, _ := userClaims.(map[string]interface{})
	userIDUint, _ := userMap["_id"].(uint)
	userIDStr, _ := userMap["_id"].(string)

	newCat := &models.CatPostgres{
		Name:  newFormCat.Name,
		Age:   newFormCat.Age,
		Owner: userIDUint,
	}

	returnedCat, err := services.CreateCatServicePostgres(newCat)
	if err != nil {
		utils.Info("Intento de creación de gato fallido del usuario con id: " + userIDStr + " desde ip: " + c.ClientIP())
		utils.JsonResponse(c, 500, err.Error(), newCat)
		return
	}

	utils.Info("Creacion exitosa de gato del usuario con id " + userIDStr + " desde ip: " + c.ClientIP())
	utils.JsonResponse(c, 201, "Gato creado con exito.", returnedCat)
}

// GetAllCatsControllerPostgres
// @Title GetAllCatsControllerPostgres
// @Description Permite obtener todos los gatos del sistema
// @Summary Obtiene todos los gatos
// @Tags Gato
// @Accept json
// @Produce json
// @Success 200 "Gatos obtenidos con exito."
// @Router /cat/ [get]
func GetAllCatsControllerPostgres(c *gin.Context) {
	returnedCats, err := services.FindAllCatServicePostgres()
	if err != nil {
		utils.JsonResponse(c, 500, "Error al obtener todos los gatos.", returnedCats)
		return
	}
	if returnedCats == nil {
		utils.JsonResponse(c, 200, "No se han encontrado gatitos :(.", returnedCats)
		return
	}
	utils.JsonResponse(c, 200, "Gatos obtenidos con exito.", returnedCats)
}
