package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/tartushkin/TSHORT.git/internal/service"
)

type Handlers struct {
	Short      *service.Short // внутриняя логика приложения
	httpServer *http.Server
}

func NewHandlers(short *service.Short) *Handlers {
	return &Handlers{Short: short}
}

// NewPersonHandlers создает новый экземпляр обработчиков запросов для Person
func (h *Handlers) newRoutes() *mux.Router {
	ht := mux.NewRouter()

	ht.HandleFunc("/", h.postHandler).Methods("POST")
	ht.HandleFunc("/{id}", h.getHandler).Methods("GET")
	return ht
}

func (h *Handlers) StartHTTP(ctx context.Context, httpPort int) error {

	h.httpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(httpPort),
		Handler: h.newRoutes(),
	}

	err := h.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("net.Listen: %s", err.Error())
	}

	return nil
}

// остнавка http сервера
func (h *Handlers) StopHTTP(ctx context.Context) {
	h.httpServer.Shutdown(ctx)
}
