package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/n1ck-kk/nicksite/internal/config"
	"github.com/n1ck-kk/nicksite/internal/handlers"
	"github.com/n1ck-kk/nicksite/internal/lab"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	lab.InitStore(cfg.SessionSecret, cfg.IsProduction())

	h := handlers.New(cfg.LabPassword)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	staticFS := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.Handle("/static/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		staticFS.ServeHTTP(w, req)
	}))

	r.Get("/health", h.HandleHealth)
	r.Get("/", h.HandleHome)

	// Lab routes
	r.Get("/lab", h.HandleLabLogin)
	r.Post("/lab", h.HandleLabLoginPost)
	r.Group(func(r chi.Router) {
		r.Use(requireLabAuth)
		r.Get("/lab/ideas", h.HandleLabIdeas)
		r.Post("/lab/logout", h.HandleLabLogout)
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server starting on :%s (env=%s)", cfg.Port, cfg.Environment)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case <-shutdown:
		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			srv.Close()
		}
	}

	fmt.Println("Server stopped")
}

func requireLabAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lab.IsAuthenticated(r) {
			http.Redirect(w, r, "/lab", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
