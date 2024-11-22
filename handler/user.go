package handler

import (
	"errors"
	"net/http"

	"github.com/webdevfuel/projectmotor/template"
)

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	sessions, err := h.SessionService.GetAllSessions(user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	sess, err := h.GetSessionStore(r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	tok, ok := sess.Values["token"].(string)
	if !ok {
		h.Error(
			w,
			errors.New("token must be present in session values"),
			http.StatusInternalServerError,
		)
		return
	}
	component := template.Profile(sessions, tok)
	component.Render(r.Context(), w)
}
