package models

import (
	"database/sql"
	"friday/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

var (
	// DBCon is the connection handle
	// for the database
	DBCon *sql.DB
)

type DatabaseInfo struct {
	Name     string
	Host     string
	Password string
	Root     string
}

func getSourceName(db DatabaseInfo) string {
	return db.Root + ":" + db.Password + "@tcp(" + db.Host + ":3306)/" + db.Name
}

func InitDB() string {
	err := godotenv.Load(".env")
	utils.FatalError{Error: err}.Handle()

	databaseInfo := DatabaseInfo{
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ROOT"),
	}

	db, err := sql.Open("mysql", getSourceName(databaseInfo))
	utils.FatalError{Error: err}.Handle()

	DBCon = db

	return "InitDB Success"
}
