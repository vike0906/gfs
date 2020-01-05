package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
)

var osType = runtime.GOOS

var Info BaseInfo

func Init() error {
	filePath, _ := os.Getwd()
	if osType == "windows" {
		filePath = filePath + "\\config\\application.yaml"
	} else if osType == "linux" {
		filePath = filePath + "/config/application.yaml"
	}
	fmt.Println(filePath)
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

type BaseInfo struct {
	Server Server `yaml:"server"`
	Mysql  Mysql  `yaml:"mysql"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Mysql struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
