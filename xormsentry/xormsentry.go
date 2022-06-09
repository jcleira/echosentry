package xormsentry

// Code based on: https://github.com/mongofs/pkg/blob/32e850fd22afb5048fb4144e9e5c32550db341a7/log/jeager/xormHook.go

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"xorm.io/builder"
	"xorm.io/xorm/contexts"
)

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

type xormHookSpan struct{}

var xormHookSpanKey = &xormHookSpan{}

func before(c *contexts.ContextHook) (context.Context, error) {

	fmt.Println("Before SPAN")

	transaction := sentry.TransactionFromContext(c.Ctx)

	fmt.Println("TRANSACTION:", transaction)

	span := sentry.StartSpan(c.Ctx, "db.query", sentry.TransactionName("/buildings"))

	c.Ctx = context.WithValue(c.Ctx, "xormspan", *span)
	return c.Ctx, nil
}

func after(c *contexts.ContextHook) error {
	fmt.Println("After SPAN:", c.SQL, c.ExecuteTime)
	span, ok := c.Ctx.Value("xormspan").(sentry.Span)
	if !ok {
		return nil
	}
	defer span.Finish()

	if c.Err != nil {

		//span.LogFields(log.Object("errors", c.Err))
		fmt.Println("HOOK errors:", c.Err)
	}
	sql, _ := builder.ConvertToBoundSQL(c.SQL, c.Args)

	fmt.Println("HOOK SQL:", sql)
	fmt.Println("HOOK Args:", c.Args)
	span.Data = map[string]interface{}{
		"sql":          sql,
		"execute_time": c.ExecuteTime.Microseconds(),
	}
	span.SetTag("execute_time", string(c.ExecuteTime.Milliseconds()))

	return nil
}

func NewTracingHook() *TracingHook {
	return &TracingHook{
		before: before,
		after:  after,
	}
}
