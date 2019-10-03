package main 

import (
	"github.com/gorilla/mux"
)

func main(){
	router := mux.NewRouter();

	router.HandleFunc("/api/blog",controllers.CreateBlog).Methods("POST")
}