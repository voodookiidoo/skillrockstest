package repository

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"skillrockstest/internal/dto"
	"skillrockstest/pkg/db"
	"strconv"
	"time"
)

type Repository struct {
	db       *pgx.Conn
	rd       *redis.Client
	migrator *migrate.Migrate
}

func (r *Repository) Close() error {
	return errors.Join(
		r.db.Close(context.Background()),
		r.migrator.Down(),
		r.rd.Close())
}

func NewRepository() *Repository {
	pg, migrator := db.MustConnect()
	rd := db.MustConnectRedis()
	return &Repository{
		db:       pg,
		migrator: migrator,
		rd:       rd,
	}
}

func (r *Repository) updateCache(ctx context.Context, task dto.Task, id int) error {
	tosave := toRedis(&task)
	if err := r.rd.HSet(ctx, strconv.Itoa(id), tosave).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]dto.Task, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (dto.Task, error) {
		return pgx.RowToStructByName[dto.Task](row)
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (r *Repository) CreateTask(ctx context.Context, req dto.TaskRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)
	now := time.Now()
	var id int
	row := tx.QueryRow(ctx,
		"INSERT INTO tasks(title, description, status, created_at, updated_at) values ($1, $2, $3, $4, $5) returning id",
		req.Title, req.Desc, req.Status, now, now,
	)
	if err = row.Scan(&id); err != nil {
		return err
	}
	task := dto.Task{
		Title:   req.Title,
		Desc:    req.Desc,
		Status:  req.Status,
		Created: now,
		Updated: now,
	}
	r.updateCache(ctx, task, id)

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
	r.clearCache(ctx, id)

	return affected.RowsAffected(), tx.Commit(ctx)
}

func (r *Repository) UpdateTask(ctx context.Context, req dto.TaskRequest, id int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	now := time.Now()

	row := tx.QueryRow(ctx,
		"UPDATE tasks SET title = $1,description = $2, status = $3, updated_at = $4 WHERE id = $5 returning created_at",
		req.Title, req.Desc, req.Status, now, id)
	created := time.Time{}
	if err = row.Scan(&created); err != nil {
		return err
	}
	task := dto.Task{
		Title:   req.Title,
		Desc:    req.Desc,
		Status:  req.Status,
		Created: created,
		Updated: now,
	}
	r.updateCache(ctx, task, id)
	return tx.Commit(ctx)
}

func (r *Repository) Get(c context.Context, id int) (*dto.Task, error) {
	res := r.rd.HMGet(c, strconv.Itoa(id), "title", "desc", "status", "c_at", "upd_at")
	if err := res.Err(); err != nil {
		return nil, res.Err()
	}
	task, err := RedisExtractor(res)
	if err == nil {
		return task, err
	}
	q := r.db.QueryRow(c, "SELECT * FROM tasks WHERE id = $1 LIMIT 1", id)
	t := new(dto.Task)

	if err = q.Scan(
		&t.Id,
		&t.Title,
		&t.Desc,
		&t.Status,
		&t.Created,
		&t.Updated); err != nil {
		return nil, err
	}
	return t, err

}

func (r *Repository) clearCache(ctx context.Context, id int) error {
	return r.rd.Del(ctx, strconv.Itoa(id)).Err()
}

type redisParser struct {
	Title   string    `redis:"title"`
	Desc    *string   `redis:"desc,omitempty"`
	Status  *string   `redis:"status"`
	Created time.Time `redis:"c_at"`
	Updated time.Time `redis:"upd_at"`
}

func toRedis(task *dto.Task) redisParser {
	return redisParser{
		Title:   task.Title,
		Desc:    task.Desc,
		Status:  task.Status,
		Created: task.Created,
		Updated: task.Created,
	}
}

func toTask(p *redisParser) dto.Task {
	return dto.Task{
		Title:   p.Title,
		Desc:    p.Desc,
		Status:  p.Status,
		Created: p.Created,
		Updated: p.Updated,
	}
}

func RedisExtractor(cmd *redis.SliceCmd) (*dto.Task, error) {
	redOut := new(redisParser)
	if err := cmd.Scan(&redOut); err != nil {
		return nil, err
	}
	t := &dto.Task{
		Title:   redOut.Title,
		Desc:    redOut.Desc,
		Status:  redOut.Status,
		Created: redOut.Created,
		Updated: redOut.Updated,
	}
	return t, nil

}
