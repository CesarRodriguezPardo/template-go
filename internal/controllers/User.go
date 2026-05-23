package controllers

import (
	"net/http"
	"CesarRodriguezPardo/template-go/internal/dto"
	"CesarRodriguezPardo/template-go/internal/middleware"
	"CesarRodriguezPardo/template-go/internal/services"
	"CesarRodriguezPardo/template-go/utils"

	logger "CesarRodriguezPardo/template-go/infra/logger"
	response "CesarRodriguezPardo/template-go/infra/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	uuid "github.com/satori/go.uuid"
)

const (
	createdUserHtml = "mailer/templates/confirmCreatedUser.html"
)

// CreateUser
// @Title CreateUser
// @Description Permite crear un usuario en el sistema
// @Summary Crea un usuario
// @Tags Usuario
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} dto.UserResponse "Usuario creado con exito."
// @Router /user [post]
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JsonResponse(c, http.StatusBadRequest, "invalid user data", nil)
		return
	}

	returnedUser, err := services.CreateUser(c, &req)
	if err != nil {
		logger.Error("failed user creation attempt", err, zap.String("ip", c.ClientIP()))
		response.JsonResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logger.Info("Created user", zap.String("email", req.Email), zap.String("ip", c.ClientIP()))
	response.JsonResponse(c, http.StatusCreated, "user created", returnedUser)
}

// GetAllUsers
// @Title GetAllUsers
// @Description Permite obtener todos los usuarios del sistema
// @Summary Obtiene todos los usuarios
// @Tags Usuario
// @Produce json
// @Param page query int false "Numero de pagina"
// @Param limit query int false "Cantidad de items por pagina"
// @Success 200 {object} response.PaginatedData "Usuarios obtenidos con exito."
// @Router /user [get]
func GetAllUsers(c *gin.Context) {
	pagination := utils.GetPaginationFromContext(c)

	users, total, err := services.GetAllUsers(c, pagination.Limit, pagination.Offset)

	if err != nil {
		response.JsonResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.PaginatedJsonResponse(c, http.StatusOK, "Usuarios obtenidos con exito.", users, pagination.Page, pagination.Limit, total)
}

// GetUserByID
// @Title GetUserByID
// @Description Permite obtener un usuario especifico del sistema
// @Summary Obtiene un usuario por ID
// @Tags Usuario
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} dto.UserResponse "Usuario obtenido con exito."
// @Router /user/{id} [get]
func GetUserByID(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := uuid.FromString(targetIDStr)
	if err != nil {
		response.JsonResponse(c, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	user, err := services.GetUserByID(c, targetID)
	if err != nil {
		response.JsonResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.JsonResponse(c, http.StatusOK, "Usuario obtenido con exito.", user)
}

// UpdateUser
// @Title UpdateUser
// @Description Permite actualizar un usuario en el sistema
// @Summary Actualiza un usuario
// @Tags Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param user body dto.CreateUserRequest true "Datos del usuario"
// @Success 200 {object} dto.UserResponse "Usuario actualizado con exito."
// @Router /user/{id} [put]
func UpdateUser(c *gin.Context) {
	targetUserIDStr := c.Param("id")
	targetUserID, err := uuid.FromString(targetUserIDStr)
	if err != nil {
		response.JsonResponse(c, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	userMap, ok := c.Get("user")
	if !ok {
		response.JsonResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	claims, ok := userMap.(middleware.UserClaims)
	var requesterID uuid.UUID
	var requesterRole string
	if ok {
		requesterID = claims.ID
		requesterRole = claims.Role
	} else {
		claimsMap, ok := userMap.(map[string]interface{})
		if !ok {
			response.JsonResponse(c, http.StatusUnauthorized, "invalid token claims", nil)
			return
		}
		requesterIDStr, _ := claimsMap["id"].(string)
		requesterRole, _ = claimsMap["role"].(string)
		requesterID, _ = uuid.FromString(requesterIDStr)
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JsonResponse(c, http.StatusBadRequest, "invalid user data", nil)
		return
	}

	returnedUser, err := services.UpdateUser(c, requesterRole, requesterID, targetUserID, &req)
	if err != nil {
		logger.Error("failed user update attempt", err, zap.String("ip", c.ClientIP()))
		response.JsonResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logger.Info("Updated user", zap.String("id", targetUserIDStr), zap.String("ip", c.ClientIP()))
	response.JsonResponse(c, http.StatusOK, "user updated", returnedUser)
}

// DeleteUser
// @Title DeleteUser
// @Description Permite eliminar un usuario en el sistema (soft delete)
// @Summary Elimina un usuario
// @Tags Usuario
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {string} string "user deleted"
// @Router /user/{id} [delete]
func DeleteUser(c *gin.Context) {
	targetUserIDStr := c.Param("id")
	targetUserID, err := uuid.FromString(targetUserIDStr)
	if err != nil {
		response.JsonResponse(c, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	userMap, ok := c.Get("user")
	if !ok {
		response.JsonResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	claims, ok := userMap.(middleware.UserClaims)
	var requesterRole string
	if ok {
		requesterRole = claims.Role
	} else {
		claimsMap, ok := userMap.(map[string]interface{})
		if !ok {
			response.JsonResponse(c, http.StatusUnauthorized, "invalid token claims", nil)
			return
		}
		requesterRole, _ = claimsMap["role"].(string)
	}

	err = services.DeleteUser(c, requesterRole, targetUserID)
	if err != nil {
		logger.Error("failed user delete attempt", err, zap.String("ip", c.ClientIP()))
		response.JsonResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logger.Info("Deleted user", zap.String("id", targetUserIDStr), zap.String("ip", c.ClientIP()))
	response.JsonResponse(c, http.StatusOK, "user deleted", nil)
}
