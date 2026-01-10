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

	mux := http.NewServeMux()

	// Swagger UI at /docs
	mux.Handle("/docs/", SwaggerUIHandler())
	mux.HandleFunc("/docs/openapi.yml", OpenAPISpecHandler)

	// API routes (Ogen handles /api/v1 prefix internally)
	mux.Handle("/", srv)

	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Fatal(err)
	}
}
