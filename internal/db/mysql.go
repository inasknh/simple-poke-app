package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/inasknh/simple-poke-app/internal/config"
	"log"
)

func NewMySql(config config.Configurations) *sql.DB {
	dbConfig := mysql.Config{
		User:      config.Database.User,
		Passwd:    config.Database.Password,
		Addr:      config.Database.Host,
		DBName:    config.Database.Name,
		ParseTime: true,
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		log.Fatalf("Failed to open database %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("DB cannot be pinged %s", err)
	}

	return db
}
