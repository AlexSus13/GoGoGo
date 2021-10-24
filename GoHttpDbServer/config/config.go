package config

import (
	"io/ioutil"
	"log"
	"github.com/go-yaml/yaml"
)

type DB struct {
	User string `yaml:"user"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Dbname string `yaml:"dbname"`
	Sslmode string `yaml:"sslmode"`
}

type Conf struct {
	Port string `yaml:"port"`
	DB DB
	Pathopenfile string `yaml:"pathopenfile"`
	Pathsavefile string `yaml:"pathsavefile"`
	Pathdeletefile string `yaml:"pathdeletefile"`
}

func (c *Conf) GetConf() *Conf {

	yamlFile, err := ioutil.ReadFile("home/ubuntu/zooProject/GoGoGo/GoHttpDbServer/etc/etc.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
