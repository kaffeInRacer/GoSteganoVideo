package server

import (
	"context"
	"errors"
	"fmt"
	"kaffein/cmd/web/routes"
	"kaffein/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultIdleTimeout = time.Minute // Waktu timeout ketika server tidak menerima permintaan
	//defaultReadTimeout    = 5 * time.Second  // Waktu timeout untuk membaca data dari client
	defaultReadTimeout    = 30 * time.Minute // Waktu timeout untuk membaca data dari client
	defaultWriteTimeout   = 2 * time.Minute  // Waktu timeout untuk menulis data ke client
	defaultShutdownPeriod = 30 * time.Second // Waktu timeout untuk proses shutdown server
)

type Application struct {
	Config config.Config
	Logger *slog.Logger
	Wg     sync.WaitGroup
}

func (app *Application) ServeHTTP() error {

	addr := fmt.Sprintf("%s:%s", app.Config.Host, app.Config.Port)
	srv := &http.Server{
		Addr:         addr,                                                    // Alamat server
		Handler:      routes.Routes(app.Logger),                               // Handler dari route, dengan logger yang diteruskan
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelWarn), // Logger untuk mencatat kesalahan
		IdleTimeout:  defaultIdleTimeout,                                      // Timeout untuk koneksi idle
		ReadTimeout:  defaultReadTimeout,                                      // Timeout untuk membaca data
		WriteTimeout: defaultWriteTimeout,                                     // Timeout untuk menulis data
	}

	shutdownErrorChan := make(chan error)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	app.Logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	app.Logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))

	app.Wg.Wait()
	return nil
}
