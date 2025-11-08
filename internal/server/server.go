package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/repository/postgres"
	"todo-api/internal/router"
)

// server que encapsula el servidor HTTP y dependencias
type Server struct {
	httpServer *http.Server
	db         *postgres.DB
	config     *config.Config
}

func New(cfg *config.Config) (*Server, error) {

	db, err := postgres.NewConnection(cfg.Database)

	if err != nil {
		return nil, err
	}

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      router.New(db),
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		db:     db,
		config: cfg,
	}, nil
}

// Start inicia el servidor HTTP y gestiona su ciclo de vida.
// Lanza el servidor en una goroutine y se bloquea esperando
// a que ocurra un error o se reciba una señal del sistema (Ctrl+C o SIGTERM)
// para apagarlo de forma controlada (graceful shutdown).
func (s *Server) Start() error {

	// Canal para comunicar errores del servidor HTTP.
	// Se usa un buffer de 1 para evitar bloqueos si se envía un error.
	errChan := make(chan error, 1)

	// Iniciar el servidor en una goroutine para no bloquear el flujo principal.
	go func() {
		log.Printf("Server starting on port %d", s.config.Server.Port)

		// Escucha y sirve peticiones HTTP. Esta llamada bloquea
		// hasta que ocurre un error o se cierra el servidor.
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Enviar cualquier error inesperado al canal de errores.
			errChan <- err
		}
	}()

	// Canal para recibir señales del sistema (interrupciones, SIGTERM, etc.).
	quit := make(chan os.Signal, 1)

	// Registrar las señales que se quieren escuchar.
	// os.Interrupt → Ctrl+C en terminal
	// syscall.SIGTERM → señal estándar de apagado en sistemas Unix/Docker
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Esperar de forma bloqueante hasta que:
	//  - llegue un error en el servidor
	//  - o se reciba una señal del sistema
	select {
	case err := <-errChan:
		// Si el servidor devuelve un error inesperado, propagarlo.
		return err

	case <-quit:
		// Si se recibe una señal de cierre, iniciar el apagado controlado.
		return s.Shutdown()
	}
}

func (s *Server) Shutdown() error {

	log.Println("Shutting down...")
	// Crear un contexto con timeout para permitir que las conexiones activas
	// finalicen correctamente antes de forzar el cierre.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Apagar el servidor de manera ordenada usando el contexto.
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	// cierra la db
	return s.db.Close()
}
