package handlers

import (
	"net/http"

	"github.com/n1ck-kk/nicksite/internal/views"
)

type Handler struct {
	labPassword string
}

func New(labPassword string) *Handler {
	return &Handler{labPassword: labPassword}
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	views.HomePage().Render(r.Context(), w)
}
