package server

import (
	"github/SirVoly/chirpy/internal/database"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             		*database.Queries
	platform       		string
	JWTsecret			string
}
