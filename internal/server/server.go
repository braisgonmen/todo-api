package server

import(
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	db			*database.DB
	config		*config.Config
}
func NewServer(
	logger *Logger
	config	*config
	commentStore *commentStore
	anotherStore *anotherStore
) http.Handler {
	mux := http.NewServerMux()
	addRoutes(
		mux,
		Logger,
		Config,
		commentStore,
		anotherStore
	)

	var handler http.Handler = mux
	handler
}
