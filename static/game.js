// Wait for the DOM to be fully loaded
document.addEventListener("DOMContentLoaded", function () {
    // Get all pieces in the grid
    const pieces = document.querySelectorAll(".piece");

    // Add style to handle both connected and non-connected pieces
    const styleElement = document.createElement("style");
    styleElement.textContent = `
    /* Remove the default hover styles from all pieces */
    .piece:hover {
      transform: none !important;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2) !important;
      border-color: rgba(255, 255, 255, 0.3) !important;
      opacity: 0.9 !important;
    }
    
    /* Apply hover effect only to pieces with connected-hover class */
    .piece.connected-hover {
      transform: scale(1.15) !important;
      box-shadow: 0 12px 20px rgba(0, 0, 0, 0.4) !important;
      border-color: rgba(255, 255, 255, 0.5) !important;
      opacity: 1 !important;
      z-index: 10;
      transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) !important;
    }
  `;

    document.head.appendChild(styleElement);

    // Add event listeners to each piece
    pieces.forEach((piece) => {
        piece.addEventListener("mouseenter", handleMouseEnter);
        piece.addEventListener("mouseleave", handleMouseLeave);
    });

    /**
     * Handles the mouseenter event on a piece
     * @param {Event} event - The mouseenter event
     */
    function handleMouseEnter(event) {
        const targetPiece = event.target;
        const connectedPieces = findConnectedPieces(targetPiece);

        // Only apply the hover effect if there are at least 2 connected pieces
        if (connectedPieces.length >= 2) {
            // Apply hover effect to all connected pieces including the hovered one
            connectedPieces.forEach((piece) => {
                document.getElementById(piece).classList.add("connected-hover");
            });
        }
    }

    /**
     * Handles the mouseleave event on a piece
     * @param {Event} event - The mouseleave event
     */
    function handleMouseLeave(event) {
        // Get the element being entered after leaving this one
        const relatedTarget = event.relatedTarget;

        // If we're not entering another piece with the connected-hover class,
        // remove the effect from all pieces
        if (
            !relatedTarget ||
            !relatedTarget.classList.contains("connected-hover")
        ) {
            document.querySelectorAll(".connected-hover").forEach((piece) => {
                piece.classList.remove("connected-hover");
            });
        }
    }

    /**
     * Finds all connected pieces of the same color
     * @param {HTMLElement} startPiece - The piece to start the search from
     * @returns {Array} - Array of all connected piece elements
     */
    function findConnectedPieces(startPiece) {
        const pieceId = startPiece.id;
        const match = pieceId.match(/piece(\d+)/);
        if (!match) return [];

        const index = parseInt(match[1], 10);
        const board = window.gamePieces;
        const width = 12;
        const height = 12;
        const totalCells = width * height;

        if (board === null || board.length !== totalCells || index < 0 || index >= totalCells) {
            return [];
        }

        const color = board[index];
        if (color === "w") return [];

        const visited = new Set();
        const stack = [index];
        const connected = [];

        const directions = [-1, 1, -width, width]; // left, right, up, down

        while (stack.length > 0) {
            const current = stack.pop();

            if (visited.has(current)) continue;
            visited.add(current);

            if (board[current] !== color) continue;

            connected.push(current);

            for (const dir of directions) {
                const neighbor = current + dir;

                if (neighbor < 0 || neighbor >= totalCells) continue;

                // Prevent wrapping around left/right edges
                if (dir === -1 && current % width === 0) continue;
                if (dir === 1 && current % width === width - 1) continue;

                if (!visited.has(neighbor)) {
                    stack.push(neighbor);
                }
            }
        }

        return connected.map((i) => "piece" + i);
    }
});
