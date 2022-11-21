package handlers

import (
	"fmt"
	"github.com/ProninIgorr/alcomarket-nearst/helpers/response"
	"github.com/ProninIgorr/alcomarket-nearst/internal/store"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"

	"github.com/d-kolpakov/logger"
)

type Handler struct {
	Db          *pgxpool.Pool
	L           *logger.Logger
	AppVersion  string
	ServiceName string
	Store       *store.Store
}

func (h *Handler) HomeRouteHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, fmt.Sprintf("Hello! This is %s. Version: %s", h.ServiceName, h.AppVersion))
}

func (h *Handler) InternalEndpoint(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Internal endpoint.")
}
