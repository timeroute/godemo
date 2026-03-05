package handler

import (
	"godemo/database"
	"godemo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListPermissions godoc
// @Summary 获取权限列表
// @Description 获取所有权限列表
// @Tags 权限
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Permission
// @Router /api/permissions [get]
func ListPermissions(c *gin.Context) {
	var permissions []models.Permission
	database.DB.Find(&permissions)
	c.JSON(http.StatusOK, permissions)
}
