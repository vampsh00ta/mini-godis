package app

import (
	"mini-godis/config"
	"mini-godis/internal/repository"
	"mini-godis/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	transport "mini-godis/internal/transport/http"
	//"go.uber.org/zap".
)

func Run(cfg *config.Config) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	rep := repository.New()
	srvc := service.New(rep)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)

	trptLogger := logger.Sugar()
	t := transport.New(srvc, trptLogger)

	logger.Info("Listening...")

	if err := http.ListenAndServe(":"+cfg.HTTP.Port, t); err != nil { //nolint:gosec
		logger.Fatal(err.Error())
	}
	select {
	case <-interrupt:
		panic("exit")
	}
}
