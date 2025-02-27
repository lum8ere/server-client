package service

import (
	"log"
	"sync"
	"sync/atomic"
)

type Service struct {
	clients       map[string]*Client
	logger        *log.Logger
	m             sync.RWMutex
	clientCounter atomic.Int64

	uploadPath string
}

func New(logger *log.Logger, uploadPath string) *Service {
	return &Service{
		clients:       make(map[string]*Client),
		logger:        logger,
		m:             sync.RWMutex{},
		clientCounter: atomic.Int64{},

		uploadPath: uploadPath,
	}
}
