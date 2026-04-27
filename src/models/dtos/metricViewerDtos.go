package modelsDTOs

import "time"

type ContainerLogRequestDTO struct {
	ContainerName string `json:"container_name"`
	HostID        string `json:"host_id"`
	Limit         int    `json:"limit"`
}

type LogResponseDTO struct {
	Time          time.Time `json:"time"`
	HostID        string    `json:"host_id"`
	ContainerName string    `json:"container_name"`
	Messages      []string  `json:"message"`
}
