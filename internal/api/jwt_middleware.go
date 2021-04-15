package api

import "net/http"

type JWTMiddleware interface {
	CheckJWT(w http.ResponseWriter, r *http.Request) error
}
