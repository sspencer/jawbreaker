// Piece constants
const TOKEN_SPACE = 1000;
const WHITE = 0;
const PURPLE = 1000;
const BLUE = 2000;
const GREEN = 3000;
const RED = 4000;
const YELLOW = 5000;
const GRAY = 6000;
const BLACK = 7000;

const GAP = 1;
const BORDER_WIDTH = 8;

// standard game pieces
const GAME_PIECES = [PURPLE, BLUE, GREEN, RED, YELLOW];

const COLOR_MAP = new Map([
    [PURPLE, "#8a2be2"],
    [BLUE, "#00a0ff"],
    [GREEN, "#00cc66"],
    [RED, "#ff3333"],
    [YELLOW, "#ffcc00"],
    [GRAY, "#888888"],
    [BLACK, "#333333"],
    [WHITE, "#ffffff"],
]);

const POWER_X = 1;
const POWER_PLUS = 2;
const POWER_CIRCLE = 3;
const POWER_RECT = 4;
const POWER_FILL = 5;
const POWER_ROTATE_RIGHT = 6;
const POWER_ROTATE_LEFT = 7;

const POWER_PIECES = [
    POWER_X,
    POWER_PLUS,
    POWER_CIRCLE,
    POWER_RECT,
];

const EXTRA_PIECES = [POWER_FILL, POWER_ROTATE_RIGHT, POWER_ROTATE_LEFT];

const SHAPE_STROKE_COLOR = "white";
const SHAPE_BORDER_COLOR = "black";
const SHAPE_LINE_WIDTH = 3; // Adjusted for potentially more complex shapes
const SHAPE_BORDER_LINE_WIDTH = 1.5; // Adjusted for potentially more complex shapes

const POWER_UP_X_DIRECTIONS = [
    {sr: -1, sc: -1, dr: -1, dc: -1}, // Top-left
    {sr: -1, sc: 1, dr: -1, dc: 1}, // Top-right
    {sr: 1, sc: -1, dr: 1, dc: -1}, // Bottom-left
    {sr: 1, sc: 1, dr: 1, dc: 1}, // Bottom-right
];

const POWER_UP_PLUS_DIRECTIONS = [
    {sr: -1, sc: 0, dr: -1, dc: 0}, // Up
    {sr: 1, sc: 0, dr: 1, dc: 0}, // Down
    {sr: 0, sc: -1, dr: 0, dc: -1}, // Left
    {sr: 0, sc: 1, dr: 0, dc: 1}, // Right
];

const POWER_UP_RECT_DIRECTIONS = [
    {sr: -2, sc: -2, dc: 0, dr: 1},
    {sr: 2, sc: -2, dc: 1, dr: 0},
    {sr: 2, sc: 2, dc: 0, dr: -1},
    {sr: -2, sc: 2, dc: -1, dr: 0},
];

const POWER_UP_CIRCLE_DIRECTIONS = [
    {sr: -1, sc: -3, dc: 0, dr: 1},
    {sr: 3, sc: -1, dc: 1, dr: 0},
    {sr: 1, sc: 3, dc: 0, dr: -1},
    {sr: -3, sc: 1, dc: -1, dr: 0},
];

// --- Global Game State ---
let gameState = {
    size: 0,
    blockSize: 0,
    cookieName: "",
    score: 0,
    lastScore: 0,
    bestScore: 0,
    board: [],
    hoverList: new Set(),
    hoverIndex: -1,
    canvas: null,
    ctx: null,
};

function clamp(value, min, max) {
    if (value === undefined || !Number.isInteger(value)) {
        return min;
    }

    return Math.min(Math.max(value, min), max);
}

function getCookie(name) {
    const cookies = document.cookie.split(";");

    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(name + "=")) {
            return cookie.substring(name.length + 1);
        }
    }
    return "";
}

function setCookie(name, value, days = 365) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = "expires=" + date.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

function lightenColor(color, percent) {
    const num = parseInt(color.replace("#", ""), 16);
    const amt = Math.round(2.55 * percent);
    const R = (num >> 16) + amt;
    const G = (num >> 8 & 0x00FF) + amt;
    const B = (num & 0x0000FF) + amt;
    return rgbToString(R, G, B);
}

function darkenColor(color, percent) {
    const num = parseInt(color.replace("#", ""), 16);
    const amt = Math.round(2.55 * percent);
    const R = (num >> 16) - amt;
    const G = (num >> 8 & 0x00FF) - amt;
    const B = (num & 0x0000FF) - amt;
    return rgbToString(R, G, B);
}

function rgbToString(R, G, B) {
    return "#" + (
        0x1000000 +
        (R < 255 ? R < 1 ? 0 : R : 255) * 0x10000 +
        (G < 255 ? G < 1 ? 0 : G : 255) * 0x100 +
        (B < 255 ? B < 1 ? 0 : B : 255)
    ).toString(16).slice(1);
}

// --- Game Logic Functions ---
function generateUniqueRandomNums(n, max) {
    if (n > max) { // n cannot be greater than the number of available spots (0 to max)
        console.warn(`Attempting to generate ${n} unique numbers from a range of ${max + 1}. Clamping to ${max + 1}.`);
        n = max + 1;
    }
    const randomNumbers = new Set();
    while (randomNumbers.size < n) {
        const randomNum = Math.floor(Math.random() * (max + 1));
        randomNumbers.add(randomNum);
    }
    return Array.from(randomNumbers);
}

function createNewBoard() {
    const boardLen = gameState.size * gameState.size;
    let board = [];
    for (let i = 0; i < boardLen; i++) {
        let color = GAME_PIECES[Math.floor(Math.random() * GAME_PIECES.length)];
        board.push(color);
    }

    const colors = [...GAME_PIECES, GRAY]; // Colors that power-ups can be on
    const ppLen = POWER_PIECES.length;
    const exLen = EXTRA_PIECES.length;
    // Max number of unique spots needed for power-ups.
    // Each of ppLen can be on any of 'colors' length. Each of exLen is typically on GRAY.
    const numPowerUpSlotsToGenerate = Math.min(boardLen, (ppLen * colors.length) + exLen);


    const uniqIndices = generateUniqueRandomNums(numPowerUpSlotsToGenerate, boardLen - 1);

    // Place EXTRA_PIECES (typically on GRAY)
    for (let i = 0; i < EXTRA_PIECES.length; i++) {
        if (uniqIndices.length === 0) break;
        board[uniqIndices.shift()] = GRAY + EXTRA_PIECES[i];
    }

    // Place standard POWER_PIECES on various colors
    for (const color of colors) {
        for (const power of POWER_PIECES) {
            if (uniqIndices.length === 0) break;
            board[uniqIndices.shift()] = color + power;
        }
        if (uniqIndices.length === 0) break;
    }
    gameState.board = board;
    return board;
}

function fillEmptySpacesOnBoard() {
    for (let i = 0; i < gameState.board.length; i++) {
        if (gameState.board[i] === WHITE) {
            gameState.board[i] = GAME_PIECES[Math.floor(Math.random() * GAME_PIECES.length)];
        }
    }
}

function getCanvasCoordinates(e) {
    const rect = gameState.canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    return {x, y};
}

function getBoardIndexFromCoordinates(x, y) {
    const adjustedX = x - BORDER_WIDTH;
    const adjustedY = y - BORDER_WIDTH;
    const col = Math.floor(adjustedX / (gameState.blockSize + GAP));
    const row = Math.floor(adjustedY / (gameState.blockSize + GAP));

    if (col < 0 || col >= gameState.size || row < 0 || row >= gameState.size) {
        return -1;
    }
    return row * gameState.size + col;
}

function getPointFromIndex(index) {
    const col = index % gameState.size;
    const row = Math.floor(index / gameState.size);
    return {col, row};
}

function getPieceColorCode(piece) {
    return COLOR_MAP.get(getBasePiece(piece));
}

function getBasePiece(piece) {
    return Math.floor(piece / TOKEN_SPACE) * TOKEN_SPACE;
}

function getPowerUpType(piece) {
    return piece % TOKEN_SPACE;
}

function getPiecePositionOnCanvas(row, col) {
    const x = BORDER_WIDTH + col * (gameState.blockSize + GAP);
    const y = BORDER_WIDTH + row * (gameState.blockSize + GAP);
    return {x, y};
}

// --- Drawing Functions ---
function drawEmptySpace(x, y, pieceSize, outlineGap, outlineSize) {
    const ctx = gameState.ctx;
    ctx.save();
    ctx.strokeStyle = darkenColor(COLOR_MAP.get(WHITE), 70);
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.rect(x + outlineGap, y + outlineGap, outlineSize, outlineSize);
    ctx.stroke();
    ctx.restore();
}

function drawColoredPiece(x, y, pieceSize, piece, cornerRadius, glow) {
    const ctx = gameState.ctx;
    ctx.save();

    const gradient = ctx.createLinearGradient(
        x + gameState.blockSize * 0.7, y + gameState.blockSize * 0.7,
        x + gameState.blockSize * 0.2, y + gameState.blockSize * 0.2,
    );

    let color = getPieceColorCode(piece);
    if (EXTRA_PIECES.includes(getPowerUpType(piece))) {
        color = darkenColor(color, 30); // Darken background for extra power-ups
    }
    gradient.addColorStop(0, lightenColor(color, 12));
    gradient.addColorStop(1, darkenColor(color, 6));
    ctx.fillStyle = gradient;

    ctx.beginPath();
    ctx.roundRect(x, y, pieceSize, pieceSize, cornerRadius);
    ctx.fill();

    // Highlights and shadows for 3D effect
    ctx.beginPath();
    ctx.rect(x + cornerRadius, y, pieceSize - (2 * cornerRadius), glow);
    ctx.rect(x, y + cornerRadius, glow, pieceSize - (2 * cornerRadius));
    ctx.fillStyle = lightenColor(color, 20);
    ctx.fill();

    ctx.beginPath();
    ctx.rect(x + pieceSize - glow, y + cornerRadius, glow, pieceSize - (2 * cornerRadius));
    ctx.rect(x + cornerRadius, y + pieceSize - glow, pieceSize - (2 * cornerRadius), glow);
    ctx.fillStyle = darkenColor(color, 20);
    ctx.fill();

    ctx.restore();
}

function drawHoverEffect(x, y, pieceSize, color, cornerRadius, bigger) {
    const ctx = gameState.ctx;
    ctx.save();
    ctx.shadowColor = color;
    ctx.shadowBlur = 15;
    ctx.strokeStyle = lightenColor(color, 40);
    ctx.lineWidth = bigger ? 8 : 4; // Thicker highlight for "bigger" (extra power-ups)

    let hoverSize = pieceSize + (bigger ? 4 : 2);
    let hoverX = x - (bigger ? 2 : 1);
    let hoverY = y - (bigger ? 2 : 1);
    if (bigger) {
        ctx.strokeStyle = "rgba(255, 255, 255, 0.7)"; // Brighter for extra
    }

    ctx.beginPath();
    ctx.roundRect(hoverX, hoverY, hoverSize, hoverSize, cornerRadius + 1);
    ctx.stroke();
    ctx.restore();
}

function prepareShapeContext() {
    const ctx = gameState.ctx;
    ctx.save();
    ctx.beginPath();
    return ctx;
}

function finalizeShapeDraw(ctx) {
    // Draw border first
    ctx.strokeStyle = SHAPE_BORDER_COLOR;
    ctx.lineWidth = SHAPE_LINE_WIDTH + SHAPE_BORDER_LINE_WIDTH * 2; // Ensure border is outside main line
    ctx.stroke();
    // Draw main shape line
    ctx.strokeStyle = SHAPE_STROKE_COLOR;
    ctx.lineWidth = SHAPE_LINE_WIDTH;
    ctx.stroke();
    ctx.restore();
}


function drawX(x, y, size) {
    const ctx = prepareShapeContext();
    const padding = size * 0.2; // Adjusted padding
    ctx.moveTo(x + padding, y + padding);
    ctx.lineTo(x + size - padding, y + size - padding);
    ctx.moveTo(x + size - padding, y + padding);
    ctx.lineTo(x + padding, y + size - padding);
    finalizeShapeDraw(ctx);
}

function drawPlus(x, y, size) {
    const ctx = prepareShapeContext();
    const padding = size * 0.2; // Adjusted padding
    ctx.moveTo(x + size / 2, y + padding);
    ctx.lineTo(x + size / 2, y + size - padding);
    ctx.moveTo(x + padding, y + size / 2);
    ctx.lineTo(x + size - padding, y + size / 2);
    finalizeShapeDraw(ctx);
}

function drawCircle(x, y, size) {
    const ctx = prepareShapeContext();
    const radius = size * 0.3; // Adjusted radius
    ctx.arc(x + size / 2, y + size / 2, radius, 0, Math.PI * 2);
    finalizeShapeDraw(ctx);
}

function drawRect(x, y, size) {
    const ctx = prepareShapeContext();
    const padding = size * 0.2; // Adjusted padding
    const rectInnerSize = size - 2 * padding;
    ctx.rect(x + padding, y + padding, rectInnerSize, rectInnerSize);
    finalizeShapeDraw(ctx);
}

function drawPowerFill(x, y, size) {
    const ctx = gameState.ctx; // Direct context for multi-fill
    ctx.save();
    const halfSize = size / 2;
    const centerSize = size * 0.3;
    const centerOffset = (size - centerSize) / 2;

    // Define colors for quadrants
    const quadColors = [BLUE, GREEN, RED, YELLOW];
    const positions = [
        {qx: x, qy: y},                           // Top-left
        {qx: x + halfSize, qy: y},                // Top-right
        {qx: x, qy: y + halfSize},                // Bottom-left
        {qx: x + halfSize, qy: y + halfSize}      // Bottom-right
    ];

    for (let i = 0; i < 4; i++) {
        ctx.beginPath();
        ctx.rect(positions[i].qx, positions[i].qy, halfSize, halfSize);
        ctx.fillStyle = COLOR_MAP.get(quadColors[i]);
        ctx.fill();
    }
    // Center piece
    ctx.beginPath();
    ctx.rect(x + centerOffset, y + centerOffset, centerSize, centerSize);
    ctx.fillStyle = COLOR_MAP.get(PURPLE);
    ctx.fill();
    ctx.restore();
}

function drawRotateRight(x, y, size) {
    const ctx = prepareShapeContext();
    const cx = x + size / 2;
    const cy = y + size / 2;
    const r = size * 0.3; // Radius of the semi-circle
    const arrowLength = size * 0.15; // Length of the arrowhead sides
    const arrowWidth = size * 0.1;  // Half-width of the arrowhead base

    // Draw the top semi-circular arc (clockwise)
    ctx.arc(cx, cy, r, Math.PI, 2 * Math.PI, false); // Start left (PI), end right (2PI or 0)

    // Arrowhead at the right end of the arc (pointing right)
    const tipX = cx + r;
    const tipY = cy;
    ctx.moveTo(tipX - arrowLength, tipY - arrowWidth);
    ctx.lineTo(tipX, tipY);
    ctx.lineTo(tipX - arrowLength, tipY + arrowWidth);

    finalizeShapeDraw(ctx);
}

function drawRotateLeft(x, y, size) {
    const ctx = prepareShapeContext();
    const cx = x + size / 2;
    const cy = y + size / 2;
    const r = size * 0.3; // Radius
    const arrowLength = size * 0.15;
    const arrowWidth = size * 0.1;

    // Draw the top semi-circular arc (counter-clockwise)
    ctx.arc(cx, cy, r, 0, Math.PI, true); // Start right (0), end left (PI)

    // Arrowhead at the left end of the arc (pointing left)
    const tipX = cx - r;
    const tipY = cy;
    ctx.moveTo(tipX + arrowLength, tipY - arrowWidth);
    ctx.lineTo(tipX, tipY);
    ctx.lineTo(tipX + arrowLength, tipY + arrowWidth);

    finalizeShapeDraw(ctx);
}

function drawPowerUp(x, y, size, powerUpType) {
    const paddingFactor = 8; // Padding from the edge of the piece to the power-up symbol
    const shapeSize = size - paddingFactor;
    const shapeX = x + (paddingFactor / 2);
    const shapeY = y + (paddingFactor / 2);

    switch (powerUpType) {
        case POWER_X:
            drawX(shapeX, shapeY, shapeSize);
            break;
        case POWER_PLUS:
            drawPlus(shapeX, shapeY, shapeSize);
            break;
        case POWER_RECT:
            drawRect(shapeX, shapeY, shapeSize);
            break;
        case POWER_FILL:
            drawPowerFill(shapeX, shapeY, shapeSize);
            break;
        case POWER_ROTATE_RIGHT:
            drawRotateRight(shapeX, shapeY, shapeSize);
            break;
        case POWER_ROTATE_LEFT:
            drawRotateLeft(shapeX, shapeY, shapeSize);
            break;
        case POWER_CIRCLE: // Fallthrough if POWER_CIRCLE is still desired, or handle explicitly
        default:
            drawCircle(shapeX, shapeY, shapeSize);
            break; // Default to circle
    }
}

function renderBoard() {
    gameState.ctx.clearRect(0, 0, gameState.canvas.width, gameState.canvas.height);
    const cornerRadius = 3; // Slightly more rounded
    const pieceSize = gameState.blockSize - GAP;
    const outlineGap = 2;
    const outlineSize = pieceSize - (2 * outlineGap);
    const glowEffectSize = 1.5; // Slightly larger glow

    // First pass: Draw all pieces
    for (let row = 0; row < gameState.size; row++) {
        for (let col = 0; col < gameState.size; col++) {
            const index = row * gameState.size + col;
            const piece = gameState.board[index];
            const {x, y} = getPiecePositionOnCanvas(row, col);

            if (piece === WHITE) {
                drawEmptySpace(x, y, pieceSize, outlineGap, outlineSize);
            } else {
                drawColoredPiece(x, y, pieceSize, piece, cornerRadius, glowEffectSize);
            }
        }
    }

    // Second pass: Draw power-ups and hover effects (to ensure they are on top)
    for (let row = 0; row < gameState.size; row++) {
        for (let col = 0; col < gameState.size; col++) {
            const index = row * gameState.size + col;
            const piece = gameState.board[index];

            if (piece !== WHITE) {
                const {x, y} = getPiecePositionOnCanvas(row, col);
                const powerUpType = getPowerUpType(piece);
                if (powerUpType > 0) {
                    drawPowerUp(x, y, pieceSize, powerUpType);
                }
                if (gameState.hoverList.has(index)) {
                    const color = getPieceColorCode(piece);
                    const isExtra = EXTRA_PIECES.includes(powerUpType);
                    drawHoverEffect(x, y, pieceSize, color, cornerRadius, isExtra);
                }
            }
        }
    }
}

// --- Connection Logic ---
function getConnectionsForPowerUp(index, targetColor, powerUpType) {
    // This function determines the area of effect for non-EXTRA power-ups like X, Plus.
    // It returns an array of indices affected by the power-up.
    switch (powerUpType) {
        case POWER_X:
            return getXConnections(index, targetColor);
        case POWER_PLUS:
            return getPlusConnections(index, targetColor);
        case POWER_RECT:
            return getRectConnections(index, targetColor);
        case POWER_CIRCLE:
            return getCircularConnections(index, targetColor);
        default:
            return [];
    }
}

function getConnectionsWithDirections(index, targetColor, directions, maxIterations) {
    const point = getPointFromIndex(index);
    let startRow = point.row;
    let startCol = point.col;
    let connectedIndices = [];

    for (let iter = 0; iter < maxIterations; iter++) {
        for (const dir of directions) {
            const r = startRow + dir.sr + dir.dr * iter;
            const c = startCol + dir.sc + dir.dc * iter;

            if (r >= 0 && r < gameState.size && c >= 0 && c < gameState.size) {
                const currentIndex = r * gameState.size + c;
                // Power-ups affect pieces of their base color, or if GRAY, any non-WHITE piece.
                const currentPieceBase = getBasePiece(gameState.board[currentIndex]);
                if (gameState.board[currentIndex] !== WHITE &&
                    (currentPieceBase === targetColor || targetColor === GRAY)) {
                    if (!connectedIndices.includes(currentIndex)) {
                        connectedIndices.push(currentIndex);
                    }
                }
            }
        }
    }
    return connectedIndices;
}

function getXConnections(index, targetColor) {
    const maxIters = gameState.size; // Iterate up to board size
    return getConnectionsWithDirections(index, targetColor, POWER_UP_X_DIRECTIONS, maxIters);
}

function getPlusConnections(index, targetColor) {
    const maxIters = gameState.size;
    return getConnectionsWithDirections(index, targetColor, POWER_UP_PLUS_DIRECTIONS, maxIters);
}

function getRectConnections(index, targetColor) {
    // Rect/Circle directions are offsets, so maxIterations is 1 for their direct application.
    // The 'dr'/'dc' in their definitions are for multi-step paths if needed, but here they are direct.
    return getConnectionsWithDirections(index, targetColor, POWER_UP_RECT_DIRECTIONS, 1);
}

function getCircularConnections(index, targetColor) {
    let c1 = getConnectionsWithDirections(index, targetColor, POWER_UP_RECT_DIRECTIONS, 1); // Inner ring
    let c2 = getConnectionsWithDirections(index, targetColor, POWER_UP_CIRCLE_DIRECTIONS, 1); // Outer ring
    return [...new Set([...c1, ...c2])]; // Combine and remove duplicates
}

function getConnectedPieces(index) {
    // This function is crucial. It determines which pieces are connected for removal or hover.
    if (index < 0 || index >= gameState.board.length || gameState.board[index] === WHITE) return [];

    const clickedPiece = gameState.board[index];
    const baseColorOfClickedPiece = getBasePiece(clickedPiece);
    const powerUpType = getPowerUpType(clickedPiece);

    // 1. Handle EXTRA_PIECES (Fill, Rotations) - for hover, they usually highlight themselves.
    // Their actual "connection" for removal is handled by their specific logic.
    if (EXTRA_PIECES.includes(powerUpType)) {
        return [index]; // For hover, highlight the power-up itself.
    }

    // 2. Handle other POWER_PIECES (X, Plus, Circle, Rect)
    if (POWER_PIECES.includes(powerUpType)) {
        let affectedIndices = getConnectionsForPowerUp(index, baseColorOfClickedPiece, powerUpType);
        // Always include the power-up piece itself in the list of pieces to be affected/removed.
        if (!affectedIndices.includes(index)) {
            affectedIndices.unshift(index);
        }
        return affectedIndices.filter(i => i >= 0 && i < gameState.board.length && gameState.board[i] !== WHITE);
    }

    // 3. Standard Flood Fill for same-colored pieces (no power-up)
    const connectedIndices = [];
    const stack = [index];
    const visited = new Array(gameState.board.length).fill(false);
    visited[index] = true;

    while (stack.length > 0) {
        const currentIndex = stack.pop();
        connectedIndices.push(currentIndex);

        const {row, col} = getPointFromIndex(currentIndex);
        const neighbors = [
            (row > 0) ? currentIndex - gameState.size : -1,             // Up
            (row < gameState.size - 1) ? currentIndex + gameState.size : -1, // Down
            (col > 0) ? currentIndex - 1 : -1,                          // Left
            (col < gameState.size - 1) ? currentIndex + 1 : -1,         // Right
        ];

        for (const neighborIndex of neighbors) {
            if (neighborIndex !== -1 && !visited[neighborIndex] &&
                gameState.board[neighborIndex] !== WHITE &&
                getBasePiece(gameState.board[neighborIndex]) === baseColorOfClickedPiece &&
                getPowerUpType(gameState.board[neighborIndex]) === 0) { // Only connect to normal pieces of same color
                visited[neighborIndex] = true;
                stack.push(neighborIndex);
            }
        }
    }
    // For normal pieces, only return if 2 or more are connected.
    return connectedIndices.length >= 2 ? connectedIndices : [];
}

// --- Game State Manipulation ---
function boardTo2DArray() {
    let boardArray = [];
    for (let r = 0; r < gameState.size; r++) {
        boardArray.push(gameState.board.slice(r * gameState.size, (r + 1) * gameState.size));
    }
    return boardArray;
}

function updateBoardFrom2DArray(boardArray) {
    gameState.board = boardArray.flat();
}

function updateCanvasDimensions() {
    const canvasSize = gameState.size * gameState.blockSize +
        (gameState.size - 1) * GAP + 2 * BORDER_WIDTH;
    gameState.canvas.width = canvasSize;
    gameState.canvas.height = canvasSize;
}

function rotateBoardClockwise() { // Rotate Right
    let boardArray = boardTo2DArray();
    let rotatedArray = Array.from({length: gameState.size}, () => Array(gameState.size).fill(WHITE));
    for (let r = 0; r < gameState.size; r++) {
        for (let c = 0; c < gameState.size; c++) {
            rotatedArray[c][gameState.size - 1 - r] = boardArray[r][c];
        }
    }
    updateBoardFrom2DArray(rotatedArray);
    applyGravityAndShiftColumns(); // Apply gravity and shift after rotation
    // renderBoard(); // renderBoard will be called by processMove or click handler
}

function rotateBoardCounterClockwise() { // Rotate Left
    let boardArray = boardTo2DArray();
    let rotatedArray = Array.from({length: gameState.size}, () => Array(gameState.size).fill(WHITE));
    for (let r = 0; r < gameState.size; r++) {
        for (let c = 0; c < gameState.size; c++) {
            rotatedArray[gameState.size - 1 - c][r] = boardArray[r][c];
        }
    }
    updateBoardFrom2DArray(rotatedArray);
    applyGravityAndShiftColumns(); // Apply gravity and shift after rotation
    // renderBoard();
}

function applyGravityAndShiftColumns() {
    let boardArray = boardTo2DArray();

    // Apply gravity (pieces fall down in each column)
    for (let c = 0; c < gameState.size; c++) {
        let writeRow = gameState.size - 1;
        for (let r = gameState.size - 1; r >= 0; r--) {
            if (boardArray[r][c] !== WHITE) {
                if (writeRow !== r) {
                    boardArray[writeRow][c] = boardArray[r][c];
                    boardArray[r][c] = WHITE;
                }
                writeRow--;
            }
        }
    }

    // Shift columns to the left if a column becomes empty
    let writeCol = 0;
    for (let c = 0; c < gameState.size; c++) {
        let isColumnEmpty = true;
        for (let r = 0; r < gameState.size; r++) {
            if (boardArray[r][c] !== WHITE) {
                isColumnEmpty = false;
                break;
            }
        }

        if (!isColumnEmpty) {
            if (writeCol !== c) {
                for (let r = 0; r < gameState.size; r++) {
                    boardArray[r][writeCol] = boardArray[r][c];
                    boardArray[r][c] = WHITE;
                }
            }
            writeCol++;
        }
    }
    updateBoardFrom2DArray(boardArray);
}

function removeTargetedPieces(index) {
    // This function handles the removal of pieces based on the clicked piece (normal or power-up)
    const clickedPieceOriginal = gameState.board[index]; // Store before modification
    const powerUpType = getPowerUpType(clickedPieceOriginal);

    let piecesToRemove = [];
    let piecesRemovedCount = 0;

    if (powerUpType === POWER_FILL) {
        gameState.board[index] = WHITE; // Remove the fill piece itself
        fillEmptySpacesOnBoard(); // Fill action
        // Gravity is applied by processMove after this
        return {count: 1, isSpecialAction: true}; // Special action, count is nominal
    } else if (powerUpType === POWER_ROTATE_RIGHT) {
        gameState.board[index] = WHITE;
        rotateBoardClockwise();
        return {count: 1, isSpecialAction: true};
    } else if (powerUpType === POWER_ROTATE_LEFT) {
        gameState.board[index] = WHITE;
        rotateBoardCounterClockwise();
        return {count: 1, isSpecialAction: true};
    }

    // For standard pieces or non-EXTRA power-ups
    piecesToRemove = getConnectedPieces(index); // getConnectedPieces now correctly identifies targets

    if (piecesToRemove.length === 0 && POWER_PIECES.includes(powerUpType)) {
        // If a non-EXTRA power-up is clicked and finds no connections, it removes itself.
        piecesToRemove = [index];
    }


    for (const i of piecesToRemove) {
        if (gameState.board[i] !== WHITE) {
            gameState.board[i] = WHITE;
            piecesRemovedCount++;
        }
    }
    return {count: piecesRemovedCount, isSpecialAction: false};
}

function calculateMoveScore(piecesRemoved) {
    if (piecesRemoved < 2) return 0;
    // Standard Jawbreaker scoring: n * (n - 1)
    return piecesRemoved * (piecesRemoved - 1);
}

function calculateRemainingPiecesScore(remainingPieces) {
    const threshold = 10; // Bonus if 10 or fewer pieces remain
    if (remainingPieces === 0) return 2000; // Bonus for clearing the board
    if (remainingPieces <= threshold) {
        return (threshold - remainingPieces + 1) * 100; // Scaled bonus
    }
    return 0; // Penalty or smaller bonus for more remaining pieces can be added here
}

function isGameOver() {
    for (let i = 0; i < gameState.board.length; i++) {
        if (gameState.board[i] === WHITE) continue;

        const powerUpType = getPowerUpType(gameState.board[i]);
        if (EXTRA_PIECES.includes(powerUpType)) return false; // Fill/Rotations are always playable

        const connections = getConnectedPieces(i); // Get potential connections
        if (connections.length > 0) return false; // If any piece can make a valid move
    }
    return true;
}

function processMove(index) {
    const removalResult = removeTargetedPieces(index);
    const n = removalResult.count;
    const piece = gameState.board[index];
    const powerUpType = getPowerUpType(piece);

    if (n > 0 && !removalResult.isSpecialAction) { // Apply gravity only if pieces were removed by non-special actions
        applyGravityAndShiftColumns();
    }
    // For special actions like Fill, gravity is handled within their functions or here if needed.
    if (powerUpType === POWER_FILL && n > 0) { // POWER_FILL already called fill, now apply gravity
        applyGravityAndShiftColumns();
    }


    if (!removalResult.isSpecialAction) { // Score only for non-special actions based on count
        gameState.score += calculateMoveScore(n);
    }
    // else: Special actions (Fill, Rotations) might have a flat score or no direct score from removal count.
    // For example: gameState.score += 10; // for using a rotation power-up.

    const gameOver = isGameOver();
    if (gameOver) {
        let remainingPieces = 0;
        gameState.board.forEach(p => {
            if (p !== WHITE) remainingPieces++;
        });
        gameState.score += calculateRemainingPiecesScore(remainingPieces);
    }

    return {
        score: gameState.score,
        gameOver: gameOver,
    };
}

// --- Event Handlers ---
function handleClick(e) {
    const {x, y} = getCanvasCoordinates(e);
    const index = getBoardIndexFromCoordinates(x, y);

    if (index >= 0 && index < gameState.board.length && gameState.board[index] !== WHITE) {
        const status = processMove(index);
        document.getElementById("current-score").innerText = status.score;

        if (status.gameOver) {
            gameState.lastScore = gameState.score; // Update last score before potential best score update
            if (gameState.score > gameState.bestScore) {
                gameState.bestScore = gameState.score;
            }
            document.getElementById("game-over-score").innerText = gameState.score;
            document.getElementById("last-score").innerText = gameState.lastScore;
            document.getElementById("best-score").innerText = gameState.bestScore;
            setCookie(gameState.cookieName, `${gameState.lastScore}|${gameState.bestScore}`);
            document.getElementById("game-over-overlay").classList.add("visible");
            return; // Stop further processing/rendering
        }

        gameState.hoverList.clear(); // Clear hover after a click
        gameState.hoverIndex = -1;
        renderBoard(); // Re-render the board after the move
    }
}

function handleMouseMove(e) {
    const {x, y} = getCanvasCoordinates(e);
    const index = getBoardIndexFromCoordinates(x, y);

    if (index === gameState.hoverIndex) return; // No change if hovering over the same piece

    gameState.hoverList.clear();
    gameState.hoverIndex = index; // Update current hover index

    if (index >= 0 && index < gameState.board.length && gameState.board[index] !== WHITE) {
        const piecesToHighlight = getConnectedPieces(index);
        if (piecesToHighlight.length > 0) {
            piecesToHighlight.forEach(i => gameState.hoverList.add(i));
        } else {
            // If getConnectedPieces returns empty (e.g. single non-powerup piece),
            // still highlight the piece under the mouse if it's not white.
            gameState.hoverList.add(index);
        }
    }
    // Always re-render on mouse move to update hover effect or clear it
    renderBoard();
}

function handleMouseLeave() {
    if (gameState.hoverList.size > 0) { // Only re-render if there was a hover to clear
        gameState.hoverList.clear();
        gameState.hoverIndex = -1;
        renderBoard();
    }
}

function resetCurrentGame() {
    createNewBoard(); // This now sets gameState.board
    gameState.hoverList.clear();
    gameState.hoverIndex = -1;
    gameState.score = 0;
    document.getElementById("current-score").innerText = gameState.score;
    document.getElementById("game-over-overlay").classList.remove("visible");
    renderBoard();
}

function registerGameEvents() {
    if (!gameState.canvas) return;
    gameState.canvas.addEventListener("click", handleClick);
    gameState.canvas.addEventListener("mousemove", handleMouseMove);
    gameState.canvas.addEventListener("mouseleave", handleMouseLeave);

    const newGameBtn = document.querySelector(".new-game-btn");
    if (newGameBtn) {
        newGameBtn.addEventListener("click", (e) => {
            e.preventDefault();
            resetCurrentGame();
        });
    }
    const gameOverRestartBtn = document.getElementById("restart-btn");
    if (gameOverRestartBtn) {
        gameOverRestartBtn.addEventListener("click", (e) => {
            e.preventDefault();
            resetCurrentGame();
        });
    }

    // Help Modal Logic
    const helpBtn = document.getElementById("help-btn");
    const helpModal = document.getElementById("help-modal");
    const closeBtn = helpModal ? helpModal.querySelector(".close") : null;

    if (helpBtn && helpModal && closeBtn) {
        helpBtn.addEventListener("click", () => {
            helpModal.style.display = "block";
        });
        closeBtn.addEventListener("click", () => {
            helpModal.style.display = "none";
        });
        window.addEventListener("click", (event) => {
            if (event.target === helpModal) {
                helpModal.style.display = "none";
            }
        });
    } else {
        console.warn("Help modal elements not found. Ensure #help-btn, #help-modal, and .close exist.");
    }
}

// --- Initialization ---
function initGame(opts) {
    gameState.size = clamp(opts.size, 8, 20); // Max size 20 for better playability
    gameState.blockSize = opts.blockSize || 36;
    gameState.cookieName = opts.cookieName || "jawbreaker_functional_scores_v3"; // Unique cookie name
    gameState.score = 0;

    gameState.canvas = document.getElementById("game-canvas");
    if (!gameState.canvas) {
        console.error("Canvas element with ID 'game-canvas' not found. Game cannot start.");
        return;
    }
    gameState.ctx = gameState.canvas.getContext("2d");

    updateCanvasDimensions(); // Set canvas size based on game size and block size
    createNewBoard(); // Initialize the board

    // Load scores from cookie
    const scoresCookie = getCookie(gameState.cookieName);
    if (scoresCookie) {
        const [last, best] = scoresCookie.split('|');
        gameState.lastScore = parseInt(last, 10) || 0;
        gameState.bestScore = parseInt(best, 10) || 0;
    } else {
        gameState.lastScore = 0;
        gameState.bestScore = 0;
    }
    // Update score displays
    document.getElementById('current-score').innerText = gameState.score;
    document.getElementById('last-score').innerText = gameState.lastScore;
    document.getElementById('best-score').innerText = gameState.bestScore;


    registerGameEvents();
    renderBoard(); // Initial render of the game board
}

