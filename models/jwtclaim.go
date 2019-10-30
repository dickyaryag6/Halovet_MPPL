package models

import (
  jwt "github.com/dgrijalva/jwt-go"
)

type TheClaims struct {
    jwt.StandardClaims
    User Pet_Owner
}
