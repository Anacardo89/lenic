package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/auth"
	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/server"
	"github.com/Anacardo89/lenic/internal/server/httphandle/api"
	"github.com/Anacardo89/lenic/internal/server/httphandle/page"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/internal/wsconnman"
	"github.com/Anacardo89/lenic/pkg/img"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/Anacardo89/lenic/pkg/mail"
)

var (
	httpServer  = &http.Server{}
	httpsServer = &http.Server{}
)

func main() {
	// Dependencies
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logg := logger.NewLogger(cfg.Log)
	dbRepo, err := initDB(cfg.DB)
	if err != nil {
		logg.Error("failed to init db", "error", err)
		os.Exit(1)
	}
	defer dbRepo.Close()
	tokenMan := auth.NewTokenManager(&cfg.Token)
	sm := session.NewSessionManager(context.Background(), cfg.Session, dbRepo)
	mailClient := mail.NewClient(cfg.Mail)
	im := img.NewImgManager(&cfg.Img)

	wsh := wshandle.NewHandler(dbRepo, logg, sm, wsconnman.NewWSConnMan())
	ah := api.NewHandler(context.Background(), logg, &cfg.Server, dbRepo, tokenMan, sm, wsh, mailClient, im)
	ph := page.NewHandler(context.Background(), logg, dbRepo, sm)
	mw := middleware.NewMiddlewareHandler(sm, logg, cfg.Server.WriteTimeout)

	srv := server.NewServer(cfg.Server, logg, ah, ph, mw, wsh)

	// Serve
	stopChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		errChan <- srv.Start()
	}()
	select {
	case sig := <-stopChan:
		logg.Info("Shutting down...", "signal", sig)
		if err := srv.Shutdown(); err != nil {
			logg.Error("Failed to shutdown server gracefully", "error", err)
			os.Exit(1)
		}
		logg.Info("Server stopped gracefully")
	case err := <-errChan:
		if err != http.ErrServerClosed {
			logg.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}
}
