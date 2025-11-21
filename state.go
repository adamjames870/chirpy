package main

import (
	"net/http"
	"sync/atomic"

	"github.com/adamjames870/chirpy/internal/database"
)

type apiState struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	mux            *http.ServeMux
	platform       string
	secret_string  string
}
