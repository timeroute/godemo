package service

import (
	"errors"
	"godemo/config"
	"godemo/database"
	"godemo/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(username, password string) (string, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return "", errors.New("用户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken(userID uint, username string) (string, error) {
	expireHours := time.Duration(config.AppConfig.JWT.ExpireHour) * time.Hour

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPermission(userID uint, permissionCode string) bool {
	var user models.User
	if err := database.DB.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return false
	}

	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Code == permissionCode {
				return true
			}
		}
	}

	return false
}
