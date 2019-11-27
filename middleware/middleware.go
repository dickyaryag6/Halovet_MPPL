package middleware

import (
	. "fmt"
	"net/http"
	"strconv"

	// "reflect"
	"context"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	// "os"
	// "github.com/joho/godotenv"
	// "log"
)

// LoginExpDuration : Durasi Token Berlaku/Valid
var LoginExpDuration = time.Duration(1) * time.Hour

// JwtSigningMethod : Metode Pembuatan Token
var JwtSigningMethod = jwt.SigningMethodHS256

// JWTAuthorization : fungsi pembuatan JWT
func JWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jwtSignKey := "notsosecret"
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
			} else if method != JwtSigningMethod {
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
	})
}

// type middleware func(http.HandlerFunc) http.HandlerFunc

// func buildChain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
// 	// if our chain is done, use the original handlerfunc
// 	if len(m) == 0 {
// 		return f
// 	}
// 	// otherwise nest the handlerfuncs
// 	return m[0](buildChain(f, m[1:cap(m)]...))
// }

// PetOwner : Middleware PetOwner
func PetOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//CEK APAKAH PET OWNER APA BUKAN

		userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
		user := userInfo["User"]
		userReal, _ := user.(map[string]interface{})
		Role, err := strconv.Atoi(Sprintf("%v", userReal["Role"]))
		if err == nil {
			Println(Role)
		}
		if Role != 1 {
			w.Write([]byte("Only Pet Owner is allowed to do this method"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Doctor : Middleware Doctor
func Doctor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//CEK APAKAH DOCTOR APA BUKAN
		userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
		user := userInfo["User"]
		userReal, _ := user.(map[string]interface{})
		Role, err := strconv.Atoi(Sprintf("%v", userReal["Role"]))
		if err == nil {
			Println(Role)
		}
		if Role != 2 {
			w.Write([]byte("Doctor is not allowed to do this method"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func userIsAdmin(role int) bool {
	if role == 3 {
		return true
	}
	return false
}

func AdminOnly(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
		user := userInfo["User"]
		userReal, _ := user.(map[string]interface{})
		Role, err := strconv.Atoi(Sprintf("%v", userReal["Role"]))
		if err != nil {
			panic("Gagal mendapatkan role user")
		}
		if !userIsAdmin(Role) {
			http.Error(w, "Admin only", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// func Admin(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		//CEK APAKAH PET OWNER APA BUKAN

// 		userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
// 		user := userInfo["User"]
// 		userReal, _ := user.(map[string]interface{})
// 		Role, err := strconv.Atoi(Sprintf("%v", userReal["Role"]))
// 		if err == nil {
// 			Println(Role)
// 		}
// 		if Role != 3 {
// 			w.Write([]byte("Only Admin is allowed to do this method"))
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
