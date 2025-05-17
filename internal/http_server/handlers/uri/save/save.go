package save

import (
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
	CreateUser(user *postgres.UserDto) error
}

func New(log *slog.Logger, userCrud UserCRUD) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.With(
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req api.Request

		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("Error decoding request body", slog.Any("err", err))
			render.JSON(writer, request, api.ErrorStatus("Failed to decode request body"))
			return
		}

		log.Info("Request body decoded", slog.Any("requestBody", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("Error validating request body", slog.Any("err", err))
			render.JSON(writer, request, api.ErrorStatus("Failed to validate request body"))
			return
		}

		err = userCrud.CreateUser(postgres.NewUser(uuid.New(), req.Firstname, req.Lastname, req.Email, req.Age))
		if err != nil {
			log.Error("Error creating user", slog.Any("err", err))
			render.JSON(writer, request, api.ErrorStatus("Failed to create user"))
			return
		}

		render.JSON(writer, request, api.OkResponse())

		log.Info("User created successfully")
	}
}
