package modelsDTOs

type DockerContainerResponseDTO struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
	Uptime    string `json:"uptime"`
}
