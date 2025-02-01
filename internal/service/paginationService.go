package service

import (
	"time"

	"Tsniper/pkg/database"
	"Tsniper/pkg/logger"

	"github.com/robfig/cron/v3"
)

type paginationService struct {
	cron   *cron.Cron
	db     database.Database
	logger logger.Logger
}

func NewPaginationService(db database.Database, logger logger.Logger) *paginationService {
	paginationService := &paginationService{
		cron:   cron.New(),
		db:     db,
		logger: logger,
	}
	err := paginationService.migrate()
	if err != nil {
		logger.Fatal("Failed to migrate pagination tables: %v", err)
	}
	return paginationService
}

func (s *paginationService) migrate() error {
	err := s.db.ExecSQLFile("./pkg/database/migrations/pagination.sql")
	return err
}

func (s *paginationService) Start() error {
	s.logger.Info("Started pagination service")
	s.Clean()
	_, err := s.cron.AddFunc("0 * * * *", s.Clean)
	if err != nil {
		return err
	}
	s.cron.Start()
	return nil
}

func (s *paginationService) Stop() error {
	s.logger.Info("Stopped pagination service")
	s.cron.Stop()
	return nil
}

func (s *paginationService) Clean() {
	threshold := time.Now().Unix() - 1200
	s.db.Exec("DELETE FROM pagination WHERE date_time < ?", threshold)
}
