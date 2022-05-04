package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockConn struct {
}

type mockRows struct {
}

var i = 0
var conn mockConn

func (r mockRows) Next() bool {
	j := i
	i++
	return (j % 2) == 0
}

func (r mockRows) Scan(dest ...interface{}) error {
	dest[0] = "Savez"
	dest[1] = "Siddiqui"
	dest[2] = "7408963464"
	dest[3] = "Some Address"
	dest[4] = true
	return nil
}

func (mockConn) Query(
	ctx context.Context,
	sql string,
	args ...interface{},
) (Rows, error) {
	var rows mockRows
	return rows, nil
}

func (mockConn) Exec(
	ctx context.Context,
	sql string,
	arguments ...interface{},
) (interface{}, error) {
	return nil, nil
}

func createMockContext() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("databaseConn", conn)
	return w, c
}

func MockJsonPost(c *gin.Context /* the test context */, content interface{}) {
	c.Request.Method = "POST" // or PUT
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

func TestGetUsers(t *testing.T) {
	w, c := createMockContext()

	GetUsers(c)
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUser(t *testing.T) {
	w, c := createMockContext()
	c.Request = &http.Request{
		Header: make(http.Header),
	}

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

func TestCreateUserWithWrongNameInput(t *testing.T) {
	w, c := createMockContext()
	c.Request = &http.Request{
		Header: make(http.Header),
	}

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

func TestCreateUserWithNoRequestBody(t *testing.T) {
	w, c := createMockContext()
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	CreateUser(c)
	if w.Code != 500 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestCreateUserWithWrongPhoneInput(t *testing.T) {
	w, c := createMockContext()
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockJsonPost(c, map[string]interface{}{
		"Fname":   "Arham",
		"Lname":   "Khan",
		"Phone":   "938939",
		"Address": "Some Random Address",
	})

	CreateUser(c)
	if w.Code != 400 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestDeactivateUser(t *testing.T) {
	w, c := createMockContext()
	c.Params = []gin.Param{{Key: "phone", Value: "7408963464"}}

	DeactivateUser(c)
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}

func TestDeactivateUserWithFailedExec(t *testing.T) {
	w, c := createMockContext()
	c.Params = []gin.Param{{Key: "phone", Value: "7408963464"}}

	DeactivateUser(c)
	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}
