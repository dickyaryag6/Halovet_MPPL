package driver

import (
	"database/sql"
	"log"
	"os"

	//driver mysql
	. "fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//Connect : Fungsi Untuk Konek Database
func Connect() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUSERNAME := os.Getenv("DB_USERNAME")
	dbNAME := os.Getenv("DB_NAME")
	dbPASSWORD := os.Getenv("DB_PASSWORD")

	dbURI := Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s",
		dbUSERNAME,
		dbPASSWORD,
		dbNAME,
	)

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}
