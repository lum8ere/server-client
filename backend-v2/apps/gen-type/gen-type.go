package main

import (
	"backed-api-v2/libs/3_infrastructure/db_manager"
	"backed-api-v2/libs/4_common/env_vars"
	"backed-api-v2/libs/4_common/smart_context"
	"os"

	"gorm.io/gen"
)

// Dynamic SQL
type Querier interface {
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	env_vars.LoadEnvVars()
	os.Setenv("LOG_LEVEL", "info")
	logger := smart_context.NewSmartContext()

	dbm, err := db_manager.NewDbManager(logger)
	if err != nil {
		logger.Fatalf("NewDbManager failed: %v", err)
	}

	logger = logger.WithDB(dbm.GetGORM())

	g := gen.NewGenerator(gen.Config{
		OutPath: "./libs/2_generated_models/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(logger.GetDB())

	g.ApplyBasic(
		g.GenerateAllTable()...,
	)

	g.Execute()
}
