package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dataSources struct {
	DB *sqlx.DB
}

// 初始化建立连接
func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")
	// 加载 env 数据
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost,
		pgPort,
		pgUser,
		pgPassword,
		pgDB,
		pgSSL,
	)

	log.Printf("connecting to Postgresql\n")
	db, err := sqlx.Open("postgres", pgConnString)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// 验证数据库是否连接正常
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return &dataSources{DB: db}, nil
}

// Close 关闭数据库
func (d *dataSources) Close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}
