package dto

import (
	"errors"
	"slices"
	"time"
)

var statusEnum = []string{"new", "in_progress", "done"}

type Task struct {
	Id      int       `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Desc    *string   `json:"description,omitempty" db:"description"`
	Status  *string   `json:"status" db:"status"`
	Created time.Time `json:"created_at" db:"created_at"`
	Updated time.Time `json:"updated_at" db:"updated_at"`
}

type TaskRequest struct {
	Title  string  `json:"title"`
	Desc   *string `json:"desc,omitempty"`
	Status *string `json:"status"`
}

func (t *TaskRequest) Validate() error {
	if t.Status != nil && !slices.Contains(statusEnum, *t.Status) {
		return errors.New("task status is invalid")
	}
	return nil
}
