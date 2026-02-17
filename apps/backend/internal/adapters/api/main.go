package api

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/persistence"
	"github.com/ChaPerx64/dobby/apps/backend/internal/config"
	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
	"github.com/rs/cors"
)

type dobbyHandler struct {
	oas.UnimplementedHandler // automatically implement all methods
	financeService           service.FinanceService
}

// Compile-time check for Handler.
var _ oas.Handler = (*dobbyHandler)(nil)

func RunServer(cfg config.Config) {
	authority := cfg.OIDCAuthority

	security := &dobbySecurity{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	clientID := cfg.OIDCBackendClientID
	clientSecret := cfg.OIDCBackendClientSecret

	// OIDC Discovery
	discoveryURL := strings.TrimSuffix(authority, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(discoveryURL)
	if err != nil {
		log.Fatalf("failed to fetch OIDC discovery: %v", err)
	}
	defer resp.Body.Close()

	var disco struct {
		IntrospectionEndpoint string `json:"introspection_endpoint"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&disco); err != nil {
		log.Fatalf("failed to decode OIDC discovery: %v", err)
	}

	if disco.IntrospectionEndpoint == "" {
		log.Fatal("OIDC discovery response missing introspection_endpoint")
	}

	security.introspectionURL = disco.IntrospectionEndpoint
	security.clientID = clientID
	security.clientSecret = clientSecret
	slog.Info("OIDC security initialized", "introspection_url", security.introspectionURL)

	repo := persistence.NewMemoryRepository()
	svc := service.NewDobbyFinancier(repo)

	srv, err := oas.NewServer(&dobbyHandler{financeService: svc}, security)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Swagger UI at /docs
	mux.Handle("/docs/", SwaggerUIHandler())
	mux.HandleFunc("/docs/openapi.yml", OpenAPISpecHandler)

	// API routes (Ogen handles /api/v1 prefix internally)
	mux.Handle("/", srv)

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	addr := ":" + cfg.BackendPort
	slog.Info("Starting server", "addr", addr)
	if err := http.ListenAndServe(addr, c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
