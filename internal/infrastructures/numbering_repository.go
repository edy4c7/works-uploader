package infrastructures

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
)

const numberingCollectionName string = "numbering"

type NumberingRepositoryImpl struct {
	ref *firestore.CollectionRef
}

func NewNumberingRepositoryImpl(client *firestore.Client) *NumberingRepositoryImpl {
	return &NumberingRepositoryImpl{
		ref: client.Collection(numberingCollectionName),
	}
}

func (r *NumberingRepositoryImpl) FindByType(ctx context.Context, typ int) (*entities.Numbering, error) {
	ds, err := r.ref.Doc(strconv.Itoa(typ)).Get(ctx)
	if err != nil {
		return nil, err
	}
	var result entities.Numbering
	ds.DataTo(&result)
	return &result, nil
}

func (r NumberingRepositoryImpl) Save(ctx context.Context, e *entities.Numbering) error {
	ref := r.ref.Doc(strconv.Itoa(e.Type))

	data := map[string]interface{}{
		"Type":      e.Type,
		"Number":    e.Number,
		"UpdatedAt": firestore.ServerTimestamp,
	}
	if e.CreatedAt.IsZero() {
		data["CreatedAt"] = firestore.ServerTimestamp
	}

	if tx, ok := ctx.Value(transactionKey).(*firestore.Transaction); ok {
		//In transaction
		if err := tx.Set(ref, data, firestore.MergeAll); err != nil {
			return err
		}
	} else {
		//Not in transaction
		if _, err := ref.Set(ctx, data, firestore.MergeAll); err != nil {
			return err
		}
	}

	return nil
}
