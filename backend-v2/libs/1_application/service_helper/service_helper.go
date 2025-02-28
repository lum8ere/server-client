package service_helper

import (
	"backed-api-v2/libs/4_infrastructure/db_manager"
	"backed-api-v2/libs/5_common/env_vars"
	"backed-api-v2/libs/5_common/shutdown"
	"backed-api-v2/libs/5_common/smart_context"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func StartService(serviceName string,
	initFunc func(sctx smart_context.ISmartContext) error,
	startFunc func(sctx smart_context.ISmartContext) error,
	closeFunc func(sctx smart_context.ISmartContext) error,
) {
	env_vars.LoadEnvVars()
	os.Setenv("LOG_LEVEL", "debug") // ставим наиболее детальный уровень чтобы далее уже юзерские настройки работали
	sctx := smart_context.NewSmartContext()
	err := internalStartService(sctx, serviceName, initFunc, startFunc, closeFunc)
	if err != nil {
		sctx.Fatalf("Error starting service: %v", err)
	}
}

func internalStartService(
	sctx smart_context.ISmartContext,
	serviceName string,
	initFunc func(sctx smart_context.ISmartContext) error,
	startFunc func(sctx smart_context.ISmartContext) error,
	closeFunc func(sctx smart_context.ISmartContext) error,
) error {
	prefix := fmt.Sprintf("Service '%s'", serviceName)
	sctx.Infof("%s: Initializing", prefix)
	defer sctx.Infof("%s: Exited", prefix)

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sctx = sctx.WithContext(ctx)
	wg := &sync.WaitGroup{}
	sctx = sctx.WithWaitGroup(wg)

	location, err := time.LoadLocation("Local")
	if err != nil {
		return fmt.Errorf("error loading local timezone: %v", err)
	}
	sctx.Infof("Local Timezone: %v", location)
	sctx.Infof("Current Time: %v", time.Now().In(location))

	dbm, err := db_manager.NewDbManager(sctx)
	if err != nil {
		return fmt.Errorf("error connecting to main (regular) database: %v", err)
	}
	sctx = sctx.WithDbManager(dbm).WithDB(dbm.GetGORM())

	// rcm := redis_cache_manager.NewRedisCacheManager(sctx)
	// if rcm != nil {
	// 	sctx.Infof("RedisCacheManager started")
	// 	sctx = sctx.WithRedisCacheManager(rcm)
	// }

	// custom init
	err = initFunc(sctx)
	if err != nil {
		return err
	}

	sctx.Infof("%s: Starting", prefix)

	err = startFunc(sctx)
	if err != nil {
		return err
	}
	sctx.Infof("%s: Started. Working...", prefix)

	defer func() {
		err := closeFunc(sctx)
		if err != nil {
			sctx.Errorf("Error closing service: %v", err)
		}
		sctx.Infof("%s: Closed", prefix)
	}()

	// тут мы зависнем до получения сигнала на завершение
	osSignal := shutdown.WaitForSignalToShutdown()
	// отменяем контекст и ждем завершения всех запросов
	sctx.Infof("%s: Received signal '%s'. Cancelling context", prefix, osSignal.String())
	cancel()
	sctx.Infof("%s: Context cancelled - no new requests will be served. Waiting for all prevously started requests to finish", prefix)
	// с этого момента все новые запросы будут получать ошибку
	wg.Wait() // тут можем ждать долго - до 2х часов - пока запросы все завершатся. новые запросы уже не будут приниматься! ни по rest ни по ws. по rest их блокирует - WithWaitGroup где проверяем serverCtx.Err(), и по ws проверяем serverCtx.Err()
	sctx.Infof("%s: All requests finished. Closing", prefix)
	return nil
}
