package infrastructures

import (
	"context"

	"github.com/edy4c7/works-uploader/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsersRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UsersRepositoryImpl {
	return &UsersRepositoryImpl{db: db}
}

func (r *UsersRepositoryImpl) Save(ctx context.Context, user *entities.User) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "nickname", "picture", "updated_at"}),
	}).Create(user).Error
}
