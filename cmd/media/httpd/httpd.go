package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.elastic.co/apm"

	"gitlab.com/slirx/newproj/internal/media"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/apmmiddleware"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/tracer"
	"gitlab.com/slirx/newproj/pkg/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("can not initialize logger: %v", err)
	}

	conf, err := NewConfig("MEDIA_HTTPD_")
	if err != nil {
		zapLogger.Fatal(err)
	}

	t := tracer.NewAPMTracer()
	responseBuilder := api.NewResponseBuilder(t)
	handler := media.NewHandler(zapLogger, responseBuilder)
	apmTracer := apm.DefaultTracer
	recoveryFunc := utils.NewRecoveryFunc(zapLogger, responseBuilder)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   conf.Server.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get(
		"/media",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Images,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/media",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Post(
		"/media",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.UploadImage,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/media",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/internal/media",
		apmmiddleware.Wrap(
			jwtmiddleware.WrapInternal(
				handler.Images,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.InternalSecrets,
			),
			"/internal/media",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web/images/media"))
	fileServer(router, "/media/image/*", filesDir)

	server := http.Server{
		Addr:    conf.Server.Addr,
		Handler: router,
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	zapLogger.Debug("started")

	<-ctx.Done()
	_ = server.Shutdown(ctx)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
