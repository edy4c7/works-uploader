// Referenced https://auth0.com/docs/quickstart/backend/golang/01-authorization
package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func NewJWTMiddleware(aud string, iss string) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token, iss+".well-known/jwks.json")
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPemCert(token *jwt.Token, url string) (string, error) {
	cert := ""
	resp, err := http.Get(url)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
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

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

type requestFilter func(*http.Request) bool

type policyFunc func(*http.Request) bool

type JWTMiddleware interface {
	CheckJWT(w http.ResponseWriter, r *http.Request) error
}

type authConfig struct {
	skipped requestFilter
}

type AuthConfigrator func(*authConfig)

func SkipAuthorization(filter requestFilter) AuthConfigrator {
	return func(c *authConfig) {
		c.skipped = filter
	}
}

func NewAutorizationMiddleware(jwtMiddleware JWTMiddleware, configrators ...AuthConfigrator) gin.HandlerFunc {
	conf := &authConfig{}
	for _, c := range configrators {
		c(conf)
	}

	return func(c *gin.Context) {
		if conf.skipped != nil && conf.skipped(c.Request) {
			//認証スキップの条件に一致する場合,終了
			return
		}

		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
}

func NewAuthenticationMiddleware(policy policyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !policy(c.Request) {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}
