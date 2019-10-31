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
func AuthenticateUser(email, password string) (bool, models.Pet_Owner) {
	//cari document dengan username dan password yg diberikan

	var user models.Pet_Owner
	sqlStatement := "SELECT id,Email,Name,Password FROM pet_owner WHERE Email=?"

	err = db.QueryRow(sqlStatement, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err == sql.ErrNoRows {
		return false, user
	}

	// check_match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	log.Println(user.Password)
	// if check_match != nil {
	if password != user.Password {
		//LOGIN FAILED, PASSWORD SALAH
		log.Println("Password Salah")
		return false, user
	}
	//LOGIN SUCCESS
	log.Println("Password Benar")
	return true, user
}
