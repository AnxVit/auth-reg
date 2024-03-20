package registr

import (
	resp "Auth-Reg/internal/lib/api/response"
	"Auth-Reg/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/render"
)

type Request struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

type Register interface {
	SaveUser(name, email, password string) error
}

var once sync.Once

func New(log *slog.Logger, register Register) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.registr.New"

		once.Do(func() {
			log = log.With(
				slog.String("op", op),
			)
		})

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err = register.SaveUser(req.Name, req.Email, req.Password)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("email", req.Email))
			render.JSON(w, r, resp.Error("user already exists"))
			return
		}

		if err != nil {
			log.Error("falied to save user", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, resp.Error("failed to save user"))
			return
		}

		log.Info("user added")
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
