package repository

import (
	"net/http"
	"server-health/model"

	"github.com/bradhe/stopwatch"
	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

type IHealthRepository interface {
	CheckHealth(path string) *model.Health
}

func NewHealthRepository(db *gorm.DB) healthRepository {
	return healthRepository{
		db: db,
	}
}

func (h healthRepository) CheckHealth(path string) *model.Health {
	watch := stopwatch.Start()
	_, err := http.Get(path)
	watch.Stop()
	health := model.Health{
		Path:         path,
		IsAlive:      err == nil,
		ResponseTime: watch.Milliseconds(),
	}
	return &health
}
