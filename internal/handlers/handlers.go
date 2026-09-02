package handlers

import (
	"net/http"

	"github.com/n1ck-kk/nicksite/internal/views"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	views.HomePage().Render(r.Context(), w)
}
