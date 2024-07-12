package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog.git/db"
	"github.com/Anacardo89/tpsi25_blog.git/fsops"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
	"github.com/gorilla/mux"
)

var (
	templates   = template.Must(template.ParseGlob("templates/*"))
	dbase       *sql.DB
	httpServer  = &http.Server{}
	httpsServer = &http.Server{}
)

func main() {
	logger.CreateLogger()

	// DB
	dbConfig, err := loadDBConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}
	dbase, err = db.LoginDB(dbConfig)
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

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/", RedirIndex).Schemes("https")
	r.HandleFunc("/home", ServeIndex).Schemes("https")
	r.HandleFunc("/login", ServeLogin).Schemes("https")
	r.HandleFunc("/register", ServeRegister).Schemes("https")

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Server
	serverConfig, err := loadServerConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}

	httpServer = &http.Server{
		Addr:    serverConfig.HttpPORT,
		Handler: http.HandlerFunc(RedirectNonSecure),
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
