package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

type Record struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_DATABASE")
    dbUser := os.Getenv("DB_USERNAME")
    dbPassword := os.Getenv("DB_PASSWORD")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    http.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w) // CORS設定を適用
        records, err := fetchRecords(db)
        if err != nil {
            http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(records)
    })

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func fetchRecords(db *sql.DB) ([]Record, error) {
    rows, err := db.Query("SELECT id, name FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var records []Record
    for rows.Next() {
        var record Record
        if err := rows.Scan(&record.ID, &record.Name); err != nil {
            return nil, err
        }
        records = append(records, record)
    }

    return records, nil
}
