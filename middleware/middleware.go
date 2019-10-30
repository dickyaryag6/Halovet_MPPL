package middleware

import (
  "net/http"
  . "fmt"
  // "reflect"
  "time"
  jwt "github.com/dgrijalva/jwt-go"
  "strings"
  "context"
  // "os"
  // "github.com/joho/godotenv"
  // "log"
)

var LOGIN_EXP_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256


func MiddlewareJWTAuthorization(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    jwtSignKey:="notsosecret"
    //
    // url:=r.URL.Path
    // if itemExists(noAuthSlice, url){
    //   next.ServeHTTP(w, r)
    //   return
    // }

    if r.URL.Path == "account/login" {
            next.ServeHTTP(w, r)
            return
    }
    //ambil data yang dikirim ke http, in this case,
    ctx := r.Context()
    if ctx == nil {
        ctx = context.Background()
    }

    // ambil token yg dikasi user
    authorizationHeader := r.Header.Get("Authorization")
    if !strings.Contains(authorizationHeader, "Bearer") {
      http.Error(w, "Invalid token", http.StatusBadRequest)
      return
    }

    tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, Errorf("Signing method invalid")
        } else if method != JWT_SIGNING_METHOD {
            return nil, Errorf("Signing method invalid")
        }

        return []byte(jwtSignKey), nil
    })

    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Jika melihat struct http.Request maka  di dalamnya
    // terdapat context. Setiap request yang datang akan diset otomatis menjadi ctx.Background().

    // ctx := context.WithValue(context.Background(), "userInfo", claims)?
    ctx = context.WithValue(ctx, "userInfo", claims)
    r = r.WithContext(ctx)

    next.ServeHTTP(w, r)
})}
