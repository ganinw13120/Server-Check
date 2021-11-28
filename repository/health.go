package repository

import (
	"net/http"
	"time"

	"github.com/bradhe/stopwatch"
	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

type IHealthRepository interface {
	CheckHealth(path string) (bool, time.Duration)
}

func NewHealthRepository(db *gorm.DB) healthRepository {
	return healthRepository{
		db: db,
	}
}

func (h healthRepository) CheckHealth(path string) (bool, time.Duration) {
	watch := stopwatch.Start()
	_, err := http.Get(path)
	if err != nil {
		return false, 0
	}
	return true, watch.Milliseconds()
}
