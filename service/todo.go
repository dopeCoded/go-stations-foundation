package service

import (
	"context"
	"database/sql"
	"time"
	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// TODO を DB に挿入
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err // 挿入時のエラーをそのまま返す
	}

	// 挿入された TODO の ID を取得
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err // ID の取得時のエラーをそのまま返す
	}

	// 挿入された TODO を取得
	row := s.db.QueryRowContext(ctx, confirm, id)

	var fetchedSubject, fetchedDescription string
	var createdAt, updatedAt time.Time

	err = row.Scan(&fetchedSubject, &fetchedDescription, &createdAt, &updatedAt)
	if err != nil {
		return nil, err // 取得時のエラーをそのまま返す
	}

	// TODO 構造体を作成
	todo := &model.TODO{
		ID:          int(id),
		Subject:     fetchedSubject,
		Description: fetchedDescription,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	return nil, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
