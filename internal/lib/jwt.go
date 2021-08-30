package lib

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/edy4c7/works-uploader/internal/common"
	"github.com/form3tech-oss/jwt-go"
)

type jsonWebKeys struct {
	Keys []struct {
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}

func GetPemCertOfJWK(token *jwt.Token, url string) (string, error) {
	cert := ""
	resp, err := http.Get(url)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = jsonWebKeys{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func CheckJWTScope(jwk string, scope string, tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, make(jwt.MapClaims), func(token *jwt.Token) (interface{}, error) {
		cert, err := GetPemCertOfJWK(token, jwk)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	if err != nil {
		log.Println(err)
	}

	var claims struct {
		Scope     string `json:"scope"`
		Audience  string `json:"aud,omitempty"`
		ExpiresAt int64  `json:"exp,omitempty"`
		Id        string `json:"jti,omitempty"`
		IssuedAt  int64  `json:"iat,omitempty"`
		Issuer    string `json:"iss,omitempty"`
		NotBefore int64  `json:"nbf,omitempty"`
		Subject   string `json:"sub,omitempty"`
	}
	if err = common.MapToStruct(token.Claims, &claims); err != nil {
		return false
	}

	hasScope := false
	if token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}

	return hasScope
}
