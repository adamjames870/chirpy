package main

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/adamjames870/chirpy/internal/database"
)

const expiryTimeRefreshToken time.Duration = time.Duration(60 * 24 * time.Hour) // 60 days
const expiryTimeAccesToken time.Duration = 60 * 60 * time.Millisecond           // 1 hour

type apiState struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	mux            *http.ServeMux
	platform       string
	secret_string  string
}
