package app

import (
	"GoHttpDbServer/config"
	"database/sql"
	"log"
)

type App struct {
	Db       *sql.DB
	errorLog *log.Logger
	infoLog  *log.Logger
	Config   *config.Conf
}

func NewApp(db *sql.DB, errorLog *log.Logger, infoLog *log.Logger, Conf *config.Conf) *App {
	return &App{
		Db:       db,
		errorLog: errorLog,
		infoLog:  infoLog,
		Config:   Conf,
	}
}
