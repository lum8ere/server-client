package main

import (
	"backed-api-v2/libs/1_application/service_helper"
	"backed-api-v2/libs/5_common/smart_context"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	var allRoutes *chi.Mux
	var webServer *http.Server
	service_helper.StartService("backend-api-v2",
		func(sctx smart_context.ISmartContext) error {
			r, err := initRoutes(sctx)
			if err != nil {
				return err
			}
			allRoutes = r

			return nil
		},
		func(sctx smart_context.ISmartContext) error {
			webServer = &http.Server{
				Addr:    ":9000",
				Handler: allRoutes,
			}
			sctx.Info("Server listening on port 9000")
			go func() {
				if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					sctx.Fatalf("Server error: %v", err)
				}
			}()

			return nil
		},
		func(sctx smart_context.ISmartContext) error {
			// Gracefully shut down the server
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second) // тут уже надо мало ждать
			defer shutdownCancel()
			if err := webServer.Shutdown(shutdownCtx); err != nil {
				sctx.Errorf("Exiting process: Error shutting down server: %v", err)
			} else {
				sctx.Info("Exiting process: Server shut down")
			}
			return nil
		},
	)
}