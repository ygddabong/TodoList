package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// NewHandler returns an http.Handler exposing the to-do REST API.
func NewHandler(store *Store) http.Handler {
	h := &handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", h.ping)
	mux.HandleFunc("GET /todos", h.list)
	mux.HandleFunc("POST /todos", h.create)
	mux.HandleFunc("GET /todos/{id}", h.get)
	mux.HandleFunc("PATCH /todos/{id}", h.update)
	mux.HandleFunc("DELETE /todos/{id}", h.delete)
	return mux
}

type handler struct {
	store *Store
}

type todoInput struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (h *handler) ping(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

func (h *handler) list(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	t, err := h.store.Get(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var in todoInput
	if !decode(w, r, &in) {
		return
	}
	if in.Title == nil || strings.TrimSpace(*in.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	t := h.store.Create(strings.TrimSpace(*in.Title))
	w.Header().Set("Location", "/todos/"+strconv.Itoa(t.ID))
	writeJSON(w, http.StatusCreated, t)
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in todoInput
	if !decode(w, r, &in) {
		return
	}
	if in.Title != nil {
		trimmed := strings.TrimSpace(*in.Title)
		if trimmed == "" {
			writeError(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		in.Title = &trimmed
	}
	t, err := h.store.Update(id, in.Title, in.Done)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.store.Delete(id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeError(w, http.StatusInternalServerError, "internal error")
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
