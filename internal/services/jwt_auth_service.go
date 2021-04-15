package services

import (
	"net/http"

	"github.com/edy4c7/darkpot-school-works/internal/api"
	"github.com/edy4c7/darkpot-school-works/internal/errors"
)

type AuthPredicate func(*http.Request) bool

type AuthConfigrator func(*JWTAuthServiceImpl)

func none(r *http.Request) bool {
	return false
}

type JWTAuthService interface {
	Authenticate(w http.ResponseWriter, r *http.Request) error
}

type JWTAuthServiceImpl struct {
	jwtMiddleware api.JWTMiddleware
	ignored       AuthPredicate
	definitions   []AuthPredicate
}

func Ignore(predicate AuthPredicate) AuthConfigrator {
	return func(s *JWTAuthServiceImpl) {
		s.ignored = predicate
	}
}

func Define(predicates ...AuthPredicate) AuthConfigrator {
	return func(s *JWTAuthServiceImpl) {
		s.definitions = predicates
	}
}

func NewJWTAuthServiceImpl(jwtmiddleware api.JWTMiddleware, configrator ...AuthConfigrator) *JWTAuthServiceImpl {
	service := &JWTAuthServiceImpl{
		jwtMiddleware: jwtmiddleware,
		ignored:       none,
	}

	for _, c := range configrator {
		c(service)
	}

	return service
}

func (r *JWTAuthServiceImpl) Authenticate(w http.ResponseWriter, req *http.Request) error {
	if r.ignored(req) {
		//認証スキップの条件に一致する場合,終了
		return nil
	}
	if err := r.jwtMiddleware.CheckJWT(w, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.NewApplicationError(errors.Code(errors.DSWE03), errors.Cause(err))
	}

	permit := false
	for _, v := range r.definitions {
		permit = v(req)
	}

	if !permit {
		return errors.NewApplicationError(errors.Code(errors.DSWE04))
	}

	return nil
}
