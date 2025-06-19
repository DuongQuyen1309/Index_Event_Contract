package router

import (
	"github.com/DuongQuyen1309/indexevent/internal/datastore"
	"github.com/DuongQuyen1309/indexevent/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store datastore.Datastore) *gin.Engine {
	router := gin.Default()
	h := &handler.Handler{EventDatastore: store}
	router.GET("/user/:address/turn-amount", h.GetTotalTurnAmountOfUser)
	router.GET("/user/:address/turn-requests", h.GetTurnsRequestsOfUser)
	router.GET("/turn-request/:request-id", h.DetailTurnRequestByRequestId)
	router.GET("/turn-request/:request-id/prizes", h.GetPrizesOfTurnRequest)
	return router
}
