package service

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type Service struct {
	clients       map[string]*Client
	logger        *log.Logger
	m             sync.RWMutex
	clientCounter atomic.Int64
    audioSources   map[string]*websocket.Conn
    audioListeners map[string][]*websocket.Conn
	uploadPath string
}

func New(logger *log.Logger, uploadPath string) *Service {
	return &Service{
		clients:       make(map[string]*Client),
		logger:        logger,
		m:             sync.RWMutex{},
		clientCounter: atomic.Int64{},
        audioSources:   make(map[string]*websocket.Conn),
        audioListeners: make(map[string][]*websocket.Conn),
		uploadPath: uploadPath,
	}
}
