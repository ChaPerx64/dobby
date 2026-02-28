package api

import (
	"testing"
	"time"

	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
	"github.com/google/uuid"
)

func TestMapPeriodSummaryToOAS(t *testing.T) {
	serviceSummary := &service.PeriodSummary{
		Period: service.Period{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now(),
		},
		TotalBudget:            1000,
		TotalSpent:             500,
		TotalRemaining:         500,
		ProjectedEndingBalance: -100,
		EnvelopeStats: []service.EnvelopeStat{
			{
				Envelope: service.Envelope{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Name: "Groceries",
				},
				Allocated: 1000,
				Spent:     500,
				Remaining: 500,
			},
		},
	}

	oasSummary := mapPeriodSummaryToOAS(serviceSummary)

	if oasSummary.TotalSpent != 500 {
		t.Errorf("expected TotalSpent to be 500, got %d", oasSummary.TotalSpent)
	}

	if len(oasSummary.EnvelopeSummaries) != 1 {
		t.Fatalf("expected 1 envelope summary, got %d", len(oasSummary.EnvelopeSummaries))
	}

	if oasSummary.EnvelopeSummaries[0].Spent != 500 {
		t.Errorf("expected EnvelopeSummaries[0].Spent to be 500, got %d", oasSummary.EnvelopeSummaries[0].Spent)
	}
}
