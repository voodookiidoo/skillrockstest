package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"skillrockstest/internal/dto"
	"skillrockstest/internal/repository"
	"skillrockstest/pkg/logger"
)

type App struct {
	repo *repository.Repository
	lg   *zap.Logger
}

func NewApp(conn *pgx.Conn) *App {
	return &App{repo: repository.NewRepository(conn), lg: logger.DefaultLogger()}
}

func (a *App) Get(c *fiber.Ctx) error {
	tasks, err := a.repo.GetAll(c.Context())
	if err != nil {
		a.lg.Error(err.Error())
		return c.SendStatus(500)
	}
	for ind, task := range tasks {
		if ind != 1 {
			c.WriteString(", ")
		}
		b := make([]byte, 0)
		if b, err = task.MarshalJSON(); err != nil {
			a.lg.Error(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Write(b)
	}
	return c.SendStatus(fiber.StatusOK)
}
func (a *App) Post(c *fiber.Ctx) error {
	req := new(dto.TaskRequest)
	if err := req.UnmarshalJSON(c.Body()); err != nil {
		if _, err = c.Status(fiber.StatusBadRequest).WriteString("invalid data format"); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return err
	}
	
	if err := req.Validate(); err != nil {
		if _, err = c.Status(fiber.StatusBadRequest).WriteString(err.Error()); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return err
	}
	if err := a.repo.CreateTask(c.Context(), *req); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (a *App) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	deleted, err := a.repo.DeleteTask(c.Context(), id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if deleted == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (a *App) Put(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	req := new(dto.TaskRequest)
	if err = req.UnmarshalJSON(c.Body()); err != nil {
		if _, err = c.Status(fiber.StatusBadRequest).WriteString("invalid data format"); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return err
	}
	if err = req.Validate(); err != nil {
		if _, err = c.Status(fiber.StatusBadRequest).WriteString(err.Error()); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return err
	}
	affected, err := a.repo.UpdateTask(c.Context(), *req, id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if affected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
