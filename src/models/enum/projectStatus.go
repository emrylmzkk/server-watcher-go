package enumModels

type ProjectStatus string

const (
	Exited  ProjectStatus = "exited"
	Running ProjectStatus = "running"
	Dead    ProjectStatus = "dead"
)
