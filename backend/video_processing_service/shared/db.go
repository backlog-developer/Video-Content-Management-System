// File: shared/db.go
package shared

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config.DBHost, Config.DBPort, Config.DBUser, Config.DBPassword, Config.DBName)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")
}
