class JB {
    // Piece constants
    static WHITE = 'w';
    static PURPLE = 'p';
    static BLUE = 'b';
    static GREEN = 'g';
    static RED = 'r';
    static YELLOW = 'y';

    // Color pieces array for random generation
    static COLOR_PIECES = [
        JB.PURPLE,
        JB.BLUE,
        JB.GREEN,
        JB.RED,
        JB.YELLOW
    ];

    // Color map for HTML class names
    static COLOR_MAP = new Map([
        [JB.PURPLE, "purple"],
        [JB.BLUE, "blue"],
        [JB.GREEN, "green"],
        [JB.RED, "red"],
        [JB.YELLOW, "yellow"]
    ]);

    /**
     * Create a new Jawbreaker game
     */
    constructor(opts) {
        this.rows = opts.rows;
        this.cols = opts.cols;
        this.cookieName = opts.cookieName;
        this.score = 0;
        this.lastScore = 0;
        this.bestScore = 0;
        this.board = opts.pieces;
        this.hoverList = new Set();

        // Read scores cookie
        const scoresCookie = this.getCookie(this.cookieName);
        if (scoresCookie) {
            const [lastScore, bestScore] = scoresCookie.split('|');
            this.lastScore = parseInt(lastScore, 10) || 0;
            this.bestScore = parseInt(bestScore, 10) || 0;

            // Update UI
            document.getElementById('last-score').innerText = this.lastScore;
            document.getElementById('best-score').innerText = this.bestScore;
        }

        this.registerEvents();
    }

    registerEvents() {
        const game = document.getElementById("game-container");
        game.addEventListener('click', (e) => {
            this.moveEvent(e);
        }, true);
        game.addEventListener('mouseover', (e) => {
            this.mouseOverEvent(e);
        }, true);
        game.addEventListener('mouseleave', (e) => {
            this.mouseLeaveEvent(e);
        }, true);

        // Add an event listener for the "New Game" button
        const newGameBtn = document.querySelector('.new-game-btn');
        if (newGameBtn) {
            newGameBtn.addEventListener('click', (e) => {
                this.resetGame();
                document.getElementById('game-over-overlay').classList.remove('visible');
                e.preventDefault();
            });
        }

        // Help button and modal
        const helpBtn = document.getElementById("help-btn");
        const helpModal = document.getElementById("help-modal");
        const closeBtn = document.querySelector(".close");

        helpBtn.addEventListener("click", () => {
            helpModal.style.display = "block";
        });

        closeBtn.addEventListener("click", () => {
            helpModal.style.display = "none";
        });

        window.addEventListener("click", (e) => {
            if (e.target === helpModal) {
                helpModal.style.display = "none";
            }
        });

    }

    /**
     * Reset the game when the "New Game" button is clicked
     */
    resetGame() {
        this.board = this.newBoard();
        this.hoverList.clear();
        this.score = 0;
        document.getElementById('current-score').innerText = this.score;
        this.renderBoard();
    }

    /**
     * Create a random board with the given dimensions
     * @returns {string} - Board as a string
     */
    newBoard() {
        let board = '';
        for (let i = 0; i < this.rows * this.cols; i++) {
            const rnd = Math.floor(Math.random() * JB.COLOR_PIECES.length);
            board += JB.COLOR_PIECES[rnd];
        }
        return board;
    }

    renderBoard() {
        let htmlString = '';

        for (let i = 0; i < this.board.length; i++) {
            let color = JB.COLOR_MAP.get(this.board[i]);
            let hover = this.hoverList.has(i) ? 'connected' : '';
            htmlString += `<div id="piece${i}" class="piece ${color} ${hover}"></div>`;
        }

        requestAnimationFrame(() => {
            document.getElementById('game').innerHTML = htmlString;
        });
    }

    moveEvent(e) {
        if (e.target.classList.contains('piece')) {
            const index = this.getIndexOfPiece(e.target.id);
            const status = this.move(index);
            document.getElementById('current-score').innerText = status.score;
            if (status.gameOver) {

                // Update lastScore and bestScore
                this.lastScore = this.score;
                if (this.score > this.bestScore) {
                    this.bestScore = this.score;
                }

                // Update UI
                document.getElementById('game-over-score').innerText = this.score;
                document.getElementById('last-score').innerText = this.lastScore;
                document.getElementById('best-score').innerText = this.bestScore;

                // Save scores to cookie
                this.setCookie(this.cookieName, `${this.lastScore}|${this.bestScore}`);

                // show game over overlay
                document.getElementById('game-over-overlay').classList.add('visible');
                return;
            }
            this.hoverList.clear();
            const pieces = this.getConnectedPieces(index);
            for (const i of pieces) {
                this.hoverList.add(i);
            }

            this.renderBoard();
        }
    }

    mouseOverEvent(e) {
        if (e.target.classList.contains('piece')) {
            const index = this.getIndexOfPiece(e.target.id);
            if (this.hoverList.has(index)) {
                return;
            }

            this.hoverList.clear();
            const pieces = this.getConnectedPieces(index);
            for (const i of pieces) {
                this.hoverList.add(i);
            }

            this.renderBoard();
        }
    }

    mouseLeaveEvent() {
        this.renderBoard();
    }

    getIndexOfPiece(piece) {
        return parseInt(piece.substring(5), 10);
    }

    /**
     * Make a move by removing connected pieces
     * @param {number} index - Index of the clicked piece
     * @returns {Object} - Status object with board, score, and gameOver
     */
    move(index) {
        const n = this.floodFill(index);
        if (n < 2) {
            return {
                board: this.board,
                score: this.score,
                gameOver: this.isGameOver()
            };
        }

        this.applyGravityAndShiftRight();
        this.score += this.calculateMoveScore(n);
        const gameOver = this.isGameOver();

        if (gameOver) {
            let remainingPieces = 0;
            for (let i = 0; i < this.board.length; i++) {
                if (this.board[i] !== JB.WHITE) {
                    remainingPieces++;
                }
            }
            this.score += this.calculateRemainingPiecesScore(remainingPieces);
        }

        return {
            board: this.board,
            score: this.score,
            gameOver: gameOver
        };
    }

    /**
     * Check if the game is over
     * @returns {boolean} - True if game is over, false otherwise
     */
    isGameOver() {
        for (let i = 0; i < this.board.length; i++) {
            if (this.board[i] === JB.WHITE) {
                continue;
            }
            if (this.getConnectedPieces(i).length > 1) {
                return false;
            }
        }
        return true;
    }

    /**
     * Calculate score for removing pieces
     * @param {number} piecesRemoved - Number of pieces removed
     * @returns {number} - Score for this move
     */
    calculateMoveScore(piecesRemoved) {
        if (piecesRemoved < 2) {
            return 0;
        }
        return piecesRemoved * (piecesRemoved - 1);
    }

    /**
     * Calculate bonus score for remaining pieces at game end
     * @param {number} remainingPieces - Number of pieces remaining
     * @returns {number} - Bonus score
     */
    calculateRemainingPiecesScore(remainingPieces) {
        const threshold = 20;
        if (remainingPieces <= threshold) {
            const p = threshold - remainingPieces;
            return p * p;
        }
        return 0;
    }


    /**
     * Get connected pieces of the same color
     * @param {number} index - Starting index
     * @returns {number[]} - Array of connected piece indices
     */
    getConnectedPieces(index) {
        if (index < 0 || index >= this.board.length) {
            return [];
        }

        const target = this.board[index];
        if (target === JB.WHITE) {
            return [];
        }

        const connectedIndices = [];
        const stack = [index];
        const visited = new Array(this.rows * this.cols).fill(false);

        while (stack.length > 0) {
            const i = stack.pop();

            if (i < 0 || i >= this.board.length || this.board[i] !== target || visited[i]) {
                continue;
            }

            visited[i] = true;
            connectedIndices.push(i);

            const row = Math.floor(i / this.cols);
            const col = i % this.cols;

            // Check all four adjacent positions (up, down, left, right)
            if (row > 0) {
                stack.push(i - this.cols); // Up
            }
            if (row < this.rows - 1) {
                stack.push(i + this.cols); // Down
            }
            if (col > 0) {
                stack.push(i - 1); // Left
            }
            if (col < this.cols - 1) {
                stack.push(i + 1); // Right
            }
        }

        if (connectedIndices.length < 2) {
            return [];
        }

        return connectedIndices;
    }

    /**
     * Perform flood fill to remove connected pieces
     * @param {number} index - Starting index
     * @returns {number} - Number of pieces removed
     */
    floodFill(index) {
        const connectedPieces = this.getConnectedPieces(index);
        if (connectedPieces.length < 2) {
            return 0;
        }

        // Convert board to array for easier manipulation
        const boardArray = this.board.split('');

        for (const i of connectedPieces) {
            boardArray[i] = JB.WHITE;
        }

        this.board = boardArray.join('');
        return connectedPieces.length;
    }

    /**
     * Apply gravity and shift columns right
     */
    applyGravityAndShiftRight() {
        // Convert board to 2D array for easier manipulation
        let boardArray = [];
        for (let row = 0; row < this.rows; row++) {
            const rowArray = [];
            for (let col = 0; col < this.cols; col++) {
                rowArray.push(this.board[row * this.cols + col]);
            }
            boardArray.push(rowArray);
        }

        // Gravity Phase: shift non-'w' characters down in each column
        for (let col = 0; col < this.cols; col++) {
            let writeRow = this.rows - 1;
            for (let row = this.rows - 1; row >= 0; row--) {
                if (boardArray[row][col] !== JB.WHITE) {
                    boardArray[writeRow][col] = boardArray[row][col];
                    if (writeRow !== row) {
                        boardArray[row][col] = JB.WHITE;
                    }
                    writeRow--;
                }
            }
        }

        // Right Shift Phase: move non-empty columns to the right
        let writeCol = this.cols - 1;
        for (let col = this.cols - 1; col >= 0; col--) {
            let isEmpty = true;
            for (let row = 0; row < this.rows; row++) {
                if (boardArray[row][col] !== JB.WHITE) {
                    isEmpty = false;
                    break;
                }
            }

            if (!isEmpty) {
                if (writeCol !== col) {
                    // Copy column to new position
                    for (let row = 0; row < this.rows; row++) {
                        boardArray[row][writeCol] = boardArray[row][col];
                        boardArray[row][col] = JB.WHITE;
                    }
                }
                writeCol--;
            }
        }

        // Convert 2D array back to string
        let newBoard = '';
        for (let row = 0; row < this.rows; row++) {
            for (let col = 0; col < this.cols; col++) {
                newBoard += boardArray[row][col];
            }
        }
        this.board = newBoard;
    }

    // Cookie functions
    getCookie(name) {
        const cookies = document.cookie.split(';');

        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(name + '=')) {
                return cookie.substring(name.length + 1);
            }
        }
        return '';
    }

    setCookie(name, value, days = 365) {
        const date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        const expires = "expires=" + date.toUTCString();
        document.cookie = name + "=" + value + ";" + expires + ";path=/";
    }
}
