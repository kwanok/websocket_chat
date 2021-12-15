package server

import (
	"database/sql"
	"friday/server/utils"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DBCon *sql.DB
var RedisClient *redis.Client
var Sqlite3 *sql.DB

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

	go initMysql()
	go initRedis()
	go initSqlite3()

	return "InitDB Success"
}

func initRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func initSqlite3() {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `	
	CREATE TABLE IF NOT EXISTS room (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		private TINYINT NULL
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	sqlStmt = `	
	CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	Sqlite3 = db
}

func initMysql() {
	databaseInfo := DatabaseInfo{
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ROOT"),
	}

	db, err := sql.Open("mysql", getSourceName(databaseInfo))
	utils.FatalError{Error: err}.Handle()

	DBCon = db
}
