package handlers

import (
	"net/http"

	"github.com/n1ck-kk/nicksite/internal/lab"
	"github.com/n1ck-kk/nicksite/internal/views"
)

func (h *Handler) HandleLabLogin(w http.ResponseWriter, r *http.Request) {
	if lab.IsAuthenticated(r) {
		http.Redirect(w, r, "/lab/ideas", http.StatusFound)
		return
	}
	views.LabLogin("").Render(r.Context(), w)
}

func (h *Handler) HandleLabLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if r.FormValue("password") != h.labPassword {
		views.LabLogin("Wrong password.").Render(r.Context(), w)
		return
	}

	if err := lab.SetAuthenticated(w, r); err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/lab/ideas", http.StatusFound)
}

func (h *Handler) HandleLabIdeas(w http.ResponseWriter, r *http.Request) {
	views.LabIdeas(lab.Ideas).Render(r.Context(), w)
}

func (h *Handler) HandleLabLogout(w http.ResponseWriter, r *http.Request) {
	lab.ClearSession(w, r)
	http.Redirect(w, r, "/lab", http.StatusFound)
}
