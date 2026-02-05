package api

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/rs/cors"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type dobbyHandler struct {
	oas.UnimplementedHandler // automatically implement all methods
}

// Compile-time check for Handler.
var _ oas.Handler = (*dobbyHandler)(nil)

func RunServer() {
	ctx := context.Background()
	issuer := os.Getenv("ZITADEL_ISSUER")
	if issuer == "" {
		slog.Warn("ZITADEL_ISSUER not set, auth will fail")
	}

	var authZ *authorization.Authorizer[*oauth.IntrospectionContext]
	if issuer != "" {
		clientID := os.Getenv("ZITADEL_CLIENT_ID")
		if clientID == "" {
			log.Fatal("ZITADEL_CLIENT_ID not set")
		}

		// Create a ZITADEL client.
		// zitadel.New expects the domain, not the URL.
		domain := strings.TrimPrefix(issuer, "https://")
		domain = strings.TrimPrefix(domain, "http://")
		z := zitadel.New(domain)

		var err error
		// Use Stateless JWT Validation (Offline)
		authZ, err = authorization.New(ctx, z, oauth.WithJWT(clientID, nil))
		if err != nil {
			log.Fatalf("failed to create authorizer: %v", err)
		}
	}

	srv, err := oas.NewServer(dobbyHandler{}, &dobbySecurity{authZ: authZ})
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
