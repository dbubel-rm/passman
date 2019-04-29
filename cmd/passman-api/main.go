package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbubel/passman/internal/mid"

	"github.com/dbubel/passman/cmd/passman-api/handlers"
	"github.com/dbubel/passman/internal/platform/db"
	"github.com/kelseyhightower/envconfig"
)

var (
	BUILD = "develop-9"
	GIT_HASH =""
	BUILD_DATE = ""

)

// TODO: hash of the service name in the DB
func main() {
	handlers.Build = BUILD
	handlers.GitHash = GIT_HASH
	handlers.BuildDate = BUILD_DATE

	log := log.New(os.Stdout, "", log.LstdFlags|log.Ltime|log.Lshortfile)

	var cfg struct {
		Web struct {
			APIHost         string        `default:":3000" envconfig:"PORT"`
			DebugHost       string        `default:":4000" envconfig:"DEBUG_HOST"`
			ReadTimeout     time.Duration `default:"5s" envconfig:"READ_TIMEOUT"`
			WriteTimeout    time.Duration `default:"5s" envconfig:"WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
			EnableTLS       string        `default:"no" envconfig:"ENABLE_TLS"`
		}
		DB struct {
			Host     string `default:"host.docker.internal" envconfig:"MYSQL_ENDPOINT"`
			Username string `default:"root" envconfig:"MYSQL_USERNAME"`
			Password string `default:"" envconfig:"MYSQL_PASSWORD"`
			Database string `default:"passman" envconfig:"MYSQL_DB"`
		}
	}

	if err := envconfig.Process("PASSMAN", &cfg); err != nil {
		log.Fatalf("Parsing Config : %v", err)
	}
	// cfgJSON, err := json.MarshalIndent(cfg, "", "    ")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	log.Println("Passman server starting", BUILD)
	// =========================================================================
	// Start MySQL

	log.Println("Initialize MySQL...")

	var masterDB *db.MySQLDB
	var err error
	for i := 0; i < 30; i++ {
		masterDB, err = db.New(connStr)
		if err != nil {
			log.Printf("main : Register DB : %s\n", err.Error())
		} else {
			log.Println("DB connect OK")
			break
		}
		time.Sleep(time.Second)
	}
	defer masterDB.Close()

	// =========================================================================
	// Start API
	api := http.Server{
		Addr:           cfg.Web.APIHost,
		Handler:        handlers.API(log, masterDB, mid.AuthHandler),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    cfg.Web.ReadTimeout,
		WriteTimeout:   cfg.Web.WriteTimeout,
	}
	defer log.Println("App Shutdown")

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("API Listening %s", cfg.Web.APIHost)
		if cfg.Web.EnableTLS == "no" {
			serverErrors <- api.ListenAndServe()
		} else {
			fmt.Println("TLS ON")
			serverErrors <- api.ListenAndServeTLS("MyCertificate.crt", "MyKey.key")
		}
	}()

	// =========================================================================
	// Start Debug Service
	// /debug/vars - Added to the default mux by the expvars package.
	// /debug/pprof - Added to the default mux by the net/http/pprof package.

	// debug := http.Server{
	// 	Addr:           cfg.Web.DebugHost,
	// 	Handler:        http.DefaultServeMux,
	// 	ReadTimeout:    cfg.Web.ReadTimeout,
	// 	WriteTimeout:   cfg.Web.WriteTimeout,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// // Not concerned with shutting this down when the
	// // application is being shutdown.
	// go func() {
	// 	log.Printf("Debug Listening %s", cfg.Web.DebugHost)
	// 	log.Printf("Debug Listener closed : %v", debug.ListenAndServe())
	// }()

	// Shutdown

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Stop API Service
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-osSignals:
		log.Println("Start shutdown...")

		// Create context for Shutdown call.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		if err := api.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Fatalf("Could not stop http server: %v", err)
			}
		}
	}
}
