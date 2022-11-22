package handler

import (
	"EchoSentry/model"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (h *Handler) Health(c echo.Context) (err error) {
	return c.String(http.StatusOK, "OK")
}

func (h *Handler) Hello(c echo.Context) (err error) {
	name := c.QueryParam("name")

	res, err := h.DB.Query("SELECT * FROM buildings;")

	fmt.Println(res, err)
	if err != nil {
		// sentry.CaptureException(err)
		captureException(err)
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("Hello %s", name))
}

func (h *Handler) NewBuilding(c echo.Context) (err error) {
	// Now on every layer we might start its own Span
	span := sentry.StartSpan(c.Request().Context(), "handler")
	defer span.Finish()

	building := new(model.Building)
	building.Name = "My Building - " + time.Now().Format("15-01-05")

	// When there is a method where we don't have control, we could use a closure
	//
	// There is no need to use a closure but it's convenient because it allows us
	// to use defer, but we could also just call StartSpan/Finish without the
	// closure.
	//
	// But let's say that the Insert method is something "ours", we need to send
	// c.Request.Context() down as a context.Context, so that method could make
	// their own sentry.StartSpan calls.
	buildingInsert := func(building *model.Building) (int64, error) {
		span := sentry.StartSpan(c.Request().Context(), "DB insert")
		defer span.Finish()

		db := h.DB.Context(c.Request().Context())
		if affected, err := db.Insert(building); err != nil {
			return affected, err
		} else {
			return affected, nil
		}
	}

	if affected, err := buildingInsert(building); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(affected)
	}

	return c.JSON(http.StatusOK, building)
}

func (h *Handler) ListBuildings(c echo.Context) (err error) {
	myContext := c.Request().Context()

	span := sentry.StartSpan(c.Request().Context(), "handler")

	session := h.DB.Context(myContext)

	buildings := make([]model.Building, 0)

	if err := session.Find(&buildings); err != nil {
		fmt.Println("There has been an error", err)
	}
	span.Finish()
	return c.JSON(http.StatusOK, buildings)
}
