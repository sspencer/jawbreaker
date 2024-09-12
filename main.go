package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DefaultPort = "5454"
)

var (
	db *sql.DB
)

type GameStats struct {
	GamesPlayedAllTime int `json:"games_played_all_time"`
	GamesPlayedToday   int `json:"games_played_today"`
	HighScore          int `json:"high_score"`
}

func main() {
	initDB()

	mount := os.Getenv("MOUNT")
	http.HandleFunc("GET "+mount+"/", indexHandler)
	http.HandleFunc("POST "+mount+"/score/{score}", scoreHandler) // Submit a score
	http.HandleFunc("GET "+mount+"/stats", statsHandler)          // Get today's highest score

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	log.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, LoggingMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

// Initialize the SQLite database
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "file:scores.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	execPragmas()
	createTables()
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		referer := r.Header.Get("Referer")
		log.Printf("%s %s, Ref: %q\n", r.Method, url, referer)

		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func scoreHandler(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("score")
	score, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, "Invalid score", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO game_scores (score) VALUES (?)", score)
	if err != nil {
		log.Printf("%s: %s\n", "scoreHandler", err.Error())
		http.Error(w, "Failed to record score", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Get UTC date in the format YYYY-MM-DD
	utc := time.Now().UTC().Format("2006-01-02")
	message := "statsHandler"

	var gamesPlayedAllTime int
	err := db.QueryRow("SELECT count(*) FROM game_scores").Scan(&gamesPlayedAllTime)
	if err != nil {
		log.Printf("%s: %s\n", message, err.Error())
	}

	var gamesPlayedToday int
	err = db.QueryRow("SELECT count(*) FROM game_scores WHERE DATE(timestamp) = ?", utc).Scan(&gamesPlayedToday)
	if err != nil {
		log.Printf("%s: %s\n", message, err.Error())
	}

	var highScore int
	err = db.QueryRow("SELECT MAX(score) FROM game_scores WHERE DATE(timestamp) = ?", utc).Scan(&highScore)
	if err != nil {
		log.Printf("%s: %s\n", message, err.Error())
	}

	stats := GameStats{
		GamesPlayedAllTime: gamesPlayedAllTime,
		GamesPlayedToday:   gamesPlayedToday,
		HighScore:          highScore,
	}

	err = sendJSON(w, http.StatusOK, stats)
	if err != nil {
		handleServerError(w, "Failed to send JSON", err)
	}
}

func sendJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)

	return err
}

func handleServerError(w http.ResponseWriter, message string, err error) {
	log.Printf("%s: %s\n", message, err.Error())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func execPragmas() {
	var err error
	// Set various pragmas to optimize the database connection.
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA cache_size=-64000;") // Set cache size.
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA busy_timeout=30000;") // Set busy timeout.
	if err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	scoresTable := `
	CREATE TABLE IF NOT EXISTS game_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		score INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	indexDate := `CREATE INDEX IF NOT EXISTS idx_game_date ON game_scores(timestamp);`
	indexScoreDate := `CREATE INDEX IF NOT EXISTS idx_game_score_date ON game_scores(score, timestamp);`

	var err error
	_, err = db.Exec(scoresTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(indexDate)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(indexScoreDate)
	if err != nil {
		log.Fatal(err)
	}
}
