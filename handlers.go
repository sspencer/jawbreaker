package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"
)

type Signals struct {
	Pieces       string `json:"pieces"`
	CurrentScore int    `json:"currentScore"`
	LastScore    int    `json:"lastScore"`
	BestScore    int    `json:"bestScore"`
	GameOver     bool   `json:"gameOver"`
}

type ScoreData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	g := NewGame(numRows, numCols)
	// Default scores
	scoreData := ScoreData{
		LastScore: 0,
		BestScore: 0,
	}

	// Try to read scores from the combined cookie
	if cookie, err := r.Cookie("scores"); err == nil {
		deserializeScoreData(&scoreData, cookie.Value)
	}

	// Create a template from the embedded HTML
	tmpl, err := template.New("index").Parse(string(indexHTML))
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Define data for the template
	data := struct {
		LastScore string
		BestScore string
		Pieces    string
		Game      template.HTML
		Hover     bool
	}{
		LastScore: strconv.Itoa(scoreData.LastScore),
		BestScore: strconv.Itoa(scoreData.BestScore),
		Pieces:    g.String(),
		Game:      template.HTML(gameToHTML(g)), // Using template.HTML to avoid escaping
		Hover:     os.Getenv("HOVER") == "1",
	}

	w.Header().Set("Content-Type", "text/html")

	// Execute the template with the data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	index := extractNumberFromPieceId(chi.URLParam(r, "id"))
	if index < 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	signals := &Signals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g, err := RestoreGame(signals.Pieces, numRows, numCols, signals.CurrentScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.Move(index)
	gameOver := g.IsGameOver()

	signals.GameOver = gameOver
	signals.CurrentScore = g.Score()
	signals.Pieces = g.String()

	if gameOver {
		signals.LastScore = signals.CurrentScore

		// Update the best score if current score is higher
		if signals.CurrentScore > signals.BestScore {
			signals.BestScore = signals.CurrentScore
		}

		// Save both scores to a single cookie
		scoreData := ScoreData{
			LastScore: signals.LastScore,
			BestScore: signals.BestScore,
		}

		scoresCookie := &http.Cookie{
			Name:     "scores",
			Value:    scoreData.serialize(),
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, scoresCookie)
	}

	sse := datastar.NewSSE(w, r)
	err = sse.MarshalAndMergeSignals(signals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sse.MergeFragments(gameToHTML(g))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	// Read current store to preserve best score
	store := &Signals{}
	if err := datastar.ReadSignals(r, store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g := NewGame(numRows, numCols)

	sse := datastar.NewSSE(w, r)

	// Send updated signals
	signals := map[string]any{
		"pieces":       g.String(),
		"currentScore": 0,
		"lastScore":    store.LastScore,
		"bestScore":    store.BestScore,
		"gameOver":     false,
	}

	err := sse.MarshalAndMergeSignals(signals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update g fragment
	err = sse.MergeFragments(gameToHTML(g))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
