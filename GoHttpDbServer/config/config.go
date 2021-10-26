package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
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
	Port           string `yaml:"port"`
	DB             DB     `yaml:"db"`
	Pathopenfile   string `yaml:"pathopenfile"`
	Pathsavefile   string `yaml:"pathsavefile"`
	Pathdeletefile string `yaml:"pathdeletefile"`
	//WriteTimeout   int    `yaml:"writetimeout"`
	//ReadTimeout    int    `yaml:"readtimeout"`
}

func GetConf() *Conf { //(c *Conf)

	var c Conf

	yamlFile, err := ioutil.ReadFile("/home/ubuntu/zooProject/GoGoGo/GoHttpDbServer/etc/etc.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}
