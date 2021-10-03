package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/common/constants"
	"github.com/edy4c7/works-uploader/internal/entities"
	myErr "github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/mocks"
	"github.com/edy4c7/works-uploader/internal/repositories"
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
		uploader := mocks.NewMockStorageClient(ctrl)

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
		uploader := mocks.NewMockStorageClient(ctrl)

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
		uploader := mocks.NewMockStorageClient(ctrl)

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
		uploader := mocks.NewMockStorageClient(ctrl)

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
		uploader := mocks.NewMockStorageClient(ctrl)

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

		offset, limit := 1, 100
		total := int64(200)

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().CountAll(gomock.Eq(ctx)).Return(total, nil)
		worksRepo.EXPECT().GetAll(gomock.Eq(ctx), offset, limit).Return(data, nil)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

		result, err := service.GetAll(ctx, offset, limit)

		pagination := &beans.PaginationBean{
			TotalItems: total,
			Offset:     offset,
		}
		for _, v := range data {
			pagination.Items = append(pagination.Items, v)
		}

		assert.Equal(t, pagination, result)
		assert.Nil(t, err)
	})

	t.Run("Is Error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)
		errExpect := errors.New("error")

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().CountAll(gomock.Eq(ctx)).Return(int64(100), nil)
		worksRepo.EXPECT().GetAll(gomock.Eq(ctx), gomock.Any(), gomock.Any()).Return(nil, errExpect)

		service := &WorksServiceImpl{
			worksRepository: worksRepo,
		}

		result, err := service.GetAll(ctx, 0, 100)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errExpect))
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
			assert.Equal(t, appErr.Code(), myErr.WUE01)
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
			assert.Equal(t, appErr.Code(), myErr.WUE99)
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
		fileUploader := mocks.NewMockStorageClient(ctrl)

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
			AuthorID:      subject,
			Description: form.Description,
			ContentURL:  form.ContentURL,
			Version:     initialVersion,
		}
		worksRepo.EXPECT().Create(gomock.Eq(ctx), work)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		actRepo.EXPECT().Create(gomock.Eq(ctx), &entities.Activity{
			Type: constants.ActivityAdded,
			UserID: subject,
			Work: work,
		})

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

		res, err := service.Create(ctx, form)
		assert.Nil(t, err)
		assert.Equal(t, work, res)
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
		fileUploader := mocks.NewMockStorageClient(ctrl)
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
			AuthorID:       subject,
			Description:  form.Description,
			ThumbnailURL: thumbnailURL,
			ContentURL:   contentURL,
			Version:      initialVersion,
		}
		worksRepo.EXPECT().Create(gomock.Eq(ctx), work)

		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		actRepo.EXPECT().Create(gomock.Eq(ctx), &entities.Activity{
			Type: constants.ActivityAdded,
			UserID: subject,
			Work: work,
		})

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

		res, err := service.Create(ctx, form)
		assert.Nil(t, err)
		assert.Equal(t, work, res)
	})

	t.Run("Fail to extract token", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		service := &WorksServiceImpl{}

		_, err := service.Create(ctx, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to extract claims", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		//lint:ignore SA1029 can use string only
		ctx = context.WithValue(ctx, userKey, &jwt.Token{})

		service := &WorksServiceImpl{}

		_, err := service.Create(ctx, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", err)
		}
	})

	t.Run("Fail to extract subject", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		//lint:ignore SA1029 can use string only
		ctx = context.WithValue(ctx, userKey, &jwt.Token{
			Claims: make(jwt.MapClaims),
		})

		service := &WorksServiceImpl{}

		_, err := service.Create(ctx, nil)

		assert.Error(t, err)
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
		fileUploader := mocks.NewMockStorageClient(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)

		service := &WorksServiceImpl{
			uuidGenerator: uuidGenerator,
			fileUploader:  fileUploader,
		}

		_, actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
		fileUploader := mocks.NewMockStorageClient(ctrl)
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("https://example.com", nil)

		uuidGenerator.EXPECT().Generate()
		fileUploader.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", expect)

		service := &WorksServiceImpl{
			uuidGenerator: uuidGenerator,
			fileUploader:  fileUploader,
		}

		_, actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
		fileUploader := mocks.NewMockStorageClient(ctrl)
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

		_, actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
		fileUploader := mocks.NewMockStorageClient(ctrl)
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
		worksRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			uuidGenerator:     uuidGenerator,
			fileUploader:      fileUploader,
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		_, actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
		fileUploader := mocks.NewMockStorageClient(ctrl)
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
		worksRepo.EXPECT().Create(gomock.Any(), gomock.Any())
		actRepo := mocks.NewMockActivitiesRepository(ctrl)
		expect := errors.New("error")
		actRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			uuidGenerator:        uuidGenerator,
			fileUploader:         fileUploader,
			transactionRunner:    tranRunner,
			worksRepository:      worksRepo,
			activitiesRepository: actRepo,
		}

		_, actual := service.Create(ctx, form)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().DeleteByID(gomock.Eq(ctx), id)

		service := &WorksServiceImpl{
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		assert.Nil(t, service.DeleteByID(ctx, id))
	})

	t.Run("Not found", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		expect := myErr.NewRecordNotFoundError("", nil)
		worksRepo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect), "%w", actual)
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE01, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Fail to find record", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		ctx = setupContext(ctx)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		expect := errors.New("Failed to delete")

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect))
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Failed to run transaction", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		worksRepo := mocks.NewMockWorksRepository(ctrl)

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		expect := errors.New("error")
		tranRunner.EXPECT().Run(gomock.Eq(ctx), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect), "%w", actual)
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
		} else {
			assert.Failf(t, "Invalid error type", "%w", actual)
		}
	})

	t.Run("Failed to delete record", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		tranRunner := mocks.NewMockTransactionRunner(ctrl)
		expect := errors.New("error")
		tranRunner.
			EXPECT().
			Run(gomock.Eq(ctx), gomock.Any()).
			DoAndReturn(func(ctx context.Context, tranFunc repositories.TransactionFunction) error {
				return tranFunc(ctx)
			})

		worksRepo := mocks.NewMockWorksRepository(ctrl)
		worksRepo.EXPECT().DeleteByID(gomock.Eq(ctx), gomock.Any()).Return(expect)

		service := &WorksServiceImpl{
			transactionRunner: tranRunner,
			worksRepository:   worksRepo,
		}

		actual := service.DeleteByID(ctx, 1)

		assert.True(t, errors.Is(actual, expect), "%w", actual)
		var appErr *myErr.ApplicationError
		if errors.As(actual, &appErr) {
			assert.Equal(t, myErr.WUE99, appErr.Code())
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
