package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func InitMySQL() {
	var err error
	connStr := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ")/" + os.Getenv("MYSQL_DB")
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		logrus.Fatalf("Unable to connect to database: %v", err)
	}

	boil.SetDB(db)
}
