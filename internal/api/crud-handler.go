package api

import (
	"context"
	"encoding/json"
	"net/http"

	"cais/internal/db/postgresql"
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

func (ih itemHandler) GetAllItems(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	items, err := ih.db.GetAllItems(ih.ctx)
	if err != nil {
		ih.log.Error(err)
		w.WriteHeader(400)
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		ih.log.Error(err)
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
