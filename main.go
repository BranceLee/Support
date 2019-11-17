package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/support/config"
	"github.com/support/coreservice"
	"go.uber.org/zap"
)

type server struct {
	env       string
	db        *gorm.DB
	sentryDSN string
	logger    *zap.Logger
}

func (s *server) connectToDB() {
	dbConfig := config.DefaultPostgresConfig()
	db, err := gorm.Open(dbConfig.Dialect(), dbConfig.ConnectionInfo())
	db.LogMode(s.env == "TEST")
	if err != nil {
		sentry.CaptureException(err)
		s.logger.Fatal("Environment variable ENV is not found")
	}
	s.db = db
}

func (s *server) initSentry() {
	// TODO: Get DNS from AWS ssmiface.SSMAPI, but in this project get in local.
	sentryDSN := config.GetSentryDSN()
	sentry.Init(sentry.ClientOptions{
		Dsn: sentryDSN,
	})
	s.sentryDSN = sentryDSN
}

func main() {
	exampaleLogger := zap.NewExample()
	env, ok := os.LookupEnv("ENV")
	if !ok {
		exampaleLogger.Fatal("ENV is not found")
	}
	serv := &server{
		logger: exampaleLogger,
		env:    env,
	}
	serv.initSentry()
	serv.connectToDB()

	undo := zap.RedirectStdLog(serv.logger)
	defer serv.logger.Sync()
	defer undo()

	handler, err := coreservice.NewHandler(serv.db)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/category/new", handler.CreateCategory).Methods("POST")
	r.HandleFunc("/api/category", handler.GetCategory).Methods("GET")
	r.HandleFunc("/api/category/blog", handler.GetAllBlogs).Methods("GET")
	r.HandleFunc("/api/category/blog/new", handler.CreateBlog).Methods("POST")

	r.HandleFunc("/api/user/new", handler.CreateUser).Methods("POST")

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		if err := server.Shutdown(ctx); err != nil {
			serv.logger.Info("HTTP server shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	serv.logger.Info("> Server runs on  8000")
	err = server.ListenAndServe()
	if err != nil {
		serv.logger.Info("HTTP server error", zap.Error(err))
	}
	<-idleConnsClosed
}
