package handler

import (
	"godemo/database"
	"godemo/models"
	"godemo/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	RoleIDs  []uint `json:"role_ids"`
}

type UpdateUserRequest struct {
	Email   string `json:"email" binding:"omitempty,email"`
	Status  *int   `json:"status"`
	RoleIDs []uint `json:"role_ids"`
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 根据ID获取用户详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Security Bearer
// @Success 200 {object} models.User
// @Router /api/users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.Preload("Roles.Permissions").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取所有用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.User
// @Router /api/users [get]
func ListUsers(c *gin.Context) {
	var users []models.User
	database.DB.Preload("Roles").Find(&users)
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "用户信息"
// @Security Bearer
// @Success 201 {object} models.User
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		database.DB.Find(&roles, req.RoleIDs)
		database.DB.Model(&user).Association("Roles").Replace(roles)
	}

	database.DB.Preload("Roles").First(&user, user.ID)
	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateUserRequest true "用户信息"
// @Security Bearer
// @Success 200 {object} models.User
// @Router /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	if req.RoleIDs != nil {
		var roles []models.Role
		database.DB.Find(&roles, req.RoleIDs)
		database.DB.Model(&user).Association("Roles").Replace(roles)
	}

	database.DB.Preload("Roles").First(&user, user.ID)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除用户（软删除）
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, _ := strconv.ParseUint(id, 10, 32)

	if userID == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "不能删除管理员账户"})
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
