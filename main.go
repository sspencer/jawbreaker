package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var ErrBodyMustNotBeEmpty = errors.New("body must not be empty")

type GameResults struct {
	Score  int `json:"score"`
	Moves  int `json:"moves"`
	Pieces int `json:"pieces"`
}

type GlobalScore struct {
	GameResults
	date string
	mu   sync.RWMutex
}

var currentScore GlobalScore

func formattedDate() string {
	return time.Now().Format("20060102")
}

func main() {
	currentScore = GlobalScore{
		GameResults: GameResults{},
		date:        formattedDate(), // "20240827"
	}

	http.HandleFunc("POST /scores", saveScoresHandler)        // via proxy
	http.HandleFunc("POST /api/scores", saveScoresHandler)    // local dev
	http.HandleFunc("GET /scores", retrieveScoresHandler)     // via proxy
	http.HandleFunc("GET /api/scores", retrieveScoresHandler) // local dev
	http.HandleFunc("GET /", indexHandler)
	http.HandleFunc("GET /index.html", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5454"
	}

	fmt.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /index.html\n")
	http.ServeFile(w, r, "index.html")
}

func retrieveScoresHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /api/scores\n")
	currentDate := formattedDate()
	result := GameResults{}

	// --------- LOCK ---------
	currentScore.mu.RLock()
	if currentDate == currentScore.date {
		result.Score = currentScore.Score
		result.Moves = currentScore.Moves
		result.Pieces = currentScore.Pieces
	}
	currentScore.mu.RUnlock()
	// -------- UNLOCK --------

	// if score wasn't set above (fell into new date), score/moves default to 0
	// by nature of Go's default initialization

	err := sendJSON(w, http.StatusOK, result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	fmt.Printf("POST /api/scores %+v\n", input)

	date := formattedDate()
	result := GameResults{}

	// --------- LOCK ---------
	currentScore.mu.Lock()

	if currentScore.date == date {
		if input.Score > currentScore.Score {
			currentScore.Score = input.Score
		}
		if currentScore.Moves == 0 || input.Moves < currentScore.Moves {
			currentScore.Moves = input.Moves
		}
		if currentScore.Pieces == 0 || input.Pieces < currentScore.Pieces {
			currentScore.Pieces = input.Pieces
		}
	} else {
		currentScore.date = date
		currentScore.Score = input.Score
		currentScore.Moves = input.Moves
		currentScore.Pieces = input.Pieces
	}

	result.Score = currentScore.Score
	result.Moves = currentScore.Moves
	result.Pieces = currentScore.Pieces

	currentScore.mu.Unlock()
	// -------- UNLOCK --------

	fmt.Printf("Saving on %q score=%d, moves=%d, pieces=%d\n", date, result.Score, result.Moves, result.Pieces)

	err = sendJSON(w, http.StatusOK, result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
			return ErrBodyMustNotBeEmpty

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
