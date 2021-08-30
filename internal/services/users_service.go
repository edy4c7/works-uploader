package services

import (
	"context"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/entities"
	"github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/repositories"
)

type UsersService interface {
	Save(context.Context, *beans.UserFormBean) error
}

type UsersServiceImpl struct {
	repository repositories.UsersRepository
}

func NewUsersServiceImpl(repo repositories.UsersRepository) *UsersServiceImpl {
	return &UsersServiceImpl{repository: repo}
}

func (r *UsersServiceImpl) Save(ctx context.Context, form *beans.UserFormBean) error {
	user := &entities.User{
		ID:       form.ID,
		Nickname: form.Nickname,
		Name:     form.Name,
	}

	if err := r.repository.Save(ctx, user); err != nil {
		return errors.NewApplicationError(errors.Code(errors.WUE99), errors.Cause(err))
	}

	return nil
}
