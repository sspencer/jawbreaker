package jawbreaker

import "core:fmt"
import "core:log"
import "core:os"
import "core:path/filepath"
import "core:strconv"
import "core:strings"

APP_NAME :: "jawbreaker"
FILENAME :: "highscore.txt"

// High score structure - you can modify this as needed
Game_State :: struct {
	best_score: int,
	last_score: int,
}

// Get the appropriate data directory for the current OS
get_app_data_dir :: proc() -> (string, bool) {
	when ODIN_OS == .Windows {
		// Use %APPDATA% on Windows
		appdata := os.get_env("APPDATA")
		if appdata == "" {
			log.error("Could not get APPDATA environment variable")
			return "", false
		}
		return filepath.join({appdata, APP_NAME}), true
	} else when ODIN_OS == .Darwin {
		// Use ~/Library/Application Support on macOS
		home := os.get_env("HOME")
		if home == "" {
			log.error("Could not get HOME environment variable")
			return "", false
		}
		return filepath.join({home, "Library", "Application Support", APP_NAME}), true
	} else when ODIN_OS == .Linux {
		// Use XDG_DATA_HOME or ~/.local/share on Linux
		xdg_data_home := os.get_env("XDG_DATA_HOME")
		if xdg_data_home != "" {
			return filepath.join({xdg_data_home, APP_NAME}), true
		}

		home := os.get_env("HOME")
		if home == "" {
			log.error("Could not get HOME environment variable")
			return "", false
		}
		return filepath.join({home, ".local", "share", APP_NAME}), true
	} else {
		// Fallback for other systems
		log.warn("Unsupported OS, using current directory")
		return APP_NAME, true
	}
}

// Ensure the app data directory exists
ensure_app_data_dir :: proc() -> (string, bool) {
	dir, ok := get_app_data_dir()
	if !ok {
		return "", false
	}

	// Create directory if it doesn't exist
	if !os.exists(dir) {
		err := os.make_directory(dir, 0o755)
		if err != nil {
			log.errorf("Failed to create directory %s: %v", dir, err)
			return "", false
		}
		log.infof("Created app data directory: %s", dir)
	}

	return dir, true
}

// Save high score to disk
save_high_score :: proc(state: Game_State) -> bool {
	data_dir, ok := ensure_app_data_dir()
	if !ok {
		return false
	}

	score_file := filepath.join({data_dir, FILENAME})

	// Format the high score data
	score_data := fmt.aprintf(
		"%d|%d\n",
		state.best_score,
		state.last_score,
	)
	defer delete(score_data)

	// Write to file
	success := os.write_entire_file(score_file, transmute([]u8)score_data)
	if !success {
		log.errorf("Failed to write high score to %s", score_file)
		return false
	}

	log.infof("High score saved to %s", score_file)
	return true
}

// Load high score from disk
load_high_score :: proc() -> (Game_State, bool) {
	data_dir, ok := get_app_data_dir()
	if !ok {
		return {}, false
	}

	score_file := filepath.join({data_dir, FILENAME})

	// Check if file exists
	if !os.exists(score_file) {
		log.infof("High score file does not exist: %s", score_file)
		return {}, false
	}

	// Read file contents
	file_data, read_ok := os.read_entire_file(score_file)
	if !read_ok {
		log.errorf("Failed to read high score file: %s", score_file)
		return {}, false
	}
	defer delete(file_data)

	// Parse the data
	content := strings.trim_space(string(file_data))
	parts := strings.split(content, "|")
	defer delete(parts)

	if len(parts) < 1 {
		log.errorf("Invalid high score file format in %s", score_file)
		return {}, false
	}

	best_score, best_ok := strconv.parse_int(parts[0])
	if !best_ok {
		log.errorf("Invalid best score value in high score file: %s", parts[0])
		return {}, false
	}

	last_score = 0
	if len(parts) > 1 {
		last_ok: bool
		last_score, last_ok = strconv.parse_int(parts[1])
		if !last_ok {
			log.errorf("Invalid score value in high score file: %s", parts[1])
			return {}, false
		}
	}

	state := Game_State {
		best_score = best_score,
		last_score = last_score,
	}

	log.infof("Score loaded from %s", score_file)
	return state, true
}

// Helper function to check if a new score is a high score
is_new_high_score :: proc(new_score: int) -> bool {
	state, ok := load_high_score()
	if !ok {
		// No existing high score, so any score is a high score
		return true
	}

	return new_score > state.best_score
}
