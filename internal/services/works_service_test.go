package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	myErr "github.com/edy4c7/darkpot-school-works/internal/errors"
	"github.com/edy4c7/darkpot-school-works/internal/mocks"
	"github.com/edy4c7/darkpot-school-works/internal/repositories"
	"github.com/form3tech-oss/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewWorksServiceImpl(t *testing.T) {
	t.Run("Is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tr := mocks.NewMockTransactionRunner(ctrl)
		workRepo := mocks.NewMockWorksRepository(ctrl)
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uploader := mocks.NewMockFileUploader(ctrl)

		service := NewWorksServiceImpl(tr, workRepo, actRepo, uuidGenerator, uploader)

		assert.Same(t, service.transactionRunner, tr)
		assert.Same(t, service.worksRepository, workRepo)
		assert.Same(t, service.activitiesRepository, actRepo)
		assert.Same(t, service.uuidGenerator, uuidGenerator)
		assert.Same(t, service.fileUploader, uploader)
	})

	t.Run("Transaction runner is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		workRepo := mocks.NewMockWorksRepository(ctrl)
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uploader := mocks.NewMockFileUploader(ctrl)

		assert.Panics(t, func() {
			NewWorksServiceImpl(nil, workRepo, actRepo, uuidGenerator, uploader)
		})
	})

	t.Run("Works repository is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tr := mocks.NewMockTransactionRunner(ctrl)
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uploader := mocks.NewMockFileUploader(ctrl)

		assert.Panics(t, func() {
			NewWorksServiceImpl(tr, nil, actRepo, uuidGenerator, uploader)
		})
	})

	t.Run("Activity repository is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tr := mocks.NewMockTransactionRunner(ctrl)
		workRepo := mocks.NewMockWorksRepository(ctrl)
		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uploader := mocks.NewMockFileUploader(ctrl)

		assert.Panics(t, func() {
			NewWorksServiceImpl(tr, workRepo, nil, uuidGenerator, uploader)
		})
	})

	t.Run("UUID generator is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tr := mocks.NewMockTransactionRunner(ctrl)
		workRepo := mocks.NewMockWorksRepository(ctrl)
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		uploader := mocks.NewMockFileUploader(ctrl)

		assert.Panics(t, func() {
			NewWorksServiceImpl(tr, workRepo, actRepo, nil, uploader)
		})
	})

	t.Run("File uploader is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tr := mocks.NewMockTransactionRunner(ctrl)
		workRepo := mocks.NewMockWorksRepository(ctrl)
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)

		assert.Panics(t, func() {
			NewWorksServiceImpl(tr, workRepo, actRepo, uuidGenerator, nil)
		})
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		data := []*entities.Work{
			{
				ID:           0,
				Title:        "hoge",
				Description:  "hogehoge",
				ThumbnailURL: "https://example.com",
				ContentURL:   "https://example.com",
			},
			{
				ID:           1,
				Title:        "hoge",
				Description:  "hogehoge",
				ThumbnailURL: "https://example.com",
				ContentURL:   "https://example.com",
			},
		}

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().GetAll(gomock.Eq(ctx)).Return(data, nil)
			},
		)

		result, err := service.GetAll(ctx)

		assert.Equal(t, data, result)
		assert.Nil(t, err)
	})

	t.Run("Is Error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)
		errExpect := errors.New("error")

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().GetAll(gomock.Eq(ctx)).Return(nil, errExpect)
			},
		)

		result, err := service.GetAll(ctx)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})
}

func TestFindByID(t *testing.T) {
	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		var id uint64 = 1
		data := &entities.Work{
			ID:           id,
			Title:        "hoge",
			Description:  "hogehoge",
			ThumbnailURL: "https://example.com",
			ContentURL:   "https://example.com",
		}

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Eq(ctx), id).Return(data, nil)
			},
		)

		result, err := service.FindByID(ctx, id)

		assert.Equal(t, data, result)
		assert.Nil(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		errExpect := myErr.NewRecordNotFoundError("", nil)
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Any()).Return(nil, errExpect)
			},
		)
		result, err := service.FindByID(ctx, 1)

		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, appErr.Code(), myErr.DSWE01)
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Is error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		errExpect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Eq(ctx), gomock.Any()).Return(nil, errExpect)
			},
		)

		result, err := service.FindByID(ctx, 1)

		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, appErr.Code(), myErr.DSWE99)
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})
}
func TestCreate(t *testing.T) {
	t.Run("New with URL", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:        constants.ContentTypeURL,
			Title:       "hoge",
			Description: "hogehoge",
			ContentURL:  "https://example.com",
		}

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Eq(ctx), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				work := &entities.Work{
					Type:        form.Type,
					Title:       form.Title,
					Author:      subject,
					Description: form.Description,
					ContentURL:  form.ContentURL,
					Version:     initialVersion,
				}
				wso.worksRepository.EXPECT().Save(gomock.Eq(ctx), work)
				wso.activittyRepository.EXPECT().Save(gomock.Eq(ctx), &entities.Activity{
					Type: constants.ActivityAdded,
					User: subject,
					Work: work,
				})
			},
		)

		assert.Nil(t, service.Create(ctx, form))
	})

	t.Run("New with file", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:        constants.ContentTypeFile,
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		thumbnailFileName := "abcde12345"
		thumbnailURL := fmt.Sprintf("https://example.com/%s", thumbnailFileName)
		contentFileName := "fghij67890"
		contentURL := fmt.Sprintf("https://example.com/%s", contentFileName)

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().Return(thumbnailFileName)
				wso.fileUploader.EXPECT().Upload(thumbnailFileName, form.Thumbnail).Return(thumbnailURL, nil)
				wso.uuidGenerator.EXPECT().Generate().Return(contentFileName)
				wso.fileUploader.EXPECT().Upload(contentFileName, form.Content).Return(contentURL, nil)
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Eq(ctx), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				work := &entities.Work{
					Type:         form.Type,
					Title:        form.Title,
					Author:       subject,
					Description:  form.Description,
					ThumbnailURL: thumbnailURL,
					ContentURL:   contentURL,
					Version:      initialVersion,
				}
				wso.worksRepository.EXPECT().Save(gomock.Eq(ctx), work)
				wso.activittyRepository.EXPECT().Save(gomock.Eq(ctx), &entities.Activity{
					Type: constants.ActivityAdded,
					User: subject,
					Work: work,
				})
			},
		)
		assert.Nil(t, service.Create(ctx, form))
	})

	t.Run("Fail to extract token", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		service := expectWorksService(ctrl)

		err := service.Create(ctx, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to extract subject", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		//lint:ignore SA1029 can use string only
		ctx = context.WithValue(ctx, userKey, &jwt.Token{})

		service := expectWorksService(ctrl)

		err := service.Create(ctx, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to upload thumbnail", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:        constants.ContentTypeFile,
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("Failed to upload")

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)
			},
		)

		actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to upload content", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:        constants.ContentTypeFile,
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("Failed to upload")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				wso.uuidGenerator.EXPECT().Generate()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)
			},
		)

		actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to run transaction", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return expect
					})
			},
		)

		actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to save work", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expect)
			},
		)

		actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to save activity", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().Save(gomock.Any(), gomock.Any())
				wso.activittyRepository.
					EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(expect)
			},
		)

		actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})
}

func TestUpdate(t *testing.T) {
	version := uint(1)

	t.Run("Update with URL", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:    constants.ContentTypeURL,
			Version: version,
		}

		var id uint64 = 12345

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), id).Return(&entities.Work{
					ID:      id,
					Version: version,
				}, nil)
				work := &entities.Work{
					Type:    form.Type,
					ID:      id,
					Author:  subject,
					Version: version + 1,
				}
				wso.worksRepository.EXPECT().Save(gomock.Any(), work)
				wso.activittyRepository.EXPECT().Save(gomock.Any(), &entities.Activity{
					Type: constants.ActivityUpdated,
					User: subject,
					Work: work,
				})
			},
		)

		assert.Nil(t, service.Update(ctx, id, form))
	})

	t.Run("Update with file", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:      constants.ContentTypeFile,
			Thumbnail: &multipart.FileHeader{},
			Content:   &multipart.FileHeader{},
			Version:   version,
		}

		var id uint64 = 12345

		thumbnailFileName := "abcde12345"
		thumbnailURL := fmt.Sprintf("https://example.com/%s", thumbnailFileName)
		contentFileName := "fghij67890"
		contentURL := fmt.Sprintf("https://example.com/%s", contentFileName)

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().Return(thumbnailFileName)
				wso.fileUploader.EXPECT().Upload(thumbnailFileName, form.Thumbnail).Return(thumbnailURL, nil)
				wso.uuidGenerator.EXPECT().Generate().Return(contentFileName)
				wso.fileUploader.EXPECT().Upload(contentFileName, form.Content).Return(contentURL, nil)
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), id).Return(&entities.Work{
					ID:      id,
					Version: version,
				}, nil)
				work := &entities.Work{
					Type:         form.Type,
					ID:           id,
					Author:       subject,
					ThumbnailURL: thumbnailURL,
					ContentURL:   contentURL,
					Version:      version + 1,
				}
				wso.worksRepository.EXPECT().Save(gomock.Any(), work)
				wso.activittyRepository.EXPECT().Save(gomock.Any(), &entities.Activity{
					Type: constants.ActivityUpdated,
					User: subject,
					Work: work,
				})
			},
		)

		assert.Nil(t, service.Update(ctx, id, form))
	})

	t.Run("conflicted", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:    constants.ContentTypeURL,
			Version: version,
		}

		var id uint64 = 12345

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), id).Return(&entities.Work{
					ID:      id,
					Version: version + 1,
				}, nil)
			},
		)

		err := service.Update(ctx, id, form)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.DSWE02, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to extract token", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		service := expectWorksService(ctrl)

		err := service.Update(ctx, 0, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to extract subject", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		//lint:ignore SA1029 can use string only
		ctx = context.WithValue(ctx, userKey, &jwt.Token{})

		service := expectWorksService(ctrl)

		err := service.Update(ctx, 0, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to upload thumbnail", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:        constants.ContentTypeFile,
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("Failed to upload")

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)
			},
		)

		actual := service.Update(ctx, 0, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to upload content", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Type:        constants.ContentTypeFile,
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("Failed to upload")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("https://example.com", nil)
				wso.uuidGenerator.EXPECT().Generate()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)
			},
		)

		actual := service.Update(ctx, 0, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to run transaction", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Title:       "hoge",
			Description: "hogehoge",
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return expect
					})
			},
		)

		actual := service.Update(ctx, 0, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := myErr.NewRecordNotFoundError("", nil)
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, expect)
			},
		)

		actual := service.Update(ctx, 0, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE01, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to save work", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
				wso.worksRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expect)
			},
		)

		actual := service.Update(ctx, 0, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to save activity", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		form := &beans.WorksFormBean{
			Thumbnail: &multipart.FileHeader{
				Filename: "thumb01",
				Size:     1,
			},
			Content: &multipart.FileHeader{
				Filename: "content01",
				Size:     1,
			},
		}

		expect := errors.New("error")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.uuidGenerator.EXPECT().Generate().AnyTimes()
				wso.fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
				wso.worksRepository.EXPECT().Save(gomock.Any(), gomock.Any())
				wso.activittyRepository.
					EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(expect)
			},
		)

		actual := service.Update(ctx, 0, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})
}

func TestDeleteByID(t *testing.T) {
	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		var id uint64 = 1
		thumbFileName := "abcde12345"
		contentFileName := "fghij67890"

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Eq(ctx), id).Return(&entities.Work{
					ThumbnailURL: thumbFileName,
					ContentURL:   contentFileName,
				}, nil)
				wso.fileUploader.EXPECT().Delete(thumbFileName)
				wso.fileUploader.EXPECT().Delete(contentFileName)
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Eq(ctx), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().DeleteByID(gomock.Eq(ctx), id)
			},
		)

		assert.Nil(t, service.DeleteByID(ctx, id))
	})

	t.Run("Not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		expect := myErr.NewRecordNotFoundError("Not found", nil)

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, expect)
			},
		)

		actual := service.DeleteByID(ctx, 1)
		assert.True(t, errors.Is(actual, expect), "%w", actual)
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE01, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to find record", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		expect := errors.New("Failed to delete")

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, expect)
			},
		)

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to delete thumbnail", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		expect := errors.New("Failed to delete")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.
					EXPECT().
					FindByID(gomock.Any(), gomock.Any()).
					Return(&entities.Work{}, nil)
				wso.fileUploader.EXPECT().Delete(gomock.Any()).Return(expect)
			},
		)

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to delete content", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		expect := errors.New("Failed to upload")
		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
				wso.fileUploader.EXPECT().Delete(gomock.Any()).Return(nil)
				wso.fileUploader.EXPECT().Delete(gomock.Any()).Return(expect)
			},
		)

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Failed to run transaction", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		expect := errors.New("error")

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
				wso.fileUploader.EXPECT().Delete(gomock.Any())
				wso.fileUploader.EXPECT().Delete(gomock.Any())
				wso.transactionRunner.EXPECT().Run(gomock.Eq(ctx), gomock.Any()).Return(expect)
			},
		)

		actual := service.DeleteByID(ctx, 1)
		assert.True(t, errors.Is(actual, expect), "%w", actual)
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Failed to delete record", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		expect := errors.New("error")

		service := expectWorksService(ctrl,
			func(wso *worksServiceOptions) {
				wso.worksRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
				wso.fileUploader.EXPECT().Delete(gomock.Any())
				wso.fileUploader.EXPECT().Delete(gomock.Any())
				wso.transactionRunner.
					EXPECT().
					Run(gomock.Eq(ctx), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
						return tranFunc(ctx)
					})
				wso.worksRepository.EXPECT().DeleteByID(gomock.Eq(ctx), gomock.Any()).Return(expect)
			},
		)

		actual := service.DeleteByID(ctx, 1)
		assert.True(t, errors.Is(actual, expect), "%w", actual)
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.DSWE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})
}

const subject string = "subABC12345"

func setupContext(ctx context.Context) context.Context {
	//lint:ignore SA1029 can use string only
	return context.WithValue(ctx, userKey, &jwt.Token{
		Claims: jwt.MapClaims{
			"sub": subject,
		},
	})
}

type worksServiceOptions struct {
	transactionRunner   *mocks.MockTransactionRunner
	worksRepository     *mocks.MockWorksRepository
	activittyRepository *mocks.MockActivitiesRepository
	uuidGenerator       *mocks.MockUUIDGenerator
	fileUploader        *mocks.MockFileUploader
}

func expectWorksService(ctrl *gomock.Controller, opts ...func(*worksServiceOptions)) *WorksServiceImpl {
	wso := &worksServiceOptions{
		transactionRunner:   mocks.NewMockTransactionRunner(ctrl),
		worksRepository:     mocks.NewMockWorksRepository(ctrl),
		activittyRepository: mocks.NewMockActivitiesRepository(ctrl),
		uuidGenerator:       mocks.NewMockUUIDGenerator(ctrl),
		fileUploader:        mocks.NewMockFileUploader(ctrl),
	}

	for _, opt := range opts {
		opt(wso)
	}

	return NewWorksServiceImpl(
		wso.transactionRunner,
		wso.worksRepository,
		wso.activittyRepository,
		wso.uuidGenerator,
		wso.fileUploader,
	)
}
