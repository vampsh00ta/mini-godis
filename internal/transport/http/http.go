package http

import (
	"github.com/go-playground/validator/v10"
	"mini-godis/internal/service"
	"net/http"

	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
	decoder  = schema.NewDecoder()
)

type transport struct {
	s service.Service
	l *zap.SugaredLogger
}

func New(t service.Service, l *zap.SugaredLogger) http.Handler {
	r := &transport{t, l}
	mux := http.NewServeMux()
	r = r

	return mux
}
