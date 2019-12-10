package coreservice_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/BranceLee/Support/config"
	"github.com/BranceLee/Support/coreservice"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	conf := config.DefaultPostgresConfig()
	conn, err := gorm.Open(conf.Dialect(), conf.ConnectionInfo())
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}
	db = conn

	defaultHandler, err := coreservice.NewHandler(db)
	if err != nil {
		log.Fatalf("Failed to create default router")
	}

	defaultRouter := mux.NewRouter()
	defaultRouter.HandleFunc("/api/category/new", defaultHandler.GetCategory).Methods("GET")

	exitCode := m.Run()
	db.Close()
	os.Exit(exitCode)
}

type requestConfig struct {
	path        string
	method      string
	params      map[string]string
	accessToken string
}

// Create a mock request
func (rc *requestConfig) build(t *testing.T) *http.Request {
	var body io.Reader
	if len(rc.params) != 0 {
		data := url.Values{}
		for key, val := range rc.params {
			data.Set(key, val)
		}
		body = strings.NewReader(data.Encode())
	}
	request, err := http.NewRequest(rc.method, rc.path, body)
	if err != nil {
		t.Errorf("Failed to create http request: %s", err)
	}

	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	if rc.accessToken != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", rc.accessToken))
	}
	return request
}
