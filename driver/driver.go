package driver

import (
	"database/sql"
	//driver mysql
	_ "github.com/go-sql-driver/mysql"
)

//Connect : Fungsi Untuk Konek Database
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/halovet_db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
