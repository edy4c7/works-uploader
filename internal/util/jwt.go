package util

import (
	"encoding/json"
	"errors"
	"net/http"

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
