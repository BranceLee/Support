package main

import (
	"log"
	"net/http"
	"os"

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

	r := mux.NewRouter()

	r.HandleFunc("/api/category/new", handler.CreateBlogCategory).Methods("POST")
	r.HandleFunc("/api/category/blog", handler.GetAllBlogs).Methods("GET")
	r.HandleFunc("/api/category/blog/new", handler.CreateBlog).Methods("POST")

	r.HandleFunc("/api/device/new", handler.CreateDevice).Methods("POST")
	r.HandleFunc("/api/user/new", handler.CreateUser).Methods("POST")

	serv.logger.Info("> Server runs on  8000")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		serv.logger.Info("HTTP server error", zap.Error(err))
	}
}
