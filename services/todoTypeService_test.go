package services_test

import (
	"context"
	"errors"
	"testing"
	"todolist/mocks"
	"todolist/models"
	"todolist/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 建立模擬資料庫的輔助函式
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

func TestTodoTypeService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	ctx := context.Background()
	service := services.NewTodoTypeService(ctx, mockRepo)

	db, sqlmock := setupMockDB(t)

	// 預期交易開始和提交
	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), "Test", 0).
		Return(false, nil).
		Times(1)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any(), gomock.AssignableToTypeOf(&models.TodoTypes{})).
		Return(nil).
		Times(1)

	result, err := service.Create(db, "Test")

	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Name)

	// 驗證所有期望都有被呼叫
	err = sqlmock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTodoTypeService_Create_NameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	ctx := context.Background()
	service := services.NewTodoTypeService(ctx, mockRepo)

	db, sqlmock := setupMockDB(t)

	// 預期交易開始和提交
	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	name := "Test"

	// 預期 IsNameExist 回傳 true，表示名稱已存在
	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), name, 0).
		Return(true, nil).
		Times(1)

	// 這時候不會呼叫 Create，因此不設定 Create 的期待

	_, err := service.Create(db, "Test")

	assert.Error(t, err)               // 期待有錯誤，因為名稱已存在
	assert.EqualError(t, err, "名稱已存在") // 錯誤訊息
}

func TestTodoTypeService_Create_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	ctx := context.Background()
	service := services.NewTodoTypeService(ctx, mockRepo)

	db, sqlmock := setupMockDB(t)

	// 預期交易開始和提交
	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	name := "Test"

	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), name, 0).
		Return(false, nil).
		Times(1)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any(), gomock.AssignableToTypeOf(&models.TodoTypes{})).
		Return(errors.New("create failed")).
		Times(1)

	_, err := service.Create(db, name)

	assert.Error(t, err)
	assert.EqualError(t, err, "create failed")
}

func TestTodoTypeService_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	db, _ := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoTypeService(ctx, mockRepo)

	page := 1
	pageSize := 10
	orderBy := []string{"id desc"}

	// 測試用資料
	expectedList := []*models.TodoTypes{
		{ID: 1, Name: "分類A"},
		{ID: 2, Name: "分類B"},
	}
	expectedTotal := int64(2)

	// 模擬 repo 回傳
	mockRepo.EXPECT().
		FindAllWithQuery(ctx, gomock.Any(), page, pageSize, gomock.Any()).
		Return(expectedList, expectedTotal, nil).
		Times(1)

	// 執行 service
	result, err := svc.Index(db, "", page, pageSize, orderBy)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTotal, result.Total)
	assert.Len(t, result.Data, 2)

	assert.Equal(t, "分類A", result.Data[0].Name)
	assert.Equal(t, "分類B", result.Data[1].Name)
}

func TestTodoTypeService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoTypeService(ctx, mockRepo)

	id := 1
	newName := "新名稱"
	existing := &models.TodoTypes{ID: id, Name: "舊名稱"}

	// 預期行為
	mockRepo.EXPECT().
		FindByID(ctx, gomock.Any(), id).
		Return(existing, nil).
		Times(1)

	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), newName, id).
		Return(false, nil).
		Times(1)

	sqlmock.ExpectBegin()
	mockRepo.EXPECT().
		Update(ctx, gomock.Any(), existing).
		Return(nil).
		Times(1)
	sqlmock.ExpectCommit()

	// 執行
	result, err := svc.Edit(db, id, newName)

	assert.NoError(t, err)
	assert.Equal(t, newName, result.Name)
}

func TestTodoTypeService_Update_Fail_NoChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoTypeService(ctx, mockRepo)

	id := 1
	newName := "新名稱"
	existing := &models.TodoTypes{ID: id, Name: "舊名稱"}

	sqlmock.ExpectBegin()

	// 找到原本資料
	mockRepo.EXPECT().
		FindByID(ctx, gomock.Any(), id).
		Return(existing, nil).
		Times(1)

	// 名稱檢查，不存在重複名稱
	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), newName, id).
		Return(false, nil).
		Times(1)

	// 嘗試更新，但回傳錯誤
	mockRepo.EXPECT().
		Update(ctx, gomock.Any(), gomock.AssignableToTypeOf(&models.TodoTypes{})).
		DoAndReturn(func(ctx context.Context, tx any, todoType *models.TodoTypes) error {
			// 確認更新參數是正確的
			assert.Equal(t, newName, todoType.Name)
			// 模擬更新失敗
			return errors.New("更新失敗")
		}).
		Times(1)

	sqlmock.ExpectRollback()

	// 執行編輯，應該失敗
	_, err := svc.Edit(db, id, newName)

	assert.Error(t, err)

	// 確認 sqlmock 期待都滿足
	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoTypeService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoTypeService(ctx, mockRepo)

	id := 1
	existingModel := &models.TodoTypes{ID: id, Name: "測試項目"}

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	mockRepo.EXPECT().
		FindByID(ctx, gomock.Any(), id).
		Return(existingModel, nil).
		Times(1)

	mockRepo.EXPECT().
		SoftDelete(ctx, gomock.Any(), existingModel).
		Return(nil).
		Times(1)

	_, err := svc.Delete(db, id)

	assert.NoError(t, err)
	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoTypeService_Delete_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoTypeRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoTypeService(ctx, mockRepo)

	id := 1
	existingModel := &models.TodoTypes{ID: id, Name: "測試項目"}

	sqlmock.ExpectBegin()
	// 預期交易 rollback 而非 commit，因刪除失敗
	sqlmock.ExpectRollback()

	mockRepo.EXPECT().
		FindByID(ctx, gomock.Any(), id).
		Return(existingModel, nil).
		Times(1)

	mockRepo.EXPECT().
		SoftDelete(ctx, gomock.Any(), existingModel).
		Return(errors.New("刪除失敗")).
		Times(1)

	result, err := svc.Delete(db, id)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "刪除失敗", err.Error())

	// 確保 sqlmock 的所有預期呼叫都完成
	assert.NoError(t, sqlmock.ExpectationsWereMet())
}
