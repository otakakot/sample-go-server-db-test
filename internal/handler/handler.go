package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/otakakot/sample-go-server-db-test/internal/gateway"
)

type Handler struct {
	gateway *gateway.Gateway
}

func New(
	gateway *gateway.Gateway,
) *Handler {
	return &Handler{
		gateway: gateway,
	}
}

func (hdl *Handler) CreateUser(
	res http.ResponseWriter,
	req *http.Request,
) {
	type RequestBody struct {
		Name string `json:"name"`
	}

	var reqBody RequestBody

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	output, err := hdl.gateway.CreateUser(req.Context(), gateway.CreateUserDAI{
		Name: reqBody.Name,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusOK)

	type ResponseBody struct {
		User struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
	}

	resBody := ResponseBody{
		User: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID:   output.User.ID,
			Name: output.User.Name,
		},
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(resBody); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, err := res.Write(buf.Bytes()); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (hdl *Handler) ReadUser(
	res http.ResponseWriter,
	req *http.Request,
) {
	id := req.PathValue("id")

	output, err := hdl.gateway.ReadUser(req.Context(), gateway.ReadUserDAI{
		ID: id,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}

	res.WriteHeader(http.StatusOK)

	type ResponseBody struct {
		User struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
	}

	resBody := ResponseBody{
		User: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID:   output.User.ID,
			Name: output.User.Name,
		},
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(resBody); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, err := res.Write(buf.Bytes()); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (hdl *Handler) UpdateUser(
	res http.ResponseWriter,
	req *http.Request,
) {
	id := req.PathValue("id")

	type RequestBody struct {
		Name string `json:"name"`
	}

	var reqBody RequestBody

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	if _, err := hdl.gateway.UpdateUser(req.Context(), gateway.UpdateUserDAI{
		ID:   id,
		Name: reqBody.Name,
	}); err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}

	res.WriteHeader(http.StatusOK)
}

func (hdl *Handler) DeleteUser(
	res http.ResponseWriter,
	req *http.Request,
) {
	id := req.PathValue("id")

	if _, err := hdl.gateway.DeleteUser(req.Context(), gateway.DeleteUserDAI{
		ID: id,
	}); err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}

	res.WriteHeader(http.StatusOK)
}
