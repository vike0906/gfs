package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
)

type BaseInfo struct {
	Server Server `yaml:"server"`
	Mysql  Mysql  `yaml:"mysql"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DataBase string `yaml:"dataBase"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
}

var Info BaseInfo

func Init() error {
	var osType = runtime.GOOS
	filePath, _ := os.Getwd()
	if osType == "windows" {
		filePath = filePath + "\\config\\application.yaml"
	} else if osType == "linux" {
		filePath = filePath + "/config/application.yaml"
	}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &Info)
	if err != nil {
		return err
	}
	return nil
}

func ConfigInfo() *BaseInfo {
	return &Info
}
