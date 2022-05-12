package main

import (
	"net/http"

	"github.com/joao-fontenele/s-rank-fiis/internal/handler"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	logger  *zap.Logger
	Version = "develop"
	DBURL   = "postgresql://usr:pwd@localhost:5432/ranks" // TODO: get from env variable
)

func main() {
	logger, _ = zap.NewProduction()

	dbGorm, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
	}

	r := repository.Stock{Conn: dbGorm}
	h := http.NewServeMux()
	h.Handle("/", handler.New(logger, r))

	logger.Info("listening on port :9999", zap.String("version", Version))
	logger.Fatal("error starting server", zap.Error(http.ListenAndServe(":9999", h)))
}
