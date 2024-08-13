package main

import (
	"log"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/config"
	"github.com/Anacardo89/tpsi25_blog/internal/pages"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
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
	r.HandleFunc("/", pages.RedirIndex).Schemes("https")
	r.HandleFunc("/home", pages.Index).Schemes("https")
	r.HandleFunc("/login", pages.Login).Schemes("https")
	r.HandleFunc("/register", pages.Register).Schemes("https")
	r.HandleFunc("/activate/{user_name}", pages.ActivateUser).Schemes("https")
	r.HandleFunc("/error", pages.Error).Schemes("https")
	r.HandleFunc("/newPost", pages.NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid}", pages.Post).Schemes("https")
	r.HandleFunc("/api/image", pages.ServeImage).Schemes("https")
	r.HandleFunc("/forgot-password", pages.ServeForgotPassword).Schemes("https")

	r.HandleFunc("/api/register", pages.RegisterPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/login", pages.LoginPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/logout", pages.LogoutPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/post", pages.PostPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/post/{post_guid}/comment", pages.CommentPOST).Methods("POST").Schemes("https")
	r.HandleFunc("/api/post/{post_guid}/comment/{comment_id}", pages.CommentPUT).Methods("PUT").Schemes("https")
	r.HandleFunc("/api/forgot-password", pages.ForgotPassword).Methods("POST").Schemes("https")

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
