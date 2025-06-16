package handler

import (
	"net/http"
	"strconv"

	"github.com/DuongQuyen1309/indexevent/internal/datastore"
	// "github.com/DuongQuyen1309/indexevent/internal/model"
	"github.com/gin-gonic/gin"
)

func GetTotalTurnAmountOfUser(c *gin.Context) {
	userAddress := c.Param("address")
	amountSum, err := datastore.GetTotalTurnAmountOfUser(userAddress, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"amount": amountSum, "address": userAddress})
}

func GetTurnsRequestsOfUser(c *gin.Context) {
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
	turns, err := datastore.GetTurnsRequestsOfUser(userAddress, limit, offset, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, turns)
}

func DetailTurnRequestByRequestId(c *gin.Context) {
	requestId := c.Param("request-id")
	turn, err := datastore.GetTurnByRequestId(requestId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, turn)
}

func GetPrizesOfHash(c *gin.Context) {
	requestId := c.Param("request-id")
	prizes, err := datastore.GetPrizesFromRequest(requestId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prizes)
}
