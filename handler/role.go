package handler

import (
	"godemo/database"
	"godemo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

// ListRoles godoc
// @Summary 获取角色列表
// @Description 获取所有角色列表
// @Tags 角色
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Role
// @Router /api/roles [get]
func ListRoles(c *gin.Context) {
	var roles []models.Role
	database.DB.Preload("Permissions").Find(&roles)
	c.JSON(http.StatusOK, roles)
}

// GetRole godoc
// @Summary 获取角色信息
// @Description 根据ID获取角色详细信息
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Security Bearer
// @Success 200 {object} models.Role
// @Router /api/roles/{id} [get]
func GetRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// CreateRole godoc
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param request body CreateRoleRequest true "角色信息"
// @Security Bearer
// @Success 201 {object} models.Role
// @Router /api/roles [post]
func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建角色失败"})
		return
	}

	if len(req.PermissionIDs) > 0 {
		var permissions []models.Permission
		database.DB.Find(&permissions, req.PermissionIDs)
		database.DB.Model(&role).Association("Permissions").Replace(permissions)
	}

	database.DB.Preload("Permissions").First(&role, role.ID)
	c.JSON(http.StatusCreated, role)
}

// UpdateRole godoc
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body CreateRoleRequest true "角色信息"
// @Security Bearer
// @Success 200 {object} models.Role
// @Router /api/roles/{id} [put]
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.Name = req.Name
	role.Description = req.Description

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新角色失败"})
		return
	}

	if req.PermissionIDs != nil {
		var permissions []models.Permission
		database.DB.Find(&permissions, req.PermissionIDs)
		database.DB.Model(&role).Association("Permissions").Replace(permissions)
	}

	database.DB.Preload("Permissions").First(&role, role.ID)
	c.JSON(http.StatusOK, role)
}

// DeleteRole godoc
// @Summary 删除角色
// @Description 删除角色（软删除）
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
