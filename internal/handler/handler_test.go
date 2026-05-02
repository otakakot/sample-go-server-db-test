package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"

	_ "modernc.org/sqlite"

	"github.com/otakakot/sample-go-server-db-test/internal/gateway"
	"github.com/otakakot/sample-go-server-db-test/internal/handler"
)

func Test(t *testing.T) {
	t.Parallel()

	file := uuid.NewString()

	db, err := sql.Open("sqlite", "file:"+file+"?cache=shared")
	if err != nil {
		t.Fatal(err)
	}

	if db.Ping() != nil {
		t.Fatal(err)
	}

	ddl := `
	CREATE TABLE users (
		id   TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`

	if _, err := db.Exec(ddl); err != nil {
		t.Fatal(err)
	}

	gw := gateway.New(db)

	hdl := handler.New(gw)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", hdl.CreateUser)
	mux.HandleFunc("GET /users/{id}", hdl.ReadUser)
	mux.HandleFunc("PUT /users/{id}", hdl.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", hdl.DeleteUser)

	srv := httptest.NewServer(mux)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}

		if err := os.Remove(file); err != nil {
			t.Error(err)
		}

		srv.Close()
	})

	var userID string

	name := uuid.NewString()

	{
		t.Log("POST /users")
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/users", bytes.NewBufferString(`{"name":"`+name+`"}`))
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
		}

		type ResponseBody struct {
			User struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
		}

		var resBody ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}

		userID = resBody.User.ID
	}

	{
		t.Log("GET /users/{id}")
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/users/"+userID, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
		}

		type ResponseBody struct {
			User struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
		}

		var resBody ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}

		if resBody.User.ID != userID {
			t.Errorf("got %s, want %s", resBody.User.ID, userID)
		}

		if resBody.User.Name != name {
			t.Errorf("got %s, want %s", resBody.User.Name, name)
		}
	}

	{
		t.Log("PUT /users/{id}")
		name := uuid.NewString()

		req, err := http.NewRequest(http.MethodPut, srv.URL+"/users/"+userID, bytes.NewBufferString(`{"name":"`+name+`"}`))
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
		}

		req, err = http.NewRequest(http.MethodGet, srv.URL+"/users/"+userID, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
		}

		type ResponseBody struct {
			User struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
		}

		var resBody ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}

		if resBody.User.ID != userID {
			t.Errorf("got %s, want %s", resBody.User.ID, userID)
		}

		if resBody.User.Name != name {
			t.Errorf("got %s, want %s", resBody.User.Name, name)
		}
	}

	{
		t.Log("DELETE /users/{id}")
		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/users/"+userID, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("got %d, want %d", res.StatusCode, http.StatusOK)
		}

		req, err = http.NewRequest(http.MethodGet, srv.URL+"/users/"+userID, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("got %d, want %d", res.StatusCode, http.StatusNotFound)
		}
	}
}
