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
	SessionID  string
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
	LastScore  int
	BestScore  int
	GameMechanics
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

	var scoreData ScoreData
	if cookie, err := r.Cookie(cookieName); err == nil {
		deserializeScoreCookie(&scoreData, cookie.Value)
	}
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pageData{
		DS:         true,
		DatastarJS: fsys.HashName("static/datastar.js"),
		StyleCSS:   fsys.HashName("static/style.css"),
		GameSize:   template.CSS(gameSize),
		SessionID:  sessionID,
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
	SessionID    string `json:"session"`
	Board        string `json:"board"`
	CurrentScore int    `json:"currentScore"`
	LastScore    int    `json:"lastScore"`
	BestScore    int    `json:"bestScore"`
}

func (app *application) clickHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, clickAction)
}

func (app *application) mouseHandler(w http.ResponseWriter, r *http.Request) {
	app.handleAction(w, r, mouseAction)
}

func (app *application) handleAction(w http.ResponseWriter, r *http.Request, action string) {
	index := extractNumberFromPieceId(chi.URLParam(r, "id"))
	if index < 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	signal := &SessionSignal{}
	if err := datastar.ReadSignals(r, signal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := serializeAction(action, index, signal.CurrentScore, signal.Board)
	slog.Info("Ready to send", "message", message)

	app.clientsMux.Lock()
	defer app.clientsMux.Unlock()

	clientChan, ok := app.clients[signal.SessionID]
	if !ok {
		//http.Error(w, "Client not found", http.StatusBadRequest)
		slog.Error("Client not found", "session", signal.SessionID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	select {
	case clientChan <- message:
		slog.Info("Sent action", "action", action, "index", index)
		w.WriteHeader(http.StatusNoContent)
	default:
		// Channel blocked, skip
		slog.Info("BLOCK action", "action", action, "index", index)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (app *application) newGameHandler(w http.ResponseWriter, r *http.Request) {
	// Read the current store to preserve the best score
	signal := &SessionSignal{}
	if err := datastar.ReadSignals(r, signal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a client channel
	clientChan := app.registerClient(signal.SessionID)
	defer app.unregisterClient(signal.SessionID)

	cfg := app.cfg
	opts := jawbreaker.GameOptions{}
	g := jawbreaker.NewGameWithOptions(cfg.size, cfg.size, opts.PowerUps())

	sse := datastar.NewSSE(w, r)

	err := sendGameStartSignals(sse, g.Board(), signal.LastScore, signal.BestScore)
	if err != nil {
		slog.Error("Error sending new game signals", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Keep connection alive and send events
	for {
		select {
		case message := <-clientChan:

			action := deserializeAction(message)
			if action == nil {
				slog.Error("Error deserializing action", "action", message)
				continue
			}

			g, err := jawbreaker.RestoreGame(action.Board, cfg.size, cfg.size, action.Score)
			if err != nil {
				slog.Error("Error restoring game", "error", err)
				continue
			}

			if action.Action == clickAction {
				status := g.Move(action.Index)
				if status.GameOver {
					err = sendGameOverSignals(sse, status)
				} else {
					err = sendGamePlaySignals(sse, status.Board, status.Score)
				}
			} else if action.Action == mouseAction {
				connected := g.GetConnectedPieces(action.Index)
				board := g.Board()
				for _, idx := range connected {
					board[idx] += jawbreaker.PieceSelected
				}

				err = sendGamePlaySignals(sse, board, g.Score())
			}

			if err != nil {
				slog.Error("Error sending game", "error", err)
				continue
			}

		case <-r.Context().Done():
			return
		}
	}
}

func sendGameStartSignals(sse *datastar.ServerSentEventGenerator, board jawbreaker.Board, lastScore, bestScore int) error {
	slog.Info("Sending game start signals")
	signals := struct {
		Board        string `json:"board"`
		CurrentScore int    `json:"currentScore"`
		LastScore    int    `json:"lastScore"`
		BestScore    int    `json:"bestScore"`
		GameOver     bool   `json:"gameOver"`
	}{
		Board:        board.Base64(),
		CurrentScore: 0,
		LastScore:    lastScore,
		BestScore:    bestScore,
		GameOver:     false,
	}

	return sse.MarshalAndMergeSignals(signals)
}

func sendGamePlaySignals(sse *datastar.ServerSentEventGenerator, board jawbreaker.Board, score int) error {
	slog.Info("Sending game play signals")
	signals := struct {
		Board        string `json:"board"`
		CurrentScore int    `json:"currentScore"`
	}{
		Board:        board.Base64(),
		CurrentScore: score,
	}

	return sse.MarshalAndMergeSignals(signals)
}

func sendGameOverSignals(sse *datastar.ServerSentEventGenerator, status jawbreaker.Status) error {
	slog.Info("Sending game over signals")
	pieceText := "pieces"
	if status.RemainingPieces == 1 {
		pieceText = "piece"
	}

	signals := struct {
		Board           string `json:"board"`
		PiecesText      string `json:"piecesText"`
		RemainingPieces int    `json:"remainingPieces"`
		BonusScore      int    `json:"bonusScore"`
		CurrentScore    int    `json:"currentScore"`
		LastScore       int    `json:"lastScore"`
		BestScore       int    `json:"bestScore"`
		GameOver        bool   `json:"gameOver"`
	}{
		Board:           status.Board.Base64(),
		CurrentScore:    status.Score,
		BonusScore:      status.Bonus,
		RemainingPieces: status.RemainingPieces,
		PiecesText:      pieceText,
		LastScore:       status.LastScore,
		BestScore:       status.BestScore,
		GameOver:        true,
	}

	return sse.MarshalAndMergeSignals(signals)
}

func oldCookieCode(w http.ResponseWriter) {
	// Save both scores to a single cookie
	scoreData := ScoreData{
		LastScore: 0,
		BestScore: 0,
	}

	scoresCookie := &http.Cookie{
		Name:     cookieName,
		Value:    scoreData.serializeScoreCookie(),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, scoresCookie)

}
