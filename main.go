package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var ErrBodyMustNotBeEmpty = errors.New("body must not be empty")

type GameResults struct {
	Date   string `json:"date"`
	Score  int    `json:"score"`
	Moves  int    `json:"moves"`
	Pieces int    `json:"pieces"`
}

type GlobalScore struct {
	GameResults
	mu sync.RWMutex
}

var currentScore GlobalScore

func formattedDate() string {
	return time.Now().Format("20060102")
}

func main() {
	currentScore = GlobalScore{
		GameResults: GameResults{
			Date: formattedDate(),
		},
	}

	mount := os.Getenv("MOUNT")
	http.HandleFunc("POST "+mount+"/scores", saveScoresHandler)
	http.HandleFunc("GET "+mount+"/scores", retrieveScoresHandler)
	http.HandleFunc("GET "+mount+"/", indexHandler)
	http.HandleFunc("GET "+mount+"/index.html", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5454"
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

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := filepath.Join("index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error parsing template: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// inline template to pass Date to index.html rand seed
	data := struct{ Date string }{
		Date: formattedDate(),
	}

	// Execute the template and write the response
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func retrieveScoresHandler(w http.ResponseWriter, r *http.Request) {
	currentDate := formattedDate()
	result := GameResults{
		Date: currentDate,
	}

	// --------- LOCK ---------
	currentScore.mu.RLock()
	if currentDate == currentScore.Date {
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

	//TODO: only consider server date here -- if input.Date != server date,
	// just return Results with new date + 0 scores

	date := input.Date
	if date == "" {
		date = formattedDate()
	}

	// --------- LOCK ---------
	currentScore.mu.Lock()

	if currentScore.Date == date {
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
		currentScore.Date = date
		currentScore.Score = input.Score
		currentScore.Moves = input.Moves
		currentScore.Pieces = input.Pieces
	}

	result := GameResults{
		Date:   currentScore.Date,
		Score:  currentScore.Score,
		Moves:  currentScore.Moves,
		Pieces: currentScore.Pieces,
	}

	currentScore.mu.Unlock()
	// -------- UNLOCK --------

	log.Printf("Saving %+v\n", result)

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
