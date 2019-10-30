package middleware

import (
  dbCon "Halovet/driver"
  models "Halovet/models"
  "database/sql"
  "log"
  // "context"
  _ "golang.org/x/crypto/bcrypt"
  . "fmt"
)

var err error
var db *sql.DB

func init() {
  // KONEK KE DATABASE
  db, err = dbCon.Connect()
  if err != nil {
    log.Fatal(err.Error())
  }
}

func AuthenticateUser(email, password string) (bool, models.Pet_Owner) {
  //cari document dengan username dan password yg diberikan

  var user models.Pet_Owner
  sql_statement:="SELECT id,Email,Name,Password FROM pet_owner WHERE Email=?"

  err = db.QueryRow(sql_statement, email).
    Scan(&user.ID,&user.Email,&user.Name,&user.Password)
  if err == sql.ErrNoRows {
      return false,user
  }


  // check_match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
  Println(user.Password)
  // if check_match != nil {
  if password != user.Password {
    //LOGIN FAILED, PASSWORD SALAH
    Print("Password Salah")
    return false, user
  }
  //LOGIN SUCCESS
  Print("Password Benar")
  return true, user
}
