package app

import (
	"bccChat/internal/handlers"
	"bccChat/internal/logger"
	_ "context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"bccChat/internal/repository"
	service "bccChat/internal/services"

	"github.com/joho/godotenv"
)

type app struct {
	log     *slog.Logger
	Handler *handlers.Handler
}

type App interface {
	GetLogger() *slog.Logger
	Start() error
}

func New() App {
	// Load environment variables

	logger := logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		logger.Error("DB_DSN not set in .env file")
	}

	db, err := repository.NewDatabase(dsn)
	if err != nil {
		logger.Error("Error connecting to the database: %v", err)
	}
	//defer db.Close()

	repo := repository.NewMessageRepository(db)
	server := service.NewChatServer(repo)
	Handler := handlers.NewHandler(server)
	return &app{log: logger, Handler: Handler}
}

func (a app) GetLogger() *slog.Logger {
	return a.log
}

func (h *app) Start() error {

	errChan := make(chan error)

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	go h.Handler.StartHandler( /*ctx,*/ errChan)

	h.GetLogger().Info("app is started")
	Register(errChan, h.GetLogger())

	return nil
}

func Register(errChan chan error, log *slog.Logger) {

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errChan:
		log.Error("gracefully shutdown error", "error", err.Error())
	case stop := <-stopChan:
		log.Error("app is finished", "signal", stop.String())
	}

}
