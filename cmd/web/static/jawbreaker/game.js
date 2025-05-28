class JawbreakerGame {
    constructor() {
        this.board = null;
        this.renderer = null;
        this.animator = null;
        this.rows = 0;
        this.cols = 0;
        this.blockSize = 0;
        this.cookieName = "";
        this.score = 0;
        this.bonus = 0;
        this.remainingPieces = 0;
        this.lastScore = 0;
        this.bestScore = 0;
        this.canvas = null;
        this.hoverList = new Set();
        this.hoverIndex = -1;
        this.touch = false;
        this.undo = null;
        this.weightedFills = true;
        this.usePowerUps = true;
        this.settingsCookieName = "jawbreaker_settings";
    }

    clamp(value, min, max) {
        return Math.max(min, Math.min(max, Number.isInteger(value) ? value : min));
    }

    async init({
                   rows = 10,
                   cols = 10,
                   blockSize = 36,
                   cookieName = "jawbreaker_scores",
                   touch = false,
                   weightedFills = true,
               }) {
        // Check for settings cookie
        const settingsStr = CookieManager.getCookie(this.settingsCookieName);
        if (settingsStr) {
            try {
                const settings = JSON.parse(settingsStr);
                rows = settings.rows || rows;
                cols = settings.cols || cols;
                this.usePowerUps = settings.usePowerUps !== undefined ? settings.usePowerUps : true;
            } catch (e) {
                console.error("Error parsing settings cookie:", e);
            }
        }

        this.rows = this.clamp(rows, 8, 20);
        this.cols = this.clamp(cols, 8, 20);
        this.blockSize = blockSize;
        this.cookieName = cookieName;
        this.touch = touch;
        this.weightedFills = weightedFills;
        this.score = 0;
        this.canvas = document.getElementById("game-canvas");
        if (!this.canvas) {
            throw new Error("Canvas element 'game-canvas' not found.");
        }
        this.canvas.width =
            this.cols * this.blockSize + (this.cols - 1) * GameConfig.GAP + 2 * GameConfig.BORDER_WIDTH;
        this.canvas.height =
            this.rows * this.blockSize + (this.rows - 1) * GameConfig.GAP + 2 * GameConfig.BORDER_WIDTH;
        this.board = new Board(this.rows, this.cols);
        this.renderer = new Renderer(
            this.canvas.getContext("2d"),
            this.board,
            this.blockSize,
            this.rows,
            this.cols
        );
        this.animator = new AnimationManager(
            this.renderer,
            this.board,
            GameConfig.ANIMATE ? GameConfig.ANIMATE_PIECE_SPEED : 0
        );
        this.initializeBoard();
        const scores = CookieManager.getCookie(this.cookieName);
        if (scores) {
            const [last, best] = scores.split("|").map((s) => parseInt(s, 10));
            this.lastScore = isNaN(last) ? 0 : last;
            this.bestScore = isNaN(best) ? 0 : best;
        }
        document.getElementById("current-score").innerText = this.score;
        document.getElementById("last-score").innerText = this.lastScore;
        document.getElementById("best-score").innerText = this.bestScore;
        this.registerEvents();
        if (GameConfig.ANIMATE) {
            await this.animateNewBoard();
        } else {
            this.renderer.render();
        }
    }

    async animateNewBoard() {
        const size = this.rows * this.cols;
        const start = Array(size).fill(GameConfig.COLORS.WHITE);
        const end = this.board.board.slice();
        const indices = this.board.createShuffledIndices(size);
        await this.animator.start(start, end, indices, GameConfig.ANIMATE_BOARD_SPEED);
    }

    initializeBoard() {
        this.board.initialize(this.usePowerUps);
    }

    getCanvasCoordinates(e) {
        const rect = this.canvas.getBoundingClientRect();
        return {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top,
        };
    }

    getBoardIndex(x, y) {
        const adjustedX = x - GameConfig.BORDER_WIDTH;
        const adjustedY = y - GameConfig.BORDER_WIDTH;
        const col = Math.floor(adjustedX / (this.blockSize + GameConfig.GAP));
        const row = Math.floor(adjustedY / (this.blockSize + GameConfig.GAP));
        if (col < 0 || col >= this.cols || row < 0 || row >= this.rows) {
            return -1;
        }
        return this.board.getIndex(row, col);
    }

    getConnectedPieces(index) {
        if (
            index < 0 ||
            index >= this.board.board.length ||
            this.board.board[index] === GameConfig.COLORS.WHITE
        ) {
            return [];
        }
        const piece = this.board.board[index];
        const color = this.board.getColor(piece);
        const powerUp = this.board.getPowerUp(piece);
        if (
            GameConfig.POWER_CONFIG.has(powerUp) &&
            !GameConfig.POWER_CONFIG.get(powerUp).connections
        ) {
            return [index];
        }
        let powerUpIndices = [];
        if (
            GameConfig.POWER_CONFIG.has(powerUp) &&
            GameConfig.POWER_CONFIG.get(powerUp).connections
        ) {
            powerUpIndices = this.getConnectedPowerUps(index, color, powerUp);
            if (powerUpIndices.length > 0) {
                powerUpIndices.unshift(index);
            }
            if (color === GameConfig.COLORS.GRAY) {
                return powerUpIndices;
            }
        }
        const connected = [];
        const stack = [index];
        const visited = new Array(this.board.board.length).fill(false);
        visited[index] = true;
        while (stack.length > 0) {
            const current = stack.pop();
            connected.push(current);
            const { row, col } = this.board.getPoint(current);
            const neighbors = [
                row > 0 ? current - this.cols : -1,
                row < this.rows - 1 ? current + this.cols : -1,
                col > 0 ? current - 1 : -1,
                col < this.cols - 1 ? current + 1 : -1,
            ];
            for (const n of neighbors) {
                if (
                    n !== -1 &&
                    !visited[n] &&
                    this.board.board[n] !== GameConfig.COLORS.WHITE &&
                    this.board.getColor(this.board.board[n]) === color
                ) {
                    visited[n] = true;
                    stack.push(n);
                }
            }
        }
        const combined = [...new Set([...powerUpIndices, ...connected])];
        return combined.length >= 2 ? combined : [];
    }

    getConnectedPowerUps(index, color, powerUp) {
        const handlers = {
            [GameConfig.POWER_UPS.X]: () =>
                this.getXConnections(index, color),
            [GameConfig.POWER_UPS.PLUS]: () =>
                this.getPlusConnections(index, color),
            [GameConfig.POWER_UPS.RECT]: () =>
                this.getRectConnections(index, color),
            [GameConfig.POWER_UPS.CIRCLE]: () =>
                this.getCircularConnections(index, color),
            [GameConfig.POWER_UPS.DISC]: () =>
                this.getDiscConnections(index, color),
        };
        return handlers[powerUp]?.() || [];
    }

    getConnectedDirections(index, color, directions, maxIterations) {
        const { row, col } = this.board.getPoint(index);
        const indices = [];
        for (let iter = 0; iter < maxIterations; iter++) {
            for (const dir of directions) {
                const r = row + dir.sr + dir.dr * iter;
                const c = col + dir.sc + dir.dc * iter;
                if (r >= 0 && r < this.rows && c >= 0 && c < this.cols) {
                    const i = this.board.getIndex(r, c);
                    const pieceColor = this.board.getColor(this.board.board[i]);
                    if (
                        (color === GameConfig.COLORS.GRAY &&
                            pieceColor !== GameConfig.COLORS.WHITE &&
                            pieceColor !== GameConfig.COLORS.GRAY) ||
                        (color !== GameConfig.COLORS.GRAY && color === pieceColor)
                    ) {
                        if (!indices.includes(i)) {
                            indices.push(i);
                        }
                    }
                }
            }
        }
        return indices;
    }

    getXConnections(index, color) {
        const maxIters = Math.ceil(
            Math.sqrt(this.rows ** 2 + this.cols ** 2)
        );
        return this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.X,
            maxIters
        );
    }

    getPlusConnections(index, color) {
        const maxIters = Math.max(this.rows, this.cols);
        return this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.PLUS,
            maxIters
        );
    }

    getRectConnections(index, color) {
        return this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.RECT,
            4
        );
    }

    getCircularConnections(index, color) {
        const c1 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.RECT,
            1
        );
        const c2 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.CIRCLE,
            3
        );
        return [...c1, ...c2];
    }

    getDiscConnections(index, color) {
        const c1 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.DISC1,
            3
        );
        const c2 = this.getConnectedDirections(
            index,
            color,
            GameConfig.POWER_UP_DIRECTIONS.DISC2,
            2
        );
        return [...c1, ...c2];
    }

    async removePieces(index) {
        const piece = this.board.board[index];
        const powerUp = this.board.getPowerUp(piece);
        let count = 0;
        let isSpecial = false;
        if (powerUp === GameConfig.POWER_UPS.FILL) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.applyGravity();
            await this.animateTransformation(() => this.board.fillSpaces());
            count = 1;
            isSpecial = true;
        } else if (powerUp === GameConfig.POWER_UPS.EXCHANGE) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.applyGravity();
            await this.animateTransformation(() => this.board.fillPieces());
            count = 1;
            isSpecial = true;
        } else if (powerUp === GameConfig.POWER_UPS.RIGHT) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.rotate(90);
            this.board.applyGravity();
            count = 1;
            isSpecial = true;
        } else if (powerUp === GameConfig.POWER_UPS.LEFT) {
            this.board.board[index] = GameConfig.COLORS.WHITE;
            this.board.rotate(270);
            this.board.applyGravity();
            count = 1;
            isSpecial = true;
        } else {
            const pieces = this.getConnectedPieces(index);
            for (const i of pieces) {
                if (this.board.board[i] !== GameConfig.COLORS.WHITE) {
                    this.board.board[i] = GameConfig.COLORS.WHITE;
                    count++;
                }
            }
        }
        return { count, isSpecial };
    }

    async animateTransformation(transformFn) {
        const start = this.board.board.slice();
        transformFn();
        const end = this.board.board.slice();
        const indices = this.board.createShuffledIndices(this.board.board.length).filter(
            (i) => start[i] !== end[i]
        );
        await this.animator.start(start, end, indices, GameConfig.ANIMATE_PIECE_SPEED);
    }

    calculateMoveScore(count) {
        return count < 2 ? 0 : count * (count - 1);
    }

    calculateBonus(remaining) {
        if (remaining === 0) return 2000;
        if (remaining <= 10) return (11 - remaining) * 100;
        return 0;
    }

    isGameOver() {
        for (let i = 0; i < this.board.board.length; i++) {
            if (this.board.board[i] === GameConfig.COLORS.WHITE) continue;
            const powerUp = this.board.getPowerUp(this.board.board[i]);
            if (
                GameConfig.POWER_CONFIG.has(powerUp) &&
                !GameConfig.POWER_CONFIG.get(powerUp).connections
            ) {
                return false;
            }
            if (this.getConnectedPieces(i).length > 0) {
                return false;
            }
        }
        return true;
    }

    async processMove(index) {
        this.undo = { board: this.board.board.slice(), score: this.score };
        document.getElementById("undo-btn").disabled = false;
        const start = this.board.board.slice();
        const { count, isSpecial } = await this.removePieces(index);
        if (count > 0 && !isSpecial) {
            const changed = this.board.board
            .map((p, i) => (p !== start[i] ? i : null))
            .filter((i) => i !== null);
            await this.animator.start(start, this.board.board, changed, 1);
            this.board.applyGravity();
        }
        if (!isSpecial) {
            this.score += this.calculateMoveScore(count);
        }
        const gameOver = this.isGameOver();
        if (gameOver) {
            this.remainingPieces = this.board.board.filter(
                (p) => p !== GameConfig.COLORS.WHITE
            ).length;
            this.bonus = this.calculateBonus(this.remainingPieces);
            this.score += this.bonus;
            this.undo = null;
            this.lastScore = this.score;
            if (this.score > this.bestScore) {
                this.bestScore = this.score;
            }
            document.getElementById("game-over-score").innerText = this.score;
            document.getElementById("last-score").innerText = this.lastScore;
            document.getElementById("best-score").innerText = this.bestScore;
            document.getElementById("bonus-points").innerText = this.bonus;
            document.getElementById("remaining-pieces").innerText = this.remainingPieces;
            document.getElementById("pieces-text").innerText =
                this.remainingPieces === 1 ? "piece" : "pieces";
            document.getElementById("game-over-overlay").classList.add("visible");
            CookieManager.setCookie(
                this.cookieName,
                `${this.lastScore}|${this.bestScore}`
            );
        }
        document.getElementById("current-score").innerText = this.score;
        this.hoverList.clear();
        this.hoverIndex = -1;
        if (this.touch) {
            const pieces = this.getConnectedPieces(index);
            if (pieces.length > 1) {
                this.hoverIndex = index;
                pieces.forEach((i) => this.hoverList.add(i));
            }
        }
        this.renderer.render(this.hoverList, this.hoverIndex);
        return { score: this.score, gameOver };
    }

    undoMove() {
        if (!this.undo) return;
        this.board.board = this.undo.board.slice();
        this.score = this.undo.score;
        this.undo = null;
        document.getElementById("undo-btn").disabled = true;
        document.getElementById("current-score").innerText = this.score;
        this.renderer.render();
    }

    async handleClick(e) {
        if (this.animator.isAnimating()) return;
        const { x, y } = this.getCanvasCoordinates(e);
        const index = this.getBoardIndex(x, y);
        if (
            index >= 0 &&
            index < this.board.board.length &&
            this.board.board[index] !== GameConfig.COLORS.WHITE
        ) {
            await this.processMove(index);
        }
    }

    handleMouseMove(e) {
        if (this.animator.isAnimating()) return;
        const { x, y } = this.getCanvasCoordinates(e);
        const index = this.getBoardIndex(x, y);
        if (index === this.hoverIndex) return;
        this.hoverList.clear();
        this.hoverIndex = index;
        if (
            index >= 0 &&
            index < this.board.board.length &&
            this.board.board[index] !== GameConfig.COLORS.WHITE
        ) {
            this.getConnectedPieces(index).forEach((i) => this.hoverList.add(i));
        }
        this.renderer.render(this.hoverList, this.hoverIndex);
    }

    handleMouseLeave() {
        if (this.animator.isAnimating() || this.hoverList.size === 0) return;
        this.hoverList.clear();
        this.hoverIndex = -1;
        this.renderer.render();
    }

    handleTouch(e) {
        if (this.animator.isAnimating()) return;
        const { x, y } = this.getCanvasCoordinates(e);
        const index = this.getBoardIndex(x, y);
        if (index === this.hoverIndex) {
            this.handleClick(e);
        } else {
            this.handleMouseMove(e);
        }
    }

    async resetGame() {
        // Update canvas size based on current rows and cols
        this.canvas.width =
            this.cols * this.blockSize + (this.cols - 1) * GameConfig.GAP + 2 * GameConfig.BORDER_WIDTH;
        this.canvas.height =
            this.rows * this.blockSize + (this.rows - 1) * GameConfig.GAP + 2 * GameConfig.BORDER_WIDTH;

        // Update renderer with new dimensions
        this.renderer = new Renderer(
            this.canvas.getContext("2d"),
            this.board,
            this.blockSize,
            this.rows,
            this.cols
        );

        // Update animator with new renderer
        this.animator = new AnimationManager(
            this.renderer,
            this.board,
            GameConfig.ANIMATE ? GameConfig.ANIMATE_PIECE_SPEED : 0
        );

        // Create a new board with the updated dimensions
        this.board = new Board(this.rows, this.cols);

        this.initializeBoard();
        this.hoverList.clear();
        this.hoverIndex = -1;
        this.score = 0;
        this.bonus = 0;
        this.remainingPieces = 0;
        this.undo = null;
        document.getElementById("undo-btn").disabled = true;
        document.getElementById("current-score").innerText = this.score;
        document.getElementById("game-over-overlay").classList.remove("visible");
        if (GameConfig.ANIMATE) {
            await this.animateNewBoard();
        } else {
            this.renderer.render();
        }
    }

    registerEvents() {
        if (this.touch) {
            this.canvas.addEventListener("click", (e) => this.handleTouch(e));
        } else {
            this.canvas.addEventListener("click", (e) => this.handleClick(e));
            this.canvas.addEventListener("mousemove", (e) => this.handleMouseMove(e));
            this.canvas.addEventListener("mouseleave", () => this.handleMouseLeave());
        }
        const newGameBtn = document.getElementById("new-game-btn");
        if (newGameBtn) {
            newGameBtn.addEventListener("click", (e) => {
                e.preventDefault();
                this.resetGame();
            });
        }
        const restartBtn = document.getElementById("restart-btn");
        if (restartBtn) {
            restartBtn.addEventListener("click", (e) => {
                e.preventDefault();
                this.resetGame();
            });
        }
        const undoBtn = document.getElementById("undo-btn");
        if (undoBtn) {
            undoBtn.disabled = true;
            undoBtn.addEventListener("click", (e) => {
                e.preventDefault();
                this.undoMove();
            });
        }
        const helpBtn = document.getElementById("help-btn");
        const helpModal = document.getElementById("help-modal");
        const helpCloseBtn = helpModal?.querySelector(".close");
        if (helpBtn && helpModal && helpCloseBtn) {
            helpBtn.addEventListener("click", () => {
                helpModal.style.display = "block";
            });
            helpCloseBtn.addEventListener("click", () => {
                helpModal.style.display = "none";
            });
            window.addEventListener("click", (e) => {
                if (e.target === helpModal) {
                    helpModal.style.display = "none";
                }
            });
        }

        // Settings modal functionality
        const settingsBtn = document.getElementById("settings-btn");
        const settingsModal = document.getElementById("settings-modal");
        const settingsCloseBtn = settingsModal?.querySelector(".close");
        const rowsSlider = document.getElementById("rows-slider");
        const colsSlider = document.getElementById("cols-slider");
        const rowsValue = document.getElementById("rows-value");
        const colsValue = document.getElementById("cols-value");
        const powerupsCheckbox = document.getElementById("powerups-checkbox");
        const settingsSaveBtn = document.getElementById("settings-save-btn");

        if (settingsBtn && settingsModal && settingsCloseBtn && rowsSlider && colsSlider && 
            rowsValue && colsValue && powerupsCheckbox && settingsSaveBtn) {

            // Initialize settings values
            rowsSlider.value = this.rows;
            colsSlider.value = this.cols;
            rowsValue.textContent = this.rows;
            colsValue.textContent = this.cols;
            powerupsCheckbox.checked = this.usePowerUps;

            // Update value displays when sliders change
            rowsSlider.addEventListener("input", () => {
                rowsValue.textContent = rowsSlider.value;
            });

            colsSlider.addEventListener("input", () => {
                colsValue.textContent = colsSlider.value;
            });

            // Open settings modal
            settingsBtn.addEventListener("click", () => {
                settingsModal.style.display = "block";
            });

            // Close settings modal
            settingsCloseBtn.addEventListener("click", () => {
                settingsModal.style.display = "none";
            });

            // Close modal when clicking outside
            window.addEventListener("click", (e) => {
                if (e.target === settingsModal) {
                    settingsModal.style.display = "none";
                }
            });

            // Save settings and restart game
            settingsSaveBtn.addEventListener("click", () => {
                const settings = {
                    rows: parseInt(rowsSlider.value, 10),
                    cols: parseInt(colsSlider.value, 10),
                    usePowerUps: powerupsCheckbox.checked
                };

                // Save settings to cookie
                CookieManager.setCookie(this.settingsCookieName, JSON.stringify(settings));

                // Update game settings
                this.rows = this.clamp(settings.rows, 8, 20);
                this.cols = this.clamp(settings.cols, 8, 20);
                this.usePowerUps = settings.usePowerUps;

                // Close modal
                settingsModal.style.display = "none";

                // Restart game with new settings
                this.resetGame();
            });
        }
    }
}

async function initGame(opts) {
    const game = new JawbreakerGame();
    await game.init(opts);
}
