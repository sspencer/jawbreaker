package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	Rows        = 12
	Cols        = 12
	DefaultPort = "5454"
)

var (
	Pieces      = []string{"🟣", "🔵", "🟢", "🔴", "🟡"}
	globalScore GlobalScore
)

type GameResults struct {
	Date       int      `json:"date"`
	DailyBoard []string `json:"daily_board,omitempty"`
	Score      int      `json:"score"`
	Moves      int      `json:"moves"`
	Pieces     int      `json:"pieces"`
}

type PageData struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
	GameResults
}

type GlobalScore struct {
	GameResults
	mu sync.RWMutex
}

func main() {
	//globalScore = GlobalScore{}
	maybeResetGlobalScores()

	mount := os.Getenv("MOUNT")
	http.HandleFunc("POST "+mount+"/scores", saveScoresHandler)
	http.HandleFunc("GET "+mount+"/scores", retrieveScoresHandler)
	http.HandleFunc("GET "+mount+"/", indexHandler)
	http.HandleFunc("GET "+mount+"/index.html", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	log.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, LoggingMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
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
	templatePath := filepath.Join("index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		handleServerError(w, "Error parsing template", err)
		return
	}

	data := getPageData()

	if err := tmpl.Execute(w, data); err != nil {
		handleServerError(w, "Error executing template", err)
		return
	}
}

func retrieveScoresHandler(w http.ResponseWriter, r *http.Request) {
	maybeResetGlobalScores()
	results := getGameResults()

	err := sendJSON(w, http.StatusOK, results)
	if err != nil {
		handleServerError(w, "Error sending response", err)
		return
	}
}

func saveScoresHandler(w http.ResponseWriter, r *http.Request) {
	var input GameResults
	err := decodeJSON(w, r, &input, true)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	results := updateGlobalScores(input)

	log.Printf("score=%d, moves=%d, pieces=%d\n", results.Score, results.Moves, results.Pieces)

	err = sendJSON(w, http.StatusOK, results)
	if err != nil {
		handleServerError(w, "Error sending response", err)
		return
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

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}, disallowUnknownFields bool) error {
	maxBytes := 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func handleServerError(w http.ResponseWriter, message string, err error) {
	log.Printf("%s: %s\n", message, err.Error())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func generateDailyBoard(seed int) []string {
	rnd := rand.New(rand.NewSource(int64(seed)))
	size := Rows * Cols
	num := len(Pieces)
	board := make([]string, size)
	for i := 0; i < size; i++ {
		board[i] = Pieces[rnd.Intn(num)]
	}

	return board
}

func currentDate() int {
	now := time.Now()
	return now.Year()*10000 + int(now.Month())*100 + now.Day()
}

func getPageData() PageData {
	data := PageData{
		Rows: Rows,
		Cols: Cols,
	}
	globalScore.mu.RLock()
	data.GameResults = globalScore.GameResults
	globalScore.mu.RUnlock()
	return data
}

func getGameResults() GameResults {
	data := GameResults{}
	globalScore.mu.RLock()
	data = globalScore.GameResults
	data.DailyBoard = nil
	globalScore.mu.RUnlock()
	return data
}

func updateGlobalScores(input GameResults) GameResults {
	now := currentDate()

	globalScore.mu.Lock()
	defer globalScore.mu.Unlock()

	if globalScore.Date == now && input.Date == now {
		if input.Score > globalScore.Score {
			globalScore.Score = input.Score
		}
		if globalScore.Moves == 0 || input.Moves < globalScore.Moves {
			globalScore.Moves = input.Moves
		}
		if globalScore.Pieces == 0 || input.Pieces < globalScore.Pieces {
			globalScore.Pieces = input.Pieces
		}
	} else {
		globalScore.Date = now
		globalScore.DailyBoard = generateDailyBoard(now)
		globalScore.Score = 0
		globalScore.Moves = 0
		globalScore.Pieces = 0
	}

	data := GameResults{}
	data = globalScore.GameResults
	return data
}

func maybeResetGlobalScores() {
	date := currentDate()
	globalScore.mu.Lock()
	defer globalScore.mu.Unlock()

	if globalScore.Date == date {
		return // If the dates are the same, do nothing
	}

	// Update the global score
	globalScore.Date = date
	globalScore.DailyBoard = generateDailyBoard(date)
	globalScore.Score = 0
	globalScore.Moves = 0
	globalScore.Pieces = 0
}
