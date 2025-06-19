package handler

import (
	// "net/http"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DuongQuyen1309/indexevent/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TotalTurnAmount struct {
	Amount int `json:"amount"`
}

type Error struct {
	Error string `json:"error"`
}

type MockDatastore struct {
	GetTotalTurnAmountOfUserFunc func(address string, c context.Context) (int, error)
	GetTurnsRequestsOfUserFunc   func(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error)
	GetTurnByRequestIdFunc       func(requestId string, c context.Context) (*model.RequestCreatedEvent, error)
	GetPrizesFromRequestFunc     func(requestId string, c context.Context) (*[]int64, error)
}

func (m *MockDatastore) GetTotalTurnAmountOfUser(address string, c context.Context) (int, error) {
	return m.GetTotalTurnAmountOfUserFunc(address, c)
}

func (m *MockDatastore) GetTurnsRequestsOfUser(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error) {
	return m.GetTurnsRequestsOfUserFunc(address, limit, offset, c)
}

func (m *MockDatastore) GetTurnByRequestId(requestId string, c context.Context) (*model.RequestCreatedEvent, error) {
	return m.GetTurnByRequestIdFunc(requestId, c)
}
func (m *MockDatastore) GetPrizesFromRequest(requestId string, c context.Context) (*[]int64, error) {
	return m.GetPrizesFromRequestFunc(requestId, c)
}
func GetDataForTestGetTotalTurnAmountOfUser(mockDatastore *MockDatastore, address string) (*TotalTurnAmount, *httptest.ResponseRecorder, error) {
	h := &Handler{EventDatastore: mockDatastore}
	router := gin.Default()
	url := "/user/:address/turn-amount"
	router.GET(url, h.GetTotalTurnAmountOfUser)
	url = strings.ReplaceAll(url, ":address", address)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	bodyResponse, err := io.ReadAll(w.Body)
	if err != nil {
		return nil, nil, err
		// t.Fatal("ReadAll failed:", err)
	}
	var totalTurnAmount TotalTurnAmount
	err = json.Unmarshal(bodyResponse, &totalTurnAmount)
	if err != nil {
		return nil, nil, err
		// t.Fatal("Unmarshal failed:", err)
	}
	return &totalTurnAmount, w, nil
}
func TestGetTotalTurnAmountOfUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTotalTurnAmountOfUserFunc: func(address string, c context.Context) (int, error) {
				return 33, nil
			},
		}
		totalTurnAmount, w, err := GetDataForTestGetTotalTurnAmountOfUser(mockDatastore, "0xAdfD8DAa41c23c18064074416d3428a3086e1621")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 33, totalTurnAmount.Amount)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success with non-exist address", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTotalTurnAmountOfUserFunc: func(address string, c context.Context) (int, error) {
				return 0, nil
			},
		}
		totalTurnAmount, w, err := GetDataForTestGetTotalTurnAmountOfUser(mockDatastore, "0xAdfD8DAa41c23c18064074416d3428a3086e1621abc")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 0, totalTurnAmount.Amount)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func GetSuccessDataForTestGetTurnsRequestsOfUser(mockDatastore *MockDatastore, address string, page int, limit int) (*[]model.RequestCreatedEvent, *httptest.ResponseRecorder, error) {
	h := &Handler{EventDatastore: mockDatastore}
	router := gin.Default()
	url := "/user/:address/turn-requests"
	router.GET(url, h.GetTurnsRequestsOfUser)
	url = strings.ReplaceAll(url, ":address", address)
	url = fmt.Sprintf("%s?limit=%d&page=%d", url, limit, page)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	bodyResponse, err := io.ReadAll(w.Body)
	if err != nil {
		return nil, nil, err
		// t.Fatal("ReadAll failed:", err)
	}
	var requestEvent []model.RequestCreatedEvent
	err = json.Unmarshal(bodyResponse, &requestEvent)
	if err != nil {
		return nil, nil, err
		// t.Fatal("Unmarshal failed:", err)
	}
	return &requestEvent, w, nil
}

func GetErrorDataForTestGetTurnsRequestsOfUser(mockDatastore *MockDatastore, address string, page string, limit string) (*Error, *httptest.ResponseRecorder, error) {
	h := &Handler{EventDatastore: mockDatastore}
	router := gin.Default()
	url := "/user/:address/turn-requests"
	router.GET(url, h.GetTurnsRequestsOfUser)
	url = strings.ReplaceAll(url, ":address", address)
	url = fmt.Sprintf("%s?limit=%s&page=%s", url, limit, page)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	bodyResponse, err := io.ReadAll(w.Body)
	if err != nil {
		return nil, nil, err
		// t.Fatal("ReadAll failed:", err)
	}
	var responseError Error
	err = json.Unmarshal(bodyResponse, &responseError)
	if err != nil {
		return nil, nil, err
	}
	return &responseError, w, nil
}
func TestGetTurnsRequestsOfUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTurnsRequestsOfUserFunc: func(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error) {
				createdAt, _ := time.Parse("2006-01-02 15:04:05.000 -0700", "2022-09-02 20:02:15.000 +0700")
				return &[]model.RequestCreatedEvent{
					{
						Id:              7,
						TransactionHash: "0xee8eaeb360d562d6d3c9d07de35b5a7b1ffaa919eaeb8882a256758499397173",
						LogIndex:        405,
						RequestId:       "4880331253831101174008392807351007390870631292530080775106683908345166554785",
						RequestOwner:    "0x735AB5B0dcC0b678E420B48a8719B1887336310e",
						Amount:          1,
						CreatedAt:       createdAt,
					},
				}, nil
			},
		}
		requestEvent, w, err := GetSuccessDataForTestGetTurnsRequestsOfUser(mockDatastore, "0x735AB5B0dcC0b678E420B48a8719B1887336310e", 1, 1)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 7, (*requestEvent)[0].Id)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Fail with invalid integer parameters", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTurnsRequestsOfUserFunc: func(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error) {
				return &[]model.RequestCreatedEvent{}, errors.New("page and limit parameters must be greater than 0")
			},
		}
		responseError, w, err := GetErrorDataForTestGetTurnsRequestsOfUser(mockDatastore, "0x735AB5B0dcC0b678E420B48a8719B1887336310e", "-1", "10")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "page and limit parameters must be greater than 0", responseError.Error)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Fail with invalid limit parameter number type", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTurnsRequestsOfUserFunc: func(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error) {
				return &[]model.RequestCreatedEvent{}, errors.New("invalid limit parameter")
			},
		}
		responseError, w, err := GetErrorDataForTestGetTurnsRequestsOfUser(mockDatastore, "0x735AB5B0dcC0b678E420B48a8719B1887336310e", "2", "a")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "invalid limit parameter", responseError.Error)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Fail with invalid page parameter number type", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTurnsRequestsOfUserFunc: func(address string, limit int, offset int, c context.Context) (*[]model.RequestCreatedEvent, error) {
				return &[]model.RequestCreatedEvent{}, errors.New("invalid page parameter")
			},
		}
		responseError, w, err := GetErrorDataForTestGetTurnsRequestsOfUser(mockDatastore, "0x735AB5B0dcC0b678E420B48a8719B1887336310e", "b", "2")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "invalid page parameter", responseError.Error)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
func TestDetailTurnRequestByRequestId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTurnByRequestIdFunc: func(requestId string, c context.Context) (*model.RequestCreatedEvent, error) {
				createdAt, _ := time.Parse("2006-01-02 15:04:05.000 -0700", "2022-09-02 20:02:15.000 +0700")
				return &model.RequestCreatedEvent{
					Id:              7,
					TransactionHash: "0xee8eaeb360d562d6d3c9d07de35b5a7b1ffaa919eaeb8882a256758499397173",
					LogIndex:        405,
					RequestId:       "4880331253831101174008392807351007390870631292530080775106683908345166554785",
					RequestOwner:    "0x735AB5B0dcC0b678E420B48a8719B1887336310e",
					Amount:          1,
					CreatedAt:       createdAt,
				}, nil
			},
		}
		h := &Handler{EventDatastore: mockDatastore}
		router := gin.Default()
		url := "/turn-request/:request-id"
		router.GET(url, h.DetailTurnRequestByRequestId)
		url = strings.ReplaceAll(url, ":request-id", "4880331253831101174008392807351007390870631292530080775106683908345166554785")
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		bodyResponse, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fatal("ReadAll failed:", err)
		}
		var requestEvent model.RequestCreatedEvent
		t.Log(string(bodyResponse))
		err = json.Unmarshal(bodyResponse, &requestEvent)
		if err != nil {
			t.Fatal("Unmarshal failed:", err)
		}
		assert.Equal(t, 7, requestEvent.Id)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Fail with non-exist address", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetTurnByRequestIdFunc: func(requestId string, c context.Context) (*model.RequestCreatedEvent, error) {
				return nil, errors.New("sql: no rows in result set")
			},
		}
		h := &Handler{EventDatastore: mockDatastore}
		router := gin.Default()
		url := "/turn-request/:request-id"
		router.GET(url, h.DetailTurnRequestByRequestId)
		url = strings.ReplaceAll(url, ":request-id", "48803312538314785")
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		bodyResponse, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fatal("ReadAll failed:", err)
		}
		var responseError Error
		err = json.Unmarshal(bodyResponse, &responseError)
		if err != nil {
			t.Fatal("Unmarshal failed:", err)
		}
		assert.Equal(t, "sql: no rows in result set", responseError.Error)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
func TestGetPrizesOfTurnRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetPrizesFromRequestFunc: func(address string, c context.Context) (*[]int64, error) {
				return &[]int64{1}, nil
			},
		}
		h := &Handler{EventDatastore: mockDatastore}
		router := gin.Default()
		url := "/turn-request/:request-id/prizes"
		router.GET(url, h.GetPrizesOfTurnRequest)
		url = strings.ReplaceAll(url, ":request-id", "4880331253831101174008392807351007390870631292530080775106683908345166554785")
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Fail with non-exist address", func(t *testing.T) {
		mockDatastore := &MockDatastore{
			GetPrizesFromRequestFunc: func(address string, c context.Context) (*[]int64, error) {
				return nil, errors.New("sql: no rows in result set")
			},
		}
		h := &Handler{EventDatastore: mockDatastore}
		router := gin.Default()
		url := "/turn-request/:request-id/prizes"
		router.GET(url, h.GetPrizesOfTurnRequest)
		url = strings.ReplaceAll(url, ":request-id", "48803312538314785")
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		bodyResponse, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fatal("ReadAll failed:", err)
		}
		var responseError Error
		err = json.Unmarshal(bodyResponse, &responseError)
		if err != nil {
			t.Fatal("Unmarshal failed:", err)
		}
		assert.Equal(t, "sql: no rows in result set", responseError.Error)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
