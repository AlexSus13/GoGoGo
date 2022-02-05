package app

import (
	"FileStorageServer/config"

	"github.com/sirupsen/logrus"

	"database/sql"
)

type App struct {
	Db       *sql.DB
	MyLogger *logrus.Logger
	Config   *config.Conf
}

func NewApp(db *sql.DB, MyLogger *logrus.Logger, Conf *config.Conf) *App {
	return &App{
		Db:       db,
		MyLogger: MyLogger,
		Config:   Conf,
	}
}
