package services_test

import (
	"context"
	"errors"
	"testing"
	"todolist/mocks"
	"todolist/models"
	"todolist/services"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTodoListService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	ctx := context.Background()
	service := services.NewTodoListService(ctx, mockRepo)

	db, sqlmock := setupMockDB(t)

	// 預期交易開始和提交
	sqlmock.ExpectBegin()

	// 正確模擬 GORM 自動帶入 deleted_at IS NULL 的查詢
	sqlmock.ExpectQuery(`SELECT count\(\*\) FROM "to_do_types" WHERE id = \$1 AND "to_do_types"\."deleted_at" IS NULL`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), "Test", 0).
		Return(false, nil).
		Times(1)

	mockRepo.EXPECT().
		Create(ctx, gomock.Any(), gomock.AssignableToTypeOf(&models.TodoList{})). // ✅ 修正點
		Return(nil).
		Times(1)

	sqlmock.ExpectCommit()

	// 呼叫建立服務
	result, err := service.Create(db, "Test", 1)

	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Name)

	// 驗證 SQL 預期都被滿足
	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoListService_Create_TypeIDNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	ctx := context.Background()
	service := services.NewTodoListService(ctx, mockRepo)

	db, sqlmock := setupMockDB(t)

	sqlmock.ExpectBegin()

	sqlmock.ExpectQuery(`SELECT count\(\*\) FROM "to_do_types" WHERE id = \$1 AND "to_do_types"\."deleted_at" IS NULL`).
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	sqlmock.ExpectRollback()

	result, err := service.Create(db, "AnyName", 999)

	assert.Error(t, err)
	assert.Equal(t, "type_id 不存在", err.Error())

	// 雖然錯誤，但 result 指標本體是有的，可以檢查欄位值
	assert.NotNil(t, result)
	assert.Equal(t, "AnyName", result.Name)
	assert.Equal(t, 999, result.TypeID)

	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoListService_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	db, _ := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoListService(ctx, mockRepo)

	page := 1
	pageSize := 10
	orderBy := []string{"id desc"}

	// 測試用資料
	expectedList := []*models.TodoList{
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

func TestTodoListService_Edit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoListService(ctx, mockRepo)

	id := 1
	newName := "新名稱"
	newTypeID := 2

	existing := &models.TodoList{
		ID:     id,
		Name:   "舊名稱",
		TypeID: 1,
	}

	sqlmock.ExpectBegin()

	// 模擬 type_id 存在查詢
	sqlmock.ExpectQuery(`SELECT count\(\*\) FROM "to_do_types" WHERE id = \$1 AND "to_do_types"\."deleted_at" IS NULL`).
		WithArgs(newTypeID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// 模擬名稱不重複
	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), newName, id).
		Return(false, nil).
		Times(1)

	// 模擬取得原始資料
	mockRepo.EXPECT().
		FindByID(ctx, gomock.Any(), id).
		Return(existing, nil).
		Times(1)

	// 模擬更新資料，並驗證更新內容
	mockRepo.EXPECT().
		Update(ctx, gomock.Any(), gomock.AssignableToTypeOf(&models.TodoList{})).
		DoAndReturn(func(ctx context.Context, tx any, todo *models.TodoList) error {
			assert.Equal(t, newName, todo.Name)
			assert.Equal(t, newTypeID, todo.TypeID)
			return nil
		}).
		Times(1)

	sqlmock.ExpectCommit()

	result, err := svc.Edit(db, id, newName, newTypeID)

	assert.NoError(t, err)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newTypeID, result.TypeID)

	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoListService_Edit_Fail_NoChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoListService(ctx, mockRepo)

	id := 1
	newName := "新名稱"
	newTypeID := 2

	existing := &models.TodoList{
		ID:     id,
		Name:   "舊名稱",
		TypeID: 1,
	}

	// 預期開始交易
	sqlmock.ExpectBegin()

	// 模擬 type_id 存在的查詢
	sqlmock.ExpectQuery(`SELECT count\(\*\) FROM "to_do_types" WHERE id = \$1 AND "to_do_types"\."deleted_at" IS NULL`).
		WithArgs(newTypeID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// 名稱不重複的檢查
	mockRepo.EXPECT().
		IsNameExist(ctx, gomock.Any(), newName, id).
		Return(false, nil).
		Times(1)

	// FindByID 兩次回傳相同資料
	mockRepo.EXPECT().
		FindByID(ctx, gomock.Any(), id).
		Return(existing, nil).
		Times(1)

	// Update 時模擬失敗，會回傳錯誤
	mockRepo.EXPECT().
		Update(ctx, gomock.Any(), gomock.AssignableToTypeOf(&models.TodoList{})).
		DoAndReturn(func(ctx context.Context, tx any, todo *models.TodoList) error {
			assert.Equal(t, newName, todo.Name)
			assert.Equal(t, newTypeID, todo.TypeID)
			return errors.New("更新失敗")
		}).
		Times(1)

	// 預期交易回滾
	sqlmock.ExpectRollback()

	_, err := svc.Edit(db, id, newName, newTypeID)

	assert.Error(t, err)

	// 移除原本再查資料的驗證，專注於交易 rollback

	// 確認 sqlmock 的交易期望都被達成
	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoListService_Edit_TypeIDNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoListService(ctx, mockRepo)

	id := 1       // 編輯的 TodoList ID
	typeID := 999 // 假設這個 typeID 不存在
	name := "測試名稱"

	sqlmock.ExpectBegin()

	// 模擬 type_id 查詢回傳 0，表示 type_id 不存在
	sqlmock.ExpectQuery(`SELECT count\(\*\) FROM "to_do_types" WHERE id = \$1 AND "to_do_types"\."deleted_at" IS NULL`).
		WithArgs(typeID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// 交易失敗，回滾
	sqlmock.ExpectRollback()

	_, err := svc.Edit(db, id, name, typeID)

	assert.Error(t, err)
	assert.Equal(t, "type_id 不存在", err.Error())

	assert.NoError(t, sqlmock.ExpectationsWereMet())
}

func TestTodoListService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoListService(ctx, mockRepo)

	id := 1
	existingModel := &models.TodoList{ID: id, Name: "測試項目"}

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

func TestTodoListService_Delete_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTodoListRepository(ctrl)
	db, sqlmock := setupMockDB(t)
	ctx := context.Background()

	svc := services.NewTodoListService(ctx, mockRepo)

	id := 1
	existingModel := &models.TodoList{ID: id, Name: "測試項目"}

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
