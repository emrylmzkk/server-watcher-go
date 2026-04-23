package modelsDTOs

type ActionOnProject struct {
	ExternalID string `json:"project_id" binding:"required"`
	// Type       enumModels.ProcessType `json:"project_type" binding:"required"`
	Action string `json:"action" binding:"required"`
}
