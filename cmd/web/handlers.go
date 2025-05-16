package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/sspencer/jawbreaker"
)

type Signals struct {
	Board           string `json:"board"`
	Rows            int    `json:"rows"`
	Cols            int    `json:"cols"`
	PiecesText      string `json:"piecesText"`
	RemainingPieces int    `json:"remainingPieces"`
	BonusScore      int    `json:"bonusScore"`
	CurrentScore    int    `json:"currentScore"`
	LastScore       int    `json:"lastScore"`
	BestScore       int    `json:"bestScore"`
	GameOver        bool   `json:"gameOver"`
}

type ScoreData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

type pageData struct {
	DS         bool
	DatastarJS string
	GameCode   template.JS
	GameSrc    string
	StyleCSS   string
	Size       int
	Block      int
	MSize      int
	MBlock     int
	CookieName string
	GameSize   template.CSS
	Game       template.HTML
	Board      string
	LastScore  int
	BestScore  int
}

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.cfg

	data := pageData{
		GameSrc:    fsys.HashName("static/jawbreaker.js"),
		StyleCSS:   fsys.HashName("static/style.css"),
		Size:       cfg.size,
		Block:      cfg.block,
		MSize:      cfg.msize,
		MBlock:     cfg.mblock,
		CookieName: cookieName,
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) oneHandler(w http.ResponseWriter, r *http.Request) {

	gameBytes, err := os.ReadFile("cmd/web/static/jawbreaker.js")
	if err != nil {
		slog.Error("jawbreaker.js not found", "error", err)
		gameBytes = gameCode
	}

	cfg := app.cfg
	data := pageData{
		GameCode:   template.JS(gameBytes),
		StyleCSS:   fsys.HashName("static/style.css"),
		Size:       cfg.size,
		Block:      cfg.block,
		MSize:      cfg.msize,
		MBlock:     cfg.mblock,
		CookieName: cookieName,
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(w, "one", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) datastarHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.cfg

	//g := jawbreaker.NewGame(cfg.size, cfg.size)
	opts := jawbreaker.GameOptions{}
	g := jawbreaker.NewGameWithOptions(cfg.size, cfg.size, opts.StartEmpty())

	var scoreData ScoreData
	if cookie, err := r.Cookie(cookieName); err == nil {
		deserializeScoreData(&scoreData, cookie.Value)
	}

	data := pageData{
		DS:         true,
		DatastarJS: fsys.HashName("static/datastar.js"),
		StyleCSS:   fsys.HashName("static/style.css"),
		GameSize:   template.CSS(gameSize),
		Game:       template.HTML(boardToHTML(g.Board(), nil)),
		Board:      g.Board().Base64(),
		Size:       cfg.size,
		Block:      cfg.block,
		MSize:      cfg.msize,
		MBlock:     cfg.mblock,
		LastScore:  scoreData.LastScore,
		BestScore:  scoreData.BestScore,
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, "datastar", data)
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

	g, err := jawbreaker.RestoreGame(signals.Board, signals.Rows, signals.Cols, signals.CurrentScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := g.Move(index)

	signals.GameOver = status.GameOver
	signals.CurrentScore = status.Score
	signals.Board = status.Board.Base64()

	if status.GameOver {
		signals.LastScore = signals.CurrentScore
		signals.BonusScore = status.Bonus
		signals.RemainingPieces = status.RemainingPieces
		if status.RemainingPieces == 1 {
			signals.PiecesText = "piece"
		} else {
			signals.PiecesText = "pieces"
		}

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
			Name:     cookieName,
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

	err = sse.MergeFragments(boardToHTML(g.Board(), g.GetConnectedPieces(index)))
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

	g, err := jawbreaker.RestoreGame(signals.Board, signals.Rows, signals.Cols, signals.CurrentScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sse := datastar.NewSSE(w, r)

	err = sse.MergeFragments(boardToHTML(g.Board(), g.GetConnectedPieces(index)))
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

	cfg := app.cfg

	opts := jawbreaker.GameOptions{}
	g := jawbreaker.NewGameWithOptions(cfg.size, cfg.size, opts.PowerUps())

	sse := datastar.NewSSE(w, r)

	// Update g fragment

	if app.cfg.animate {
		size := len(g.Board())
		piecesPerIter := 6
		maxIters := (size / piecesPerIter) + 1
		indices := jawbreaker.ShuffledIndices(size)
		for n := 0; n < maxIters; n++ {
			piecesToShow := (n + 1) * piecesPerIter
			board, err := g.AnimateBoard(indices, piecesToShow)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = sse.MergeFragments(boardToHTML(board, nil))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			time.Sleep(17 * time.Millisecond)
		}
	} else {
		err := sse.MergeFragments(boardToHTML(g.Board(), nil))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Send updated signals
	signals := map[string]any{
		"board":        g.Board().Base64(),
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

}
