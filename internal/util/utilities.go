package util

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Logger functions to write to a file for debugginf purposes.
func logger() *os.File {
	logFile, err := os.OpenFile("api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return logFile
}

func WriteLog(message interface{}, arg ...interface{}) {
	logFile := logger()
	defer logFile.Close() // Close the log file when the application exits

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Println(message, arg)
}

// Functions to send respoonses back to the client.
func ParseUserID(vars map[string]string) (int64, error) {
	idStr, ok := vars["id"]
	if !ok {
		return 0, errors.New("id is missing")
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}
	return int64(userID), nil
}
func RespondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithSuccess(w http.ResponseWriter, payload interface{}) {
	respondWithJSON(w, http.StatusOK, payload)
}
