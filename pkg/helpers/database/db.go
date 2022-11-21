package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

type Config struct {
	Host               string
	DBName             string
	Username           string
	Password           string
	Port               string
	DBType             string
	Timeout            int
	MaxOpenConnections int
	MaxIdleConnections int
}

func (c *Config) GetConnectionString() string {
	if c.Timeout == 0 {
		c.Timeout = 5
	}

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		c.DBType,
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.DBName,
		c.Timeout)

	return connStr
}

func New(dbConfig Config) (*sql.DB, error) {

	connectionString := dbConfig.GetConnectionString()
	db, err := sql.Open(dbConfig.DBType, connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	db.SetMaxIdleConns(dbConfig.MaxIdleConnections)

	// Проверять подключение к базе данных нужно через Ping.
	// При вызове sql.Open не происходит реальной инициализации подключения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error %s database connection %e ", connectionString, err)
	}

	return db, nil
}
