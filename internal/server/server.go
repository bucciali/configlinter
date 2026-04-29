package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"configlinter/internal/config"
	"configlinter/internal/engine"
	"configlinter/internal/grpcserver"
	"configlinter/internal/parser"
	"configlinter/internal/server/handlers"
	"configlinter/internal/server/router"
	pb "configlinter/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Start(cfg *config.ServerConfig, reg *parser.Registry, eng *engine.Engine) error {
	h := handlers.New(reg, eng)
	r := router.New(h)

	httpAddr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         httpAddr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on %s", httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listen: %v", err)
		}
	}()

	grpcAddr := ":" + cfg.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterConfigLinterServer(grpcSrv, grpcserver.New(reg, eng))
	reflection.Register(grpcSrv)

	go func() {
		log.Printf("gRPC server listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	grpcSrv.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}

	log.Println("All servers stopped gracefully")
	return nil
}
