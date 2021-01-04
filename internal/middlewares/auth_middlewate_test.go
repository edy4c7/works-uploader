package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edy4c7/darkpot-school-works/internal/mocks"
	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const authorizationURL string = "http://localhost:9099/identitytoolkit.googleapis.com/v1/%v?key=%v"

var firebaseScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/datastore",
	"https://www.googleapis.com/auth/devstorage.full_control",
	"https://www.googleapis.com/auth/firebase",
	"https://www.googleapis.com/auth/identitytoolkit",
	"https://www.googleapis.com/auth/userinfo.email",
}

type credential struct {
	Email             string
	Password          string
	ReturnSecureToken bool
}

type signInWithPasswordResponse struct {
	IDToken      string
	EMail        string
	RefreshToken string
	ExpiresIn    string
	localID      string
	Registered   bool
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("Is valid", func(t *testing.T) {
		ctrl, service := setupMockService(t)
		defer ctrl.Finish()
		tokenStr := "tokentoken"
		service.EXPECT().VerifyToken(gomock.Any(), tokenStr).Return(nil)
		r := setupRouter(service)
		called := false
		r.GET("/", func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tokenStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		//認証OK時は後続のハンドラが実行されること
		assert.True(t, called)
	})

	t.Run("Is invalid", func(t *testing.T) {
		ctrl, service := setupMockService(t)
		defer ctrl.Finish()
		tokenStr := "tokentoken"
		service.EXPECT().VerifyToken(gomock.Any(), tokenStr).Return(errors.New("error"))
		r := setupRouter(service)
		called := false
		r.GET("/", func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tokenStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		//認証NG時
		//ステータスコードが401 Unauthorizedとなること
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		//後続のハンドラが実行されないこと
		assert.False(t, called)
	})
}

func setupMockService(t *testing.T) (*gomock.Controller, *mocks.MockAuthService) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockAuthService(ctrl)

	return ctrl, service
}

func setupRouter(service services.AuthService) *gin.Engine {
	r := gin.New()

	r.Use(NewAuthenticationMiddleware(service))

	return r
}

func signIn(credential *credential) *signInWithPasswordResponse {
	reqBody, _ := json.Marshal(credential)
	url := fmt.Sprintf(authorizationURL, "accounts:signUp", "AIzaSyDpuLWMl8KFRwP1ZcFUxuYhHAqTbcQibCw")
	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	var resObj signInWithPasswordResponse
	json.Unmarshal(resBody, &resObj)

	return &resObj
}
