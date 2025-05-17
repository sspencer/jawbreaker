package main

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/sspencer/jawbreaker"
)

const (
	clickAction = "click"
	mouseAction = "mouse"
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
	DS         bool
	DatastarJS string
	Session    string
	GameCode   template.JS
	GameSrc    string
	StyleCSS   string
	Size       int
	Block      int
	MSize      int
	MBlock     int
	CookieName string
	GameStyle  template.CSS
	Game       template.HTML
	LastScore  int
	BestScore  int
	GameMechanics
}

type Signals struct {
	Session   string `json:"session"`
	Board     string `json:"board"`
	Size      int    `json:"size"`
	Score     int    `json:"score"`
	Bonus     int    `json:"bonus"`
	Last      int    `json:"last"`
	Best      int    `json:"best"`
	Remaining int    `json:"remaining"`
	Pieces    string `json:"pieces"`
	GameOver  bool   `json:"gameOver"`
}

func (app *Application) indexHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) oneHandler(w http.ResponseWriter, r *http.Request) {

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

func (app *Application) datastarHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.cfg

	var scoreData ScoreData
	if cookie, err := r.Cookie(cookieName); err == nil {
		deserializeCookie(&scoreData, cookie.Value)
	}
	session, err := generateSessionID()
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data := pageData{
		DS:         true,
		DatastarJS: fsys.HashName("static/datastar.js"),
		StyleCSS:   fsys.HashName("static/style.css"),
		GameStyle:  template.CSS(gameSizeCSS(cfg.size, cfg.block)),
		CookieName: cookieName,
		Session:    session,
		Size:       cfg.size,
		Block:      cfg.block,
		MSize:      cfg.msize,
		MBlock:     cfg.mblock,
		LastScore:  scoreData.LastScore,
		BestScore:  scoreData.BestScore,
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
	err = tmpl.ExecuteTemplate(w, "datastar", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

type SessionSignal struct {
	Session string `json:"session"`
	Board   string `json:"board"`
	Size    int    `json:"size"`
	Score   int    `json:"score"`
	Last    int    `json:"last"`
	Best    int    `json:"best"`
}

func (app *Application) clickHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, clickAction)
	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) mouseHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, mouseAction)
	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) handleAction(w http.ResponseWriter, r *http.Request, actionName string) {
	index := extractNumberFromPieceId(chi.URLParam(r, "id"))
	if index < 0 && actionName == clickAction {
		return
	}

	signal := &SessionSignal{}
	if err := datastar.ReadSignals(r, signal); err != nil {
		slog.Error("Error reading signal", "error", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	//slog.Info("Read signal", "signal", signal)

	app.clientsMux.Lock()
	defer app.clientsMux.Unlock()

	clientChan, ok := app.clients[signal.Session]
	if !ok {
		slog.Error("Client not found", "session", signal.Session)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	action := Action{
		Action:  actionName,
		Index:   index,
		Session: signal.Session,
		Board:   signal.Board,
	}

	select {
	case clientChan <- action:
		//slog.Info("Sent action", "action", action)
		return
	default:
		// Channel blocked, skip
		//slog.Info("BLOCK action", "action", action)
		return
	}
}

func (app *Application) newGameHandler(w http.ResponseWriter, r *http.Request) {
	// Read the current store to preserve the best score
	signal := &SessionSignal{}
	if err := datastar.ReadSignals(r, signal); err != nil {
		slog.Error("Error reading signals", "error", err)
		http.Error(w, "Error reading signals", http.StatusInternalServerError)
		return
	}

	// Create a client channel
	clientChan := app.registerClient(signal.Session)
	defer app.unregisterClient(signal.Session)

	opts := jawbreaker.GameOptions{}
	g := jawbreaker.NewGameWithOptions(signal.Size, signal.Size, opts.PowerUps())
	app.gamesMux.Lock()
	app.games[signal.Session] = g
	app.gamesMux.Unlock()

	sse := datastar.NewSSE(w, r)

	err := sendGameStartSignals(sse, g.Board(), signal.Last, signal.Best)
	if err != nil {
		slog.Error("Error sending new game signals", "error", err)
		http.Error(w, "Error sending new game signals", http.StatusInternalServerError)
		return
	}

	// Keep connection alive and send events
	for {
		select {
		case action := <-clientChan:
			err := app.processAction(sse, &action)
			if err != nil {
				slog.Error("Error processing action", "error", err)
			}

		case <-r.Context().Done():
			slog.Debug("Client disconnected", "session", signal.Session)
			return
		}
	}
}

func (app *Application) processAction(sse *datastar.ServerSentEventGenerator, action *Action) error {

	app.gamesMux.Lock()
	g := app.games[action.Session]
	app.gamesMux.Unlock()

	if g == nil {
		return errors.New("Game not found")
	}

	if action.Action == clickAction {
		status := g.Move(action.Index)
		if status.GameOver {
			if err := sendGameOverSignals(sse, status); err != nil {
				return err
			}
		} else {
			if err := sendGamePlaySignals(sse, status.Board, status.Score); err != nil {
				return err
			}
		}
	} else if action.Action == mouseAction {
		board := g.BoardWithConnections(action.Index)
		if err := sendGamePlaySignals(sse, board, g.Score()); err != nil {
			return err
		}
	}

	return nil
}

func sendGameStartSignals(sse *datastar.ServerSentEventGenerator, board jawbreaker.Board, lastScore, bestScore int) error {
	signals := map[string]any{
		"board":    board.Base64(),
		"score":    0,
		"last":     lastScore,
		"best":     bestScore,
		"gameOver": false,
	}

	return sse.MarshalAndMergeSignals(signals)
}

func sendGamePlaySignals(sse *datastar.ServerSentEventGenerator, board jawbreaker.Board, score int) error {
	signals := map[string]any{
		"board": board.Base64(),
		"score": score,
	}

	return sse.MarshalAndMergeSignals(signals)
}

func sendGameOverSignals(sse *datastar.ServerSentEventGenerator, status jawbreaker.Status) error {
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

	return sse.MarshalAndMergeSignals(signals)
}
