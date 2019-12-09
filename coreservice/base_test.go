package coreservice_test

import (
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/support/config"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	conf := config.DefaultPostgresConfig()
	conn, err := gorm.Open(conf.Dialect(), conf.ConnectionInfo())
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}
	db = conn

	exitCode := m.Run()
	db.Close()
	os.Exit(exitCode)
}
