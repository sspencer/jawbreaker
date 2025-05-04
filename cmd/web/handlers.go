package main

import (
	_ "embed"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/sspencer/jawbreaker"
)

//go:embed index.html
var indexHTML []byte

//go:embed breaker.html
var breakerHTML []byte

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
	Style     string
}

type ScoreData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

func breakerHandler(w http.ResponseWriter, r *http.Request) {
	// Create a template from the embedded HTML
	tmpl, err := template.New("index").Parse(string(breakerHTML))
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	g := jawbreaker.NewGame(numRows, numCols)
	data := IndexData{
		Style:  fsys.HashName("static/style.css"),
		Pieces: g.Board().String(),
		Game:   template.HTML(gameToHTML(g, nil)),
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	g := jawbreaker.NewGame(numRows, numCols)
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
		Pieces:    g.Board().String(),
		Game:      template.HTML(gameToHTML(g, nil)),
		Datastar:  fsys.HashName("static/datastar.js"),
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

	g, err := jawbreaker.RestoreGame(signals.Pieces, numRows, numCols, signals.CurrentScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gs := g.Move(index)

	signals.GameOver = gs.GameOver
	signals.CurrentScore = gs.Score
	signals.Pieces = string(gs.Board)

	if gs.GameOver {
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

	err = sse.MergeFragments(gameToHTML(g, g.GetConnectedPieces(index)))
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

	g, err := jawbreaker.RestoreGame(signals.Pieces, numRows, numCols, signals.CurrentScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sse := datastar.NewSSE(w, r)

	err = sse.MergeFragments(gameToHTML(g, g.GetConnectedPieces(index)))
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

	g := jawbreaker.NewGame(numRows, numCols)

	sse := datastar.NewSSE(w, r)

	// Send updated signals
	signals := map[string]any{
		"pieces":       g.Board().String(),
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
	err = sse.MergeFragments(gameToHTML(g, nil))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
