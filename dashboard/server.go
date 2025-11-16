package dashboard

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
)

//go:embed public/*
var StaticFS embed.FS

var (
	StaticSys = hashfs.NewFS(StaticFS)
)

type Server struct {
	env    *EnvVars
	mux    chi.Router
	server *http.Server
}

type EnvVars struct {
	IsProd        bool
	LogLocation   string
	SensorAddress string
}

var ServerEnv EnvVars

func NewServer(address string, env *EnvVars) *Server {

	mux := chi.NewMux()

	return &Server{
		env: env,
		mux: mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func StaticPath(format string, args ...any) string {
	return "/" + StaticSys.HashName(fmt.Sprintf("public/"+format, args...))
}

func (s *Server) Start() error {

	slog.Info("Starting webserver", "address", s.server.Addr)

	s.mux.Group(func(router chi.Router) {

		router.Handle("/public/*", hashfs.FileServer(StaticSys))

		webRoutes(router, s.env)

		router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("route does not exist"))
			slog.Info("route does not exist" + r.URL.Path)
		})

	})

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	slog.Info("Stopping the http server")
	// ensure shutdown doesnt hang
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Stopped the http server")
	return nil
}

func logAndError(w http.ResponseWriter, err error) {
	slog.Error("App error", "error", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
}

func formatError(prefix string, r *http.Request, err error) error {
	return fmt.Errorf(prefix+" - url:%v error:%v", chi.RouteContext(r.Context()).RoutePattern(), err)
}
