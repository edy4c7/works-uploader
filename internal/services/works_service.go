package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	myErr "github.com/edy4c7/darkpot-school-works/internal/errors"
	"github.com/edy4c7/darkpot-school-works/internal/repositories"
	"github.com/edy4c7/darkpot-school-works/internal/tools"
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

//WorksService は、作品管理機能のインターフェースを定義する
type WorksService interface {
	GetAll(context.Context) ([]*entities.Work, error)
	FindByID(context.Context, uint64) (*entities.Work, error)
	Save(context.Context, uint64, *beans.WorksFormBean) error
	DeleteByID(context.Context, uint64) error
}

//WorksServiceImpl は、作品管理機能を実装する
type WorksServiceImpl struct {
	transactionRunner    repositories.TransactionRunner
	worksRepository      repositories.WorksRepository
	activitiesRepository repositories.ActivitiesRepository
	uuidGenerator        tools.UUIDGenerator
	fileUploader         tools.FileUploader
}

//NewWorksServiceImpl は、TransuctionRunner、リポジトリオブジェクトを指定し、WorksServiceImplの新しいインスタンスを生成する
func NewWorksServiceImpl(
	tranRnr repositories.TransactionRunner,
	worksRepo repositories.WorksRepository,
	activitiesRepo repositories.ActivitiesRepository,
	uuidGenerator tools.UUIDGenerator,
	fileUploader tools.FileUploader,
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
func (r *WorksServiceImpl) GetAll(ctx context.Context) ([]*entities.Work, error) {
	result, err := r.worksRepository.GetAll(ctx)
	if err != nil {
		return nil, myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}
	return result, err
}

//FindByID は、指定したIDの作品を取得する
func (r *WorksServiceImpl) FindByID(ctx context.Context, id uint64) (*entities.Work, error) {
	result, err := r.worksRepository.FindByID(ctx, id)

	if err != nil {
		var dbErr *myErr.RecordNotFoundError
		if errors.As(err, &dbErr) {
			return nil, myErr.NewApplicationError(myErr.Code(myErr.DSWE01), myErr.Cause(err))
		}

		return nil, myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	return result, nil
}

//Save は、作品の投稿及び更新を行う
func (r *WorksServiceImpl) Save(ctx context.Context, id uint64, bean *beans.WorksFormBean) error {
	token, ok := ctx.Value(userKey).(*jwt.Token)
	if !ok {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99))
	}
	clm, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99))
	}
	author, ok := clm[subjectKey].(string)
	if !ok {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99))
	}

	w := &entities.Work{
		ID:          id,
		Type:        bean.Type,
		Author:      author,
		Title:       bean.Title,
		Description: bean.Description,
	}

	if bean.Type == constants.ContentTypeFile {
		thumbURL, err := r.fileUploader.Upload(r.uuidGenerator.Generate(), bean.Thumbnail)
		if err != nil {
			return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
		}
		w.ThumbnailURL = thumbURL

		contentURL, err := r.fileUploader.Upload(r.uuidGenerator.Generate(), bean.Content)
		if err != nil {
			return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
		}
		w.ContentURL = contentURL
	} else {
		w.ContentURL = bean.ContentURL
	}

	err := r.transactionRunner.Run(ctx, func(ctx context.Context) error {
		if err := r.worksRepository.Save(ctx, w); err != nil {
			return err
		}

		actType := constants.ActivityType(0)

		if id == 0 {
			actType = constants.ActivityAdded
		} else {
			actType = constants.ActivityUpdated
		}

		act := &entities.Activity{
			Type: actType,
			User: author,
			Work: w,
		}
		if err := r.activitiesRepository.Save(ctx, act); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		var dbErr *myErr.RecordNotFoundError
		if errors.As(err, &dbErr) {
			return myErr.NewApplicationError(myErr.Code(myErr.DSWE01), myErr.Cause(err))
		}
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	return nil
}

//DeleteByID は、指定したIDの作品を削除する
func (r *WorksServiceImpl) DeleteByID(ctx context.Context, id uint64) error {
	w, err := r.worksRepository.FindByID(ctx, id)
	if err != nil {
		var dbErr *myErr.RecordNotFoundError
		if errors.As(err, &dbErr) {
			return myErr.NewApplicationError(myErr.Code(myErr.DSWE01), myErr.Cause(err))
		}

		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	if err := r.fileUploader.Delete(w.ThumbnailURL); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	if err := r.fileUploader.Delete(w.ContentURL); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	err = r.transactionRunner.Run(ctx, func(c context.Context) error {
		return r.worksRepository.DeleteByID(ctx, id)
	})

	if err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	return nil
}
