/*
 * PROJECT:  Jawbreaker
 * VERSION:  1.0
 * LICENSE:  BSD (revised)
 * AUTHOR:   (c) 2005 Steven Spencer
 * LINK:     http://www.bigfrog.net/jawbreaker/
 *
 * This script can be used freely as long as all
 * copyright messages are intact.
 *
 * The html layout, css and javascript "style"  are derived
 * from JsTetris.  The copyright for JsTetris appears below:
 */

/*
 * PROJECT:  JsTetris
 * VERSION:  1.1.0
 * LICENSE:  BSD (revised)
 * AUTHOR:   (c) 2004 Cezary Tomczak
 * LINK:     http://gosu.pl/dhtml/JsTetris.html
 *
 * This script can be used freely as long as all
 * copyright messages are intact.
 */

Array.prototype.copy = function()
{
	var arr = new Array(this.length);
	for (var i = 0; i < this.length; i++) {
		arr[i] = this[i];
	}

	return arr;
}

Array.prototype.equals = function(arr)
{
	if (arr == null) return false;

	if (this.length != arr.length) return false;

	for (var i = 0; i < this.length; i++) {
		if (this[i] != arr[i]) return false;
	}

	return true;
}


function JawBreaker() {
    var self = this;
    
    this.board          = null;  // 2 dimensional array of colors (red, green, etc)
    this.undoBoard      = null;  // 2 dimensional array of colors (red, green, etc) for last move
    this.undoScore      = 0;     // so we can get the last score
    this.undoColor      = null;
    this.undoPieces     = 0;
    this.areaW          = 11;    // area width
    this.areaH          = 12;    // area height
    this.imageH         = 20;    // image height
    this.imageW         = 20;    // image width
    this.numColors      = 5;
    this.colors         = Array("white", "red", "blue", "green", "yellow", "purple");
    this.bonus          = Array(200, 150, 120, 100, 80, 60, 40, 20, 10);
    this.gameScore      = 0;
    this.gamesPlayed    = 0;
    this.bestScore      = 0;
    this.totalScore     = 0;
    this.selectedBlocks = null;
    this.uuid = null;

    this.pcount         = Array();

	// constructor
    this.start = function() {

		if (self.board == null)
		{
			// init game array
			self.board     = new Array(self.areaH);
			self.undoBoard = new Array(self.areaH);

			// create rows/cols in html table
            var tbody = $('boardbody');

			var x, y;
			for (y = 0; y < self.areaH; y++)
			{
				self.board[y]     = new Array(self.areaW);
				self.undoBoard[y] = new Array(self.areaW);
                var row = document.createElement('TR');

				for (x = 0; x < self.areaW; x++)
				{
					var cell = document.createElement('TD');
					var link = document.createElement('A');
					var img = document.createElement('IMG');

					link.onclick = function(e) { self.mouseclick(e); }

					img.width = self.imageW;
					img.height = self.imageH;
					img.id = self.cid(y, x);

					link.appendChild(img);
					cell.appendChild(link);
                    row.appendChild(cell);
                }
                
                tbody.appendChild(row);
            }

            self.saveScore();
		}

        self.pcount["purple"] = 0;    
        self.pcount["blue"] = 0;    
        self.pcount["green"] = 0;    
        self.pcount["red"] = 0;    
        self.pcount["yellow"] = 0;    

        var color;
		for (y = 0; y < self.areaH; y++)
		{
			for (x = 0; x < self.areaW; x++)
			{
                color = self.rnd(self.numColors);
				self.board[y][x] = color;
                self.pcount[color] += 1;
				self.draw(y, x);
			}
		}

		self.gameScore = 0;
		self.score(0);
		self.disable('menu-undo');

        // x=1 to skip white
        var max = 0;
        for (x = 1; x < self.colors.length; x++) {
            color = self.colors[x];
            if (self.pcount[color] > max) {
                max = self.pcount[color];
            }
        }

        for (x = 1; x < self.colors.length; x++) {
            color = self.colors[x];
            if (self.pcount[color] == max) {
		if (max > 39) 
  	              $(color).innerHTML = "<b style='color: red'>" + self.pcount[color] + "</b>";
		else 
	              $(color).innerHTML = "<b>" + self.pcount[color] + "</b>";
            } else
                $(color).innerHTML = self.pcount[color];
        }
 
    }

	this.disable = function(id) {
		$(id).disabled = true;
		$(id).style.color = '#999';
	}
	this.enable = function(id) {
		$(id).disabled = false;
		$(id).style.color = '#000';
	}
		
	// ajax call to store the user's score, retrieve updated stats
    this.saveScore = function(score) {
        var url = "score.php";
        var p = "uuid=" + self.uuid + "&board=" + self.areaH + "x" + self.areaW;
        if (score) p += "&score=" + score;

        new Ajax.Request(url, {method: 'get', parameters: p, onComplete: self.showScore});
    }

	// called on ajax return, update the scores
    this.showScore = function(req) {
        var arr = req.responseText.split("|");
        for(var i =0; i < arr.length; i+=2) {
            try {
                $(arr[i]).innerHTML = arr[i+1];

				// score updated both on left menu and game stats
                if (arr[i] == "usergames" || arr[i] == "useraverage" || arr[i] == "userbest" ) {
                    $(arr[i]+"1").innerHTML = arr[i+1];
                }
            } catch(e) {}
        }
    }
    
	// recursively compute the selected blocks
	this.getNeighbors = function(list, y, x, c) {

		// position is valid?
		if (y < 0 || y >= self.areaH || x < 0 && x >= self.areaW) {
			return list;
		}

		// position is correct color
		if (self.board[y][x] != c || self.board[y][x] == "white") {
			return list;
		}

		// position hasn't been seen
		for (var i = 0; i < list.length; i++) {
			if (list[i] == self.cid(y, x)) {
				return list;
			}
		}

		// save new position
		list.push(self.cid(y, x));

		// recursively count neighbors
		self.getNeighbors(list, y-1, x, c);
		self.getNeighbors(list, y+1, x, c);
		self.getNeighbors(list, y, x-1, c);
		self.getNeighbors(list, y, x+1, c);

		return list;
	}

	// user clicked on a block
	this.mouseclick = function(e) {
		var el;
		if (!e) var e = window.event;
		if (e.target) el = e.target;
		else if (e.srcElement) el = e.srcElement;
		if (el.nodeType == 3) // defeat Safari bug
		el = el.parentNode;

		var id = el.id;
		if (id.match(/r\d+c\d+/)) {
			var rc = self.rc(id);
			list = self.getNeighbors(new Array(), rc[0], rc[1], self.board[rc[0]][rc[1]]);
			list.sort();

			if (list.equals(self.selectedBlocks)) {
                var color = self.board[rc[0]][rc[1]];
                self.pcount[color] -= list.length;
                $(color).innerHTML = self.pcount[color];
                self.undoColor = color;
                self.undoPieces = list.length;
				self.copyBoard(self.board, self.undoBoard);
				self.undoScore = self.removeBlocks(list);
				self.selectedBlocks = null;
				if (self.gameScore > 0) self.enable('menu-undo');
			} else {
				self.deselectBlocks()
				self.selectBlocks(list);
			}
		}
	}

    // copy game board for undo
	this.copyBoard = function(src, dst) {
		for (y = 0; y < self.areaH; y++) {
			for (x = 0; x < self.areaW; x++) {
				dst[y][x] = src[y][x];
			}
		}
	}
	
	// given a rYcX id, return the (y,x)
	this.rc = function(id) {
		var rc = id.match(/r0?(\d+)c0?(\d+)/);
		//if (!rc) alert("id=" + id);
		return new Array(parseInt(rc[1]), parseInt(rc[2]))
	}

	// given a (x,y) position, returns its html id
	this.cid = function (y, x) {
		if (y < 10) y = "0" + y;
		if (x < 10) x = "0" + x;
		return  "r" + y + "c" + x;
	}


	// draw selected blocks and update block score
	this.selectBlocks = function(list) {
		if (list.length < 2) return;

		for (var i = 0; i < list.length; i++) {
			var rc = self.rc(list[i]);
			self.draw(rc[0], rc[1], true);
		}

		self.blockScore(list.length);
		self.selectedBlocks = list.copy();
	}

	// user clicks outside of selected blocks ... deselect them
	this.deselectBlocks = function() {
		if (self.selectedBlocks != null) {
			for (var i = 0; i < self.selectedBlocks.length; i++) {
				var rc = self.rc(self.selectedBlocks[i]);
				self.draw(rc[0], rc[1]);
			}
		}
		self.selectedBlocks = null;
		self.blockScore(0);
	}

	// remove the selected blocks
	this.removeBlocks = function(list) {
		if (list.length < 2) {
			return 0;
		}

		// marked removed blocks
		var minx = 999999;
		var maxx = -1;
		var x;
		var y;
		for (var i = 0; i < list.length; i++) {
			var rc = self.rc(list[i]);
			y = rc[0];
			x = rc[1];

			if (x < minx) minx = x;
			if (x > maxx) maxx = x;
			self.board[y][x] = "remove";
		}

		// shift blocks down
		for (y = 0; y < self.areaH; y++)
		{
			for (x = minx; x <= maxx; x++)
			{
				if (self.board[y][x] == "remove") {

					for (var i = y; i > 0; i--) {
						self.board[i][x] = self.board[i-1][x];
						self.draw(i, x);
					}

					self.board[0][x] = "white";
					self.draw(0, x);
				}
			}
		}

		// look for empty columns and remove them

		for (x = 1; x < self.areaW; x++)
		{
			if (self.isColEmpty(x)) {
				for (var xx = x; xx > 0; xx--) {
					self.moveCols(xx);
				}
			}
		}

		var score = self.score(list.length);

		if (self.isGameOver()) {
			self.endGame(true);
		}

		return score;
	}

	// compute the total score (there may be a bonus), optionally display a message
	// and restart the game
	this.endGame = function(showGameOverDialog) {
		// count # blocks left over for bonus
		var count = 0;
		for (var y = 0; y < self.areaH; y++) {
			for (var x = 0; x < self.areaW; x++) {
				if (self.board[y][x] != "white")
					count++;
			}
		}

		if (self.bonus[count]) {
			bonus = self.bonus[count];
		} else {
			bonus = 0;
		}

		self.gameScore += bonus;
		if (showGameOverDialog) {
			var msg = "Game Over!  You scored " + self.gameScore + " points";
			if (bonus > 0) {
				var blockstr = "block";
				if (count != 1) blockstr += "s";
				msg += " which includes a bonus of " + bonus + " for leaving only " + count + " " +
					blockstr + ".";
			} else {
				msg += ".";
			}
			
			alert(msg);
		}

		self.saveScore(self.gameScore);
		self.start();
	}
		
	
	// fill empty columns with adjacent column
	this.moveCols = function(x) {
		for(var y=0; y < self.areaH; y++)
		{
			self.board[y][x] = self.board[y][x-1];
			self.draw(y, x);
			self.board[y][x-1] = "white";
			self.draw(y, x-1);
		}
	}
	
	// is a (vertical) column empty
	this.isColEmpty = function(c) {
		for (var r = 0; r < self.areaH; r++) {
			if (self.board[r][c] != "white") {
				return false;
			}
		}

		return true;
	}
	
	// see if there's at least one move (two contiguous pieces) left
	this.isGameOver = function() {
		var list = new Array();
		for (var y = 0; y < self.areaH; y++) {
			for (var x = 0; x < self.areaW; x++) {
				if (self.board[y][x] != "white") {
					list = self.getNeighbors(list, y, x, self.board[y][x]);
					if (list.length > 1)
						return false;
					// clear the 1 and only array element
					list.pop();
				}
			}
		}
		return true;
	}
	
	// draw a plain or selected block
	this.draw = function(y, x, selected) {
		var type;
		if (selected) {
			type = "s_";
		} else {
			type = "p_";
		}
		$(self.cid(y, x)).src = "img/" + type + self.board[y][x] + ".gif";
	}

	// redraw entire board
	this.redrawBoard = function() {
		for (var y = 0; y < self.areaH; y++) {
			for (var x = 0; x < self.areaW; x++) {
				// save 100's of function calls by drawing inline
				//draw(y, x);
				$(self.cid(y, x)).src = "img/p_" + self.board[y][x] + ".gif";
			}
		}
	}

	// return a random color
	this.rnd = function (max) {
		var n = Math.ceil(max * Math.random());
		return self.colors[n];
	}

	// compute the current score and update the text
	this.score = function(n) {
		var score = 0;
		if (n > 1) score = n*(n-1);

		self.gameScore += score;
		
		$('userscore').innerHTML = self.gameScore;
		self.blockScore(0);
		return score;
	}

	// update the currently selected block score and count
	this.blockScore = function(n) {
		var score = 0;
		if (n > 1) score = n*(n-1);

		$('blockscore').innerHTML = score;
		$('blockcount').innerHTML = n;
	}
	
	// undo the last move -- we only allow 1 undo level
	this.undo = function() {
		self.copyBoard(self.undoBoard, self.board);
		self.redrawBoard();
		self.gameScore -= self.undoScore;

        var color = self.undoColor;
        self.pcount[color] += self.undoPieces;
        $(color).innerHTML = self.pcount[color];
		self.blockScore(0);
		$('userscore').innerHTML = self.gameScore;
		self.disable('menu-undo');
	}
	
	this.game = function() {
		$('board').style.display='';
		$('about').style.display='none';
		$('help').style.display='none';
		
	}
	this.help = function() {
		$('board').style.display='none';
		$('about').style.display='none';
		$('help').style.display='';
	}
	
	this.about = function() {
		$('board').style.display='none';
		$('about').style.display='';
		$('help').style.display='none';
	}



    // windows
    var helpwindow  = new Window("game-help");
    var statswindow = new Window("game-stats");


	$("menu-start").onclick = function() { helpwindow.close(); statswindow.close(); self.start(); this.blur(); }
	$("menu-undo").onclick = function() { self.undo(); this.blur(); }
	$("menu-help").onclick = function() { statswindow.close(); helpwindow.activate(); this.blur(); }
	$("menu-stats").onclick = function() { helpwindow.close(); statswindow.activate(); this.blur(); }


    $("game-help-close").onclick  = helpwindow.close;
    $("game-stats-close").onclick = statswindow.close;

    /**
     * Window replaces game area, for example help window
     * @param string id
     */
    function Window(id) {
        
        this.id = id;
        this.el = document.getElementById(this.id);
        var self = this;
        
        /**
         * Activate or deactivate a window - update html
         * @return void
         * @access event
         */
        this.activate = function() {
            self.el.style.display = (self.el.style.display == "block" ? "none" : "block");
        }
        
        /**
         * Close window - update html
         * @return void
         * @access event
         */
        this.close = function() {
            self.el.style.display = "none";
        }
        
        /**
         * @return bool
         * @access public
         */
        this.isActive = function() {
            return (self.el.style.display == "block");
        }
    }
}
