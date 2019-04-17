package middleware

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"
	"net/http"
)

//RouterWare is a function used to modify the mux
type HandlerWare func(h http.Handler) http.Handler

func WithJWT(secret string) HandlerWare {
	return func(h http.Handler) http.Handler {
		j := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})
		j.Handler(h)
		return h
	}
}

func WithCORS(options cors.Options) HandlerWare {
	return func(h http.Handler) http.Handler {
		cors.New(options).Handler(h)
		return h
	}
}
