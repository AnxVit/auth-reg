package smtpmid

import (
	"Auth-Reg/internal/http-server/handlers/registr"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/smtp"
	"os"
	"regexp"

	"github.com/go-chi/render"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) bool {
	isValid := emailRegex.MatchString(email)
	if !isValid {
		return isValid
	}

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASS")
	fmt.Println(from, password)

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	message := []byte("Your email is correct")

	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, []string{email}, message)
	return err == nil
}

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Info("SMTP Logger enabled")
			var req registr.Request

			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			err := render.DecodeJSON(r.Body, &req)
			if err != nil {
				log.Error("failed to decode request body")
				w.WriteHeader(http.StatusBadRequest)

				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			log.Info("request body decoded", slog.Any("request", req))

			isValid := validateEmail(req.Email)
			if !isValid {
				log.Error("Not correct email")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Info("email is correct", slog.String("email", req.Email))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
