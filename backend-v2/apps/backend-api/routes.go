package main

import (
	"backed-api-v2/libs/1_application/ws_server"
	"backed-api-v2/libs/2_domain_methods/handlers"
	"backed-api-v2/libs/2_domain_methods/handlers/applications"
	"backed-api-v2/libs/2_domain_methods/handlers/auth"
	"backed-api-v2/libs/2_domain_methods/handlers/devices"
	"backed-api-v2/libs/2_domain_methods/handlers/dicts"
	"backed-api-v2/libs/2_domain_methods/handlers/metrics"
	"backed-api-v2/libs/2_domain_methods/handlers/test_handlers"
	"backed-api-v2/libs/2_domain_methods/handlers/users"
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

	r.Get("/rnd2", run_processor.WrapRestApiSmartHandler(sctx, test_handlers.RndHandler2))

	// Запрос для обработки команд
	r.Post("/send_command", rest_middleware.RoleMiddleware("ADMIN",
		run_processor.WrapRestApiSmartHandler(sctx, handlers.SendCommandHandler)))

	// запросы для фронта
	r.Get("/api/dicts/roles", run_processor.WrapRestApiSmartHandler(sctx, dicts.GetRoleDictsHandler))

	r.Get("/api/devices", run_processor.WrapRestApiSmartHandler(sctx, devices.GetDevicesHandler))
	r.Get("/api/users", rest_middleware.RoleMiddleware("ADMIN",
		run_processor.WrapRestApiSmartHandler(sctx, users.GetUsersHandler)))
	r.Get("/api/devices/{id}", run_processor.WrapRestApiSmartHandler(sctx, devices.GetDevicesByIDHandler))
	r.Get("/api/metrics", run_processor.WrapRestApiSmartHandler(sctx, metrics.GetMetricsHandler))
	r.Get("/api/metrics/{id}", run_processor.WrapRestApiSmartHandler(sctx, metrics.GetMetricsByDeviceIDHandler))         // тут id это id девайса
	r.Get("/api/apps/{id}", run_processor.WrapRestApiSmartHandler(sctx, applications.GetApplicationsByDevicesIDHandler)) // тут id это id девайса

	// запросы на регистрацию и авторизацию
	r.Post("/api/auth/register", rest_middleware.RoleMiddleware("ADMIN",
		run_processor.WrapRestApiSmartHandler(sctx, auth.RegisterHandler)))
	r.Post("/api/auth/login", run_processor.WrapRestApiSmartHandler(sctx, auth.LoginHandler))
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
