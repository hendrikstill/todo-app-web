package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type TodoAppConfig struct {
	HealthCheckTime int
	DBDriver        string
	DBConfig        map[string]string
	ReleaseMode     string
	SecurityConfig  SecurityConfig
}

type SecurityConfig struct {
	EnableJWTAuthentification bool
	JWTSharedKey              string
}

func readConfig(configFile string) (*TodoAppConfig, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &TodoAppConfig{
			DBDriver:       "redis",
			DBConfig:       map[string]string{},
			ReleaseMode:    gin.DebugMode,
			SecurityConfig: SecurityConfig{EnableJWTAuthentification: false, JWTSharedKey: "secrete"},
		}, err
	}
	config := &TodoAppConfig{}
	json.Unmarshal(file, config)

	if config.DBDriver == "" {
		log.Println("Use redis as default")
		config.DBDriver = "redis"
	}

	if config.DBConfig == nil {
		config.DBConfig = map[string]string{}
	}

	if config.ReleaseMode == "" {
		config.ReleaseMode = gin.DebugMode
	}

	if config.SecurityConfig.EnableJWTAuthentification {
		log.Println("Enabled authentification via JWT")
	}

	if config.SecurityConfig.EnableJWTAuthentification && config.SecurityConfig.JWTSharedKey == "" {
		log.Println("Use default JWT sharedSecrete")
		config.SecurityConfig = SecurityConfig{JWTSharedKey: "secrete"}
	}

	return config, err
}
