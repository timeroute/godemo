package database

import (
	"godemo/config"
	"godemo/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("✅ 数据库连接成功:", config.AppConfig.Database.DSN)

	// 自动迁移
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.RequestLog{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化默认数据
	initDefaultData()
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func initDefaultData() {
	// 创建默认权限
	permissions := []models.Permission{
		{Name: "查看用户", Code: "user:view", Description: "查看用户列表和详情"},
		{Name: "创建用户", Code: "user:create", Description: "创建新用户"},
		{Name: "编辑用户", Code: "user:edit", Description: "编辑用户信息"},
		{Name: "删除用户", Code: "user:delete", Description: "删除用户"},
		{Name: "查看角色", Code: "role:view", Description: "查看角色列表和详情"},
		{Name: "创建角色", Code: "role:create", Description: "创建新角色"},
		{Name: "编辑角色", Code: "role:edit", Description: "编辑角色信息"},
		{Name: "删除角色", Code: "role:delete", Description: "删除角色"},
	}

	for _, perm := range permissions {
		DB.Where(models.Permission{Code: perm.Code}).FirstOrCreate(&perm)
	}

	// 创建默认角色
	var adminRole models.Role
	DB.Where(models.Role{Name: "管理员"}).FirstOrCreate(&adminRole, models.Role{
		Name:        "管理员",
		Description: "系统管理员，拥有所有权限",
	})

	// 为管理员角色分配所有权限
	var allPerms []models.Permission
	DB.Find(&allPerms)
	DB.Model(&adminRole).Association("Permissions").Replace(allPerms)

	// 创建默认管理员用户
	var adminUser models.User
	result := DB.Where(models.User{Username: "admin"}).First(&adminUser)
	if result.Error != nil {
		// 用户不存在，创建新用户
		hashedPassword, _ := hashPassword("admin123")
		adminUser = models.User{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@example.com",
			Status:   1,
		}
		DB.Create(&adminUser)
	}

	// 为管理员用户分配管理员角色
	DB.Model(&adminUser).Association("Roles").Replace([]models.Role{adminRole})
}
