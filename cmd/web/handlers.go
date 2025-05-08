package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/sspencer/jawbreaker"
)

type Signals struct {
	Pieces       string `json:"pieces"`
	CurrentScore int    `json:"currentScore"`
	LastScore    int    `json:"lastScore"`
	BestScore    int    `json:"bestScore"`
	GameOver     bool   `json:"gameOver"`
}

type IndexData struct {
	LastScore    int
	BestScore    int
	Pieces       string
	Game         template.HTML
	GameSize     template.CSS
	DS           bool
	DatastarJS   string
	JawbreakerJS string
	StyleCSS     string
	Rows         int
	Cols         int
	Block        int
	Border       int
	Gap          int
	CookieName   string
}

type ScoreData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

func (app *application) jsHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.getConfig()
	block := cfg.block
	if isMobile(r) {
		block = cfg.mblock
	}
	data := IndexData{
		JawbreakerJS: fsys.HashName("static/jawbreaker.js"),
		StyleCSS:     fsys.HashName("static/style.css"),
		Rows:         cfg.rows,
		Cols:         cfg.cols,
		Block:        block,
		Gap:          cfg.gap,
		Border:       cfg.border,
		CookieName:   cfg.cookie,
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, "js", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.getConfig()

	g := jawbreaker.NewGame(cfg.rows, cfg.cols)
	var scoreData ScoreData
	if cookie, err := r.Cookie(cfg.cookie); err == nil {
		deserializeScoreData(&scoreData, cookie.Value)
	}

	block := cfg.block
	if isMobile(r) {
		block = cfg.mblock
	}

	gameSize := getGameSize(cfg.rows, cfg.cols, block, cfg.gap*2)

	data := IndexData{
		DS:         true,
		DatastarJS: fsys.HashName("static/datastar.js"),
		StyleCSS:   fsys.HashName("static/style.css"),
		GameSize:   template.CSS(gameSize),
		Game:       template.HTML(gameToHTML(g, nil)),
		Pieces:     g.Board().String(),
		LastScore:  scoreData.LastScore,
		BestScore:  scoreData.BestScore,
		Rows:       cfg.rows,
		Cols:       cfg.cols,
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) clickHandler(w http.ResponseWriter, r *http.Request) {
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

	cfg := app.getConfig()
	g, err := jawbreaker.RestoreGame(signals.Pieces, cfg.rows, cfg.cols, signals.CurrentScore)
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

		// Update the best score if the current score is higher
		if signals.CurrentScore > signals.BestScore {
			signals.BestScore = signals.CurrentScore
		}

		// Save both scores to a single cookie
		scoreData := ScoreData{
			LastScore: signals.LastScore,
			BestScore: signals.BestScore,
		}

		scoresCookie := &http.Cookie{
			Name:     cfg.cookie,
			Value:    scoreData.serialize(),
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
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

func (app *application) mouseHandler(w http.ResponseWriter, r *http.Request) {
	index := extractNumberFromPieceId(chi.URLParam(r, "id"))

	signals := &Signals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := app.getConfig()
	g, err := jawbreaker.RestoreGame(signals.Pieces, cfg.rows, cfg.cols, signals.CurrentScore)
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

func (app *application) newGameHandler(w http.ResponseWriter, r *http.Request) {
	// Read the current store to preserve the best score
	store := &Signals{}
	if err := datastar.ReadSignals(r, store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := app.getConfig()
	g := jawbreaker.NewGame(cfg.rows, cfg.cols)

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
