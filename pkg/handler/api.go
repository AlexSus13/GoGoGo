package handler

import (
	"net/http"
)

func(h *Handler) Hello(w http.ResponseWriter, r *http.Request) {}

func(h *Handler) HandlerGET(w http.ResponseWriter, r *http.Request) {}

func(h *Handler) HandlerPOST(w http.ResponseWriter, r *http.Request) {}

func(h *Handler) HandlerFILE(w http.ResponseWriter, r *http.Request) {}
