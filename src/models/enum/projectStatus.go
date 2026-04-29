package enumModels

type ProjectStatus string

const (
	Exited   ProjectStatus = "exited"
	Running  ProjectStatus = "running"
	Dead     ProjectStatus = "dead"
	Errored  ProjectStatus = "errored"
	Stopping ProjectStatus = "stopping"
	Deleted  ProjectStatus = "deleted"
)
