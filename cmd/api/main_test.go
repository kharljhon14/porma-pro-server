package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
	"github.com/stretchr/testify/require"
)

func newTestingServer(t *testing.T, store db.Store) *Server {
	server, err := NewServer(store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
