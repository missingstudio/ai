package api

import (
	"log/slog"

	"github.com/missingstudio/studio/backend/core/connection"
	"github.com/missingstudio/studio/backend/core/prompt"
	"github.com/missingstudio/studio/backend/internal/ingester"
	"github.com/missingstudio/studio/backend/internal/providers"
	"github.com/missingstudio/studio/backend/internal/ratelimiter"
)

type Deps struct {
	Logger            *slog.Logger
	Ingester          ingester.Ingester
	RateLimiter       *ratelimiter.RateLimiter
	ProviderService   *providers.Service
	ConnectionService *connection.Service
	PromptService     *prompt.Service
}

func NewDeps(
	logger *slog.Logger,
	ingester ingester.Ingester,
	ratelimiter *ratelimiter.RateLimiter,
	ps *providers.Service,
	cs *connection.Service,
	pms *prompt.Service,
) *Deps {
	return &Deps{
		Logger:            logger,
		Ingester:          ingester,
		RateLimiter:       ratelimiter,
		ProviderService:   ps,
		ConnectionService: cs,
		PromptService:     pms,
	}
}
