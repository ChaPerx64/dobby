package api

import (
	"context"

	"github.com/ChaPerx64/dobby/apps/backend/internal/oas"
)

type DobbySecurity struct {
	// You might hold a JWT verifier or database connection here
}

func (s *DobbySecurity) HandleBearerAuth(ctx context.Context, operationName oas.OperationName, t oas.BearerAuth) (context.Context, error) {
	return ctx, nil
}
