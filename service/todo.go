package service

import (
	"context"
	"database/sql"
	"fmt"
    "strings"
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
		ID:          int64(id),
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

	// size が 0 以下の場合、空スライスを返す
	if size <= 0 {
			return []*model.TODO{}, nil
	}

	var rows *sql.Rows
	var err error

	// prevID の有無でクエリを選択
	if prevID > 0 {
			rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
			rows, err = s.db.QueryContext(ctx, read, size)
	}
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var todos []*model.TODO

	// 取得した行をスキャンしてスライスに追加
	for rows.Next() {
			var todo model.TODO
			if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
					return nil, err
			}
			todos = append(todos, &todo)
	}

	// イテレーション中にエラーが発生したか確認
	if err := rows.Err(); err != nil {
			return nil, err
	}

	// 空の結果の場合は空スライスを返す
	if todos == nil {
			todos = []*model.TODO{}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
			update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
			confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// トランザクションの開始
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
			return nil, err
	}
	defer tx.Rollback()

	// Prepared Statement の作成 (UPDATE)
	stmtUpdate, err := tx.PrepareContext(ctx, update)
	if err != nil {
			return nil, err
	}
	defer stmtUpdate.Close()

	// UPDATE クエリの実行
	result, err := stmtUpdate.ExecContext(ctx, subject, description, id)
	if err != nil {
			return nil, err
	}

	// 更新された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
			return nil, err
	}
	if rowsAffected == 0 {
			return nil, &model.ErrNotFound{}
	}

	// Prepared Statement の作成 (SELECT)
	stmtSelect, err := tx.PrepareContext(ctx, confirm)
	if err != nil {
			return nil, err
	}
	defer stmtSelect.Close()

	// 更新された TODO の取得
	row := stmtSelect.QueryRowContext(ctx, id)
	var todo model.TODO
	if err := row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
	}

	// トランザクションのコミット
	if err := tx.Commit(); err != nil {
			return nil, err
	}

	// 更新された TODO を返す
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
    // ids が空の場合、何もせずに nil を返す
    if len(ids) == 0 {
        return nil
    }

    // ids の数だけ '?' を作成し、カンマで区切る
    placeholders := strings.TrimLeft(strings.Repeat(",?", len(ids)), ",")

    // DELETE クエリを作成
    query := fmt.Sprintf("DELETE FROM todos WHERE id IN (%s)", placeholders)

    // ids を []interface{} 型に変換
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        args[i] = id
    }

    // クエリを実行
    result, err := s.db.ExecContext(ctx, query, args...)
    if err != nil {
        return err
    }

    // 削除された行数を取得
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    // 削除された行数が 0 の場合、ErrNotFound を返す
    if rowsAffected == 0 {
        return &model.ErrNotFound{}
    }

    // 正常に削除された場合は nil を返す
    return nil
}
