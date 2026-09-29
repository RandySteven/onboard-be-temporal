package servers

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RandySteven/onboard-be/apps"
	"github.com/RandySteven/onboard-be/configs"
	"github.com/RandySteven/onboard-be/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Server struct {
	Host string
	Port string

	ServerTimeout time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
}

func init() {
	if err := godotenv.Load("./files/env/.env"); err != nil {
		log.Println("no .env file loaded, using process environment")
	}
}

func getConfig() *configs.Config {
	configPath, err := configs.ParseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	config, err := configs.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}

func initServer(config *configs.Config) *Server {
	serverConfig := config.Configs.Server
	return &Server{
		Host:          serverConfig.Host,
		Port:          serverConfig.Port,
		ServerTimeout: serverConfig.Timeout.Server,
		ReadTimeout:   serverConfig.Timeout.Read,
		WriteTimeout:  serverConfig.Timeout.Write,
		IdleTimeout:   serverConfig.Timeout.Idle,
	}
}

func (s *Server) run(ctx context.Context, r *mux.Router) error {
	srv := &http.Server{
		Addr:         s.Host + ":" + s.Port,
		Handler:      r,
		ReadTimeout:  s.ReadTimeout * time.Second,
		WriteTimeout: s.WriteTimeout * time.Second,
		IdleTimeout:  s.IdleTimeout * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("Server exiting")
	return nil
}

func Serve() {
	config := getConfig()
	server := initServer(config)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := apps.NewApp(config)
	if err != nil {
		log.Fatalln("Error starting app ", err)
	}
	defer app.MySQL.Close()

	apis := app.PrepareHttpHandler(ctx)
	r := mux.NewRouter()
	router := routes.NewEndpointRouters(apis)
	routes.InitRouter(router, r)

	if err = app.Temporal.Start(); err != nil {
		log.Fatalln("Failed to start Temporal worker:", err)
	}
	defer app.Temporal.Stop()

	if err := server.run(ctx, r); err != nil {
		log.Fatal(err)
	}
}
