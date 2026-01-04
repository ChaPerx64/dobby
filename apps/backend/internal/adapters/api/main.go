package api

import (
	"log"
	"net/http"

	"github.com/ChaPerx64/dobby/apps/backend/internal/oas"
)

type dobbyHandler struct {
	oas.UnimplementedHandler // automatically implement all methods
}

// Compile-time check for Handler.
var _ oas.Handler = (*dobbyHandler)(nil)

func RunServer() {
	srv, err := oas.NewServer(dobbyHandler{}, &dobbySecurity{})
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe("0.0.0.0:8080", srv); err != nil {
		log.Fatal(err)
	}
}
