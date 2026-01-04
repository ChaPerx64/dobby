package api

import (
	"context"
	"github.com/ChaPerx64/dobby/apps/backend/internal/oas"
	"github.com/google/uuid"
	"log"
)

func (h dobbyHandler) GetCurrentUser(ctx context.Context) (r *oas.User, _ error) {
	log.Println("Got a request @/me")
	return &oas.User{
		ID:             uuid.UUID{},
		Name:           "mock name",
		CurrentBalance: 10000,
	}, nil
}
