package services

import (
	"context"
	"errors"
	"testing"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/entities"
	"github.com/edy4c7/works-uploader/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		form := &beans.UserFormBean{
			ID:       "hogehoge",
			Name: "Fugafugao",
			Nickname: "hogetaro",
		}

		user := &entities.User{
			ID:       form.ID,
			Name: form.Name,
			Nickname: form.Nickname,
		}

		repo := mocks.NewMockUsersRepository(ctrl)
		repo.EXPECT().Save(ctx, user).Return(nil)
		service := &UsersServiceImpl{repository: repo}

		err := service.Save(ctx, form)

		assert.Nil(t, err)
	})

	t.Run("Is error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		form := &beans.UserFormBean{
			ID:       "hogehoge",
			Name: "Fugafugao",
			Nickname: "hogetaro",
		}

		repo := mocks.NewMockUsersRepository(ctrl)
		errExpect := errors.New("error")
		repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errExpect)
		service := &UsersServiceImpl{repository: repo}

		errActual := service.Save(ctx, form)

		if errActual != nil {
			assert.True(t, errors.Is(errActual, errExpect))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}
