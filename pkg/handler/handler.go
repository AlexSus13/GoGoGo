package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)
type Handler struct {

}

func(h *Handler) NewRouter() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/", h.Hello).Methods("GET")
	router.HandleFunc("/getfiles/{filename}", h.HandlerGET).Methods("GET")
	router.HandleFunc("/postfile", h.HandlerPOST).Methods("POST")
	router.HandleFunc("/getfile/{filename}", h.HandlerFILE).Methods("GET")

	return router
}
