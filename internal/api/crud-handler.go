package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cais/internal/db/postgresql"
	"cais/internal/entities"
	"cais/pkg/logger"

	"github.com/julienschmidt/httprouter"
)

type itemHandler struct {
	ctx context.Context
	db  db.Storage
	log *logger.Logger
}

func NewItemHandler(ctx context.Context, db db.Storage, log *logger.Logger) *itemHandler {
	return &itemHandler{
		ctx: ctx,
		db:  db,
		log: log,
	}
}

// GetAllItems получить список всех item (GET)
func (ih *itemHandler) GetAllItems(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	items, err := ih.db.GetAllItems(ih.ctx)
	if err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ih.writeJSONResponse(w, http.StatusOK, items)
}

// CreateNewItem создает новый item, возвращает message с информацией о статусе выполнения (POST)
func (ih *itemHandler) CreateNewItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var item entities.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := ih.db.ItemsHandleCUD(ih.ctx, "create", item); err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Entity with id: %d and name '%s' has been added to the database", item.Id, item.Name),
	}

	ih.writeJSONResponse(w, http.StatusOK, resp)
}

// DeleteItem удаляет сущность item (DELETE)
func (ih *itemHandler) DeleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := ih.db.ItemsHandleCUD(ih.ctx, "delete", entities.Item{Id: int64(id)}); err != nil {
		if err != nil {
			ih.log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Entity with id: %d was successfully deleted", id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// UpdateItem обновляет сущность item (PATH) TODO:
func (ih *itemHandler) UpdateItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var item entities.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := ih.db.ItemsHandleCUD(ih.ctx, "update", item); err != nil {
		ih.log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Entity with id: %d and name '%s' has been updated to the database", item.Id, item.Name),
	}

	ih.writeJSONResponse(w, http.StatusOK, resp)
}
