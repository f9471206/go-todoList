package repositories

import (
	"todolist/models"
	"todolist/repositories/base"
)

type TodoListDetailsRepository struct {
	*base.BaseRepository[*models.TodoListDetails]
}

func NewTodoListDetailsRepository() *TodoListDetailsRepository {
	return &TodoListDetailsRepository{
		BaseRepository: base.NewBaseRepository[*models.TodoListDetails](),
	}
}
