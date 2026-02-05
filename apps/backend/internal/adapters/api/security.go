package api

import (
	"context"
	"fmt"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
)

type dobbySecurity struct {
	authZ *authorization.Authorizer[*oauth.IntrospectionContext]
}

func (s *dobbySecurity) HandleBearerAuth(ctx context.Context, operationName oas.OperationName, t oas.BearerAuth) (context.Context, error) {
	if s.authZ == nil {
		return ctx, fmt.Errorf("authorization not configured")
	}

	authCtx, err := s.authZ.CheckAuthorization(ctx, t.Token)
	if err != nil {
		return nil, err
	}

	// Extract the subject (user ID) from the authorized context
	// If no error, we are authorized.
	return context.WithValue(ctx, userIDKey, authCtx.UserID), nil
}

// GetUserID retrieves the user ID from the context if it exists.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
