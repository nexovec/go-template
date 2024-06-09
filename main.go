package main

import (
	"app"
	"configuration"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
)

func main() {
	var err error
	// validating the environment variables
	{
		lo.Must0(godotenv.Load("./configs/envfiles/.app.dev.env"))
		lo.Must(configuration.GetDeploymentConfig())
	}

	// initialize submodules
	environ := lo.Must(configuration.GetDeploymentConfig())
	cfg := lo.Must(configuration.GetGeneralConfig())
	lo.Must0(utils.Initialize())

	// initialize default logger
	{
		logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		slogger := slog.New(logHandler)
		slog.SetDefault(slogger)
	}

	// create the server
	rootHandler := fiber.New(fiber.Config{
		AppName: lo.Must(configuration.GetDeploymentConfig()).AppName,
	})

	// initialize middlewares
	{
		// recover after panic
		rootHandler.Use(recover.New())

		// logging
		// PERFORMANCE: this code is for fiber's default logger, but it's making the whole server ~70% slower
		rootHandler.Use(logger.New(logger.Config{
			TimeFormat: time.RFC3339Nano,
			TimeZone:   "UTC",
			Done: func(c fiber.Ctx, logString []byte) {
				if c.Response().StatusCode() != fiber.StatusOK {
					slog.Warn(string(logString))
				}
			},
		}))

		// TODO: rate limiting

		// security
		rootHandler.Use(helmet.New())

		// profiling
		if configuration.EnumDeploymentProd != lo.Must(configuration.GetDeploymentConfig()).Deployment {
			rootHandler.Use(pprof.New(pprof.Config{}))
		}
	}

	// log current configuration
	{
		// get the git commit id
		dir := lo.Must(os.Getwd())
		cmd := exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = dir
		output, err := cmd.Output()
		if err != nil {
			output = []byte("unknown")
		}
		commitID := strings.TrimSpace(string(output))
		slog.Info("Launching application server", "Git commit ID", commitID)
	}

	// set up routes
	{
		rootHandler.Get("/favicon.ico", func(c fiber.Ctx) error {
			return c.Redirect().To("/static/favicon.ico")
		})
		rootHandler.Get("/favicon.ico", func(c fiber.Ctx) error { return c.Redirect().To("https://bulma.io/images/bulma-logo.png") })

		rootHandler.Get("/health", func(c fiber.Ctx) error {
			return c.Redirect().To("/healthcheck")
		})
		rootHandler.All("/healthcheck", func(c fiber.Ctx) error {
			return c.SendString("ALIVE")
		})

		rootHandler.Get("/security.txt", func(c fiber.Ctx) error {
			return c.Redirect().To("/.well-known/security.txt")
		})
		rootHandler.Get("/.well-known/security.txt", func(c fiber.Ctx) error {
			return c.SendFile("./static/security.txt")
		})
		rootHandler.Get("/robots.txt", func(c fiber.Ctx) error {
			return c.SendFile("./static/robots.txt")
		})

		rootHandler.Static("/static", "./static")

		app.RegisterRoutes(rootHandler)

		// default error pages
		rootHandler.All("*", func(c fiber.Ctx) error {
			// TODO:
			return c.SendStatus(fiber.StatusNotImplemented)
		})
	}

	// launch server
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() { // start the server with retry
		retry := 0
		retry_delay := 1 * time.Second
		err := configuration.ErrServerRefusedToStart
		for err != nil && retry < cfg.MaxServerStartRetries {
			host := environ.AppHost
			port := environ.AppPort
			slog.Warn("Server start", "host", host, "port", port)
			err = rootHandler.Listen(fmt.Sprintf("%s:%s", host, port), fiber.ListenConfig{
				EnablePrefork: false,
			})
			if err == nil {
				break
			}
			retry++
			slog.Error("Failed to start the server", "error", err, "retrying in", retry_delay)
			retry_delay *= 2
			retry_delay = max(retry_delay, 10*time.Second)
			time.Sleep(retry_delay)
		}
		slog.Error("Failed to start the server", "error", err)
		os.Exit(1)
	}()

	// do cleanup
	slog.Info("received signal", "signal", <-sigc)
	gracefulShutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()
	err = rootHandler.ShutdownWithContext(gracefulShutdownCtx)
	if err != nil {
		slog.Error("Failed to shutdown the server", "error", err)
	}
	slog.Info("Server is gracefully shutdown")
	os.Exit(0)
}
