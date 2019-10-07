package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/support/config"
	"github.com/support/coreservice"
	"log"
	"net/http"
	"os"
)

func main() {
	dbConfig := config.DefaultPostgresConfig()
	db, err := gorm.Open(dbConfig.Dialect(), dbConfig.ConnectionInfo())
	env, _ := os.LookupEnv("ENV")
	db.LogMode(env == "DEV")
	if err != nil {
		log.Fatal(err)
		return
	}
	handler, err := coreservice.NewHandler(db)
	if err != nil {
		log.Fatal(err)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/blog", handler.CreateBlog).Methods("POST")

	log.Println("> Server runs on  8000")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		log.Println(err)
	}
}
