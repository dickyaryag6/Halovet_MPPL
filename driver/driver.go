package driver

import(
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
)

func Connect() (*sql.DB, error) {
  db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/halovet_db")
  if err != nil {
      return nil, err
  }

  return db, nil
}
