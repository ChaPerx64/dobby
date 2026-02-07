package api

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/rs/cors"
)

type dobbyHandler struct {
	oas.UnimplementedHandler // automatically implement all methods
}

// Compile-time check for Handler.
var _ oas.Handler = (*dobbyHandler)(nil)

func RunServer() {
	authority := os.Getenv("OIDC_AUTHORITY")
	if authority == "" {
		slog.Warn("OIDC_AUTHORITY not set, auth will fail")
	}

	security := &dobbySecurity{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	if authority != "" {
		clientID := os.Getenv("OIDC_BACKEND_CLIENT_ID")
		clientSecret := os.Getenv("OIDC_BACKEND_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			log.Fatal("OIDC_BACKEND_CLIENT_ID or OIDC_BACKEND_CLIENT_SECRET not set")
		}

		// OIDC Discovery
		discoveryURL := strings.TrimSuffix(authority, "/") + "/.well-known/openid-configuration"
		resp, err := http.Get(discoveryURL)
		if err != nil {
			log.Fatalf("failed to fetch OIDC discovery: %v", err)
		}
		defer resp.Body.Close()

		var config struct {
			IntrospectionEndpoint string `json:"introspection_endpoint"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			log.Fatalf("failed to decode OIDC discovery: %v", err)
		}

		if config.IntrospectionEndpoint == "" {
			log.Fatal("OIDC discovery response missing introspection_endpoint")
		}

		security.introspectionURL = config.IntrospectionEndpoint
		security.clientID = clientID
		security.clientSecret = clientSecret
		slog.Info("OIDC security initialized", "introspection_url", security.introspectionURL)
	}

	srv, err := oas.NewServer(dobbyHandler{}, security)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Swagger UI at /docs
	mux.Handle("/docs/", SwaggerUIHandler())
	mux.HandleFunc("/docs/openapi.yml", OpenAPISpecHandler)

	// API routes (Ogen handles /api/v1 prefix internally)
	mux.Handle("/", srv)

	allowedOrigins := []string{"*"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = []string{frontendURL}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe("0.0.0.0:8080", c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
