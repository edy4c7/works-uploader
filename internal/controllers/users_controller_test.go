package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	endpoint := "/"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		form := &beans.UserFormBean{
			ID:       "hogehoge",
			Name: "Fugafugao",
			Nickname: "hogetaro",
		}

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)
		json, _ := json.Marshal(form)
		req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(json))
		req = req.WithContext(ctx)
		ginCtx.Request = req

		service := mocks.NewMockUsersService(ctrl)
		service.EXPECT().Save(ctx, form).Return(nil)

		userCtrl := &UsersController{service: service}
		r.POST(endpoint, userCtrl.Save)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Is Fail", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		form := &beans.UserFormBean{
			ID:       "hogehoge",
			Name: "Fugafugao",
			Nickname: "hogetaro",
		}

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)
		json, _ := json.Marshal(form)
		req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(json))
		req = req.WithContext(ctx)
		ginCtx.Request = req

		errExpect := errors.New("error")
		service := mocks.NewMockUsersService(ctrl)
		service.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errExpect)

		userCtrl := &UsersController{service: service}
		r.POST(endpoint, userCtrl.Save)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}
