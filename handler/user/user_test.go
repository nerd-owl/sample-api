package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"example/web-service-gin/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func createMockContext(
	conn *mocks.MockDbConn,
) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Set("databaseConn", conn)
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

func TestDeactivateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConn := mocks.NewMockDbConn(mockCtrl)
	w, c := createMockContext(mockConn)
	mockConn.EXPECT().
		Exec(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	c.Params = []gin.Param{{Key: "phone", Value: "7408963464"}}

	DeactivateUser(c)

	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestDeactivateUserWithError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConn := mocks.NewMockDbConn(mockCtrl)
	mockConn.EXPECT().
		Exec(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("Dummy Error"))

	w, c := createMockContext(mockConn)
	c.Params = []gin.Param{{Key: "phone", Value: "7408963464"}}

	DeactivateUser(c)

	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConn := mocks.NewMockDbConn(mockCtrl)
	mockConn.EXPECT().
		Exec(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	w, c := createMockContext(mockConn)
	MockJsonPost(c, map[string]interface{}{
		"Fname":   "Arham",
		"Lname":   "Khan",
		"Phone":   "9389474439",
		"Address": "Some Random Address",
	})

	CreateUser(c)
	if w.Code != 201 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConn := mocks.NewMockDbConn(mockCtrl)
	mockConn.EXPECT().
		Exec(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("Dummy Error"))

	w, c := createMockContext(mockConn)
	MockJsonPost(c, map[string]interface{}{
		"Fname":   "Arham",
		"Lname":   "Khan",
		"Phone":   "9389474439",
		"Address": "Some Random Address",
	})

	CreateUser(c)
	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithInvalidName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConn := mocks.NewMockDbConn(mockCtrl)

	w, c := createMockContext(mockConn)
	MockJsonPost(c, map[string]interface{}{
		"Fname":   "Arham ",
		"Lname":   "Khan",
		"Phone":   "9389474439",
		"Address": "Some Random Address",
	})

	CreateUser(c)
	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithInvalidPhone(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConn := mocks.NewMockDbConn(mockCtrl)

	w, c := createMockContext(mockConn)
	MockJsonPost(c, map[string]interface{}{
		"Fname":   "Arham",
		"Lname":   "Khan",
		"Phone":   "93439",
		"Address": "Some Random Address",
	})

	CreateUser(c)
	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestGetUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRows := mocks.NewMockRows(mockCtrl)
	mockConn := mocks.NewMockDbConn(mockCtrl)
	gomock.InOrder(
		mockConn.EXPECT().
			Query(gomock.Any(), gomock.Any()).
			Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true).Times(1),
		mockRows.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1),
		mockRows.EXPECT().Next().Return(false).Times(1),
	)

	w, c := createMockContext(mockConn)

	GetUsers(c)

	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestGetUsersWithScanError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRows := mocks.NewMockRows(mockCtrl)
	mockConn := mocks.NewMockDbConn(mockCtrl)
	gomock.InOrder(
		mockConn.EXPECT().
			Query(gomock.Any(), gomock.Any()).
			Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true).Times(1),
		mockRows.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("Dummy Error")).
			Times(1),
	)

	w, c := createMockContext(mockConn)

	GetUsers(c)

	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestGetUsersWithQueryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRows := mocks.NewMockRows(mockCtrl)
	mockConn := mocks.NewMockDbConn(mockCtrl)
	gomock.InOrder(
		mockConn.EXPECT().
			Query(gomock.Any(), gomock.Any()).
			Return(mockRows, errors.New("Dummy Error")),
	)

	w, c := createMockContext(mockConn)

	GetUsers(c)

	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}
