package handler

import (
	"net/http"
	"strconv"

	"github.com/DuongQuyen1309/indexevent/internal/datastore"
	// "github.com/DuongQuyen1309/indexevent/internal/model"
	"github.com/gin-gonic/gin"
)

// type Config struct {
// 	R              *gin.Engine
// 	EventDatastore model.EventDatastore
// }

//	func NewHandler(c *Config) {
//		h := &Handler{
//			EventDatastore: c.EventDatastore,
//		}
//		g := c.R.Group("/api/v1")
//		g.GET("/user/:address/turn-amount", h.GetTotalTurnAmountOfUser)
//		g.GET("/user/:address/turn-requests", h.GetTurnsRequestsOfUser)
//		g.GET("/turn-request/:request-id", h.DetailTurnRequestByRequestId)
//		g.GET("/turn-request/:request-id/prizes", h.GetPrizesOfTurnRequest)
//	}
type Handler struct {
	EventDatastore datastore.Datastore
}

func (h *Handler) GetTotalTurnAmountOfUser(c *gin.Context) {
	userAddress := c.Param("address")
	amountSum, err := h.EventDatastore.GetTotalTurnAmountOfUser(userAddress, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"amount": amountSum})
}
func (h *Handler) GetTurnsRequestsOfUser(c *gin.Context) {
	userAddress := c.Param("address")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}
	if page <= 0 || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page and limit parameters must be greater than 0"})
		return
	}
	offset := (page - 1) * limit
	turns, err := h.EventDatastore.GetTurnsRequestsOfUser(userAddress, limit, offset, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, turns)
}

func (h *Handler) DetailTurnRequestByRequestId(c *gin.Context) {
	requestId := c.Param("request-id")
	turn, err := h.EventDatastore.GetTurnByRequestId(requestId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, turn)
}

func (h *Handler) GetPrizesOfTurnRequest(c *gin.Context) {
	requestId := c.Param("request-id")
	prizes, err := h.EventDatastore.GetPrizesFromRequest(requestId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prizes)
}
