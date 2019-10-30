package handler

import (
  "net/http"
  models "Halovet/models"
  dbCon "Halovet/driver"
  "log"
  . "fmt"
  "time"
  jwt "github.com/dgrijalva/jwt-go"
  mid "Halovet/middleware"
  "encoding/json"
  // "io/ioutil"
  // "encoding/json"
  "github.com/kataras/go-sessions"
  "database/sql"
  _ "golang.org/x/crypto/bcrypt"
  _ "github.com/go-sql-driver/mysql"
  "strings"
)

var err error
var db *sql.DB
var LOGIN_EXP_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256


type  M map[string]interface{}

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

func QueryUser(email string) (bool, models.Pet_Owner){

	var user models.Pet_Owner

  sql_statement := "SELECT id,Email,Name,Password FROM pet_owner WHERE Email=?"
	err = db.QueryRow(sql_statement, email).
		Scan(&user.ID,&user.Email,&user.Name,&user.Password)
  if err == sql.ErrNoRows {
      return false,user
  }
  return true,user

}

func Login(w http.ResponseWriter, r *http.Request) {

    Println("Endpoint Hit: Login")
    var response struct {
      Status bool
      Message string
      Data    map[string]interface{}
    }
    // var response models.JwtResponse
    //CEK UDAH ADA YANG LOGIN ATAU BELUM
    // session := sessions.Start(w, r)
  	// if len(session.GetString("email")) != 0 && checkErr(w, r, err) {
  	// 	http.Redirect(w, r, "/", 302)
  	// }
    jwtSignKey:="notsosecret"
    appName:="Halovet"
    var message string
  	// email := r.FormValue("email")
  	// password := r.FormValue("password")

    //dapetin informasi dari form login
    email, password, ok := r.BasicAuth()
    if !ok {
      message = "Invalid email or password"
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      response.Status   = true
      response.Message  = message
      json.NewEncoder(w).Encode(response)
      return
    }

    ok, userInfo := mid.AuthenticateUser(email, password)
    if !ok {
      message = "Invalid email or password"
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      response.Status   = true
      response.Message  = message
      json.NewEncoder(w).Encode(response)
      return
    }

    claims := models.TheClaims{
            StandardClaims: jwt.StandardClaims{
            Issuer:    appName,
            ExpiresAt: time.Now().Add(LOGIN_EXP_DURATION).Unix(),
        },
        User : userInfo,
    }

    token := jwt.NewWithClaims(
    JWT_SIGNING_METHOD,
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
      "jwtToken" : signedToken,
      "user"     : userInfo,
    }

    //RESPON JSON
    message = "Login Succesfully"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    response.Status   = true
    response.Message  = message
    response.Data  = data
    json.NewEncoder(w).Encode(response)

    http.SetCookie(w, &http.Cookie{
      Name    : "token",
      Value   : signedToken,
      Expires : time.Now().Add(LOGIN_EXP_DURATION),
    })
}

func ValidateEmail(email string) (string, bool) {
  if !strings.Contains(email, "@") {
		return "Email address is required", false
	}
  message:="Email Valid"
  return message, true
}

func ValidatePassword(password string) (string, bool) {
  if len(password) < 6 {
    return "Password should be more than 6 character", false
  }
  message:="Password Valid"
  return message, true
}

func Register(w http.ResponseWriter, r *http.Request) {

  Println("Endpoint Hit: Register")

  var response struct {
    Status bool
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
  if message, status := ValidateEmail(email) ; status != true {
    Println(message)
  }
  password := r.FormValue("password")
  if message, status := ValidatePassword(password) ; status != true {
    Println(message)
  }
  name := r.FormValue("name")
  //cek apakah email user sudah ada di database
  //query user dengan email tersebut
  status,_ := QueryUser(email)

  // kalo status false , berarti register
  // kalo status true, berarti print email terdaftar

  if status{
    Println("Email sudah terdaftar")
  } else {
    // hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
    hashedPassword := password

    // Println(hashedPassword)
    if len(hashedPassword) != 0 && checkErr(w, r, err) {
      stmt, err := db.Prepare("INSERT INTO pet_owner (Email, Name, Password) VALUES (?,?,?)")
      if err == nil {
        _, err := stmt.Exec(&email, &name, &hashedPassword)
        if err != nil {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
        }
        message := "Register Succesfull"
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        response.Status   = true
        response.Message  = message
        json.NewEncoder(w).Encode(response)
        http.Redirect(w, r, "/login", 201)
        return
      }
    } else {
      message := "Registration Failed"
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      response.Status   = true
      response.Message  = message
      json.NewEncoder(w).Encode(response)

      http.Redirect(w, r, "/register", 302)
      }
  }
}

func Logout(w http.ResponseWriter, r *http.Request) {
  Println("Endpoint Hit: Logout")

	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}
