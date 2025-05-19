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
const BORDER_WIDTH = 4;

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

const GRADIENT_MAP = new Map([
    [PURPLE, ["#d442ff", "#8a2be2"]],
    [BLUE, ["#00f0ff", "#0070ff"]],
    [GREEN,["#00ff99", "#00cc66"]],
    [RED, ["#ff5500", "#dd0000"]],
    [YELLOW, ["#ffee00", "#ff8800"]],
    [GRAY, ["#777777", "#cccccc"]],
]);

const POWER_X = 1;
const POWER_PLUS = 2;
const POWER_CIRCLE = 3;
const POWER_RECT = 4;
const POWER_FILL = 5;
const POWER_ROTATE_RIGHT = 6;
const POWER_ROTATE_LEFT = 7;
const POWER_BUBBLE = 8;
const POWER_DISC = 9; // filled circle

const POWER_PIECES = [
    POWER_X,
    POWER_PLUS,
    POWER_CIRCLE,
    POWER_RECT,
];

const CONNECTION_PIECES = [...POWER_PIECES, POWER_DISC];

const EXTRA_PIECES = [
    POWER_FILL,
    POWER_ROTATE_RIGHT,
    POWER_ROTATE_LEFT,
];

const SPECIAL_PIECES= [ POWER_DISC ];


const SHAPE_STROKE_COLOR = "white";
const SHAPE_BORDER_COLOR = "black";
const SHAPE_LINE_WIDTH = 3; // Adjusted for potentially more complex shapes
const SHAPE_BORDER_LINE_WIDTH = 1; // Adjusted for potentially more complex shapes

const POWER_UP_X_DIRECTIONS = [
    {sr: -1, sc: -1, dr: -1, dc: -1}, // Top-left
    {sr: -1, sc:  1, dr: -1, dc:  1}, // Top-right
    {sr:  1, sc: -1, dr:  1, dc: -1}, // Bottom-left
    {sr:  1, sc:  1, dr:  1, dc:  1}, // Bottom-right
];

const POWER_UP_PLUS_DIRECTIONS = [
    {sr: -1, sc:  0, dr: -1, dc:  0}, // Up
    {sr:  1, sc:  0, dr:  1, dc:  0}, // Down
    {sr:  0, sc: -1, dr:  0, dc: -1}, // Left
    {sr:  0, sc:  1, dr:  0, dc:  1}, // Right
];

const POWER_UP_RECT_DIRECTIONS = [
    {sr: -2, sc: -2, dc:  0, dr:  1},
    {sr:  2, sc: -2, dc:  1, dr:  0},
    {sr:  2, sc:  2, dc:  0, dr: -1},
    {sr: -2, sc:  2, dc: -1, dr:  0},
];

const POWER_UP_CIRCLE_DIRECTIONS = [
    {sr: -1, sc: -3, dc:  0, dr:  1},
    {sr:  3, sc: -1, dc:  1, dr:  0},
    {sr:  1, sc:  3, dc:  0, dr: -1},
    {sr: -3, sc:  1, dc: -1, dr: 0},
];
const POWER_UP_DISC1_DIRECTIONS = [
    {sr: -1, sc: -2, dc:  0, dr:  1},
    {sr:  2, sc: -1, dc:  1, dr:  0},
    {sr:  1, sc:  2, dc:  0, dr: -1},
    {sr: -2, sc:  1, dc: -1, dr: 0},
];

const POWER_UP_DISC2_DIRECTIONS = [
    {sr: -1, sc: -1, dc:  0, dr:  1},
    {sr:  1, sc: -1, dc:  1, dr:  0},
    {sr:  1, sc:  1, dc:  0, dr: -1},
    {sr: -1, sc:  1, dc: -1, dr: 0},
];

// --- Global Game State ---
let gameState = {
    rows: 0,
    cols: 0,
    blockSize: 0,
    cookieName: "",
    score: 0,
    bonus: 0,
    remainingPieces: 0,
    lastScore: 0,
    bestScore: 0,
    board: [],
    animateStart: [],
    animateEnd: [],
    animateIndices: [],
    animatePosition: 0,
    animatePieces: 8, // draw (n) animated pieces per frame
    animateId: null,
    animateResolve: null, // Promise resolve function for animation
    hoverList: new Set(),
    hoverIndex: -1,
    mobile: false,
    canvas: null,
    ctx: null,
    undo: null,
    weightedFills: true,
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

/**
 * Creates a new game board filled with random pieces and power-ups
 * @returns {Array} The newly created board
 */
function createNewBoard() {
    const size = gameState.rows * gameState.cols;

    // Start with an empty board
    let board = Array(size).fill(WHITE);

    // Fill with random game pieces
    for (let i = 0; i < size; i++) {
        board[i] = randomItem(GAME_PIECES);
    }

    const colors = [...GAME_PIECES, GRAY]; // Colors that power-ups can be on
    const ppLen = POWER_PIECES.length;
    const exLen = EXTRA_PIECES.length;
    const spLen = SPECIAL_PIECES.length;
    // Max number of unique spots needed for power-ups.
    // Each of ppLen can be on any of 'colors' length. Each of exLen is typically on GRAY.
    const numPowerUpSlotsToGenerate = Math.min(size, (ppLen * colors.length) + exLen + spLen);

    const uniqIndices = generateUniqueRandomNums(numPowerUpSlotsToGenerate, size - 1);

    for (let i = 0; i < EXTRA_PIECES.length; i++) {
        if (uniqIndices.length === 0) break;
        board[uniqIndices.shift()] = GRAY + EXTRA_PIECES[i];
    }

    for (let i = 0; i < SPECIAL_PIECES.length; i++) {
        if (uniqIndices.length === 0) break;
        board[uniqIndices.shift()] = GRAY + SPECIAL_PIECES[i];
    }

    // Place standard POWER_PIECES on various colors
    for (const color of colors) {
        for (const power of POWER_PIECES) {
            if (uniqIndices.length === 0) break;
            board[uniqIndices.shift()] = color + power;
        }
        if (uniqIndices.length === 0) break;
    }

    // Set the game board and return it
    gameState.board = board;
    return board;
}

/**
 * Sets up animation for a new board appearing
 * @returns {Promise} Resolves when animation completes
 */
async function animateNewBoard() {
    const size = gameState.rows * gameState.cols;

    // Set up animation from empty board to current board
    gameState.animateStart =  Array(size).fill(WHITE);
    gameState.animateEnd = gameState.board.slice();

    // For new board animation, we want to animate all indices in random order
    gameState.animateIndices = createIndexArray(size);
    shuffleArray(gameState.animateIndices);

    // Run the animation
    return renderBoardAnimated(8);
}

function randomItem(arr) {
    return arr[Math.floor(Math.random() * arr.length)];
}

function fillSpaces(connections) {
    if (connections === undefined) {
        connections = [];
        for (let i = 0; i < gameState.board.length; i++) {
            if (gameState.board[i] === WHITE) {
                connections.push(i);
            }
        }
    }

    if (gameState.weightedFills) {
        // Use an object to count occurrences
        const counts = {};
        for (const item of gameState.board) {
            const piece = getPieceBase(item);
            if (piece === WHITE || piece === GRAY) continue;
            counts[piece] = (counts[piece] || 0) + 1;
        }

        const choices = Object.entries(counts).map(([value, count]) => ({value, count}));

        for (let c in connections) {
            gameState.board[connections[c]] = weightedRandom(choices); // randomItem(GAME_PIECES);
        }
    } else {
        for (let c in connections) {
            gameState.board[connections[c]] = randomItem(GAME_PIECES);
        }

    }
}

function weightedRandom(choices) {
    // Sum up all counts
    const total = choices.reduce((sum, obj) => sum + obj.count, 0);
    // Get a random number in [0, total)
    let r = Math.random() * total;
    // Walk through array, subtracting counts, find where it lands
    for (let i = 0; i < choices.length; i++) {
        if (r < choices[i].count) {
            return choices[i].value;
        }
        r -= choices[i].count;
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

    if (col < 0 || col >= gameState.cols || row < 0 || row >= gameState.cols) {
        return -1;
    }
    return row * gameState.cols + col;
}

function getPointFromIndex(index) {
    const col = index % gameState.cols;
    const row = Math.floor(index / gameState.cols);
    return {col, row};
}

function getPieceColor(piece) {
    return COLOR_MAP.get(getPieceBase(piece));
}

function getPieceGradient(piece) {
    return GRADIENT_MAP.get(getPieceBase(piece));
}

function getPieceBase(piece) {
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


function createIndexArray(n) {
    return Array.from({ length: n }, (_, index) => index);
}

// Fisher-Yates shuffle algorithm
function shuffleArray(array) {
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
}

function findChangedIndices(array1, array2) {
    if (array1.length !== array2.length) {
        throw new Error("Arrays must be of the same length.");
    }

    const changed = [];

    for (let i = 0; i < array1.length; i++) {
        if (array1[i] !== array2[i]) {
            changed.push(i);
        }
    }

    return changed;
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

    const size = gameState.blockSize;

    const angle = 225 * Math.PI / 180;
    const diagonal = Math.sqrt(size * size + size * size);

    const startX = x + size / 2 + (Math.cos(angle) * diagonal) / 2;
    const startY = y + size / 2 + (Math.sin(angle) * diagonal) / 2;
    const endX = x + size / 2 - (Math.cos(angle) * diagonal) / 2;
    const endY = y + size / 2 - (Math.sin(angle) * diagonal) / 2;

    const gradBody = ctx.createLinearGradient(startX, startY, endX, endY);
    const gradHighlight = ctx.createLinearGradient(x, y, x + size, y + size);
    //const gradHighlight = ctx.createLinearGradient(startX, startY, endX, endY);
    const color = getPieceGradient(piece);

    // highlight
    ctx.beginPath();
    gradHighlight.addColorStop(0, color[1]);
    gradHighlight.addColorStop(1, lightenColor(color[0], 5));
    ctx.fillStyle = gradHighlight;
    ctx.roundRect(x, y, pieceSize, pieceSize, cornerRadius);
    ctx.fill();

    // inner fill
    ctx.beginPath();
    gradBody.addColorStop(0, darkenColor(color[0], 4));
    gradBody.addColorStop(1, lightenColor(color[1], 8));
    ctx.fillStyle = gradBody;
    ctx.roundRect(x+2, y+2, pieceSize-4, pieceSize-4, cornerRadius);
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

function finalizeShapeDraw(ctx, opts) {
    // Draw border first
    ctx.strokeStyle = SHAPE_BORDER_COLOR;
    ctx.lineWidth = SHAPE_LINE_WIDTH + SHAPE_BORDER_LINE_WIDTH * 2; // Ensure border is outside main line
    ctx.stroke();
    // Draw the main shape line
    if (opts) {
        ctx.strokeStyle = SHAPE_BORDER_COLOR;
    } else {
        ctx.strokeStyle = SHAPE_STROKE_COLOR;
    }
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

function drawDisc(x, y, size) {
    const ctx = prepareShapeContext();
    const radius = size * 0.3; // Adjust this value as needed
    const cx = x + size / 2;
    const cy = y + size / 2;

    ctx.beginPath();
    ctx.arc(cx, cy, radius, 0, Math.PI * 2);
    ctx.fillStyle = SHAPE_STROKE_COLOR
    ctx.fill();
    ctx.strokeStyle = SHAPE_BORDER_COLOR;
    ctx.lineWidth = 2;
    ctx.stroke();

    finalizeShapeDraw(ctx);
}

function drawSpike(x, y, size) {
    const ctx = prepareShapeContext();
    const padding = size * 0.2;
    const left = x + padding;
    const right = x + size - padding;
    const centerY = y + size / 2;
    const spikeHeight = size * 0.2;

    // Main horizontal line (with spike up and down at 1/3 and 2/3 points)
    ctx.beginPath();
    ctx.moveTo(left, centerY);

    // 1st segment
    const spike1x = left + (right - left) / 3;
    ctx.lineTo(spike1x, centerY);

    // Spike up
    ctx.lineTo(spike1x + (right - left) * 0.05, centerY - spikeHeight);
    ctx.lineTo(spike1x + (right - left) * 0.10, centerY);

    // 2nd segment
    const spike2x = left + 2 * (right - left) / 3;
    ctx.lineTo(spike2x, centerY);

    // Spike down
    ctx.lineTo(spike2x + (right - left) * 0.05, centerY + spikeHeight);
    ctx.lineTo(spike2x + (right - left) * 0.10, centerY);

    // Last segment
    ctx.lineTo(right, centerY);

    ctx.strokeStyle = SHAPE_STROKE_COLOR;
    ctx.lineWidth = SHAPE_LINE_WIDTH;
    ctx.stroke();

    // Border line (draw over main line for higher-line border effect)
    ctx.save();
    ctx.strokeStyle = SHAPE_BORDER_COLOR;
    ctx.lineWidth = SHAPE_BORDER_LINE_WIDTH;
    ctx.stroke();
    ctx.restore();

    finalizeShapeDraw(ctx);
}

function drawRect(x, y, size, opts) {
    const ctx = prepareShapeContext();
    const padding = size * 0.2; // Adjusted padding
    const rectInnerSize = size - 2 * padding;
    ctx.rect(x + padding, y + padding, rectInnerSize, rectInnerSize);
    finalizeShapeDraw(ctx, opts);
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
    // Centerpiece
    ctx.beginPath();
    ctx.rect(x + centerOffset, y + centerOffset, centerSize, centerSize);
    ctx.fillStyle = COLOR_MAP.get(PURPLE);
    ctx.fill();
    ctx.restore();

    drawRect(x-8, y-8, size+16, "reverse");
}

function drawPowerBubble(x, y, size) {
    const ctx = gameState.ctx;
    ctx.save();

    // Main circle for clipping and outline
    const centerX = x + size / 2;
    const centerY = y + size / 2;
    const radius = size * 0.3;

    // Clip to circle, so color doesn't bleed out
    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
    ctx.clip();

    // Draw colored pie slices inside clip region
    const quadColors = [BLUE, GREEN, RED, YELLOW];
    for (let i = 0; i < 4; i++) {
        ctx.beginPath();
        ctx.moveTo(centerX, centerY);
        ctx.arc(
            centerX,
            centerY,
            radius,
            (Math.PI / 2) * i,
            (Math.PI / 2) * (i + 1)
        );
        ctx.closePath();
        ctx.fillStyle = COLOR_MAP.get(quadColors[i]);
        ctx.fill();
    }

    // Optional highlight
    ctx.beginPath();
    ctx.arc(centerX, centerY - radius / 2, radius / 3, 0, Math.PI * 2);
    ctx.fillStyle = "rgba(255,255,255,0.28)";
    ctx.fill();

    ctx.restore();

    const f = gameState.blockSize / 9;
    drawCircle(x-f, y-f, size+f*2);

}


function drawRotateRight(x, y, size) {
    const ctx = prepareShapeContext();
    const cx = x + size / 2;
    const cy = y + size / 2;
    const r = size * 0.3;
    const arrowLength = size * 0.25; // Bigger arrowhead
    const arrowWidth = size * 0.15;

    // Draw the top semicircular arc (clockwise)
    ctx.arc(cx, cy, r, Math.PI, 2 * Math.PI, false);

    // Arrowhead at the right end of the arc (pointing right)
    const tipX = cx + r + 2;
    const tipY = cy + 3;
    ctx.moveTo(tipX - arrowLength, tipY - arrowWidth);
    ctx.lineTo(tipX, tipY);
    ctx.lineTo(tipX - arrowLength, tipY + arrowWidth);
    ctx.closePath();

    finalizeShapeDraw(ctx);
}

function drawRotateLeft(x, y, size) {
    const ctx = prepareShapeContext();
    const cx = x + size / 2;
    const cy = y + size / 2;
    const r = size * 0.3;
    const arrowLength = size * 0.25;
    const arrowWidth = size * 0.15;

    // Draw the top semicircular arc (counter-clockwise)
    ctx.arc(cx, cy, r, 0, Math.PI, true);

    // Arrowhead at the left end of the arc (pointing left)
    const tipX = cx - r - 2;
    const tipY = cy + 3;
    ctx.moveTo(tipX + arrowLength, tipY - arrowWidth);
    ctx.lineTo(tipX, tipY);
    ctx.lineTo(tipX + arrowLength, tipY + arrowWidth);
    ctx.closePath();

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
        case POWER_BUBBLE:
            drawPowerBubble(shapeX, shapeY, shapeSize);
            break;
        case POWER_ROTATE_RIGHT:
            drawRotateRight(shapeX, shapeY, shapeSize);
            break;
        case POWER_ROTATE_LEFT:
            drawRotateLeft(shapeX, shapeY, shapeSize);
            break;
        case POWER_DISC:
            drawDisc(shapeX, shapeY, shapeSize);
            break;
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
    for (let row = 0; row < gameState.rows; row++) {
        for (let col = 0; col < gameState.cols; col++) {
            const index = row * gameState.cols + col;
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
    for (let row = 0; row < gameState.rows; row++) {
        for (let col = 0; col < gameState.cols; col++) {
            const index = row * gameState.cols + col;
            const piece = gameState.board[index];

            if (piece !== WHITE) {
                const {x, y} = getPiecePositionOnCanvas(row, col);
                const powerUpType = getPowerUpType(piece);
                if (powerUpType > 0) {
                    drawPowerUp(x, y, pieceSize, powerUpType);
                }
                if (gameState.hoverList.has(index)) {
                    const color = getPieceColor(piece);
                    const isExtra = EXTRA_PIECES.includes(powerUpType);
                    drawHoverEffect(x, y, pieceSize, color, cornerRadius, isExtra);
                }
            }
        }
    }
}

async function animateTransformation(transformFn, transformOpts) {
    // Save the starting state
    const startBoard = gameState.board.slice();

    // Apply the transformation
    if (typeof transformFn === 'function') {
        transformFn(transformOpts);
    }

    // Set up animation state
    gameState.animateStart = startBoard;
    gameState.animateEnd = gameState.board.slice();
    let changes = findChangedIndices(gameState.animateStart, gameState.animateEnd);
    shuffleArray(changes);
    gameState.animateIndices = changes;
    // Run the animation
    await renderBoardAnimated(1);
    return gameState.board;
}

function renderBoardAnimated(piecesToAnimate = 8) {
    console.log("renderBoardAnimated:", piecesToAnimate);
    gameState.animatePieces = piecesToAnimate;
    gameState.animatePosition = 0;
    return new Promise((resolve) => {
        gameState.animateResolve = resolve;
        gameState.animateId = requestAnimationFrame(doRenderBoardAnimated);
    });
}

function doRenderBoardAnimated(timestamp) {
    gameState.board = gameState.animateStart.slice();

    const newPos = gameState.animatePosition + gameState.animatePieces;
    const maxPos = gameState.animateIndices.length;
    const endPos = Math.min(newPos, maxPos);

    for (let iter = 0; iter < endPos; iter++) {
        const index =gameState.animateIndices[iter];
        gameState.board[index] = gameState.animateEnd[index];
    }
    renderBoard();

    gameState.animateStart = gameState.board.slice();
    gameState.animatePosition = endPos;

    if (gameState.animatePosition < maxPos) {
        gameState.animateId = requestAnimationFrame(doRenderBoardAnimated)
        return;
    }

    gameState.animateStart = [];
    gameState.animateEnd = [];
    gameState.animatePosition = 0;
    gameState.animateId = null;

    // Resolve the Promise if it exists
    if (gameState.animateResolve) {
        gameState.animateResolve();
        gameState.animateResolve = null;
    }
}

function getConnectedPieces(index) {
    // This function is crucial. It determines which pieces are connected for removal or hover.
    if (index < 0 || index >= gameState.board.length || gameState.board[index] === WHITE) return [];

    const clickedPiece = gameState.board[index];
    const baseColorOfClickedPiece = getPieceBase(clickedPiece);
    const powerUpType = getPowerUpType(clickedPiece);

    // 1. Handle EXTRA_PIECES (Fill, Rotations) - for hover, they usually highlight themselves.
    // Their actual "connection" for removal is handled by their specific logic.
    if (EXTRA_PIECES.includes(powerUpType)) {
        return [index]; // For hover, highlight the power-up itself.
    }

    let powerUpIndices = [];
    // 2. Handle other POWER_PIECES (X, Plus, Circle, Rect)

    if (CONNECTION_PIECES.includes(powerUpType)) {
        powerUpIndices = getConnectionsForPowerUp(index, baseColorOfClickedPiece, powerUpType);
        if (powerUpIndices.length > 0) {
            powerUpIndices.unshift(index);
        }
        //return affectedIndices.filter(i => i >= 0 && i < gameState.board.length && gameState.board[i] !== WHITE);
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
            (row > 0) ? currentIndex - gameState.cols : -1,             // Up
            (row < gameState.cols - 1) ? currentIndex + gameState.cols : -1, // Down
            (col > 0) ? currentIndex - 1 : -1,                          // Left
            (col < gameState.cols - 1) ? currentIndex + 1 : -1,         // Right
        ];

        for (const neighborIndex of neighbors) {
            if (neighborIndex !== -1 && !visited[neighborIndex] &&
                gameState.board[neighborIndex] !== WHITE &&
                getPieceBase(gameState.board[neighborIndex]) === baseColorOfClickedPiece) {
                visited[neighborIndex] = true;
                stack.push(neighborIndex);
            }
        }
    }

    // combineAndRemoveDuplicates
    const combinedArray = [...powerUpIndices, ...connectedIndices];
    const uniqueArray = [...new Set(combinedArray)];

    // For normal pieces, only return if 2 or more are connected.
    return uniqueArray.length >= 2 ? uniqueArray : [];
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
        case POWER_DISC:
            return getDiscConnections(index, targetColor);
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

            if (r >= 0 && r < gameState.rows && c >= 0 && c < gameState.cols) {
                const currentIndex = r * gameState.cols + c;
                // Power ups affect pieces of their base color, or if GRAY, any non-WHITE piece.
                const currentPieceBase = getPieceBase(gameState.board[currentIndex]);
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
    const maxIters = Math.ceil(Math.sqrt(Math.pow(gameState.rows, 2) + Math.pow(gameState.cols, 2)));
    return getConnectionsWithDirections(index, targetColor, POWER_UP_X_DIRECTIONS, maxIters);
}

function getPlusConnections(index, targetColor) {
    const maxIters = Math.max(gameState.rows, gameState.cols);
    return getConnectionsWithDirections(index, targetColor, POWER_UP_PLUS_DIRECTIONS, maxIters);
}

function getRectConnections(index, targetColor) {
    return getConnectionsWithDirections(index, targetColor, POWER_UP_RECT_DIRECTIONS, 4);
}

function getCircularConnections(index, targetColor) {
     const c1 = getConnectionsWithDirections(index, targetColor, POWER_UP_RECT_DIRECTIONS, 1); // Inner ring
     const c2 = getConnectionsWithDirections(index, targetColor, POWER_UP_CIRCLE_DIRECTIONS, 3); // Outer ring
     return [...c1, ...c2];
}

function getDiscConnections(index, targetColor) {
    const c1 = getConnectionsWithDirections(index, targetColor, POWER_UP_DISC1_DIRECTIONS, 3);
    const c2 = getConnectionsWithDirections(index, targetColor, POWER_UP_DISC2_DIRECTIONS, 2);
    return [...c1, ...c2];
}


// --- Game State Manipulation ---
function boardTo2DArray() {
    let boardArray = [];
    for (let r = 0; r < gameState.rows; r++) {
        boardArray.push(gameState.board.slice(r * gameState.cols, (r + 1) * gameState.cols));
    }
    return boardArray;
}

function updateBoardFrom2DArray(boardArray) {
    gameState.board = boardArray.flat();
}

function updateCanvasDimensions() {
    gameState.canvas.width = gameState.cols * gameState.blockSize + (gameState.cols - 1) * GAP + 2 * BORDER_WIDTH
    gameState.canvas.height = gameState.rows * gameState.blockSize + (gameState.rows - 1) * GAP + 2 * BORDER_WIDTH;
}

function rotateBoard(degrees) {
    const arr = gameState.board;
    const n = Math.sqrt(arr.length);
    if (!Number.isInteger(n))
        throw new Error("Array length must be a perfect square.");

    // Normalize degrees to 0, 90, 180, or 270
    let rotation = ((degrees % 360) + 360) % 360;

    // Convert to 2D row-major grid
    const grid = Array.from({ length: n }, (_, row) =>
        arr.slice(row * n, row * n + n)
    );

    let rotated = Array.from({ length: n }, () => Array(n));

    switch (rotation) {
        case 90:
            for (let row = 0; row < n; row++) {
                for (let col = 0; col < n; col++) {
                    rotated[col][n - 1 - row] = grid[row][col];
                }
            }
            break;
        case 270:
            for (let row = 0; row < n; row++) {
                for (let col = 0; col < n; col++) {
                    rotated[n - 1 - col][row] = grid[row][col];
                }
            }
            break;
        default:
            throw new Error("Rotation must be a multiple of 90 degrees.");
    }

    // Flatten 2D grid back to 1D row-major array
    gameState.board = rotated.flat();
}

function applyGravityAndShiftColumns() {
    let boardArray = boardTo2DArray();

    // Apply gravity (pieces fall down in each column)
    for (let c = 0; c < gameState.cols; c++) {
        let writeRow = gameState.rows - 1;
        for (let r = gameState.rows - 1; r >= 0; r--) {
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
    let writeCol = gameState.cols - 1;
    for (let c = gameState.cols - 1; c >= 0; c--) {
        let isEmpty = true;
        for (let r = 0; r < gameState.rows; r++) {
            if (boardArray[r][c] !== WHITE) {
                isEmpty = false;
                break;
            }
        }

        if (!isEmpty) {
            if (writeCol !== c) {
                for (let row = 0; row < gameState.rows; row++) {
                    boardArray[row][writeCol] = boardArray[row][c];
                    boardArray[row][c] = WHITE;
                }
            }
            writeCol--;
        }
    }

    updateBoardFrom2DArray(boardArray);
}

async function removeTargetedPieces(index) {
    // This function handles the removal of pieces based on the clicked piece (normal or power-up)
    const clickedPieceOriginal = gameState.board[index]; // Store before modification
    const powerUpType = getPowerUpType(clickedPieceOriginal);

    let piecesToRemove = [];
    let piecesRemovedCount = 0;

    if (powerUpType === POWER_FILL) {
        gameState.board[index] = WHITE; // randomItem(GAME_PIECES); // Remove the fill piece itself
        applyGravityAndShiftColumns();
        // Animate filling empty spaces
        await animateTransformation(fillSpaces);
        return {count: 1, isSpecialAction: true}; // Special action, count is nominal
    } else if (powerUpType === POWER_ROTATE_RIGHT) {
        gameState.board[index] = WHITE; //randomItem(GAME_PIECES);
        rotateBoard(90);
        applyGravityAndShiftColumns();
        return {count: 1, isSpecialAction: true};
    } else if (powerUpType === POWER_ROTATE_LEFT) {
        gameState.board[index] = WHITE; // randomItem(GAME_PIECES);
        rotateBoard(-90);
        applyGravityAndShiftColumns();
        return {count: 1, isSpecialAction: true};
    } /*else if (powerUpType === POWER_BUBBLE) {
        getConnectionsForPowerUp(index, getPieceBase(gameState.board[index]), POWER_BUBBLE);
        applyGravityAndShiftColumns();
        await animateTransformation(fillSpaces, getConnectedPieces(index));
        return {count: 1, isSpecialAction: true};
    }*/

    // For standard pieces or non-EXTRA power-ups
    piecesToRemove = getConnectedPieces(index);

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

async function processMove(index) {
    // Save the starting board state for animation
    const startBoard = gameState.board.slice();

    // Process the piece removal
    const removalResult = await removeTargetedPieces(index);
    const n = removalResult.count;

    if (n > 0 && !removalResult.isSpecialAction) { // Apply gravity only if pieces were removed by non-special actions
        // First animate the piece removal
        gameState.animateStart = startBoard;
        gameState.animateEnd = gameState.board.slice();
        gameState.animateIndices = findChangedIndices(gameState.animateStart, gameState.animateEnd, index);

        applyGravityAndShiftColumns();
        renderBoard();
    }

    if (!removalResult.isSpecialAction) { // Score only for non-special actions based on count
        gameState.score += calculateMoveScore(n);
    }

    const gameOver = isGameOver();
    if (gameOver) {
        let remainingPieces = 0;
        gameState.board.forEach(p => {
            if (p !== WHITE) remainingPieces++;
        });
        gameState.bonus = calculateRemainingPiecesScore(remainingPieces);
        gameState.remainingPieces = remainingPieces;
        gameState.score += gameState.bonus;
    }

    return {
        score: gameState.score,
        gameOver: gameOver,
    };
}

function undoMove() {
    if (gameState.undo === null) return;
    const {board, score} = gameState.undo;
    gameState.board = board.slice();
    gameState.score = score;
    gameState.undo = null;
    document.getElementById('undo-btn').disabled = true;
    document.getElementById("current-score").innerText = ""+score;

    renderBoard();
}

function handleTouch(e) {
    if (gameState.animateId !== null) {
        return;
    }

    const {x, y} = getCanvasCoordinates(e);
    const index = getBoardIndexFromCoordinates(x, y);
    if (index === gameState.hoverIndex) {
        handleClick(e).then(r => {});
        return
    }

    handleMouseMove(e);
}


// --- Event Handlers ---
async function handleClick(e) {
    if (gameState.animateId !== null) {
        return;
    }

    const {x, y} = getCanvasCoordinates(e);
    const index = getBoardIndexFromCoordinates(x, y);

    if (index >= 0 && index < gameState.board.length && gameState.board[index] !== WHITE) {


        gameState.undo = {
            board: gameState.board.slice(),
            score: gameState.score,
        };
        document.getElementById('undo-btn').disabled = false;


        const status = await processMove(index);
        document.getElementById("current-score").innerText = ""+status.score;

        if (status.gameOver) {
            gameState.undo = null;
            gameState.lastScore = gameState.score; // Update last score before potential best score update
            if (gameState.score > gameState.bestScore) {
                gameState.bestScore = gameState.score;
            }

            document.getElementById("game-over-score").innerText = gameState.score;
            document.getElementById("last-score").innerText = gameState.lastScore;
            document.getElementById("best-score").innerText = gameState.bestScore;
            document.getElementById("bonus-points").innerText = ""+gameState.bonus;
            document.getElementById("remaining-pieces").innerText = ""+gameState.remainingPieces;
            document.getElementById("pieces-text").innerText =  gameState.remainingPieces === 1 ? "piece" : "pieces";
            document.getElementById("game-over-overlay").classList.add("visible");

            setCookie(gameState.cookieName, `${gameState.lastScore}|${gameState.bestScore}`);
            return; // Stop further processing/rendering
        }

        gameState.hoverList.clear(); // Clear hover after a click
        gameState.hoverIndex = -1;
        if (gameState.mobile) {
            const pieces = getConnectedPieces(index);
            if (pieces.length > 1) {
                gameState.hoverIndex = index;
                for (const i of pieces) {
                    gameState.hoverList.add(i);
                }
            }
        }

        renderBoard(); // Re-render the board after the move
    }
}

function handleMouseMove(e) {
    if (gameState.animateId !== null) {
        return;
    }

    const {x, y} = getCanvasCoordinates(e);
    const index = getBoardIndexFromCoordinates(x, y);

    if (index === gameState.hoverIndex) return; // No change if hovering over the same piece

    gameState.hoverList.clear();
    gameState.hoverIndex = index;

    if (index >= 0 && index < gameState.board.length && gameState.board[index] !== WHITE) {
        const piecesToHighlight = getConnectedPieces(index);
        if (piecesToHighlight.length > 0) {
            piecesToHighlight.forEach(i => gameState.hoverList.add(i));
        }
    }
    // Always rerender on mouse move to update the hover effect or clear it
    renderBoard();
}

function handleMouseLeave() {
    if (gameState.animateId !== null) {
        return;
    }

    if (gameState.hoverList.size > 0) { // Only re-render if there was a hover to clear
        gameState.hoverList.clear();
        gameState.hoverIndex = -1;
        renderBoard();
    }
}

/**
 * Resets the game with a new board and clears game state
 */
async function resetCurrentGame() {
    // Create a new board
    createNewBoard();

    // Reset game state
    gameState.hoverList.clear();
    gameState.hoverIndex = -1;
    gameState.score = 0;
    gameState.bonus = 0;
    gameState.remainingPieces = 0;
    gameState.undo = null;
    document.getElementById('undo-btn').disabled = true;
    document.getElementById("current-score").innerText = gameState.score;
    document.getElementById("game-over-overlay").classList.remove("visible");

    // Animate the new board appearing
    await animateNewBoard();
}

function registerGameEvents() {
    if (!gameState.canvas) return;
    if (gameState.mobile) {
        gameState.canvas.addEventListener("click", handleTouch);
    } else {
        gameState.canvas.addEventListener("click", async (e) => {
            await handleClick(e);
        });
        gameState.canvas.addEventListener("mousemove", handleMouseMove);
        gameState.canvas.addEventListener("mouseleave", handleMouseLeave);
    }

    const newGameBtn = document.getElementById("new-game-btn");
    if (newGameBtn) {
        newGameBtn.addEventListener("click", async (e) => {
            e.preventDefault();
            await resetCurrentGame();
        });
    }
    const gameOverRestartBtn = document.getElementById("restart-btn");
    if (gameOverRestartBtn) {
        gameOverRestartBtn.addEventListener("click", async (e) => {
            e.preventDefault();
            await resetCurrentGame();
        });
    }

    const undoBtn = document.getElementById("undo-btn");
    if (undoBtn) {
        undoBtn.disabled = true;
        undoBtn.addEventListener("click", async (e) => {
            console.log("undo button pressed");
            e.preventDefault();
            undoMove()
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
/**
 * Initializes the game with the given options
 * @param {Object} opts - Game options
 * @param {number} opts.size - Board size (number of cells per side)
 * @param {number} opts.blockSize - Size of each cell in pixels
 * @param {string} opts.cookieName - Name of the cookie to store scores
 * @param {boolean} opts.mobile - Whether the game is running on a mobile device
 */
async function initGame(opts) {
    // Initialize game state
    gameState.rows = clamp(opts.rows, 8, 20); // Max size 20 for better playability
    gameState.cols = clamp(opts.cols, 8, 20); // Max size 20 for better playability
    gameState.blockSize = opts.blockSize || 36;
    gameState.cookieName = opts.cookieName || "jawbreaker_functional_scores_v3"; // Unique cookie name
    gameState.mobile = opts.mobile || false;
    gameState.score = 0;

    // Set up canvas
    gameState.canvas = document.getElementById("game-canvas");
    if (!gameState.canvas) {
        console.error("Canvas element with ID 'game-canvas' not found. Game cannot start.");
        return;
    }
    gameState.ctx = gameState.canvas.getContext("2d");
    updateCanvasDimensions(); // Set canvas size based on game size and block size

    // Create the initial game board
    createNewBoard();

    // Load scores from the cookie
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

    // Register event handlers
    registerGameEvents();

    // Animate the initial board appearance
    await animateNewBoard();
}
