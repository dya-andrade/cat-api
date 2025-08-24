package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/dya-andrade/cat-api/internal/domain"
	"github.com/dya-andrade/cat-api/internal/service"
)

type CatsHandler struct {
	svc       service.CatService
	validator *validator.Validate
}

// Construtor do handler. Recebe o serviço e cria o validador.
func NewCatsHandler(svc service.CatService) *CatsHandler {
	return &CatsHandler{
		svc:       svc,
		validator: validator.New(),
	}
}

// List: lista gatos com paginação.
// - Lê o parâmetro "limit" da URL, define limite de itens (padrão 20, máximo 100).
// - Lê o parâmetro "cursor" da URL, converte para time.Time se existir.
// - Chama o serviço para buscar os gatos.
// - Se houver erro, retorna erro 500.
// - Monta resposta JSON com os gatos e o próximo cursor (se existir).
func (h *CatsHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if lstr := r.URL.Query().Get("limit"); lstr != "" {
		if l, err := strconv.Atoi(lstr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	var cursor *time.Time
	if cstr := r.URL.Query().Get("cursor"); cstr != "" {
		if t, err := time.Parse(time.RFC3339, cstr); err == nil {
			cursor = &t
		}
	}

	cats, next, err := h.svc.List(r.Context(), limit, cursor)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}

	resp := map[string]any{
		"items": cats,
	}
	if next != nil {
		resp["next_cursor"] = next.Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, resp)
}

// Create: cria um novo gato.
// - Decodifica o corpo da requisição para struct CatCreate.
// - Valida os dados recebidos.
// - Chama o serviço para criar o gato.
// - Se houver erro, retorna erro apropriado.
// - Retorna o gato criado em JSON.
func (h *CatsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in domain.CatCreate
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.validator.Struct(in); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err)
		return
	}
	cat, err := h.svc.Create(r.Context(), in)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, cat)
}

// GetByID: busca um gato pelo ID.
// - Lê o parâmetro "id" da URL.
// - Converte para inteiro e valida.
// - Chama o serviço para buscar o gato.
// - Se não encontrar, retorna erro 404.
// - Se houver outro erro, retorna erro 500.
// - Retorna o gato encontrado em JSON.
func (h *CatsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, errors.New("id inválido"))
		return
	}
	cat, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpError(w, http.StatusNotFound, err)
			return
		}
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

/*
func (h *CatsHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, errors.New("id inválido"))
		return
	}
	var in domain.UpdateCatInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.validator.Struct(in); err != nil {
		httpError(w, http.StatusUnprocessableEntity, err)
		return
	}
	cat, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpError(w, http.StatusNotFound, err)
			return
		}
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

func (h *CatsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, errors.New("id inválido"))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpError(w, http.StatusNotFound, err)
			return
		}
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
*/

/*** helpers ***/
// Helpers para resposta JSON e erro.
// writeJSON: escreve resposta JSON com status.
// httpError: escreve resposta de erro em JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{
		"error": err.Error(),
	})
}
