package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/handlers"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error in main")

			os.Exit(1)
		}
	}()

	// cfg is the struct type that contains fields that stores the necessary configuration
	// gathered from the environment.
	var cfg struct {
		DaemonPort int `envconfig:"DAEMON_PORT" default:"3000"`

		DBUser string `envconfig:"DB_USER" default:"root"`
		DBPass string `envconfig:"DB_PASS" default:"root"`
		DBName string `envconfig:"DB_NAME" default:"list"`
		DBHost string `envconfig:"DB_HOST" default:"db"`
		DBPort int    `envconfig:"DB_PORT" default:"5432"`

		ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
	}
	if err := envconfig.Process("LIST", &cfg); err != nil {
		err = errors.Wrap(err, "parse environment variables")
		return
	}

	dbCfg := db.Config{
		User: cfg.DBUser,
		Pass: cfg.DBPass,
		Name: cfg.DBName,
		Host: cfg.DBHost,
		Port: cfg.DBPort,
	}
	dbc, err := db.NewConnection(dbCfg)
	if err != nil {
		err = errors.Wrap(err, "connect to postgres db")
		return
	}

	defer func(){
		if err := dbc.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	server := http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.DaemonPort),
		Handler:        handlers.NewApplication(dbc),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Start listening for requests made to the daemon and create a channel
	// to collect non-HTTP related server errors on.
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server started, listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown of the daemon.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Waiting for an osSignal or a non-HTTP related server error.
	select {
	case e := <-serverErrors:
		err = fmt.Errorf("server failed to start: %+v", e)
		return

	case <-osSignals:
	}

	// Gracefully shutdown server once an exit signal or error is received.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", cfg.ShutdownTimeout, err)

		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}
}
