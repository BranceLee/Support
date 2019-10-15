package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/support/config"
	"github.com/support/coreservice"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

type server struct {
	env       string
	db        *gorm.DB
	sentryDSN string
	logger    *zap.Logger
}

func newServer(logger *zap.Logger) *server {
	return &server{
		logger: logger,
	}
}

func (s *server) connectToDB() {
	dbConfig := config.DefaultPostgresConfig()
	db, err := gorm.Open(dbConfig.Dialect(), dbConfig.ConnectionInfo())
	env, ok := os.LookupEnv("ENV")
	if !ok {
		s.logger.Fatal("ENV Can not be null")
	}
	db.LogMode(env == "DEV")
	if err != nil {
		sentry.CaptureException(err)
		s.logger.Fatal("Environment variable ENV is not found")
	}
	s.db = db
}

func (s *server) initSentry() {
	// TODO: Get DNS from AWS ssmiface.SSMAPI, but in this project get in local.
	// sentryDSN, ok := os.LookupEnv("sentry")
	//  if ok  {
	// 	 s.logger.Fatal("Failed to initialize sentry")
	//  }
	sentryDSN := config.GetSentryDSN()
	sentry.Init(sentry.ClientOptions{
		Dsn: sentryDSN,
	})
	s.sentryDSN = sentryDSN
}

func main() {
	exampaleLogger := zap.NewExample()
	serv := newServer(exampaleLogger)
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

	router := mux.NewRouter()

	router.HandleFunc("/api/blog", handler.CreateBlog).Methods("POST")
	router.HandleFunc("/api/blog/all", handler.GetAllBlogs).Methods("GET")
	router.HandleFunc("/api/device/new", handler.CreateDevice).Methods("POST")
	router.HandleFunc("/api/user/new", handler.CreateUser).Methods("POST")
	serv.logger.Info("> Server runs on  8000")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		serv.logger.Info("HTTP server error", zap.Error(err))
	}
}
