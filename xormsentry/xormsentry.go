package xormsentry

// Code based on: https://github.com/mongofs/pkg/blob/32e850fd22afb5048fb4144e9e5c32550db341a7/log/jeager/xormHook.go

import (
	"context"

	"github.com/getsentry/sentry-go"
	"xorm.io/builder"
	"xorm.io/xorm/contexts"
)

func NewTracingHook() *TracingHook {
	return &TracingHook{
		before: before,
		after:  after,
	}
}

type TracingHook struct {
	before func(c *contexts.ContextHook) (context.Context, error)
	after  func(c *contexts.ContextHook) error
}

func (h *TracingHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return h.before(c)
}

func (h *TracingHook) AfterProcess(c *contexts.ContextHook) error {
	return h.after(c)
}

var _ contexts.Hook = &TracingHook{}

func before(c *contexts.ContextHook) (context.Context, error) {
	span := sentry.StartSpan(c.Ctx, "xorm.query")
	c.Ctx = context.WithValue(c.Ctx, "xormspan", *span)
	return c.Ctx, nil
}

func after(c *contexts.ContextHook) error {
	span, ok := c.Ctx.Value("xormspan").(sentry.Span)
	if !ok {
		return nil
	}
	defer span.Finish()

	sql, _ := builder.ConvertToBoundSQL(c.SQL, c.Args)
	span.Data = map[string]interface{}{
		"sql":          sql,
		"execute_time": c.ExecuteTime.Microseconds(),
	}

	return nil
}
