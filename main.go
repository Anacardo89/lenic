package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog.git/api"
	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/db"
	"github.com/Anacardo89/tpsi25_blog.git/fsops"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
	"github.com/Anacardo89/tpsi25_blog.git/rabbit"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	templates   = template.Must(template.ParseGlob("templates/*"))
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
	sessConfig, err := loadSessionConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}
	auth.SessionStore = sessions.NewCookieStore([]byte(sessConfig.Pass))

	// RabbitMQ
	rabbit.RabbitMQ, err = loadRabbitConfig()
	if err != nil {
		logger.Error.Fatal(err)
	}

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/", RedirIndex).Schemes("https")
	r.HandleFunc("/home", Index).Schemes("https")
	r.HandleFunc("/login", Login).Schemes("https")
	r.HandleFunc("/register", Register).Schemes("https")
	r.HandleFunc("/error", Error).Schemes("https")
	r.HandleFunc("/newPost", NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid:[0-9a-zA-Z\\-=]+}", Post).Schemes("https")

	r.HandleFunc("/api/register", api.RegisterPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/login", api.LoginPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/logout", api.LogoutPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/post", api.PostPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/post/{post_guid:[0-9a-zA-Z\\-=]+}/comment", api.CommentPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/post/{post_guid:[0-9a-zA-Z\\-=]+}/comment/{comment_id:[0-9]+}", api.CommentPUT).Methods("PUT").Schemes("https")

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
