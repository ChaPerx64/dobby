package oas

import (
	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
)

// ToLogicModel converts CreateEnvelope DTO to logic parameters.
// We only return name because ID and UserID are handled by service/handler.
func (req *CreateEnvelope) ToLogicModel() string {
	return req.Name
}

// ToLogicModel converts CreateTransaction DTO to logic model.
func (req *CreateTransaction) ToLogicModel() service.Transaction {
	t := service.Transaction{
		EnvelopeID: req.EnvelopeId,
		Amount:     req.Amount,
	}

	if v, ok := req.Description.Get(); ok {
		t.Description = v
	}

	if v, ok := req.Date.Get(); ok {
		t.Date = v
	}

	if v, ok := req.Category.Get(); ok {
		t.Category = v
	}

	return t
}
