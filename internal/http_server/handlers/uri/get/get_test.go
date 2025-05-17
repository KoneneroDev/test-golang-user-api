package get

import (
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
	getFunc func(id uuid.UUID) (*postgres.UserDto, error)
}

func (m *mockUserCRUD) GetUser(id uuid.UUID) (*postgres.UserDto, error) {
	return m.getFunc(id)
}

func TestGetUserHandler(t *testing.T) {
	t.Run("successfully retrieves user", func(t *testing.T) {
		//given
		id := uuid.New()
		r := chi.NewRouter()
		user := &postgres.UserDto{
			ID:        id,
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@gmail.com",
			Age:       30,
		}

		mockCrud := &mockUserCRUD{
			getFunc: func(uid uuid.UUID) (*postgres.UserDto, error) {
				assert.Equal(t, id, uid)
				return user, nil
			},
		}

		handler := New(slog.Default(), mockCrud)
		r.Get("/users/{id}", handler)

		req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(user)
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error for invalid UUID", func(t *testing.T) {
		//given
		r := chi.NewRouter()
		handler := New(slog.Default(), &mockUserCRUD{})
		r.Get("/users/{id}", handler)

		req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("Invalid UUID"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		//given
		id := uuid.New()
		r := chi.NewRouter()
		mockCrud := &mockUserCRUD{
			getFunc: func(uid uuid.UUID) (*postgres.UserDto, error) {
				return nil, errors.New("not found")
			},
		}

		handler := New(slog.Default(), mockCrud)
		r.Get("/users/{id}", handler)

		req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		expected, _ := json.Marshal(api.ErrorStatus("User not found"))
		require.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, string(expected), resp.Body.String())
	})
}
