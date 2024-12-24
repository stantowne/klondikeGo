package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
	"os"
)

type DBConfig struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

func LoadConfig() (*DBConfig, error) {
	err := godotenv.Load("sqlDbConnect.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &DBConfig{
		Server:   os.Getenv("DB_SERVER"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}

	// Validate that all necessary environment variables are set
	if config.Server == "" || config.Port == "" || config.User == "" || config.Password == "" || config.Database == "" {
		return nil, fmt.Errorf("missing one or more required environment variables")
	}

	return config, nil
}

func Connect(config *DBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s", config.Server, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating a SQL Server database connection: %w", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error verifying connection to the SQL Server: %w", err)
	}

	return db, nil
}
