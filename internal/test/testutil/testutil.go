package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func BindFormToObject(req *http.Request, obj interface{}) error {
	return (&gin.Context{Request: req}).ShouldBind(obj)
}

func ExecuteHandler(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
}
