package repositories

import (
	"context"

	"github.com/edy4c7/works-uploader/internal/entities"
)

type UsersRepository interface {
	Save(context.Context, *entities.User) error
}
