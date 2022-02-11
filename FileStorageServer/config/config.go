package config

import (
//	"github.com/go-yaml/yaml"
//	"github.com/pkg/errors"

//	"io/ioutil"
	"os"
)

type DB struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

type Conf struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	DB          DB     `yaml:"db"`
	PathToFile  string `yaml:"pathtofile"`
	KeyPassword string `yaml:"keypassword"`
	KeyToken    string `yaml:"keytoken"`
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
		PathToFile:  os.Getenv("PATHTOFILE"),
		KeyPassword: os.Getenv("KEYPASSWORD"),
		KeyToken:    os.Getenv("KEYTOKEN"),
	}
	return &dconf
}

/*func Get() (*Conf, error) {

	var dconf Conf

	yamlFile, err := ioutil.ReadFile("/home/ubuntu/HttpServer/GoGoGo/FileStorageServer/etc/etc.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "Read .yaml File, func Get") //error handling
	}

	err = yaml.Unmarshal(yamlFile, &dconf)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal .yaml File, func Get") //error handling
	}

	return &dconf, nil
}*/
