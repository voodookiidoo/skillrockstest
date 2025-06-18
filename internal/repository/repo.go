package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"skillrockstest/internal/dto"
	"time"
)

type Repository struct {
	db           *pgx.Conn
	postTaskStmt *pgconn.StatementDescription
	updTaskStmt  *pgconn.StatementDescription
}

func NewRepository(db *pgx.Conn) *Repository {
	postStmt, err := db.Prepare(context.Background(), "newtask",
		"INSERT INTO tasks(title, description, status, created_at, updated_at) values ($1, $2, $3, $4, $5)",
	)
	if err != nil {
		panic(err)
	}
	updStmt, err := db.Prepare(context.Background(), "updtask",
		"UPDATE tasks SET title = $1,description = $2, status = $3, updated_at = $4 WHERE id = $5")
	if err != nil {
		panic(err)
	}
	
	return &Repository{
		db:           db,
		postTaskStmt: postStmt,
		updTaskStmt:  updStmt,
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]dto.Task, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make([]dto.Task, 0, 1024)
	for rows.Next() {
		t := new(dto.Task)
		if err = rows.Scan(
			&t.Id,
			&t.Title,
			&t.Desc,
			&t.Status,
			&t.Created,
			&t.Updated,
		); err != nil {
			return nil, err
		}
		data = append(data, *t)
	}
	return data, nil
}
func (r *Repository) CreateTask(ctx context.Context, req dto.TaskRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	now := time.Now()
	if _, err = tx.Exec(ctx, r.postTaskStmt.Name, req.Title, req.Desc, req.Status, now, now); err != nil {
		
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) DeleteTask(ctx context.Context, id int) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	
	affected, err := tx.Exec(ctx, "DELETE FROM tasks where id = $1", id)
	if err != nil {
		return 0, err
	}
	
	return affected.RowsAffected(), tx.Commit(ctx)
}

func (r *Repository) UpdateTask(ctx context.Context, req dto.TaskRequest, id int) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	now := time.Now()
	
	affected, err := tx.Exec(ctx, r.updTaskStmt.Name, req.Title, req.Desc, req.Status, now, id)
	if err != nil {
		return 0, err
	}
	
	return affected.RowsAffected(), tx.Commit(ctx)
}
