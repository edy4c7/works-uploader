package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/common/constants"
	"github.com/edy4c7/works-uploader/internal/entities"
	myErr "github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/lib"
	"github.com/edy4c7/works-uploader/internal/repositories"
	"github.com/form3tech-oss/jwt-go"
)

const userKey string = "user"
const subjectKey string = "sub"
const cannotBeNullMessage = "%s can't be null"
const msgTransactionRunner = "transaction runner"
const msgWorksRepository = "works repository"
const msgActivitiesRepository = "activities repository"
const msgUUIDGenerator = "UUID generator"
const msgFileUploader = "file uploader"
const initialVersion uint = 1

//WorksService は、作品管理機能のインターフェースを定義する
type WorksService interface {
	GetAll(ctx context.Context, offset int, limit int) (*beans.PaginationBean, error)
	FindByID(context.Context, uint64) (*entities.Work, error)
	Create(context.Context, *beans.WorksFormBean) (*entities.Work, error)
	Update(context.Context, uint64, *beans.WorksFormBean) (*entities.Work, error)
	DeleteByID(context.Context, uint64) error
}

//WorksServiceImpl は、作品管理機能を実装する
type WorksServiceImpl struct {
	transactionRunner    repositories.TransactionRunner
	worksRepository      repositories.WorksRepository
	activitiesRepository repositories.ActivitiesRepository
	uuidGenerator        lib.UUIDGenerator
	fileUploader         lib.FileUploader
}

//NewWorksServiceImpl は、TransuctionRunner、リポジトリオブジェクトを指定し、WorksServiceImplの新しいインスタンスを生成する
func NewWorksServiceImpl(
	tranRnr repositories.TransactionRunner,
	worksRepo repositories.WorksRepository,
	activitiesRepo repositories.ActivitiesRepository,
	uuidGenerator lib.UUIDGenerator,
	fileUploader lib.FileUploader,
) *WorksServiceImpl {

	if tranRnr == nil {
		panic(fmt.Sprintf(cannotBeNullMessage, msgTransactionRunner))
	}
	if worksRepo == nil {
		panic(fmt.Sprintf(cannotBeNullMessage, msgWorksRepository))
	}
	if activitiesRepo == nil {
		panic(fmt.Sprintf(cannotBeNullMessage, msgActivitiesRepository))
	}
	if uuidGenerator == nil {
		panic(fmt.Sprintf(cannotBeNullMessage, msgUUIDGenerator))
	}
	if fileUploader == nil {
		panic(fmt.Sprintf(cannotBeNullMessage, msgFileUploader))
	}

	return &WorksServiceImpl{
		transactionRunner:    tranRnr,
		worksRepository:      worksRepo,
		activitiesRepository: activitiesRepo,
		uuidGenerator:        uuidGenerator,
		fileUploader:         fileUploader,
	}
}

//GetAll は、作品の全件取得を行う
func (r *WorksServiceImpl) GetAll(ctx context.Context, offset int, limit int) (*beans.PaginationBean, error) {
	count, err := r.worksRepository.CountAll(ctx)
	if err != nil {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	result, err := r.worksRepository.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	pagination := &beans.PaginationBean{
		TotalItems: count,
		Offset:     offset,
		Items:      make([]interface{}, 0, limit),
	}

	for _, v := range result {
		pagination.Items = append(pagination.Items, v)
	}

	return pagination, err
}

//FindByID は、指定したIDの作品を取得する
func (r *WorksServiceImpl) FindByID(ctx context.Context, id uint64) (*entities.Work, error) {
	result, err := r.worksRepository.FindByID(ctx, id)

	if err != nil {
		var dbErr *myErr.RecordNotFoundError
		if errors.As(err, &dbErr) {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE01), myErr.Cause(err))
		}

		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	return result, nil
}

func (r *WorksServiceImpl) Create(ctx context.Context, bean *beans.WorksFormBean) (*entities.Work, error) {
	token, ok := ctx.Value(userKey).(*jwt.Token)
	if !ok {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99))
	}
	clm, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99))
	}
	author, ok := clm[subjectKey].(string)
	if !ok {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99))
	}

	w := &entities.Work{
		Type:        bean.Type,
		Author:      author,
		Title:       bean.Title,
		Description: bean.Description,
		Version:     initialVersion,
	}

	if bean.Type == constants.ContentTypeFile {
		thumbURL, err := r.fileUploader.Upload(r.uuidGenerator.Generate(), bean.Thumbnail)
		if err != nil {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}
		w.ThumbnailURL = thumbURL

		contentURL, err := r.fileUploader.Upload(r.uuidGenerator.Generate(), bean.Content)
		if err != nil {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}
		w.ContentURL = contentURL
	} else {
		w.ContentURL = bean.ContentURL
	}

	err := r.transactionRunner.Run(ctx, func(ctx context.Context) error {
		if err := r.worksRepository.Save(ctx, w); err != nil {
			return err
		}

		act := &entities.Activity{
			Type: constants.ActivityAdded,
			User: author,
			Work: w,
		}
		if err := r.activitiesRepository.Create(ctx, act); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	return w, nil
}

//Update は、作品の投稿及び更新を行う
func (r *WorksServiceImpl) Update(ctx context.Context, id uint64, bean *beans.WorksFormBean) (*entities.Work, error) {
	token, ok := ctx.Value(userKey).(*jwt.Token)
	if !ok {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99))
	}
	clm, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99))
	}
	author, ok := clm[subjectKey].(string)
	if !ok {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99))
	}

	var thumbURL string
	var contentURL string
	if bean.Type == constants.ContentTypeFile {
		var err error
		thumbURL, err = r.fileUploader.Upload(r.uuidGenerator.Generate(), bean.Thumbnail)
		if err != nil {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}

		contentURL, err = r.fileUploader.Upload(r.uuidGenerator.Generate(), bean.Content)
		if err != nil {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}
	} else {
		contentURL = bean.ContentURL
	}

	var work *entities.Work

	err := r.transactionRunner.Run(ctx, func(ctx context.Context) error {
		w, err := r.worksRepository.FindByID(ctx, id)

		if err != nil {
			var rnfErr *myErr.RecordNotFoundError
			if errors.As(err, &rnfErr) {
				return myErr.NewApplicationError(myErr.Code(myErr.WUE01), myErr.Cause(err))
			}

			return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}

		if w.Version != bean.Version {
			return myErr.NewApplicationError(myErr.Code(myErr.WUE02))
		}

		w.Type = bean.Type
		w.Author = author
		w.Title = bean.Title
		w.Description = bean.Description
		w.ThumbnailURL = thumbURL
		w.ContentURL = contentURL
		w.Version = w.Version + 1

		if err := r.worksRepository.Save(ctx, w); err != nil {
			return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}

		act := &entities.Activity{
			Type: constants.ActivityUpdated,
			User: author,
			Work: w,
		}
		if err := r.activitiesRepository.Create(ctx, act); err != nil {
			return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
		}

		work = w

		return nil
	})

	if err != nil {
		var appErr *myErr.ApplicationError
		if errors.As(err, &appErr) {
			return nil, err
		}
		return nil, myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	return work, nil
}

//DeleteByID は、指定したIDの作品を削除する
func (r *WorksServiceImpl) DeleteByID(ctx context.Context, id uint64) error {
	w, err := r.worksRepository.FindByID(ctx, id)
	if err != nil {
		var dbErr *myErr.RecordNotFoundError
		if errors.As(err, &dbErr) {
			return myErr.NewApplicationError(myErr.Code(myErr.WUE01), myErr.Cause(err))
		}

		return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	err = r.transactionRunner.Run(ctx, func(ctx context.Context) error {
		return r.worksRepository.DeleteByID(ctx, id)
	})

	if err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	if err := r.fileUploader.Delete(w.ThumbnailURL); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	if err := r.fileUploader.Delete(w.ContentURL); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.WUE99), myErr.Cause(err))
	}

	return nil
}
