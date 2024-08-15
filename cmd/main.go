package main

import (
	"log"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/config"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/pages"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
	"github.com/Anacardo89/tpsi25_blog/internal/routes"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	httpServer  = &http.Server{}
	httpsServer = &http.Server{}
)

func main() {
	logger.CreateLogger()

	// DB
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}
	db.Dbase, err = db.LoginDB(dbConfig)
	if err != nil {
		logger.Error.Fatal(err)
	}

	// Certificate
	cert, err := fsops.MakePaths()
	if err != nil {
		logger.Error.Fatal(err)
	}
	tlsConf, err := fsops.LoadCertificates(cert)
	if err != nil {
		logger.Error.Fatal(err)
	}

	// Session Store
	sessConfig, err := config.LoadSessionConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}
	auth.SessionStore = sessions.NewCookieStore([]byte(sessConfig.Pass))

	// RabbitMQ
	rabbit.RabbitMQ, err = config.LoadRabbitConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}

	// Router
	r := mux.NewRouter()
	routes.DeclareRoutes(r)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

	// Server
	serverConfig, err := config.LoadServerConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}

	httpServer = &http.Server{
		Addr:    serverConfig.HttpPORT,
		Handler: http.HandlerFunc(pages.RedirectNonSecure),
	}

	httpsServer = &http.Server{
		Addr:      serverConfig.HttpsPORT,
		TLSConfig: tlsConf,
	}

	// Work
	errChan := make(chan error, 2)

	go func() {
		log.Println("Starting HTTP server on :8081")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	go func() {
		log.Println("Starting HTTPS server on :8082")
		if err := httpsServer.ListenAndServeTLS(cert.CertPath, cert.KeyPath); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	logger.Error.Fatalf("Server error: %v", <-errChan)

}
