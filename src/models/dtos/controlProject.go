package modelsDTOs

type ActionOnProject struct {
	ExternalID string `json:"project_id" binding:"required"`
	Action     string `json:"action" binding:"required"`
}
