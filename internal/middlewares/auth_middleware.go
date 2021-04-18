// Referenced https://auth0.com/docs/quickstart/backend/golang/01-authorization
package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

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

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")

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
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

type AuthPredicate func(*http.Request) bool

type Authorizer func(w http.ResponseWriter, r *http.Request) error

type authenticationRule struct {
	match AuthPredicate
	rule  AuthPredicate
}

type authConfig struct {
	ignored      AuthPredicate
	authorize    Authorizer
	authenticate []*authenticationRule
}

type AuthConfigrator func(*authConfig)

func IgnoreAuth(filter AuthPredicate) AuthConfigrator {
	return func(c *authConfig) {
		c.ignored = filter
	}
}

func Authorize(authorizer Authorizer) AuthConfigrator {
	return func(c *authConfig) {
		c.authorize = authorizer
	}
}

func Authenticate(match AuthPredicate, rule AuthPredicate) AuthConfigrator {
	return func(c *authConfig) {
		c.authenticate = append(c.authenticate, &authenticationRule{
			match: match,
			rule:  rule,
		})
	}
}

func Path(method string, path string) AuthPredicate {
	return func(r *http.Request) bool {
		return r.Method == method && r.URL.Path == path
	}
}

func PermitAll() AuthPredicate {
	return func(r *http.Request) bool { return true }
}

func HasScope(scope string) AuthPredicate {
	return func(r *http.Request) bool {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		tokenStr := authHeaderParts[1]

		token, _ := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			cert, err := getPemCert(token)
			if err != nil {
				return nil, err
			}
			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		})

		claims, ok := token.Claims.(*CustomClaims)

		hasScope := false
		if ok && token.Valid {
			result := strings.Split(claims.Scope, " ")
			for i := range result {
				if result[i] == scope {
					hasScope = true
				}
			}
		}

		return hasScope
	}
}

func NewAuthenticationMiddleware(configrators ...AuthConfigrator) gin.HandlerFunc {
	conf := &authConfig{}
	for _, c := range configrators {
		c(conf)
	}

	return func(c *gin.Context) {
		if conf.ignored != nil && conf.ignored(c.Request) {
			//認証スキップの条件に一致する場合,終了
			return
		}

		if err := conf.authorize(c.Writer, c.Request); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		permit := false
		for _, r := range conf.authenticate {
			req := c.Request
			if r.match(req) {
				permit = r.rule(req)
			}
		}

		if !permit {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
