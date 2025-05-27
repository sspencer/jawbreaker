package main

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/sspencer/jawbreaker"
)

const (
	clickAction = "click"
	mouseAction = "mouse"
	newAction   = "new"
)

type ScoreData struct {
	LastScore int `json:"lastScore"`
	BestScore int `json:"bestScore"`
}

type GameMechanics struct {
	PieceSpace       int
	White            int
	Purple           int
	Blue             int
	Green            int
	Red              int
	Yellow           int
	PowerX           int
	PowerPlus        int
	PowerCircle      int
	PowerRect        int
	PowerFill        int
	PowerRotateRight int
	PowerRotateLeft  int
	PieceSelected    int
}

type pageData struct {
	DS          bool
	DatastarUrl string
	GameSrc     template.JS
	GameUrl     string
	GameStyle   template.CSS
	StyleUrl    string
	StyleSrc    template.CSS
	Rows        int
	Cols        int
	Block       int
	MRows       int
	MCols       int
	MBlock      int
	CookieName  string
	LastScore   int
	BestScore   int
	GameMechanics
}

func (app *Application) indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.cfg

	data := pageData{
		GameUrl:    jsys.HashName("js/jawbreaker.js"),
		StyleUrl:   fsys.HashName("static/style.css"),
		Rows:       cfg.rows,
		Cols:       cfg.cols,
		Block:      cfg.block,
		MRows:      cfg.mrows,
		MCols:      cfg.mcols,
		MBlock:     cfg.mblock,
		CookieName: app.getCookieName(),
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) oneHandler(w http.ResponseWriter, r *http.Request) {

	code, err := bundleLiveJavascript()
	if err != nil {
		slog.Error("Error bundling javascript", "error", err)
		http.Error(w, "Error loading javascript", http.StatusInternalServerError)
		return
	}

	css, err := staticFiles.ReadFile("static/style.css")
	if err != nil {
		slog.Error("Error reading style.css", "error", err)
		http.Error(w, "Error loading style.css", http.StatusInternalServerError)
		return
	}

	cfg := app.cfg
	data := pageData{
		GameSrc:    template.JS(code),
		StyleSrc:   template.CSS(css),
		Rows:       cfg.rows,
		Cols:       cfg.cols,
		Block:      cfg.block,
		MRows:      cfg.mrows,
		MCols:      cfg.mcols,
		MBlock:     cfg.mblock,
		CookieName: app.getCookieName(),
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(w, "one", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) datastarHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.cfg

	var scoreData ScoreData
	if cookie, err := r.Cookie(app.getCookieName()); err == nil {
		deserializeCookie(&scoreData, cookie.Value)
	}

	data := pageData{
		DS:          true,
		DatastarUrl: fsys.HashName("static/datastar.js"),
		StyleUrl:    fsys.HashName("static/style.css"),
		GameStyle:   template.CSS(app.gameSizeCSS()),
		CookieName:  app.getCookieName(),
		Rows:        cfg.rows,
		Cols:        cfg.cols,
		Block:       cfg.block,
		MRows:       cfg.mrows,
		MCols:       cfg.mcols,
		MBlock:      cfg.mblock,
		LastScore:   scoreData.LastScore,
		BestScore:   scoreData.BestScore,
		GameMechanics: GameMechanics{
			PieceSpace:       jawbreaker.PieceSpace,
			White:            jawbreaker.White,
			Purple:           jawbreaker.Purple,
			Blue:             jawbreaker.Blue,
			Green:            jawbreaker.Green,
			Red:              jawbreaker.Red,
			Yellow:           jawbreaker.Yellow,
			PowerX:           jawbreaker.PowerX,
			PowerPlus:        jawbreaker.PowerPlus,
			PowerCircle:      jawbreaker.PowerCircle,
			PowerRect:        jawbreaker.PowerRect,
			PowerFill:        jawbreaker.PowerFill,
			PowerRotateRight: jawbreaker.PowerRotateRight,
			PowerRotateLeft:  jawbreaker.PowerRotateLeft,
			PieceSelected:    jawbreaker.PieceSelected,
		},
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.ExecuteTemplate(w, "datastar", data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) newGameHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, newAction)
}

func (app *Application) clickHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, clickAction)
}

func (app *Application) mouseHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, mouseAction)
}

type Signals struct {
	Board string `json:"board"`
	Rows  int    `json:"rows"`
	Cols  int    `json:"cols"`
	Score int    `json:"score"`
	Last  int    `json:"last"`
	Best  int    `json:"best"`
}

func createGame(actionName string, s *Signals) (*jawbreaker.Game, error) {
	if actionName == newAction {
		opts := jawbreaker.GameOptions{}
		g := jawbreaker.NewGameWithOptions(s.Rows, s.Cols, opts.PowerUps())
		return g, nil
	}

	return jawbreaker.RestoreGame(s.Board, s.Rows, s.Cols, s.Score)
}

func (app *Application) handleAction(w http.ResponseWriter, r *http.Request, actionName string) {
	index := extractNumberFromPieceId(chi.URLParam(r, "id"))

	signal := &Signals{}
	if err := datastar.ReadSignals(r, signal); err != nil {
		slog.Error("Error reading signals", "error", err)
		//w.WriteHeader(http.StatusNoContent)
		return
	}

	g, err := createGame(actionName, signal)
	if err != nil {
		slog.Error("Error restoring game in click", "error", err)
		return
	}

	sse := datastar.NewSSE(w, r)
	switch actionName {
	case newAction:
		sendGameStartSignals(sse, g, signal.Last, signal.Best)

	case clickAction:
		if index >= 0 && index < len(g.Board()) {
			status := g.Move(index)
			if status.GameOver {
				sendGameOverSignals(sse, status)
			} else {
				sendGamePlaySignals(sse, status.Board, status.Score)
			}
		}

	case mouseAction:
		board := g.BoardWithConnections(index)
		sendGamePlaySignals(sse, board, g.Score())
	}
}

func sendGameStartSignals(sse *datastar.ServerSentEventGenerator, game *jawbreaker.Game, lastScore, bestScore int) {
	signals := map[string]any{
		"board":     game.Board().Base64(),
		"score":     0,
		"rows":      game.Rows(), // always square
		"cols":      game.Cols(), // always square
		"last":      lastScore,
		"best":      bestScore,
		"remaining": 0,
		"gameOver":  false,
	}

	if err := sse.MarshalAndMergeSignals(signals); err != nil {
		slog.Error("Error sending game start signals", "error", err)
	}
}

func sendGamePlaySignals(sse *datastar.ServerSentEventGenerator, board jawbreaker.Board, score int) {
	signals := map[string]any{
		"board": board.Base64(),
		"score": score,
	}

	if err := sse.MarshalAndMergeSignals(signals); err != nil {
		slog.Error("Error sending game play signals", "error", err)
	}
}

func sendGameOverSignals(sse *datastar.ServerSentEventGenerator, status jawbreaker.Status) {
	pieceText := "pieces"
	if status.Remaining == 1 {
		pieceText = "piece"
	}

	signals := map[string]any{
		"board":     status.Board.Base64(),
		"score":     status.Score,
		"bonus":     status.Bonus,
		"remaining": status.Remaining,
		"pieces":    pieceText,
		"last":      status.Last,
		"best":      status.Best,
		"gameOver":  true,
	}

	if err := sse.MarshalAndMergeSignals(signals); err != nil {
		slog.Error("Error sending game over signals", "error", err)
	}
}
