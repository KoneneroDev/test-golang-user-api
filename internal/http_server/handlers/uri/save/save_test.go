package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"test_golang_user_api/internal/api"
	"test_golang_user_api/internal/storage/postgres"
	"testing"
)

type mockUserCRUD struct {
	createFunc func(user *postgres.UserDto) error
}

func (m *mockUserCRUD) CreateUser(user *postgres.UserDto) error {
	return m.createFunc(user)
}

func TestSaveUserHandler(t *testing.T) {
	t.Run("successfully creates user", func(t *testing.T) {
		//given
		r := http.NewServeMux()
		reqBody := api.Request{
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@example.com",
			Age:       30,
		}
		mockCrud := &mockUserCRUD{
			createFunc: func(u *postgres.UserDto) error {
				assert.Equal(t, reqBody.Firstname, u.Firstname)
				return nil
			},
		}

		handler := New(slog.Default(), mockCrud)
		r.Handle("/users", handler)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.OkResponse())
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error for invalid JSON body", func(t *testing.T) {
		//given
		r := http.NewServeMux()
		handler := New(slog.Default(), &mockUserCRUD{})
		r.Handle("/users", handler)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid-json")))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Failed to decode request body"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error for validation failure", func(t *testing.T) {
		//given
		r := http.NewServeMux()
		handler := New(slog.Default(), &mockUserCRUD{})
		r.Handle("/users", handler)

		invalidBody := api.Request{Firstname: "", Lastname: "", Email: "", Age: 0}
		body, _ := json.Marshal(invalidBody)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Failed to validate request body"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error when CreateUser fails", func(t *testing.T) {
		//given
		r := http.NewServeMux()
		reqBody := api.Request{
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@example.com",
			Age:       30,
		}
		mockCrud := &mockUserCRUD{
			createFunc: func(u *postgres.UserDto) error {
				return errors.New("db error")
			},
		}

		handler := New(slog.Default(), mockCrud)
		r.Handle("/users", handler)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Failed to create user"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})
}
