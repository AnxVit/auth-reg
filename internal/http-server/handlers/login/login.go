package login

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"Auth-Reg/internal/domain/models"
	resp "Auth-Reg/internal/lib/api/response"
	"Auth-Reg/internal/lib/jwt"
	"Auth-Reg/internal/storage"

	"github.com/go-chi/render"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
	Token string
}

type Loginer interface {
	User(email string) (models.User, error)
}

var once sync.Once

func New(log *slog.Logger, loginer Loginer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.login.New"

		once.Do(func() {
			log = log.With(
				slog.String("op", op),
			)
		})

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decoded request")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		user, err := loginer.User(req.Email)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found")
			render.JSON(w, r, resp.Error("user not found"))
			return
		}
		if err != nil {
			log.Error("falied to find user", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, resp.Error("failed to find user"))
			return
		}

		if req.Password != string(user.Password) {
			log.Info("invalid password")
			render.JSON(w, r, resp.Error("invalid password"))
			return
		}

		log.Info("user logged in successfully")

		token, err := jwt.NewToken(user)
		if err != nil {
			log.Error("failed to generate token")
			render.JSON(w, r, resp.Error("filed to generate token"))
			return
		}
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Token:    token,
		})

	}
}
