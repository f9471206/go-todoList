package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todolist/config"
	"todolist/models"
	"todolist/repositories"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(config.JWTSecret)

type AuthService struct {
	ctx  context.Context
	repo *repositories.AuthRepository
}

func NewAuthService(ctx context.Context) *AuthService {
	return &AuthService{
		ctx:  ctx,
		repo: &repositories.AuthRepository{},
	}
}

func (s *AuthService) Login(db *gorm.DB, account, password string) (string, error) {
	user := &models.User{}

	// 查帳號
	if err := db.Where("account = ?", account).First(user).Error; err != nil {
		return "", errors.New("account not found")
	}

	// 密碼比對
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("incorrect password")
	}

	// 產生 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小時過期
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Register(db *gorm.DB, account, password string) (*models.User, error) {
	// 檢查帳號是否已存在
	var count int64
	db.Model(&models.User{}).Where("account = ?", account).Count(&count)
	if count > 0 {
		return nil, errors.New("account already exists")
	}

	// 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 查詢預設角色 (遊客)
	var guestRole models.Role
	if err := db.Where("name = ?", "Guest").First(&guestRole).Error; err != nil {
		return nil, fmt.Errorf("default role '遊客' not found %w", err)
	}

	user := &models.User{
		Account:  account,
		Password: string(hashedPassword),
		Roles:    []models.Role{guestRole},
	}

	// 新增到資料庫
	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
