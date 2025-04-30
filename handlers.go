package main

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"
)

//go:embed tmpl/index.html
var indexHTML []byte
var noConnections = []int{}

type Signals struct {
	Pieces       string `json:"pieces"`
	CurrentScore int    `json:"currentScore"`
	LastScore    int    `json:"lastScore"`
	BestScore    int    `json:"bestScore"`
	GameOver     bool   `json:"gameOver"`
}

type IndexData struct {
	LastScore int
	BestScore int
	Pieces    string
	Game      template.HTML
	Datastar  string
	GameJS    string
	Style     string
}

type ScoreData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

func envTrue(env string) bool {
	val := strings.ToLower(os.Getenv(env))
	return val == "1" || val == "true"
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

	w.Header().Set("Content-Type", "text/html")

	data := IndexData{
		LastScore: scoreData.LastScore,
		BestScore: scoreData.BestScore,
		Pieces:    g.String(),
		Game:      template.HTML(gameToHTML(g, noConnections)),
		Datastar:  fsys.HashName("static/datastar.js"),
		GameJS:    fsys.HashName("static/game.js"),
		Style:     fsys.HashName("static/style.css"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func clickHandler(w http.ResponseWriter, r *http.Request) {
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

	err = sse.MergeFragments(gameToHTML(g, g.getConnectedPieces(index)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func mouseHandler(w http.ResponseWriter, r *http.Request) {
	index := extractNumberFromPieceId(chi.URLParam(r, "id"))

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

	sse := datastar.NewSSE(w, r)

	err = sse.MergeFragments(gameToHTML(g, g.getConnectedPieces(index)))
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
	err = sse.MergeFragments(gameToHTML(g, noConnections))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
