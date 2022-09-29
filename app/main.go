package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/ztrue/shutdown"

	iotDelivery "github.com/RacoWireless/iot-gw-stresser/implementation/_start/http"
	iotService "github.com/RacoWireless/iot-gw-stresser/implementation/service"
	iotUsecase "github.com/RacoWireless/iot-gw-stresser/implementation/usecase"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/spf13/viper"

	Log "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	lecho "github.com/ziflex/lecho/v3"
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Panic().Err(err).Msg("")
	}
	log.Info().Msg(`path: ` + path)
	viper.SetConfigType(`json`)
	viper.SetConfigName(`config`)
	viper.AddConfigPath(`./`)
	viper.AddConfigPath(`../`)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic().Err(err).Msg("config file not found")
		}
		time.Sleep(10 * time.Second)
		log.Panic().Err(err).Msg("config file not found")
	}

	if viper.GetBool(`debug`) {
		log.Info().Msg("MQTT Stresser Service RUN on DEBUG mode")
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// @title IOT Model Management API
// @version 1.0
// @description This is a Iot Device Management  server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.korewireless.com
// @contact.email support@korewireless.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host iot.korewireless.com
// @BasePath /
func main() {

	log.Info().Msg("Go Time")

	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvPrefix(viper.GetString("ENV"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	e := echo.New()
	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(Log.DEBUG),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Logger = logger
	e.Use(lecho.Middleware(lecho.Config{
		Logger: logger}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	//timeoutContext := time.Duration(viper.GetInt("CONTEXT.TIMEOUT")) * time.Second
	brokerUrl := viper.GetString("ENV_BROKER_URL")
	if brokerUrl == "" {
		log.Panic().Msg("Configuration Error: BROKER URL String not available")

	}
	StressService := iotService.NewStresserService(brokerUrl)
	Usecase := iotUsecase.NewUsecase(StressService)

	iotDelivery.NewIoTtHandler(e, Usecase)
	shutdown.Add(func() {
		log.Info().Msg("Stopping...")
		time.Sleep(time.Second)
		log.Info().Msg("Stopped")
	})
	go func() {
		log.Error().Err(e.Start(viper.GetString("ENV_AUTH_SERVER"))).Msg("")
	}()
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

}
