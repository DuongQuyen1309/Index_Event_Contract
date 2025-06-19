package datastore

import (
	"context"

	"github.com/DuongQuyen1309/indexevent/internal/model"
)

// "context"

type Datastore interface {
	GetTotalTurnAmountOfUser(address string, c context.Context) (int, error)
	GetTurnsRequestsOfUser(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error)
	GetTurnByRequestId(requestId string, c context.Context) (*model.RequestCreatedEvent, error)
	GetPrizesFromRequest(requestId string, c context.Context) (*[]int64, error)
}
type DataStoreInDB struct {
}
