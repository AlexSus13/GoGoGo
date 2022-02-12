package config

import (
	"os"
)

type DB struct {
	User     string
	Host     string
	Port     string
	Password string
	Dbname   string
	Sslmode  string
}

type Conf struct {
	Host        string
	Port        string
	DB          DB
	PathToFile  string
	KeyPassword string
	KeyToken    string
}

func Get() *Conf {

	dconf := Conf{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
		DB: DB{
			User:     os.Getenv("POSTGRES_USER"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Dbname:   os.Getenv("POSTGRES_DB"),
			Sslmode:  os.Getenv("Sslmode_DB"),
		},
		PathToFile:  os.Getenv("PATHTOSAVEFILE"),
		KeyPassword: os.Getenv("KEYPASSWORD"),
		KeyToken:    os.Getenv("KEYTOKEN"),
	}
	return &dconf
}
