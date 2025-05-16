package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"test_golang_user_api/internal/api"
)

type UserCRUD interface {
	DeleteUser(id uuid.UUID) error
}

func New(log *slog.Logger, crud UserCRUD) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.With(
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		idStr := chi.URLParam(request, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("Invalid UUID", err)
			render.JSON(writer, request, api.ErrorStatus("Invalid UUID"))
			return
		}

		err = crud.DeleteUser(id)
		if err != nil {
			log.Error("Error deleting user", err)
			render.JSON(writer, request, api.ErrorStatus("Failed to delete user"))
			return
		}

		render.JSON(writer, request, api.OkResponse())

		log.Info("User deleted successfully")
	}

}
