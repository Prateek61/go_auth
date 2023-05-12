package postgres

import (
	"github.com/Prateek61/go_auth/graph/model"
	"github.com/go-pg/pg/v10"
)

type TodosRepo struct {
	DB *pg.DB
}

func (t *TodosRepo) GetTodos() ([]*model.Todo, error) {
	var todos []*model.Todo
	err := t.DB.Model(&todos).Select()

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *TodosRepo) CreateTodo(todo *model.Todo) error {
	_, err := t.DB.Model(todo).Insert()

	if err != nil {
		return err
	}

	return nil
}

func (t *TodosRepo) GetTodosByUserID(userID string) ([]*model.Todo, error) {
	var todos []*model.Todo
	err := t.DB.Model(&todos).Where("user_id = ?", userID).Select()

	if err != nil {
		return nil, err
	}

	return todos, nil
}