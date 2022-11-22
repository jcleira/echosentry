package main

// https://github.com/fonysaputra/go-utils/blob/9b79ea7f79a7b9f63fa005d908b6050b0c1b0a17/apm/sentry/sentry.go
// https://github.com/hlmn/senyum-go-utils/blob/12befb105256991a52b77ff1c6e428b8e3de88f3/apm/sentry/sentry.go
// https://github.com/maiaaraujo5/gostart/blob/5bdf031f72d53fae97c7a263a2e9623e7696da58/rest/handler/wrapper/wrapper.go

import (
	"EchoSentry/handler"
	"EchoSentry/middlewares"
	"EchoSentry/model"
	"EchoSentry/xormsentry"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"xorm.io/xorm"

	sentryecho "github.com/getsentry/sentry-go/echo"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	sentry.CaptureException(err)

	fmt.Println("Capturing:", code)
	//errorPage := fmt.Sprintf("%d.html", code)
	// if err := c.File(errorPage); err != nil {
	//	c.Logger().Error(err)
	// }
	err = c.JSON(code, "Undefined error")
}

func main() {
	config := loadConfig()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDsn,
		TracesSampleRate: 1.0,
		Environment:      "local",
		Release:          "1.0",
		Debug:            true,
		AttachStacktrace: true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryecho.Options{
		Repanic:         true,
		Timeout:         1 * time.Second,
		WaitForDelivery: false,
	}))

	e.Use(middlewares.SentryTransaction())

	db, err := xorm.NewEngine("sqlite3", config.Database)
	if err != nil {
		e.Logger.Fatal(err)
	}

	db.AddHook(xormsentry.NewTracingHook())

	if err := db.Sync2(new(model.Building)); err != nil {
		fmt.Println(err)
	}

	h := &handler.Handler{DB: db}

	e.GET("/buildings/new", h.NewBuilding)
	e.GET("/buildings", h.ListBuildings)
	e.GET("/health", h.Health)
	e.GET("/hello", h.Hello)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/health")
	})

	e.HTTPErrorHandler = customHTTPErrorHandler

	sentry.CaptureMessage("Starting APP")

	// address := ":" + strconv.Itoa(config.Port)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}

type Config struct {
	Port      int
	Database  string
	SentryDsn string
}

func loadConfig() *Config {
	cfg := &Config{
		Port:      1234,
		Database:  "myfile.db",
		SentryDsn: "https://565ab10db08448289861fe107cb0867b@o913183.ingest.sentry.io/6469771",
	}

	viper.SetConfigName("front-test")
	viper.SetEnvPrefix("ft")

	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE")

	//Flags
	viper.AutomaticEnv()

	viper.ReadInConfig()

	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Printf("cannot unmarshal config: %v\n", err)
	}

	return cfg
}
