package handler

import (
	dbCon "Halovet/driver"
	mid "Halovet/middleware"
	models "Halovet/models"
	"encoding/json"
	. "fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/go-sessions"
	_ "golang.org/x/crypto/bcrypt"
)

var err error
var db *sql.DB

// LoginExpDuration : ...
var LoginExpDuration = time.Duration(1) * time.Hour

// JwtSigningMethod : ...
var JwtSigningMethod = jwt.SigningMethodHS256

type M map[string]interface{}

func init() {
	// KONEK KE DATABASE
	db, err = dbCon.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	// if err != nil {
	//
	//   Println(r.Host + r.URL.Path)
	//
	// 	http.Redirect(w, r, r.Host+r.URL.Path, 301)
	// 	return false
	// }

	return true
}

// QueryUser : ngecek apakah sudah ada user dengan email tersebut
func QueryUser(email string) (bool, models.Account) {

	var user models.Account
	// Println(email)
	sqlStatement := "SELECT email,name,password,role FROM account WHERE email=?"
	err = db.QueryRow(sqlStatement, email).
		Scan(&user.Email, &user.Name, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return false, user
	}
	return true, user

}

// ValidateEmail : ngecek apakah format email benar
func ValidateEmail(email string) (string, bool) {
	if !strings.Contains(email, "@") || email == "" {
		return "Email address is required", false
	}
	message := "Email Valid"
	return message, true
}

// ValidatePassword : ngecek apakah password sesuai
func ValidatePassword(password string) (string, bool) {
	if len(password) < 6 {
		return "Password should be more than 6 character", false
	}
	message := "Password Valid"
	return message, true
}

// Login : fungsi Login
func Login(w http.ResponseWriter, r *http.Request) {

	Println("Endpoint Hit: Login")
	var response struct {
		Status  bool
		Message string
		Data    map[string]interface{}
	}
	// var response models.JwtResponse
	//CEK UDAH ADA YANG LOGIN ATAU BELUM
	// session := sessions.Start(w, r)
	// if len(session.GetString("email")) != 0 && checkErr(w, r, err) {
	// 	http.Redirect(w, r, "/", 302)
	// }
	jwtSignKey := "notsosecret"
	appName := "Halovet"
	var message string

	//dapetin informasi dari Basic Auth
	// email, password, ok := r.BasicAuth()
	// if !ok {
	// 	message = "Invalid email or password"
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(200)
	// 	response.Status = false
	// 	response.Message = message
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }
	//dapetin informasi dari form
	email := r.FormValue("email")
	if _, status := ValidateEmail(email); status != true {
		message := "Format Email Salah atau Kosong"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return
	}
	password := r.FormValue("password")
	if _, status := ValidatePassword(password); status != true {
		message := "Format Password Salah atau Kosong, Minimal 6 Karakter"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return
	}

	ok, userInfo := mid.AuthenticateUser(email, password)
	if !ok {
		message = "Invalid email or password"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return
	}

	claims := models.TheClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    appName,
			ExpiresAt: time.Now().Add(LoginExpDuration).Unix(),
		},
		User: userInfo,
	}

	token := jwt.NewWithClaims(
		JwtSigningMethod,
		claims,
	)

	signedToken, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//tokennya dijadiin json
	// tokenString, _ := json.Marshal(M{ "token": signedToken })
	// w.Write([]byte(tokenString))
	data := map[string]interface{}{
		"jwtToken": signedToken,
		"user":     userInfo,
	}

	//RESPON JSON
	message = "Login Succesfully"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	response.Status = true
	response.Message = message
	response.Data = data
	json.NewEncoder(w).Encode(response)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   signedToken,
		Expires: time.Now().Add(LoginExpDuration),
	})
}

// Register : fungsi register
func Register(w http.ResponseWriter, r *http.Request) {

	Println("Endpoint Hit: Register")

	var response struct {
		Status  bool
		Message string
	}

	// reqBody, _ := ioutil.ReadAll(r.Body)
	// var register_user models.Pet_Owner
	// json.Unmarshal(reqBody, &register_user)

	// email := register_user.Email
	// password := register_user.Password
	// name := register_user.Name

	//BIKIN VALIDATION

	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

	if len(name) == 0 {
		message := "Ada Kolom Yang Kosong"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, status := ValidateEmail(email); status != true {
		message := "Format Email Kosong atau Salah"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, status := ValidatePassword(password); status != true {
		message := "Format Password Kosong atau Salah, Minimal 6 Karakter"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return
	}

	//cek apakah email user sudah ada di database
	//query user dengan email tersebut
	status, _ := QueryUser(email)

	// kalo status false , berarti register
	// kalo status true, berarti print email terdaftar

	if status {
		message := "Email sudah terdaftar"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		response.Status = false
		response.Message = message
		json.NewEncoder(w).Encode(response)
		return

	} else {
		// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		hashedPassword := password
		Role := 1

		// Println(hashedPassword)
		if len(hashedPassword) != 0 && checkErr(w, r, err) {
			stmt, err := db.Prepare("INSERT INTO account (Email, Name, Password, Role) VALUES (?,?,?,?)")
			if err == nil {
				_, err := stmt.Exec(&email, &name, &hashedPassword, &Role)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				message := "Register Succesfull"
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				response.Status = true
				response.Message = message
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			message := "Registration Failed"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			response.Status = false
			response.Message = message
			json.NewEncoder(w).Encode(response)

		}
	}
}

// Logout : fungsi logout
func Logout(w http.ResponseWriter, r *http.Request) {
	Println("Endpoint Hit: Logout")

	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}
