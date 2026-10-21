package article

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"article/internal/httpx"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc} }

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /article/{$}", h.create)
	mux.HandleFunc("GET /article/{limit}/{offset}", h.list)
	mux.HandleFunc("GET /article/{id}", h.get)
	mux.HandleFunc("PUT /article/{id}", h.update)
	mux.HandleFunc("DELETE /article/{id}", h.delete)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	a, ok := decode(w, r)
	if !ok {
		return
	}
	if err := h.svc.Create(a); err != nil {
		respondErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, struct{}{})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, e1 := strconv.Atoi(r.PathValue("limit"))
	offset, e2 := strconv.Atoi(r.PathValue("offset"))
	if e1 != nil || e2 != nil || limit < 0 || offset < 0 {
		httpx.Error(w, http.StatusBadRequest, "limit and offset must be non-negative integers")
		return
	}
	out, err := h.svc.List(limit, offset)
	if err != nil {
		respondErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, out)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	a, err := h.svc.Get(id)
	if err != nil {
		respondErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, a)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	a, ok := decode(w, r)
	if !ok {
		return
	}
	if err := h.svc.Update(id, a); err != nil {
		respondErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, struct{}{})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		respondErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, struct{}{})
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func decode(w http.ResponseWriter, r *http.Request) (Article, bool) {
	var a Article
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json")
		return a, false
	}
	return a, true
}

// respondErr memetakan error domain ke status HTTP.
func respondErr(w http.ResponseWriter, err error) {
	var ve ValidationError
	switch {
	case errors.As(err, &ve):
		httpx.JSON(w, http.StatusBadRequest, map[string]any{"errors": []string(ve)})
	case errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "not found")
	default:
		httpx.Error(w, http.StatusInternalServerError, err.Error())
	}
}
