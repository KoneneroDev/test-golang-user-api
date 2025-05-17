package patch

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	editFunc func(user *postgres.UserDto) (*postgres.UserDto, error)
}

func (m *mockUserCRUD) EditUser(user *postgres.UserDto) (*postgres.UserDto, error) {
	return m.editFunc(user)
}

func TestPatchUserHandler(t *testing.T) {
	t.Run("successfully edits user", func(t *testing.T) {
		//given
		id := uuid.New()
		r := chi.NewRouter()
		request := api.Request{
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@example.com",
			Age:       30,
		}
		updatedUser := postgres.NewUser(id, request.Firstname, request.Lastname, request.Email, request.Age)
		mockCrud := &mockUserCRUD{
			editFunc: func(u *postgres.UserDto) (*postgres.UserDto, error) {
				assert.Equal(t, updatedUser.Firstname, u.Firstname)
				return updatedUser, nil
			},
		}
		handler := New(slog.Default(), mockCrud)
		r.Patch("/users/{id}", handler)

		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPatch, "/users/"+id.String(), bytes.NewReader(body))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(updatedUser)
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error for invalid UUID", func(t *testing.T) {
		//given
		r := chi.NewRouter()
		handler := New(slog.Default(), &mockUserCRUD{})
		r.Patch("/users/{id}", handler)

		body, _ := json.Marshal(api.Request{
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@example.com",
			Age:       30,
		})
		req := httptest.NewRequest(http.MethodPatch, "/users/not-a-uuid", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Invalid UUID"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error for invalid body", func(t *testing.T) {
		//given
		r := chi.NewRouter()
		handler := New(slog.Default(), &mockUserCRUD{})
		r.Patch("/users/{id}", handler)

		id := uuid.New()
		req := httptest.NewRequest(http.MethodPatch, "/users/"+id.String(), bytes.NewReader([]byte("invalid-json")))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Failed to decode request body"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error when update fails", func(t *testing.T) {
		//given
		id := uuid.New()
		r := chi.NewRouter()
		mockCrud := &mockUserCRUD{
			editFunc: func(u *postgres.UserDto) (*postgres.UserDto, error) {
				return nil, errors.New("db error")
			},
		}
		handler := New(slog.Default(), mockCrud)
		r.Patch("/users/{id}", handler)

		body, _ := json.Marshal(api.Request{
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@example.com",
			Age:       30,
		})
		req := httptest.NewRequest(http.MethodPatch, "/users/"+id.String(), bytes.NewReader(body))
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Failed to edit user"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})
}
