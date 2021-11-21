package server

import (
	"database/sql"
	"friday/tools"
	_ "github.com/go-sql-driver/mysql"
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

func InitDB() {
	databaseInfo := DatabaseInfo{
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ROOT"),
	}

	db, err := sql.Open("mysql", getSourceName(databaseInfo))
	tools.ErrorHandler(err)

	DBCon = db
}
