package user

import (
	"bytes"
	"encoding/json"
	"errors"
	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func createMockContext(
	querier db.Querier,
) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Set("querier", querier)
	return w, c
}

func MockJsonPost(c *gin.Context /* the test context */, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestGetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)
	mockdb.EXPECT().ListUser(c.Request.Context()).Return([]db.Kuser{{
		Firstname: "Savez",
		Lastname:  "Siddiqui",
		Phone:     "7408963464",
		Addr:      "Some Addr",
		Active:    true,
	}}, nil)

	GetUsers(c)

	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestGetUserWithDbError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)
	mockdb.EXPECT().
		ListUser(c.Request.Context()).
		Return(nil, errors.New("Dummy Error"))

	GetUsers(c)

	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)
	mockdb.EXPECT().CreateUser(c.Request.Context(), db.CreateUserParams{
		Firstname: "Savez",
		Lastname:  "Siddiqui",
		Phone:     "7408963464",
		Addr:      "Some Addr",
	}).Return(nil)

	MockJsonPost(c, map[string]interface{}{
		"Firstname": "Savez",
		"Lastname":  "Siddiqui",
		"Phone":     "7408963464",
		"Addr":      "Some Addr",
	})

	CreateUser(c)

	if w.Code != 201 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithInvalidNameInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)

	MockJsonPost(c, map[string]interface{}{
		"Firstname": "Savez ",
		"Lastname":  "Siddiqui",
		"Phone":     "7408963464",
		"Addr":      "Some Addr",
	})

	CreateUser(c)

	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithInvalidPhoneInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)

	MockJsonPost(c, map[string]interface{}{
		"Firstname": "Savez",
		"Lastname":  "Siddiqui",
		"Phone":     "74064",
		"Addr":      "Some Addr",
	})

	CreateUser(c)

	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithDbError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)
	mockdb.EXPECT().CreateUser(c.Request.Context(), db.CreateUserParams{
		Firstname: "Savez",
		Lastname:  "Siddiqui",
		Phone:     "7408963464",
		Addr:      "Some Addr",
	}).Return(errors.New("Dummy Error"))

	MockJsonPost(c, map[string]interface{}{
		"Firstname": "Savez",
		"Lastname":  "Siddiqui",
		"Phone":     "7408963464",
		"Addr":      "Some Addr",
	})

	CreateUser(c)

	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestDeactivateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)
	c.Params = []gin.Param{{Key: "phone", Value: "7408963464"}}
	mockdb.EXPECT().
		DeactivateUser(c.Request.Context(), "7408963464").
		Return(nil)

	DeactivateUser(c)

	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestDeactivateUserWithDbError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockdb := mocks.NewMockQuerier(mockCtrl)
	w, c := createMockContext(mockdb)
	c.Params = []gin.Param{{Key: "phone", Value: "7408963464"}}
	mockdb.EXPECT().
		DeactivateUser(c.Request.Context(), "7408963464").
		Return(errors.New("Dummy Error"))

	DeactivateUser(c)

	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}
