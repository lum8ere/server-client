package main

import (
	"backed-api-v2/libs/1_application/ws_server"
	"backed-api-v2/libs/2_domain_methods/handlers"
	"backed-api-v2/libs/2_domain_methods/handlers/auth"
	"backed-api-v2/libs/2_domain_methods/handlers/test_handlers"
	"backed-api-v2/libs/2_domain_methods/run_processor"
	"backed-api-v2/libs/5_common/rest_middleware"
	"backed-api-v2/libs/5_common/smart_context"
	"runtime"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func initRoutes(sctx smart_context.ISmartContext) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "X-Request-Id", "X-Session-Id", "X-Api-Key", "X-Auth-Provider"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/auth/login", rest_middleware.WithRestApiSmartContext(sctx, auth.LoginHandler))
	r.Get("/rnd2", run_processor.WrapRestApiSmartHandler(sctx, test_handlers.RndHandler2))
	r.Post("/send_command", run_processor.WrapRestApiSmartHandler(sctx, handlers.SendCommandHandler))

	// r.Post("/devices/register", rest_middleware.WithRestApiSmartContext(sctx, devices.RegisterDeviceHandler))
	// // Прием метрик
	// r.Post("/metrics", rest_middleware.WithRestApiSmartContext(sctx, metrics.ReceiveMetricsHandler))

	// pprof
	runtime.SetMutexProfileFraction(1)
	r.Mount("/debug", chi_middleware.Profiler())

	// Create a subrouter for other routes that require middleware
	protectedRoutes := chi.NewRouter()

	wsupgrader := ws_server.NewWsUpgrader(sctx)

	// Register the WebSocket handler
	protectedRoutes.Get("/ws", rest_middleware.WithWsApiSmartContext(sctx, wsupgrader.HandleWebSocket))

	protectedRoutes.NotFound(wsupgrader.HandleNotFound)

	// Mount the subrouter on the main router
	r.Mount("/", protectedRoutes)

	return r, nil
}
