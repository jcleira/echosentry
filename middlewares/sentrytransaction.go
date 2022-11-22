package middlewares

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// SentryTransaction send transaction into sentry.
func SentryTransaction() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// This call is:
			// - Starting a transaction.
			// - Starting a initial spam.
			// - Writing both in the existing echo.Context (within Request's context)
			span := sentry.StartSpan(ctx.Request().Context(), "http",
				sentry.TransactionName(ctx.Request().URL.Path),
			)
			defer span.Finish()
			return next(ctx)
		}
	}
}
