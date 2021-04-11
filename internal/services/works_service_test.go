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

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().GetAll(gomock.Eq(ctx)).Return(data, nil)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

		result, err := service.GetAll(ctx)

		assert.Equal(t, data, result)
		assert.Nil(t, err)
	})

	t.Run("Is Error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)
		errExpect := errors.New("error")

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().GetAll(gomock.Eq(ctx)).Return(nil, errExpect)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

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

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Eq(ctx), id).Return(data, nil)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

		result, err := service.FindByID(ctx, id)

		assert.Equal(t, data, result)
		assert.Nil(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		errExpect := myErr.NewRecordNotFoundError("", nil)

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Any()).Return(nil, errExpect)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

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

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Any()).Return(nil, errExpect)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		fileUploader := mocks.NewMockFileUploader(ctrl)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		work := &entities.Work{
			Type:        form.Type,
			Title:       form.Title,
			Author:      subject,
			Description: form.Description,
			ContentURL:  form.ContentURL,
			Version:     initialVersion,
		}
		worksRepo.EXPECT().Save(gomock.Eq(ctx), work)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		actRepo.EXPECT().Save(gomock.Eq(ctx), &entities.Activity{
			Type: constants.ActivityAdded,
			User: subject,
			Work: work,
		})

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().Return(thumbnailFileName)
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(thumbnailFileName, form.Thumbnail).Return(thumbnailURL, nil)

		uuidGenerator.EXPECT().Generate().Return(contentFileName)
		fileUploader.EXPECT().Upload(contentFileName, form.Content).Return(contentURL, nil)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		work := &entities.Work{
			Type:         form.Type,
			Title:        form.Title,
			Author:       subject,
			Description:  form.Description,
			ThumbnailURL: thumbnailURL,
			ContentURL:   contentURL,
			Version:      initialVersion,
		}
		worksRepo.EXPECT().Save(gomock.Eq(ctx), work)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		actRepo.EXPECT().Save(gomock.Eq(ctx), &entities.Activity{
			Type: constants.ActivityAdded,
			User: subject,
			Work: work,
		})

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

		assert.Nil(t, service.Create(ctx, form))
	})

	t.Run("Fail to extract token", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		service := &WorksServiceImpl{}

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

		service := &WorksServiceImpl{}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)

		service := &WorksServiceImpl{
			uuidGenerator: uuidGenerator,
			fileUploader:  fileUploader,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("https://example.com", nil)

		uuidGenerator.EXPECT().Generate()
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)

		service := &WorksServiceImpl{
			uuidGenerator: uuidGenerator,
			fileUploader:  fileUploader,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return expect
			})

		service := &WorksServiceImpl{
			uuidGenerator:     uuidGenerator,
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		fileUploader := mocks.NewMockFileUploader(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		expect := errors.New("error")
		worksRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			uuidGenerator:     uuidGenerator,
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		fileUploader := mocks.NewMockFileUploader(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().Save(gomock.Any(), gomock.Any())
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		expect := errors.New("error")
		actRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		fileUploader := mocks.NewMockFileUploader(ctrl)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(ctx, id).Return(&entities.Work{
			ID:      id,
			Version: version,
		}, nil)
		work := &entities.Work{
			Type:    form.Type,
			ID:      id,
			Author:  subject,
			Version: version + 1,
		}
		worksRepo.EXPECT().Save(gomock.Eq(ctx), work)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		actRepo.EXPECT().Save(gomock.Eq(ctx), &entities.Activity{
			Type: constants.ActivityUpdated,
			User: subject,
			Work: work,
		})

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().Return(thumbnailFileName)
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(thumbnailFileName, form.Thumbnail).Return(thumbnailURL, nil)

		uuidGenerator.EXPECT().Generate().Return(contentFileName)
		fileUploader.EXPECT().Upload(contentFileName, form.Content).Return(contentURL, nil)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Eq(ctx), gomock.Any()).Return(&entities.Work{
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
		worksRepo.EXPECT().Save(gomock.Eq(ctx), work)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		actRepo.EXPECT().Save(gomock.Eq(ctx), &entities.Activity{
			Type: constants.ActivityUpdated,
			User: subject,
			Work: work,
		})

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		fileUploader := mocks.NewMockFileUploader(ctrl)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), id).Return(&entities.Work{
			ID:      id,
			Version: version + 1,
		}, nil)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

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

		service := &WorksServiceImpl{}

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

		service := &WorksServiceImpl{}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)

		service := &WorksServiceImpl{
			uuidGenerator: uuidGenerator,
			fileUploader:  fileUploader,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("https://example.com", nil)

		uuidGenerator.EXPECT().Generate()
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)

		service := &WorksServiceImpl{
			uuidGenerator: uuidGenerator,
			fileUploader:  fileUploader,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return expect
			})

		service := &WorksServiceImpl{
			uuidGenerator:     uuidGenerator,
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		expect := myErr.NewRecordNotFoundError("", nil)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, expect)

		service := &WorksServiceImpl{
			uuidGenerator:     uuidGenerator,
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		expect := errors.New("error")
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
		worksRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			uuidGenerator:     uuidGenerator,
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}
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

		uuidGenerator := mocks.NewMockUUIDGenerator(ctrl)
		uuidGenerator.EXPECT().Generate().AnyTimes()
		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
		worksRepo.EXPECT().Save(gomock.Any(), gomock.Any())
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		expect := errors.New("error")
		actRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

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

		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Delete(thumbFileName)
		fileUploader.EXPECT().Delete(contentFileName)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Eq(ctx), id).Return(&entities.Work{
			ThumbnailURL: thumbFileName,
			ContentURL:   contentFileName,
		}, nil)
		worksRepo.EXPECT().DeleteByID(gomock.Eq(ctx), id)

		service := &WorksServiceImpl{
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		assert.Nil(t, service.DeleteByID(ctx, id))
	})

	t.Run("Not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).AnyTimes()

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		expect := myErr.NewRecordNotFoundError("", nil)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, expect)

		service := &WorksServiceImpl{
			fileUploader:    fileUploader,
			worksRepository: worksRepo,
		}

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

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, expect)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

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

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)

		fileUploader := mocks.NewMockFileUploader(ctrl)
		expect := errors.New("Failed to delete")
		fileUploader.EXPECT().Delete(gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			fileUploader:    fileUploader,
			worksRepository: worksRepo,
		}

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

		fileUploader := mocks.NewMockFileUploader(ctrl)
		expect := errors.New("Failed to delete")
		fileUploader.EXPECT().Delete(gomock.Any()).Return(nil)
		fileUploader.EXPECT().Delete(gomock.Any()).Return(expect)

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)

		service := &WorksServiceImpl{
			fileUploader:    fileUploader,
			worksRepository: worksRepo,
		}

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

		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Delete(gomock.Any()).AnyTimes()

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		expect := errors.New("error")
		tranRunner.EXPECT().Run(gomock.Eq(ctx), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

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

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)

		fileUploader := mocks.NewMockFileUploader(ctrl)
		fileUploader.EXPECT().Delete(gomock.Any()).AnyTimes()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		expect := errors.New("error")
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo.EXPECT().DeleteByID(gomock.Eq(ctx), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

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
