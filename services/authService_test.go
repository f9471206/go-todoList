package services

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用來模擬 GORM DB
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestLogin_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	// 模擬 GORM 的 SELECT SQL（含 soft delete 過濾與 LIMIT）
	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `users` WHERE account = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"),
	).
		WithArgs("admin", 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "account", "password"}).
				AddRow(1, "admin", string(hashedPwd)),
		)

	auth := &AuthService{ctx: context.TODO()}
	token, err := auth.Login(db, "admin", "123456")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 檢查 JWT 格式與內容
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	assert.NoError(t, err)

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		assert.Equal(t, float64(1), claims["user_id"])
	} else {
		t.Fatal("JWT parsing failed")
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_AccountNotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `users` WHERE account = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"),
	).
		WithArgs("notfound", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account", "password"})) // 空結果

	auth := &AuthService{ctx: context.TODO()}
	token, err := auth.Login(db, "notfound", "123")

	assert.EqualError(t, err, "account not found")
	assert.Equal(t, "", token)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_WrongPassword(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// 正確的密碼 hash
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `users` WHERE account = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"),
	).
		WithArgs("admin", 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "account", "password"}).
				AddRow(2, "admin", string(hashedPwd)),
		)

	auth := &AuthService{ctx: context.TODO()}
	token, err := auth.Login(db, "admin", "wrongpassword")

	assert.EqualError(t, err, "incorrect password")
	assert.Equal(t, "", token)

	assert.NoError(t, mock.ExpectationsWereMet())
}
