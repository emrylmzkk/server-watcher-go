package models

import (
	"sync"

	modelsDTOs "server-watcher-app/src/models/dtos"
)

type DockerCache struct {
	Data []modelsDTOs.DockerStats
	Mu   sync.RWMutex
}
