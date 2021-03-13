package services

import (
	"context"
	"errors"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	myErr "github.com/edy4c7/darkpot-school-works/internal/errors"
	"github.com/edy4c7/darkpot-school-works/internal/repositories"
	"github.com/edy4c7/darkpot-school-works/internal/tools"
	"github.com/form3tech-oss/jwt-go"
)

type userKeyType string

const userKey userKeyType = "user"
const subjectKey string = "sub"
const cannotBeNullMessage = "%s can't be null"

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
		panic("Transaction runner can't be nil")
	}
	if worksRepo == nil {
		panic("Works repository can't be nil")
	}
	if activitiesRepo == nil {
		panic("Activities repository can't be nil")
	}
	if uuidGenerator == nil {
		panic("UUID generator can't be nil")
	}
	if fileUploader == nil {
		panic("File uploader can't be nil")
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
	return r.worksRepository.GetAll(ctx)
}

//FindByID は、指定したIDの作品を取得する
func (r *WorksServiceImpl) FindByID(ctx context.Context, id uint64) (*entities.Work, error) {
	result, err := r.worksRepository.FindByID(ctx, id)

	if err != nil {
		var dbErr *myErr.DataBaseAccessError
		if errors.As(err, &dbErr) {
			switch dbErr.Code() {
			case myErr.RecordNotFound:
				return nil, myErr.NewApplicationError(myErr.Code(myErr.DSWE01), myErr.Cause(err))
			}
		}

		return nil, myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	return result, nil
}

//Save は、作品の投稿及び更新を行う
func (r *WorksServiceImpl) Save(ctx context.Context, id uint64, bean *beans.WorksFormBean) error {
	token, ok := ctx.Value(userKey).(*jwt.Token)
	if !ok {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99),
			myErr.Message("Failed to load JWT"))
	}
	clm, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99),
			myErr.Message("Failed to load Claims"))
	}
	author, ok := clm[subjectKey].(string)

	thumbFileName := r.uuidGenerator.Generate()
	if err := r.fileUploader.Upload(thumbFileName, bean.Thumbnail); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	contentFileName := r.uuidGenerator.Generate()
	if err := r.fileUploader.Upload(contentFileName, bean.Content); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	err := r.transactionRunner.Run(ctx, func(ctx context.Context) error {
		w := &entities.Work{
			ID:                id,
			Author:            author,
			Title:             bean.Title,
			Description:       bean.Description,
			ThumbnailFileName: thumbFileName,
			ContentFileName:   contentFileName,
		}
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
		var dbErr *myErr.DataBaseAccessError
		if errors.As(err, &dbErr) {
			switch dbErr.Code() {
			case myErr.RecordNotFound:
				return myErr.NewApplicationError(myErr.Code(myErr.DSWE01), myErr.Cause(err))
			}
		}
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	return nil
}

//DeleteByID は、指定したIDの作品を削除する
func (r *WorksServiceImpl) DeleteByID(ctx context.Context, id uint64) error {
	w, err := r.worksRepository.FindByID(ctx, id)
	if err != nil {
		var dbErr *myErr.DataBaseAccessError
		if errors.As(err, &dbErr) {
			switch dbErr.Code() {
			case myErr.RecordNotFound:
				return myErr.NewApplicationError(myErr.Code(myErr.DSWE01), myErr.Cause(err))
			}
		}

		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	if err := r.fileUploader.Delete(w.ThumbnailFileName); err != nil {
		return myErr.NewApplicationError(myErr.Code(myErr.DSWE99), myErr.Cause(err))
	}

	if err := r.fileUploader.Delete(w.ContentFileName); err != nil {
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
