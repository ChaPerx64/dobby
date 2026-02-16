package api

import (
	_ "embed"
	"net/http"

	"github.com/swaggest/swgui/v5emb"
)

//go:embed openapi_copy.yml
var openapiSpec []byte

func SwaggerUIHandler() http.Handler {
	return v5emb.New("Dobby API", "/docs/openapi.yml", "/docs/")
}

func OpenAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(openapiSpec)
}
