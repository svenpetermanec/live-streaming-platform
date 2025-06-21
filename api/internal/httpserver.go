package internal

import (
	"context"
	"fmt"
	"net/http"

	"transcoder/internal/logging"
)

type HttpServer struct {
	controller *Controller
	port       int
	logger     *logging.Logger
}

func NewHttpServer(controller *Controller, port int, logger *logging.Logger) *HttpServer {
	return &HttpServer{
		controller: controller,
		port:       port,
		logger:     logger,
	}
}

func (h *HttpServer) Start(ctx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc("/streams/", h.controller.ServeStream)
	mux.HandleFunc("/account/create", h.controller.CreateAccount)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", h.port),
		Handler: mux,
	}

	serverErr := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			serverErr <- err
		}
	}()

	h.logger.Info("Started listening", nil)

	select {
	case <-ctx.Done():
		h.logger.Info("Shutting down", nil)
		server.Shutdown(ctx)
	case err := <-serverErr:
		h.logger.Error("Server failed to start", logging.Data{"error": err})
		return
	}
}
