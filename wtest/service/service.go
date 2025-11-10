// service/service.go
package service

import (
	"databaseAi/wtest/db"
	"log"
)

type Service struct {
	DB  *db.DB
	Log *log.Logger
}

func NewService(d *db.DB, l *log.Logger) *Service {
	return &Service{DB: d, Log: l}
}
