package middlewares

import (
	"context"

	"biinge-api/internal/app/models"
)

type Claim struct{}
type Token struct{}
type TraceId struct{}
type CurrentUser struct{}

type Modifier interface {
	WithTraceId(traceId string) Modifier
	WithCurrentUser(user *models.User) Modifier
	Context() context.Context
}

type modifier struct {
	ctx context.Context
}

func NewContextModifier(ctx context.Context) Modifier {
	return &modifier{ctx: ctx}
}

func (m *modifier) WithTraceId(traceId string) Modifier {
	m.ctx = context.WithValue(m.ctx, TraceId{}, traceId)
	return m
}

func (m *modifier) WithCurrentUser(user *models.User) Modifier {
	m.ctx = context.WithValue(m.ctx, CurrentUser{}, user)
	return m
}

func (m *modifier) Context() context.Context {
	return m.ctx
}
