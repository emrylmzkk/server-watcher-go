package enumModels

type ProjectStatus string

const (
	Exited  ProjectStatus = "exited"
	Running ProjectStatus = "closed"
	Dead    ProjectStatus = "dead"
)
