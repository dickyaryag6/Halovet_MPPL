package middleware

import (
	dbCon "Halovet/driver"
	models "Halovet/models"
	"database/sql"
	"log"
	// "context"
	// . "fmt"
	// _ "golang.org/x/crypto/bcrypt"
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

// AuthenticateUser : ngecek apakah user tersebut terdaftar
func AuthenticateUser(email, password string) (bool, models.Account) {
	//cari document dengan username dan password yg diberikan

	var user models.Account
	sqlStatement := "SELECT email,name,password,role FROM account WHERE email=?"

	err = db.QueryRow(sqlStatement, email).
		Scan(&user.Email, &user.Name, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return false, user
	}

	// check_match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	log.Println(user.Password)
	// if check_match != nil {
	if password != user.Password {
		//LOGIN FAILED, PASSWORD SALAH
		log.Println("Password atau email salah")
		return false, user
	}
	//LOGIN SUCCESS
	log.Println("Password atau email salah")
	return true, user
}
