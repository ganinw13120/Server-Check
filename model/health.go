package model

import "time"

type Health struct {
	Path         string
	IsAlive      bool
	ResponseTime time.Duration
}
