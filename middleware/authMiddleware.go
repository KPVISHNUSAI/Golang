package middleware

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 30*time.Second {
			expirationTime := time.Now().Add(24 * time.Hour)
			claims.ExpiresAt = expirationTime.Unix()
			tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := tkn.SignedString(jwtKey)
			if err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    tokenString,
					Expires:  expirationTime,
					HttpOnly: true,
				})
			}
		}

		next.ServeHTTP(w, r)
	})
}
