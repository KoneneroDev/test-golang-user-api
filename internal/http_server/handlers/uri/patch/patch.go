package patch

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"test_golang_user_api/internal/api"
	"test_golang_user_api/internal/storage/postgres"
)

type UserCRUD interface {
	EditUser(user *postgres.UserDto) (*postgres.UserDto, error)
}

func New(log *slog.Logger, userCrud UserCRUD) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.With(
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req api.Request

		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("Error decoding request body", err)
			render.JSON(writer, request, api.ErrorStatus("Failed to decode request body"))
			return
		}

		log.Info("Request body decoded", slog.Any("requestBody", req))

		idStr := chi.URLParam(request, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("Invalid UUID", err)
			render.JSON(writer, request, api.ErrorStatus("Invalid UUID"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("Error validating request body", err)
			render.JSON(writer, request, api.ErrorStatus("Failed to validate request body"))
			return
		}

		user, err := userCrud.EditUser(postgres.NewUser(id, req.Firstname, req.Lastname, req.Email, req.Age))
		if err != nil {
			log.Error("Error edit user", err)
			render.JSON(writer, request, api.ErrorStatus("Failed to edit user"))
			return
		}

		render.JSON(writer, request, user)

		log.Info("User edit successfully")
	}
}
