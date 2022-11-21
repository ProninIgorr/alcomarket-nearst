package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ProninIgorr/alcomarket-nearst/internal/migrations"
	"github.com/ProninIgorr/alcomarket-nearst/internal/routes"
	"github.com/ProninIgorr/alcomarket-nearst/internal/store"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/pg"
	"github.com/d-kolpakov/logger"
	"github.com/dhnikolas/configo"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const ServiceName = "pronin-location"

var AppVersion string

func main() {
	fmt.Println("Starting service " + ServiceName)

	//Logger initialization
	logDriver := &logger.STDOUTDriver{}
	loggerConfig := logger.LoggerConfig{
		ServiceName: ServiceName,
		Level:       configo.EnvInt("logging-level", 2),
		Buffer:      configo.EnvInt("app-logger-buffer-size", 10000),
		Output:      []logger.LogDriver{logDriver},
	}
	l, err := logger.GetLogger(loggerConfig)
	if err != nil {
		panic(err)
	}

	cfg := &pg.Config{}
	cfg.Host = configo.EnvString("db-host", "")
	cfg.Username = configo.EnvString("db-username", "")
	cfg.Password = configo.EnvString("db-password", "")
	cfg.Port = configo.EnvString("db-port", "")
	cfg.DbName = configo.EnvString("db-name", "")
	cfg.Timeout = 5
	config, err := pg.NewPoolConfig(cfg)
	if err != nil {
		l.Err(err)
		panic(err)
	}
	config.MaxConns = 100
	c, err := pg.NewConnection(config)
	if err != nil {
		l.Err(err)
		panic(err)
	}

	mdb, err := sql.Open("postgres", config.ConnString())
	if err != nil {
		l.Err(err)
		panic(err)
	}
	err = goose.Up(mdb, "/var")
	if err != nil {
		l.Err(err)
		panic(err)
	}

	s := store.New(c)
	err = s.CacheShops()
	if err != nil {
		l.Err(err)
		panic(err)
	}

	routes.Start(ServiceName, AppVersion, l, c, s)
}
