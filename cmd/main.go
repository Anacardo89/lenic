package main

import (
	"log"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/config"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/routes"
	"github.com/Anacardo89/tpsi25_blog/internal/server"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/rabbitmq"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	httpServer  = &http.Server{}
	httpsServer = &http.Server{}
)

func main() {
	logger.CreateLogger()
	logger.Info.Println("System start")
	fsops.MakeImgDir()

	// DB
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		logger.Error.Fatalln("Could not load dbConfig:", err)
	}
	db.Dbase, err = db.LoginDB(dbConfig)
	if err != nil {
		logger.Error.Fatalln("Could not connect to DB: ", err)
	}
	orm.Da.Db = db.Dbase
	logger.Info.Println("Connecting to DB OK")

	// Certificate
	cert := fsops.MakePaths()
	tlsConf, err := fsops.LoadCertificates(cert)
	if err != nil {
		logger.Error.Fatalln("Could not load SSL Certificates:", err)
	}
	logger.Info.Println("Loading SSL Certificates OK")

	// Session Store
	sessConfig, err := config.LoadSessionConfig()
	if err != nil {
		logger.Error.Fatalln("Could not load sessConfig:", err)
	}
	auth.SessionStore = sessions.NewCookieStore([]byte(sessConfig.Pass))
	logger.Info.Println("Creating SessionStore OK")

	// RabbitMQ
	rabbitmq.RMQ, err = config.LoadRabbitConfig()
	if err != nil {
		logger.Error.Fatalln("Could not load rabbitConfig:", err)
	}
	rconn, err := rabbitmq.RMQ.Connect()
	if err != nil {
		logger.Error.Fatalln("Could not connect to RabbitMQ:", err)
	}
	rabbitmq.RCh, err = rconn.Channel()
	if err != nil {
		logger.Error.Fatalln("Could not create Rabbit Channel:", err)
	}
	defer rabbitmq.RCh.Close()
	logger.Info.Println("Connecting to RabbitMQ OK")

	err = rabbitmq.RMQ.DeclareQueues(rabbitmq.RCh)
	if err != nil {
		logger.Error.Fatalln("Could not declare Rabbit Queue:", err)
	}
	logger.Info.Println("Declaring Rabbit Queues OK")

	// Router
	r := mux.NewRouter()
	routes.DeclareRoutes(r)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Websocket
	wsocket.WSConnMan = wsocket.NewWSConnManager()

	// Server
	server.Server, err = config.LoadServerConfig()
	if err != nil {
		logger.Error.Fatalln("Could not load serverConfig:", err)
	}
	logger.Info.Println("Loading serverConfig OK")

	httpServer = &http.Server{
		Addr:    ":" + server.Server.HttpPORT,
		Handler: http.HandlerFunc(redirect.RedirectNonSecure),
	}

	httpsServer = &http.Server{
		Addr:      ":" + server.Server.HttpsPORT,
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
