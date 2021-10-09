package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// GetEntries Test Case
func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newUserHandlers().getusers)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":45,"name":"Ramu","email":"ram1@yahoo.com","password":"\ufffdH\ufffdm\u0026\ufffd\u0006Z\ufffd\ufffd\u0000\u0004\ufffd\ufffdwᣧ\ufffd-z\ufffd\r\ufffd\ufffdR\ufffd~\u0007_\ufffd\ufffd\ufffd\ufffd]\ufffd\ufffd\ufffd"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "45")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newUserHandlers().getUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":45,"name":"Ramu","email":"ram1@yahoo.com","password":"\ufffdH\ufffdm\u0026\ufffd\u0006Z\ufffd\ufffd\u0000\u0004\ufffd\ufffdwᣧ\ufffd-z\ufffd\r\ufffd\ufffdR\ufffd~\u0007_\ufffd\ufffd\ufffd\ufffd]\ufffd\ufffd\ufffd"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "100")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newUserHandlers().getUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestUserEntry(t *testing.T) {

	var jsonStr = []byte(`{"id":11,"name":"kiran","email":"kiran@kap.com","password":"11452369870"}`)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newUserHandlers().postusers)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":11,"name":"kiran","email":"kiran@kap.com","password":"11452369870"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPosts(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newPostHandlers().getposts)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":21,"userid":0,"caption":"scenary","imageurl":"http://c:/documents/j.jpg","postedtimestamp":"2021-10-09T15:28:01.439Z"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/posts/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "21")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newPostHandlers().getPost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":21,"userid":0,"caption":"scenary","imageurl":"http://c:/documents/j.jpg","postedtimestamp":"2021-10-09T15:28:01.439Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostByIDNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "123")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newPostHandlers().getPost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostEntry(t *testing.T) {

	var jsonStr = []byte(`{"id":11,"userid":21,"caption":"smile face","imageurl":"http://c:/imag.png"}`)

	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newPostHandlers().postposts)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":11,"userid":21,"caption":"smile face","imageurl":"http://c:/imag.png"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
