package main

import (
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
	env				string
	db				*gorm.DB
	sentryDSN		string
	logger			*zap.Logger
}

func newServer(logger *zap.Logger) *server {
	return &server{
		logger:logger,
	}
}

func (s *server) connectToDB(){
	dbConfig := config.DefaultPostgresConfig()
	db, err := gorm.Open(dbConfig.Dialect(), dbConfig.ConnectionInfo())
	env, _ := os.LookupEnv("ENV")
	db.LogMode(env == "DEV")
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
}

func (s *server) initSentry(){
	// Or get DNS from AWS ssmiface.SSMAPI 
	 sentryDSN, ok := os.LookupEnv("sentry")
	 if !ok  {
		 s.logger.Fatal("Failed to initialize sentry")
	 }
	 s.sentryDSN=sentryDSN
}

func main() {
	exampaleLogger := zap.NewExample()
	serv := newServer(exampaleLogger)
	serv.connectToDB()
	serv.initSentry()

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

	serv.logger.Info("> Server runs on  8000")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		serv.logger.Info("HTTP server error", zap.Error(err))
	}
}
