package seed

import (
	"log"
	"todolist/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var user models.User
	// 檢查 admin 是否已存在
	if err := db.Where("account = ?", "admin").First(&user).Error; err == nil {
		log.Println("🟡 admin 使用者已存在，跳過建立")
		return
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ 密碼加密失敗: %v", err)
	}

	// 取得 admin 角色
	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		log.Fatalf("❌ 找不到 Admin 角色，請先執行角色 seed: %v", err)
	}

	// 建立使用者
	newUser := models.User{
		Account:  "admin",
		Password: string(hashedPassword),
		Roles:    []models.Role{adminRole},
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("❌ 建立 admin 使用者失敗: %v", err)
	}

	log.Println("✅ 管理員帳號 admin 建立完成")
}
