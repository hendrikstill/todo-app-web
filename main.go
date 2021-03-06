package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/johscheuer/todo-app-web/tododb"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
)

var (
	appVersion  string
	showVersion bool
	database    tododb.TodoDB
)

func main() {
	configFile := flag.String("config-file", "./default.config", "Path to the configuration file")
	flag.BoolVar(&showVersion, "version", false, "Shows the version")
	flag.Parse()

	if showVersion {
		log.Printf("Version: %s\n", appVersion)
		return
	}

	config, err := readConfig(*configFile)
	if err != nil {
		log.Println(err)
	}

	gin.SetMode(config.ReleaseMode)
	if strings.ToLower(config.DBDriver) == "mysql" {
		database = tododb.NewMySQLDB(config.DBConfig, appVersion)
	} else if strings.ToLower(config.DBDriver) == "redis" {
		database = tododb.NewRedisDB(config.DBConfig, appVersion)
	}

	p := ginprometheus.NewPrometheus("gin")
	database.RegisterMetrics()

	// Iniitialize metrics
	quit := make(chan struct{})
	defer close(quit)
	if config.HealthCheckTime > 0 {
		healthCheckTimer := time.NewTicker(time.Duration(config.HealthCheckTime) * time.Second)
		go func() {
			for {
				select {
				case <-healthCheckTimer.C:
					log.Println("Called Health check")
					database.GetHealthStatus()
				case <-quit:
					healthCheckTimer.Stop()
					return
				}
			}
		}()
	}

	router := gin.New()

	p.Use(router)

	router.GET("/health", healthCheckHandler)
	router.GET("/whoami", whoAmIHandler)
	router.GET("/version", versionHandler)

	authorized := router.Group("/", JwtAuthRequired(config.SecurityConfig))
	{
		authorized.GET("/read/todo", readTodoHandler)
		authorized.GET("/insert/todo/:value", insertTodoHandler)
		authorized.GET("/delete/todo/:value", deleteTodoHandler)
	}

	router.Use(static.Serve("/", static.LocalFile("./public", true)))
	router.Run(":3000")
}
