package get

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"test_golang_user_api/internal/api"
	"test_golang_user_api/internal/storage/postgres"
)

type UserCRUD interface {
	GetUser(id uuid.UUID) (postgres.UserDto, error)
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

		user, err := crud.GetUser(id)

		if err != nil {
			log.Error("User not found", err)
			render.JSON(writer, request, api.ErrorStatus("User not found"))
			return
		}

		render.JSON(writer, request, user)

		log.Info("User successfully retrieved")
	}
}
