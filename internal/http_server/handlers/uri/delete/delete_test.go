package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserCRUD struct {
	deleteFunc func(id uuid.UUID) error
}

func (m *mockUserCRUD) DeleteUser(id uuid.UUID) error {
	return m.deleteFunc(id)
}

func TestDeleteUserHandler(t *testing.T) {
	t.Run("successfully deletes user", func(t *testing.T) {
		//given
		id := uuid.New()
		r := chi.NewRouter()
		mockCrud := &mockUserCRUD{
			deleteFunc: func(uid uuid.UUID) error {
				assert.Equal(t, id, uid)
				return nil
			},
		}
		handler := New(slog.Default(), mockCrud)
		r.Delete("/users/{id}", handler)

		req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("returns error for invalid UUID", func(t *testing.T) {
		//given
		r := chi.NewRouter()
		handler := New(slog.Default(), &mockUserCRUD{})
		r.Delete("/users/{id}", handler)

		req := httptest.NewRequest(http.MethodDelete, "/users/not-a-uuid", nil)
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)

		//then
		require.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid UUID")
	})

	t.Run("returns error when deletion fails", func(t *testing.T) {
		//given
		id := uuid.New()
		r := chi.NewRouter()
		mockCrud := &mockUserCRUD{
			deleteFunc: func(uuid.UUID) error {
				return errors.New("db error")
			},
		}
		handler := New(slog.Default(), mockCrud)
		r.Delete("/users/{id}", handler)

		req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
		resp := httptest.NewRecorder()
		//when
		r.ServeHTTP(resp, req)
		//then
		require.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Failed to delete user")
	})
}
