package config

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DBCon *sql.DB
var JwtRedis *redis.Client
var PubSubRedis *redis.Client

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
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	go initMysql()
	go initRedis()

	return "InitDB Success"
}

func initRedis() {
	//Initializing redis
	jwtDsn := os.Getenv("REDIS_JWT_DSN")
	if len(jwtDsn) == 0 {
		jwtDsn = "localhost:6379"
	}
	JwtRedis = redis.NewClient(&redis.Options{
		Addr: jwtDsn, //redis port
	})

	psDsn := os.Getenv("REDIS_PUB_SUB_DSN")
	if len(psDsn) == 0 {
		psDsn = "localhost:6380"
	}
	PubSubRedis = redis.NewClient(&redis.Options{
		Addr: psDsn, //redis port
	})
}

func initMysql() {
	databaseInfo := DatabaseInfo{
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ROOT"),
	}

	db, err := sql.Open("mysql", getSourceName(databaseInfo))
	if err != nil {
		log.Fatal(err)
	}

	DBCon = db
}
