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

	res, err := h.DB.Query("SELECT * FROM table;")

	fmt.Println(res, err)
	if err != nil {
		// sentry.CaptureException(err)
		captureException(err)
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("Hello %s", name))
}

func (h *Handler) NewBuilding(c echo.Context) (err error) {
	building := new(model.Building)
	building.Name = "My Building - " + time.Now().Format("15-01-05")
	if affected, err := h.DB.Insert(building); err != nil {
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
