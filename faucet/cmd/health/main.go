package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/gorilla/handlers"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/filecoin-project/faucet/internal/failure"
	app "github.com/filecoin-project/faucet/internal/http"
	"github.com/filecoin-project/faucet/internal/platform/lotus"
	"github.com/filecoin-project/faucet/pkg/version"
	"github.com/filecoin-project/lotus/api/client"
)

func main() {
	logger := logging.Logger("CALIBRATION-HEALTH")

	lvl, err := logging.LevelFromString("info")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	if err := run(logger); err != nil {
		logger.Fatalln("main: error:", err)
	}
}

func run(log *logging.ZapEventLogger) error {
	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:60s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			Host            string        `conf:"default:0.0.0.0:9000"`
		}
		TLS struct {
			Disable  bool   `conf:"default:true"`
			CertFile string `conf:"default:nocert.pem"`
			KeyFile  string `conf:"default:nokey.pem"`
		}
		Lotus struct {
			APIHost   string `conf:"default:127.0.0.1:1230"`
			AuthToken string
		}
	}{
		Version: conf.Version{
			Build: version.Version(),
			Desc:  "Calibration Health Service",
		},
	}

	help, err := conf.Parse("HEALTH", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return err
	}

	// =========================================================================
	// App Starting

	ctx := context.Background()

	log.Infow("starting service", "version", version.Version())
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	expvar.NewString("build").Set(version.Version())

	// =========================================================================
	// Initialize authentication support

	log.Infow("startup", "status", "initializing authentication support")

	var authToken string

	if cfg.Lotus.AuthToken == "" {
		authToken, err = lotus.GetToken()
		if err != nil {
			return fmt.Errorf("error getting authentication token: %w", err)
		}
	} else {
		authToken = cfg.Lotus.AuthToken
	}
	header := http.Header{"Authorization": []string{"Bearer " + authToken}}

	// =========================================================================
	// Start Lotus client

	log.Infow("startup", "status", "initializing Lotus support", "host", cfg.Lotus.APIHost)

	lotusClient, lotusCloser, err := client.NewFullNodeRPCV1(ctx, "ws://"+cfg.Lotus.APIHost+"/rpc/v1", header)
	if err != nil {
		return fmt.Errorf("connecting to Lotus failed: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping Lotus client support")
		lotusCloser()
	}()

	log.Infow("Successfully connected to Lotus node")

	// =========================================================================
	// Start Detector Service

	d := failure.NewDetector(log, lotusClient, time.Minute, 3*time.Minute)

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing HTTP API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      handlers.RecoveryHandler()(app.HealthHandler(log, lotusClient, d, version.Version())),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		switch cfg.TLS.Disable {
		case true:
			serverErrors <- api.ListenAndServe()
		case false:
			serverErrors <- api.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		}
	}()

	// =========================================================================
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		d.Stop()

		if err := api.Shutdown(ctx); err != nil {
			api.Close() // nolint
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
