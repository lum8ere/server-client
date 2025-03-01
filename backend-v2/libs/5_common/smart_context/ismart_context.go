package smart_context

import (
	"backed-api-v2/libs/5_common/types"
	"context"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ISmartContext interface {
	// Методы логирования
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	// Методы управления полями логирования
	LogField(key string, value interface{}) ISmartContext  // будут заполнять поля в логгере
	LogFields(fields types.Fields) ISmartContext           // будут заполнять поля в логгере
	WithField(key string, value interface{}) ISmartContext // будут заполнять поля в DataField
	WithFields(fields types.Fields) ISmartContext          // будут заполнять поля в DataField

	GetLogger() *zap.Logger

	GetDataFields() types.Fields

	WithDbManager(db IDbManager) ISmartContext
	GetDbManager() IDbManager

	WithDB(db *gorm.DB) ISmartContext
	GetDB() *gorm.DB

	WithWaitGroup(wg *sync.WaitGroup) ISmartContext
	GetWaitGroup() *sync.WaitGroup

	// Метод для получения стандартного context.Context
	WithContext(ctx context.Context) ISmartContext
	GetContext() context.Context

	WithSessionId(session string) ISmartContext
	GetSessionId() string

	// WithMinioManager(minioManager IMinioManager) ISmartContext
	// GetMinioManager() IMinioManager
}
