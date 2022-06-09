package handler

import (
	"github.com/getsentry/sentry-go"
	"xorm.io/xorm"
)

type (
	Handler struct {
		DB *xorm.Engine
		// DataFiles embed.FS
	}
)

func captureException(err error) {
	sentry.CaptureException(err)
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
